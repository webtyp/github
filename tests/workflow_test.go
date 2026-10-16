package github_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"webtyp.com/deploy"
	"webtyp.com/github"
)

func TestProvision(t *testing.T) {
	tmpDir := t.TempDir()
	r := newFakeRunner()
	ci, _ := github.New(github.Config{RootDir: tmpDir, Runner: r})

	p := deploy.Plan{Name: "misitio"}
	job := deploy.Job{
		Uses:    "some-action@v1",
		With:    []deploy.Input{{Name: "k", Value: "v"}},
		Secrets: []string{"SEC1"},
	}

	err := ci.Provision(p, job)
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(tmpDir, ".github", "workflows", "deploy.yml")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile("testdata/deploy.yml")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(content, expected) {
		t.Errorf("content does not match testdata. Got:\n%s", content)
	}

	info1, _ := os.Stat(path)

	// W2: Provision does not rewrite an identical file
	err = ci.Provision(p, job)
	if err != nil {
		t.Fatal(err)
	}
	info2, _ := os.Stat(path)

	if info1.ModTime() != info2.ModTime() {
		t.Errorf("expected mtime to be unchanged")
	}
}
