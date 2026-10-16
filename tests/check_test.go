package github_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"webtyp.com/deploy"
	"webtyp.com/github"
)

func samplePlanAndJob() (deploy.Plan, deploy.Job) {
	return deploy.Plan{Name: "misitio"}, deploy.Job{
		Uses:    "some-action@v1",
		With:    []deploy.Input{{Name: "k", Value: "v"}},
		Secrets: []string{"SEC1"},
	}
}

func TestCheckG1(t *testing.T) {
	r := newFakeRunner()
	ci, _ := github.New(github.Config{RootDir: t.TempDir(), Runner: r})
	p, job := samplePlanAndJob()

	r.err["gh --version"] = errors.New("fail")
	missing, err := ci.Check(p, job)
	if err != nil || len(missing) != 1 || missing[0].Item != "GitHub CLI" || missing[0].Fix != deploy.FixManual {
		t.Errorf("expected only G1 missing, got %v, %v", missing, err)
	}
	if len(r.cmds) != 1 {
		t.Errorf("expected no further probes after G1 failure, got %d cmds", len(r.cmds))
	}
}

func TestCheckG2(t *testing.T) {
	r := newFakeRunner()
	ci, _ := github.New(github.Config{RootDir: t.TempDir(), Runner: r})
	p, job := samplePlanAndJob()

	r.out["gh --version"] = []byte("gh version")
	r.err["gh auth status"] = errors.New("fail")
	missing, err := ci.Check(p, job)
	if err != nil || len(missing) != 1 || missing[0].Item != "GitHub CLI auth" || missing[0].Fix != deploy.FixManual {
		t.Errorf("expected only G2 missing, got %v, %v", missing, err)
	}
	if len(r.cmds) != 2 {
		t.Errorf("expected no further probes after G2 failure, got %d cmds", len(r.cmds))
	}
}

func TestCheckG3(t *testing.T) {
	tmpDir := t.TempDir()
	r := newFakeRunner()
	ci, _ := github.New(github.Config{RootDir: tmpDir, Runner: r})
	p, job := samplePlanAndJob()

	if err := ci.Provision(p, job); err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(filepath.Join(tmpDir, ".github", "workflows", "deploy.yml"))
	if err != nil {
		t.Fatal(err)
	}

	r.out["gh --version"] = []byte("gh version")
	r.out["gh auth status"] = []byte("logged in")
	r.err["gh repo view --json nameWithOwner -q .nameWithOwner"] = errors.New("no remote")
	r.out["git show origin/main:.github/workflows/deploy.yml"] = expected

	missing, err := ci.Check(p, job)
	if err != nil {
		t.Fatal(err)
	}
	if len(missing) != 1 || missing[0].Item != "GitHub repository" || missing[0].Fix != deploy.FixManual {
		t.Errorf("expected only G3 missing, got %v", missing)
	}
}

func TestCheckG4(t *testing.T) {
	tmpDir := t.TempDir()
	r := newFakeRunner()
	ci, _ := github.New(github.Config{RootDir: tmpDir, Runner: r})
	p, job := samplePlanAndJob()

	expected, err := os.ReadFile("testdata/deploy.yml")
	if err != nil {
		t.Fatal(err)
	}

	r.out["gh --version"] = []byte("gh version")
	r.out["gh auth status"] = []byte("logged in")
	r.out["gh repo view --json nameWithOwner -q .nameWithOwner"] = []byte("org/repo\n")
	r.out["git show origin/main:.github/workflows/deploy.yml"] = expected

	// Absent workflow file → G4 missing (FixAutomatic)
	missing, err := ci.Check(p, job)
	if err != nil {
		t.Fatal(err)
	}
	if len(missing) != 1 || missing[0].Item != "Deploy workflow" || missing[0].Fix != deploy.FixAutomatic {
		t.Errorf("expected only G4 missing when file absent, got %v", missing)
	}

	// Outdated local workflow file → G4 missing (FixAutomatic)
	wfPath := filepath.Join(tmpDir, ".github", "workflows", "deploy.yml")
	if err := os.MkdirAll(filepath.Dir(wfPath), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(wfPath, []byte("outdated"), 0644); err != nil {
		t.Fatal(err)
	}
	missing, err = ci.Check(p, job)
	if err != nil {
		t.Fatal(err)
	}
	if len(missing) != 1 || missing[0].Item != "Deploy workflow" || missing[0].Fix != deploy.FixAutomatic {
		t.Errorf("expected only G4 missing when file differs, got %v", missing)
	}
}

func TestCheckG5(t *testing.T) {
	tmpDir := t.TempDir()
	r := newFakeRunner()
	ci, _ := github.New(github.Config{RootDir: tmpDir, Runner: r})
	p, job := samplePlanAndJob()

	if err := ci.Provision(p, job); err != nil {
		t.Fatal(err)
	}
	expected, err := os.ReadFile(filepath.Join(tmpDir, ".github", "workflows", "deploy.yml"))
	if err != nil {
		t.Fatal(err)
	}

	r.out["gh --version"] = []byte("gh version")
	r.out["gh auth status"] = []byte("logged in")
	r.out["gh repo view --json nameWithOwner -q .nameWithOwner"] = []byte("org/repo\n")
	r.out["git show origin/main:.github/workflows/deploy.yml"] = []byte("old remote content")

	missing, err := ci.Check(p, job)
	if err != nil {
		t.Fatal(err)
	}
	if len(missing) != 1 || missing[0].Item != "Deploy workflow outdated" || missing[0].Fix != deploy.FixManual {
		t.Errorf("expected only G5 missing, got %v", missing)
	}

	// When origin/main matches, nothing is missing
	r.out["git show origin/main:.github/workflows/deploy.yml"] = expected
	missing, err = ci.Check(p, job)
	if err != nil {
		t.Fatal(err)
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing items when all probes pass, got %v", missing)
	}
}

