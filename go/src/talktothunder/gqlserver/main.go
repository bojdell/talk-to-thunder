package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

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

	DefaultNumTokensToGenerate = 50
)

type Server struct {
	db *livesql.LiveDB
}

type Snippet struct {
	Id            int64 `sql:",primary" graphql:",key"`
	CreatedAt     time.Time
	DeletedAt     *time.Time
	State         SnippetState
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
		err := s.db.Query(ctx, &rows, sqlgen.Filter{}, &sqlgen.SelectOptions{
			Where:   "deleted_at IS NULL",
			OrderBy: "id DESC",
		})
		return rows, err
	})

	object.FieldFunc("deletedSnippets", func(ctx context.Context) ([]*Snippet, error) {
		var rows []*Snippet
		err := s.db.Query(ctx, &rows, sqlgen.Filter{}, &sqlgen.SelectOptions{
			Where:   "deleted_at IS NOT NULL",
			OrderBy: "deleted_at DESC",
		})
		return rows, err
	})

	object.FieldFunc("snippet", func(ctx context.Context, args struct{ Id int64 }) (*Snippet, error) {
		var snippet *Snippet
		if err := s.db.QueryRow(ctx, &snippet, sqlgen.Filter{"id": args.Id}, nil); err != nil {
			if err != sql.ErrNoRows {
				return nil, oops.Wrapf(err, "")
			}
			return nil, nil
		}

		return snippet, nil
	})
}

func int64OrElse(a *int64, b int64) int64 {
	if a == nil {
		return b
	}
	return *a
}

func (s *Server) registerMutationRoot(schema *schemabuilder.Schema) {
	object := schema.Mutation()

	// Challenge 3a: add a mutation that will cause your computer to read text out loud.
	// The text should be provided as a GraphQL argument.
	// hint: try pasting this into your terminal: say hello world!
	// https://www.idownloadblog.com/2016/02/24/make-mac-talk-terminal/

	object.FieldFunc("createSnippet", func(ctx context.Context, args struct {
		Text             string
		FinalTokenLength *int64
	}) (int64, error) {
		currentSnippet, err := s.getCurrentSnippet(ctx)
		if err != nil {
			return -1, oops.Wrapf(err, "")
		}
		if currentSnippet != nil {
			return -1, nil
		}
		return s.createSnippet(ctx, args.Text, int64OrElse(args.FinalTokenLength, DefaultNumTokensToGenerate))
	})
	object.FieldFunc("deleteSnippet", func(ctx context.Context, args struct {
		Id int64
	}) error {
		ctx, tx, err := s.db.WithTx(ctx)
		if err != nil {
			return oops.Wrapf(err, "")
		}
		defer tx.Rollback()

		var row *Snippet
		if err := s.db.QueryRow(ctx, &row, sqlgen.Filter{"id": args.Id}, &sqlgen.SelectOptions{
			Where: "deleted_at IS NULL",
		}); err != nil {
			if err != sql.ErrNoRows {
				return oops.Wrapf(err, "")
			}
			return nil
		}

		timeNow := time.Now()
		row.DeletedAt = &timeNow

		if err := s.db.UpdateRow(ctx, row); err != nil {
			return oops.Wrapf(err, "")
		}

		if err := tx.Commit(); err != nil {
			return oops.Wrapf(err, "")
		}

		return nil
	})

	object.FieldFunc("undeleteSnippet", func(ctx context.Context, args struct {
		Id int64
	}) error {
		ctx, tx, err := s.db.WithTx(ctx)
		if err != nil {
			return oops.Wrapf(err, "")
		}
		defer tx.Rollback()

		var row *Snippet
		if err := s.db.QueryRow(ctx, &row, sqlgen.Filter{"id": args.Id}, &sqlgen.SelectOptions{
			Where: "deleted_at IS NOT NULL",
		}); err != nil {
			if err != sql.ErrNoRows {
				return oops.Wrapf(err, "")
			}
			return nil
		}

		row.DeletedAt = nil

		if err := s.db.UpdateRow(ctx, row); err != nil {
			return oops.Wrapf(err, "")
		}

		if err := tx.Commit(); err != nil {
			return oops.Wrapf(err, "")
		}

		return nil
	})

	object.FieldFunc("emptyTrash", func(ctx context.Context) error {
		if _, err := s.db.QueryExecer(ctx).ExecContext(ctx, "DELETE FROM snippets WHERE deleted_at IS NOT NULL"); err != nil && err != sql.ErrNoRows {
			return oops.Wrapf(err, "")
		}
		return nil
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

func panicIfErr(err error) {
	if err == nil {
		return
	}
	panic(err.Error())
}

func (s *Server) getCurrentSnippet(ctx context.Context) (*Snippet, error) {
	var snippets []*Snippet
	if err := s.db.Query(ctx, &snippets, sqlgen.Filter{"state": SnippetStateInProgress}, nil); err != nil {
		if err != sql.ErrNoRows {
			return nil, oops.Wrapf(err, "")
		}
		return nil, nil
	}

	if len(snippets) == 0 {
		return nil, nil
	}

	return snippets[0], nil
}

func (s *Server) createSnippet(ctx context.Context, seedText string, numTokensToGenerate int64) (int64, error) {
	row := &Snippet{
		State:    SnippetStateInProgress,
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
	// TODO: config file?
	InputFileName  = "/tmp/talktothunder/input.txt"
	OutputFileName = "/tmp/talktothunder/output.txt"
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
			currentSnippet.State = SnippetStateCompleted
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

	graphqlSchema := server.Schema()
	introspection.AddIntrospectionToSchema(graphqlSchema)

	http.Handle("/graphql", graphql.Handler(graphqlSchema))
	http.Handle("/graphiql/", http.StripPrefix("/graphiql/", graphiql.Handler()))
	log.Println("== STARTED ==")
	panicIfErr(http.ListenAndServe(":3030", nil))
}
