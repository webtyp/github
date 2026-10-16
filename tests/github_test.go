package github_test

import (
	"errors"
	"testing"

	"webtyp.com/deploy"
	"webtyp.com/github"
	"webtyp.com/lang"
)

func TestNew(t *testing.T) {
	_, err := github.New(github.Config{RootDir: "", Runner: newFakeRunner()})
	if err != github.ErrNoRootDir {
		t.Errorf("expected ErrNoRootDir, got %v", err)
	}

	_, err = github.New(github.Config{RootDir: "/dir", Runner: nil})
	if err != github.ErrNoRunner {
		t.Errorf("expected ErrNoRunner, got %v", err)
	}

	ci, err := github.New(github.Config{RootDir: "/dir", Runner: newFakeRunner()})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if ci.ID() != github.Actions {
		t.Errorf("expected Actions, got %v", ci.ID())
	}
}

func TestLang(t *testing.T) {
	lang.Load(github.Dictionary())
	lang.OutLang("es")
	t.Cleanup(func() { lang.OutLang("en") })

	r := newFakeRunner()
	r.err["gh --version"] = errors.New("not found")
	ci, err := github.New(github.Config{RootDir: "/dir", Runner: r})
	if err != nil {
		t.Fatal(err)
	}

	missing, err := ci.Check(deploy.Plan{}, deploy.Job{})
	if err != nil || len(missing) != 1 {
		t.Fatalf("expected 1 missing item, got %v, %v", missing, err)
	}

	if got := lang.Translate(missing[0].Item).String(); got != "CLI de GitHub" {
		t.Errorf("expected translated Item 'CLI de GitHub', got %q", got)
	}
	if got := lang.Translate(missing[0].Why).String(); got != "webtyp lo usa para conectarse a GitHub con tu propia sesión" {
		t.Errorf("expected translated Why, got %q", got)
	}
}

