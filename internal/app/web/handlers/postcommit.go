package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/flow-ci/flow-ci/internal/ci"
)

func PostCheckItWorks(rb RequestBody, primaryLang string) error {
	// Create workspace - SSH auth is handled by ci.NewWorkspaceFromGit
	fmt.Println("Creating workspace")
	var ws ci.Workspace
	ws, err := ci.NewWorkspaceFromGit("./tmp", rb.Url, rb.Branch)
	if err != nil {
		return fmt.Errorf("failed to create workspace: %v", err)
	}

	// Run the pipeline with language-specific configuration
	executor := ci.NewExecutor(ws)
	ctx := context.Background()
	var output string
	switch strings.ToLower(primaryLang) {
	case "golang":
		fmt.Println("Running goLang pipeline")
		output, err = executor.Run(ctx, &ci.Pipeline{
			Steps: []ci.Step{{Name: "BUILD", Commands: []string{"go", "build", "./..."}}},
		})
	case "python":
		fmt.Println("Running python pipeline")
		output, err = executor.Run(ctx, &ci.Pipeline{
			Steps: []ci.Step{{Name: "BUILD", Commands: []string{"python", "-m", "compileall", "."}}},
		})
	case "javascript", "typescript":
		fmt.Println("Running NodeJS pipeline")

		output, err = executor.Run(ctx, &ci.Pipeline{
			Steps: []ci.Step{{Name: "BUILD", Commands: []string{"npm", "install"}}},
		})
	default:
		fmt.Println("no pipeline found for primary language")
		//output, err = executor.RunDefault(ctx)
	}
	if err != nil {
		return fmt.Errorf("execution failed for %s: %s - %v", primaryLang, output, err)
	}

	// Print the result
	result := fmt.Sprintf(
		"Successfully executed pipeline for %s.\n%s\n\nFrom branch: %s\nCommit: %s\nIn directory: %s\n",
		primaryLang,
		output,
		ws.Branch(),
		ws.Commit(),
		ws.Dir(),
	)
	fmt.Printf(result)
	return nil
}
