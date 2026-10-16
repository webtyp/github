package github

import "errors"

var (
	ErrNoRootDir        = errors.New("github: missing RootDir")
	ErrNoRunner         = errors.New("github: missing Runner")
	ErrNoRemote         = errors.New("github: no remote")
	ErrRunNotFound      = errors.New("github: run not found")
	ErrUnknownRunStatus = errors.New("github: unknown run status")
)
