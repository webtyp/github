package github

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"webtyp.com/deploy"
)

// Actions is the CI id a project declares in config/deploy.go.
const Actions deploy.CIID = "webtyp.com/github.Actions"

// Runner runs an external command in dir, writing stdin to it, and returns
// its stdout. Secret values travel only through stdin.
type Runner interface {
	Run(dir, stdin, name string, args ...string) ([]byte, error)
}

// ExecRunner runs commands with os/exec.
type ExecRunner struct{}

func (e ExecRunner) Run(dir, stdin, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		if stderr.Len() > 0 {
			return nil, errors.New(err.Error() + ": " + strings.TrimSpace(stderr.String()))
		}
		return nil, err
	}
	return out, nil
}

// Config configures a CI. Both fields are required.
type Config struct {
	RootDir string // the project root (contains .git and go.mod)
	Runner  Runner
}

// CI implements deploy.CI for GitHub Actions.
type CI struct {
	rootDir string
	runner  Runner
	repo    string // cached repository name
}

// New returns a CI, or an error if a Config field is missing.
func New(cfg Config) (*CI, error) {
	if cfg.RootDir == "" {
		return nil, ErrNoRootDir
	}
	if cfg.Runner == nil {
		return nil, ErrNoRunner
	}
	return &CI{
		rootDir: cfg.RootDir,
		runner:  cfg.Runner,
	}, nil
}

func (c *CI) ID() deploy.CIID {
	return Actions
}




func (c *CI) Provision(p deploy.Plan, job deploy.Job) error {
	path := filepath.Join(c.rootDir, ".github", "workflows", "deploy.yml")
	content := render(p, job)

	existing, err := os.ReadFile(path)
	if err == nil && bytes.Equal(existing, content) {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0644)
}

var _ deploy.CI = (*CI)(nil)
