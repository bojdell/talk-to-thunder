package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/samsarahq/go/oops"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/graphiql"
	"github.com/samsarahq/thunder/graphql/introspection"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"github.com/samsarahq/thunder/livesql"
	"github.com/samsarahq/thunder/sqlgen"
)

type SnippetState int

const (
	DbName = "talktothunder"

	SnippetStateInvalid    SnippetState = 0
	SnippetStateInProgress              = 1
	SnippetStateCompleted               = 2

	DefaultNumTokensToGenerate = 100
)

type Server struct {
	db *livesql.LiveDB
}

type Snippet struct {
	Id            int64 `sql:",primary" graphql:",key"`
	CreatedAt     time.Time
	State2        SnippetState
	SeedText      string
	GeneratedText string
}

func (s *Server) registerQueryRoot(schema *schemabuilder.Schema) {
	object := schema.Query()

	// TODO: why doesn't this work??
	object.FieldFunc("currentSnippet", func(ctx context.Context) (*Snippet, error) {
		return s.getCurrentSnippet(ctx)
	})

	object.FieldFunc("allSnippets", func(ctx context.Context) ([]*Snippet, error) {
		var rows []*Snippet
		return rows, s.db.Query(ctx, &rows, sqlgen.Filter{}, nil)
	})

	object.FieldFunc("snippet", func(ctx context.Context, args struct{ Id int64 }) (*Snippet, error) {
		var snippet *Snippet
		if err := s.db.QueryRow(ctx, &snippet, sqlgen.Filter{"id": args.Id}, nil); err != nil && err != sql.ErrNoRows {
			return nil, oops.Wrapf(err, "")
		}

		return snippet, nil
	})
}

func (s *Server) registerMutationRoot(schema *schemabuilder.Schema) {
	object := schema.Mutation()
	object.FieldFunc("createSnippet", func(ctx context.Context, args struct{ Text string }) (int64, error) {
		return s.createSnippet(ctx, args.Text, DefaultNumTokensToGenerate)
	})
}

func (s *Server) SchemaBuilderSchema() *schemabuilder.Schema {
	schema := schemabuilder.NewSchema()

	s.registerQueryRoot(schema)
	s.registerMutationRoot(schema)

	return schema
}

func (s *Server) Schema() *graphql.Schema {
	return s.SchemaBuilderSchema().MustBuild()
}

type executionLogger struct{}

func (e *executionLogger) StartExecution(ctx context.Context, tags map[string]string, initial bool) {}
func (e *executionLogger) FinishExecution(ctx context.Context, tags map[string]string, delay time.Duration) {
}
func (e *executionLogger) Error(ctx context.Context, err error, tags map[string]string) {
	log.Printf("error:%v\n%s", tags, err)
}

type subscriptionLogger struct {
	server *Server
}

func panicIfErr(err error) {
	if err == nil {
		return
	}
	panic(err.Error())
}

func (l *subscriptionLogger) Subscribe(ctx context.Context, id string, tags map[string]string) {
	intId, err := strconv.ParseInt(id, 10, 64)
	panicIfErr(err)

	log.Println("~~ Subscribe", intId)
}

func (l *subscriptionLogger) Unsubscribe(ctx context.Context, id string) {
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		panic("error parsing subscription id")
	}
	log.Println("~~ Unsubscribe", intId)
}

func (s *Server) getCurrentSnippet(ctx context.Context) (*Snippet, error) {
	var snippet *Snippet
	if err := s.db.QueryRow(ctx, &snippet, sqlgen.Filter{"state2": int64(SnippetStateInProgress)}, nil); err != nil {
		if err != sql.ErrNoRows {
			return nil, oops.Wrapf(err, "")
		}
		return nil, nil
	}

	return snippet, nil
}

func (s *Server) createSnippet(ctx context.Context, seedText string, numTokensToGenerate int) (int64, error) {
	row := &Snippet{
		State2:   SnippetStateInProgress,
		SeedText: seedText,
	}

	response, err := s.db.InsertRow(ctx, row)
	if err != nil {
		return 0, oops.Wrapf(err, "")
	}

	fileBytes := []byte(fmt.Sprintf("%d %s", numTokensToGenerate, seedText))
	if err := ioutil.WriteFile(InputFileName, fileBytes, 0644); err != nil {
		return 0, oops.Wrapf(err, "")
	}
	log.Println("wrote input file", InputFileName, string(fileBytes))

	return response.LastInsertId()
}

const (
	// TODO: fix these
	InputFileName  = "/Users/bo/dev/gpt-2/input/input.txt"
	OutputFileName = "/Users/bo/dev/gpt-2/output/output.txt"
)

func (s *Server) pollCurrentSnippetLoop(ctx context.Context) error {
	var currentSnippet *Snippet
	var err error
	for {
		currentSnippet, err = s.updateCurrentSnippet(ctx, currentSnippet)
		if err != nil {
			return oops.Wrapf(err, "")
		}
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (s *Server) updateCurrentSnippet(ctx context.Context, currentSnippet *Snippet) (*Snippet, error) {
	_, err := os.Stat(InputFileName)
	if err != nil && os.IsNotExist(err) {
		if currentSnippet != nil {
			currentSnippet.State2 = SnippetStateCompleted
			if err := s.db.UpdateRow(ctx, currentSnippet); err != nil {
				return nil, oops.Wrapf(err, "")
			}
		}
		return nil, nil
	}
	if err != nil {
		return nil, oops.Wrapf(err, "")
	}
	_, err = os.Stat(OutputFileName)
	if err != nil && os.IsNotExist(err) {
		return currentSnippet, nil
	}
	if err != nil {
		return nil, oops.Wrapf(err, "")
	}

	currentSnippet, err = s.getCurrentSnippet(ctx)
	if err != nil {
		return nil, oops.Wrapf(err, "")
	}
	if currentSnippet == nil {
		return nil, nil
	}

	// TODO: perf optimization, only read diff since file is append only
	file, err := os.Open(OutputFileName)
	if err != nil {
		return nil, oops.Wrapf(err, "")
	}
	defer file.Close()

	text, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, oops.Wrapf(err, "")
	}
	currentSnippet.GeneratedText = string(text)

	if err := s.db.UpdateRow(ctx, currentSnippet); err != nil {
		return nil, oops.Wrapf(err, "")
	}

	return currentSnippet, nil
}

func main() {
	sqlgenSchema := sqlgen.NewSchema()
	sqlgenSchema.MustRegisterType("snippets", sqlgen.AutoIncrement, Snippet{})

	db, err := livesql.Open("localhost", 3307, "root", "", DbName, sqlgenSchema)
	panicIfErr(err)

	server := &Server{
		db: db,
	}

	ctx := context.Background()
	go func() {
		panicIfErr(server.pollCurrentSnippetLoop(ctx))
		return
	}()

	// _, err = server.createSnippet(ctx, "The best IOT company is", DefaultNumTokensToGenerate)
	// panicIfErr(err)

	graphqlSchema := server.Schema()
	introspection.AddIntrospectionToSchema(graphqlSchema)

	s, err := server.getCurrentSnippet(ctx)
	log.Println(s, err)

	http.Handle("/graphql", graphql.Handler(graphqlSchema))
	http.Handle("/graphiql/", http.StripPrefix("/graphiql/", graphiql.Handler()))
	log.Println("== STARTED ==")
	panicIfErr(http.ListenAndServe(":3030", nil))
}
