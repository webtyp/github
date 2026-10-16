package github

import (
	"encoding/json"
	"strings"

	"webtyp.com/deploy"
)

type githubSecrets struct {
	ci *CI
}

func (s githubSecrets) Has(key string) (bool, error) {
	repo, err := s.ci.getRepo()
	if err != nil {
		return false, err
	}

	out, err := s.ci.runner.Run(s.ci.rootDir, "", cmdGh, "secret", "list", "--repo", repo, "--json", "name")
	if err != nil {
		return false, err
	}

	var list []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(out, &list); err != nil {
		return false, err
	}

	for _, item := range list {
		if item.Name == key {
			return true, nil
		}
	}
	return false, nil
}

func (s githubSecrets) Put(key, value string) error {
	repo, err := s.ci.getRepo()
	if err != nil {
		return err
	}
	_, err = s.ci.runner.Run(s.ci.rootDir, value, cmdGh, "secret", "set", key, "--repo", repo)
	return err
}

func (c *CI) getRepo() (string, error) {
	if c.repo != "" {
		return c.repo, nil
	}
	out, err := c.runner.Run(c.rootDir, "", cmdGh, "repo", "view", "--json", "nameWithOwner", "-q", ".nameWithOwner")
	if err != nil {
		return "", ErrNoRemote
	}
	name := strings.TrimSpace(string(out))
	if name == "" {
		return "", ErrNoRemote
	}
	c.repo = name
	return name, nil
}

func (c *CI) Secrets(p deploy.Plan) deploy.Destination {
	return githubSecrets{ci: c}
}
