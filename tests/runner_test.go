package github_test

import (
	"errors"
	"strings"
)

type command struct {
	dir   string
	stdin string
	name  string
	args  []string
}

type fakeRunner struct {
	cmds []command
	out  map[string][]byte
	err  map[string]error
}

func newFakeRunner() *fakeRunner {
	return &fakeRunner{
		out: make(map[string][]byte),
		err: make(map[string]error),
	}
}

func (f *fakeRunner) Run(dir, stdin, name string, args ...string) ([]byte, error) {
	cmdStr := name + " " + strings.Join(args, " ")
	f.cmds = append(f.cmds, command{dir, stdin, name, args})

	if err, ok := f.err[cmdStr]; ok {
		return nil, err
	}
	if out, ok := f.out[cmdStr]; ok {
		return out, nil
	}
	return nil, errors.New("unexpected command: " + cmdStr)
}
