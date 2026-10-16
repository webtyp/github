package github_test

import (
	"strings"
	"testing"

	"webtyp.com/deploy"
	"webtyp.com/github"
)

func TestSecretsHas(t *testing.T) {
	r := newFakeRunner()
	ci, _ := github.New(github.Config{RootDir: "/dir", Runner: r})

	r.out["gh repo view --json nameWithOwner -q .nameWithOwner"] = []byte("org/repo")
	r.out["gh secret list --repo org/repo --json name"] = []byte(`[{"name":"K"}]`)

	dest := ci.Secrets(deploy.Plan{})
	has, err := dest.Has("K")
	if err != nil {
		t.Fatal(err)
	}
	if !has {
		t.Error("expected true")
	}

	has, err = dest.Has("OTHER")
	if err != nil {
		t.Fatal(err)
	}
	if has {
		t.Error("expected false")
	}
}

func TestSecretsPut(t *testing.T) {
	r := newFakeRunner()
	ci, _ := github.New(github.Config{RootDir: "/dir", Runner: r})

	r.out["gh repo view --json nameWithOwner -q .nameWithOwner"] = []byte("org/repo")
	r.out["gh secret set K --repo org/repo"] = []byte{}

	dest := ci.Secrets(deploy.Plan{})
	err := dest.Put("K", "secret_value")
	if err != nil {
		t.Fatal(err)
	}

	for _, cmd := range r.cmds {
		if cmd.name == "gh" && len(cmd.args) > 0 && cmd.args[0] == "secret" && cmd.args[1] == "set" {
			if cmd.stdin != "secret_value" {
				t.Errorf("expected stdin 'secret_value', got %v", cmd.stdin)
			}
			for _, arg := range cmd.args {
				if strings.Contains(arg, "secret_value") {
					t.Errorf("value found in args: %v", cmd.args)
				}
			}
		}
	}
}
