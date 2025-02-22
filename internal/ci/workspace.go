package ci

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

// Load SSH key

func NewWorkspaceFromGit(root string, url string, branch string) (*workspaceImpl, error) {
	sshKeyPath := os.Getenv("SSH_KEY_PATH")
	if sshKeyPath == "" {
		fmt.Println("check .env dumbass sshkey got the problem")
	}
	sshAuth, err := ssh.NewPublicKeysFromFile("git", sshKeyPath, "")
	if err != nil {
		fmt.Println("Error loading SSH key:", err)
		return nil, fmt.Errorf("failed to load SSH key: %w", err)
	}

	dir, err := os.MkdirTemp(root, "workspace")
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil, fmt.Errorf("workspace root directory %s does not exist", root)
	}

	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:               url,
		ReferenceName:     plumbing.NewBranchReferenceName(branch),
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Depth:             1,
		Auth:              sshAuth,
	})
	if err != nil {
		return nil, err
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	return &workspaceImpl{
		dir:    dir,
		branch: branch,
		commit: ref.Hash().String(),
		env:    []string{},
	}, nil
}

func NewWorkspaceFromDir(dir string) (*workspaceImpl, error) {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil, err
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	return &workspaceImpl{
		dir:    dir,
		branch: ref.Name().Short(),
		commit: ref.Hash().String(),
		env:    []string{},
	}, nil
}

type workspaceImpl struct {
	branch      string
	commit      string
	dir         string
	env         []string
	nodeChecked bool
}

func (ws *workspaceImpl) Branch() string {
	return ws.branch
}

func (ws *workspaceImpl) Commit() string {
	return ws.commit
}

func (ws *workspaceImpl) Dir() string {
	return ws.dir
}

func (ws *workspaceImpl) Env() []string {
	return ws.env
}

/*
func (ws *workspaceImpl) LoadPipeline() (*Pipeline, error) {
	path := filepath.Join(ws.dir, "flow-ci.yaml")
	fmt.Println("Loading pipeline from:", path)

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading pipeline file:", err)
		return nil, err
	}

	fmt.Println("Pipeline file contents:\n", string(data))

	var pipeline Pipeline
	err = yaml.Unmarshal(data, &pipeline)
	if err != nil {
		fmt.Println("Error parsing YAML:", err)
		return nil, err
	}

	fmt.Println("Pipeline loaded successfully:", pipeline.Name)
	return &pipeline, nil
}*/

func (ws *workspaceImpl) ExecuteCommand(ctx context.Context, cmd string, args []string) ([]byte, error) {
	fmt.Println("Executing:", cmd, args) // Log command

	// Prepare command
	command := exec.CommandContext(ctx, cmd, args...)
	command.Dir = ws.dir
	command.Env = append(command.Environ(), ws.Env()...)
	// ✅ Check Node.js installation for npm/yarn commands
	if !ws.nodeChecked {
		fmt.Println("Checking Framework installation...")
		checkNode := exec.Command(cmd, "-v")
		if _, err := checkNode.Output(); err != nil {
			return nil, fmt.Errorf(cmd + "is not installed, please install it first")
		}
		ws.nodeChecked = true
	}

	// 💡 Execute command and capture combined stdout/stderr
	output, err := command.CombinedOutput()
	fmt.Println("Command Output:", string(output)) // Log command output

	// 🔥 Error handling
	if err != nil {
		fmt.Println("Execution Error:", err)
		return output, err
	}

	return output, nil
}
