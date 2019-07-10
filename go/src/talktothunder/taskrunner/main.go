package main

import (
	"context"
	"fmt"
	"os"

	"github.com/samsarahq/taskrunner"
	"github.com/samsarahq/taskrunner/cache"
	"github.com/samsarahq/taskrunner/clireporter"
	"github.com/samsarahq/taskrunner/goextensions"
	"github.com/samsarahq/taskrunner/shell"
	"mvdan.cc/sh/interp"
)

var builder = goextensions.NewGoBuilder()
var cacher = cache.New(shellEnv)

func main() {
	addTasks()

	defaultOptions := []taskrunner.RunOption{
		taskrunner.ExecutorOptions(taskrunner.ShellRunOptions(shellEnv)),
		cacher.Option,
	}
	if os.Getenv("EXPERIMENTAL_TASKRUNNER_TUI") != "" {
		// defaultOptions = append(defaultOptions, tui.Option)
	} else {
		builder.LogToStdout = true
		defaultOptions = append(defaultOptions, clireporter.StdoutOption)
	}
	taskrunner.Run(defaultOptions...)
}

func shellEnv(r *interp.Runner) {
	path, _ := r.Env.Get("PATH")
	gopath, _ := r.Env.Get("GOPATH")
	r.Env.Set("PATH", fmt.Sprintf("%s:%s:%s", path, fmt.Sprintf("%s/bin", gopath), "./node_modules/.bin/"))
}

func addTasks() {
	build := taskrunner.Add(&taskrunner.Task{
		Name: "gqlserver/build",
	}, builder.WrapWithGoBuild("go/src/talktothunder/gqlserver"))

	taskrunner.Add(&taskrunner.Task{
		Name:         "gqlserver",
		Dependencies: []*taskrunner.Task{build},
		Run: func(ctx context.Context, shellRun shell.ShellRun) error {
			return shellRun(ctx, "gqlserver")
		},
	})
}
