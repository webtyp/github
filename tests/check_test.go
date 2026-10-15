package tests

import (
	"errors"
	"testing"
	"webtyp.com/deploy"
	"webtyp.com/github"
)

func TestCheckG1(t *testing.T) {
	r := newFakeRunner()
	ci, _ := github.New(github.Config{RootDir: "/dir", Runner: r})
	p := deploy.Plan{}
	job := deploy.Job{}

	r.err["gh --version"] = errors.New("fail")
	missing, err := ci.Check(p, job)
	if err != nil || len(missing) != 1 || missing[0].Item != "GitHub CLI" {
		t.Errorf("expected G1 missing, got %v, %v", missing, err)
	}
}

func TestCheckG2(t *testing.T) {
	r := newFakeRunner()
	ci, _ := github.New(github.Config{RootDir: "/dir", Runner: r})
	p := deploy.Plan{}
	job := deploy.Job{}

	r.out["gh --version"] = []byte("gh version")
	r.err["gh auth status"] = errors.New("fail")
	missing, err := ci.Check(p, job)
	if err != nil || len(missing) != 1 || missing[0].Item != "GitHub CLI auth" {
		t.Errorf("expected G2 missing, got %v, %v", missing, err)
	}
}

func TestCheckG3G4G5(t *testing.T) {
	r := newFakeRunner()
	ci, _ := github.New(github.Config{RootDir: "/dir", Runner: r})
	p := deploy.Plan{}
	job := deploy.Job{}

	r.out["gh --version"] = []byte("gh version")
	r.out["gh auth status"] = []byte("logged in")
	r.err["gh repo view --json nameWithOwner -q .nameWithOwner"] = errors.New("fail")

	missing, err := ci.Check(p, job)
	if err != nil {
		t.Fatal(err)
	}

	hasG3 := false
	for _, m := range missing {
		if m.Item == "GitHub repository" {
			hasG3 = true
		}
	}
	if !hasG3 {
		t.Errorf("expected G3 missing, got %v", missing)
	}
}
