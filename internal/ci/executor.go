package ci

import (
	"context"
	"fmt"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load() // Load .env once at startup
}

type Executor struct {
	ws Workspace
}

type Workspace interface {
	Branch() string
	Commit() string
	Dir() string
	Env() []string
	//	LoadPipeline() (*Pipeline, error)
	ExecuteCommand(ctx context.Context, cmd string, args []string) ([]byte, error)
}

func NewExecutor(ws Workspace) *Executor {
	return &Executor{
		ws: ws,
	}
}

/*
func (e *Executor) RunDefault(ctx context.Context) (string, error) {
	pipeline, err := e.ws.LoadPipeline()
	if err != nil {
		return "", err
	}
	return e.Run(ctx, pipeline)
} */

func (e *Executor) Run(ctx context.Context, pipeline *Pipeline) (string, error) {
	output := strings.Builder{}
	for _, step := range pipeline.Steps {
		fmt.Println("Running step: ")
		fmt.Println(step.Name)
		fmt.Println("\n")
		cmd := strings.Fields(step.Commands[0])[0]
		args := strings.Fields(step.Commands[0])[1:]
		_, err := e.ws.ExecuteCommand(ctx, cmd, args)
		if err != nil {
			return "",
				fmt.Errorf("error executing command %s %s: %w", cmd, strings.Join(args, " "), err)
		}
	}
	return output.String(), nil
}
