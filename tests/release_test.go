package github_test

import (
	"errors"
	"testing"

	"webtyp.com/deploy"
	"webtyp.com/github"
)

func TestRelease(t *testing.T) {
	r := newFakeRunner()
	ci, _ := github.New(github.Config{RootDir: "/dir", Runner: r})

	r.out["gh repo view --json nameWithOwner -q .nameWithOwner"] = []byte("org/repo")
	r.out["gh workflow run deploy.yml --repo org/repo --ref main"] = []byte{}
	r.out["gh run list --workflow deploy.yml --repo org/repo --event workflow_dispatch --limit 1 --json databaseId,url,createdAt"] = []byte(`[{"databaseId": 123, "url": "http://run", "createdAt": "2099-01-01T00:00:00Z"}]`)

	run, err := ci.Release(deploy.Plan{})
	if err != nil {
		t.Fatal(err)
	}

	if run.ID != "123" || run.URL != "http://run" {
		t.Errorf("expected 123/http://run, got %v", run)
	}
}

func TestStatus(t *testing.T) {
	r := newFakeRunner()
	ci, _ := github.New(github.Config{RootDir: "/dir", Runner: r})

	r.out["gh repo view --json nameWithOwner -q .nameWithOwner"] = []byte("org/repo")

	tests := []struct {
		json string
		want deploy.RunState
	}{
		{`{"status":"queued"}`, deploy.RunQueued},
		{`{"status":"waiting"}`, deploy.RunQueued},
		{`{"status":"pending"}`, deploy.RunQueued},
		{`{"status":"in_progress"}`, deploy.RunInProgress},
		{`{"status":"completed","conclusion":"success"}`, deploy.RunSucceeded},
		{`{"status":"completed","conclusion":"failure"}`, deploy.RunFailed},
		{`{"status":"completed","conclusion":"cancelled"}`, deploy.RunFailed},
	}

	for _, tt := range tests {
		r.out["gh run view 123 --repo org/repo --json status,conclusion"] = []byte(tt.json)
		state, err := ci.Status(deploy.Run{ID: "123"})
		if err != nil {
			t.Fatal(err)
		}
		if state != tt.want {
			t.Errorf("expected %v, got %v for %s", tt.want, state, tt.json)
		}
	}

	r.out["gh run view 123 --repo org/repo --json status,conclusion"] = []byte(`{"status":"mystery"}`)
	if _, err := ci.Status(deploy.Run{ID: "123"}); !errors.Is(err, github.ErrUnknownRunStatus) {
		t.Errorf("expected ErrUnknownRunStatus, got %v", err)
	}
}

