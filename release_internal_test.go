// Root-level test (justified): exercises timeNow and timeSleep — the polling clock and sleep are unexported seams so multi-poll and timeout tests do not wait in real time.
package github

import (
	"errors"
	"strings"
	"testing"
	"time"

	"webtyp.com/deploy"
)

type seqRunner struct {
	listCalls int
	listOuts  [][]byte
}

func (s *seqRunner) Run(_, _, name string, args ...string) ([]byte, error) {
	cmd := name + " " + strings.Join(args, " ")
	switch cmd {
	case "gh repo view --json nameWithOwner -q .nameWithOwner":
		return []byte("org/repo\n"), nil
	case "gh workflow run deploy.yml --repo org/repo --ref main":
		return nil, nil
	case "gh run list --workflow deploy.yml --repo org/repo --event workflow_dispatch --limit 1 --json databaseId,url,createdAt":
		idx := s.listCalls
		s.listCalls++
		if idx < len(s.listOuts) {
			return s.listOuts[idx], nil
		}
		if len(s.listOuts) > 0 {
			return s.listOuts[len(s.listOuts)-1], nil
		}
		return []byte(`[]`), nil
	default:
		return nil, errors.New("unexpected command: " + cmd)
	}
}

func TestReleaseIgnoresOlderRunAndPolls(t *testing.T) {
	now := time.Date(2026, 10, 15, 12, 0, 0, 0, time.UTC)
	origNow, origSleep := timeNow, timeSleep
	defer func() {
		timeNow = origNow
		timeSleep = origSleep
	}()
	timeNow = func() time.Time { return now }
	timeSleep = func(d time.Duration) { now = now.Add(d) }

	r := &seqRunner{
		listOuts: [][]byte{
			[]byte(`[{"databaseId": 100, "url": "https://github.com/org/repo/actions/runs/100", "createdAt": "2026-10-15T11:59:00Z"}]`),
			[]byte(`[{"databaseId": 200, "url": "https://github.com/org/repo/actions/runs/200", "createdAt": "2026-10-15T12:00:01Z"}]`),
		},
	}
	ci, err := New(Config{RootDir: "/dir", Runner: r})
	if err != nil {
		t.Fatal(err)
	}

	run, err := ci.Release(deploy.Plan{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if run.ID != "200" || run.URL != "https://github.com/org/repo/actions/runs/200" {
		t.Errorf("expected run 200, got %+v", run)
	}
	if r.listCalls != 2 {
		t.Errorf("expected 2 list calls, got %d", r.listCalls)
	}
}

func TestReleaseTimeout(t *testing.T) {
	now := time.Date(2026, 10, 15, 12, 0, 0, 0, time.UTC)
	origNow, origSleep := timeNow, timeSleep
	defer func() {
		timeNow = origNow
		timeSleep = origSleep
	}()
	timeNow = func() time.Time { return now }
	timeSleep = func(d time.Duration) { now = now.Add(d) }

	r := &seqRunner{
		listOuts: [][]byte{
			[]byte(`[{"databaseId": 100, "url": "https://github.com/org/repo/actions/runs/100", "createdAt": "2026-10-15T11:59:00Z"}]`),
		},
	}
	ci, err := New(Config{RootDir: "/dir", Runner: r})
	if err != nil {
		t.Fatal(err)
	}

	_, err = ci.Release(deploy.Plan{})
	if !errors.Is(err, ErrRunNotFound) {
		t.Errorf("expected ErrRunNotFound, got %v", err)
	}
	if r.listCalls != 15 {
		t.Errorf("expected 15 polls (30s / 2s), got %d", r.listCalls)
	}
}
