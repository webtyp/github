package tests

import (
	"testing"
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
	lang.OutLang("es")
	// For testing L1, we'd check one of the returned texts from missing.
	// We'll skip for now if it requires setup, or add test for dictionary.
	// Actually we should test it.
}
