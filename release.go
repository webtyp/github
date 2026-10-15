package github

import (
	"encoding/json"
	"fmt"
	"time"

	"webtyp.com/deploy"
)

var (
	pollInterval = 2 * time.Second
	pollTimeout  = 30 * time.Second
	timeNow      = time.Now
	timeSleep    = time.Sleep
)

func (c *CI) Release(p deploy.Plan) (deploy.Run, error) {
	repo, err := c.getRepo()
	if err != nil {
		return deploy.Run{}, err
	}

	dispatchTime := timeNow()

	_, err = c.runner.Run(c.rootDir, "", "gh", "workflow", "run", "deploy.yml", "--repo", repo, "--ref", "main")
	if err != nil {
		return deploy.Run{}, err
	}

	deadline := dispatchTime.Add(pollTimeout)

	for timeNow().Before(deadline) {
		out, err := c.runner.Run(c.rootDir, "", "gh", "run", "list", "--workflow", "deploy.yml", "--repo", repo, "--event", "workflow_dispatch", "--limit", "1", "--json", "databaseId,url,createdAt")
		if err == nil {
			var runs []struct {
				DatabaseID int64    `json:"databaseId"`
				URL        string `json:"url"`
				CreatedAt  string `json:"createdAt"`
			}
			if err := json.Unmarshal(out, &runs); err == nil && len(runs) > 0 {
				createdAt, err := time.Parse(time.RFC3339, runs[0].CreatedAt)
				if err == nil && createdAt.After(dispatchTime) {
					return deploy.Run{
						ID:  fmt.Sprintf("%d", runs[0].DatabaseID),
						URL: runs[0].URL,
					}, nil
				}
			}
		}
		timeSleep(pollInterval)
	}

	return deploy.Run{}, ErrRunNotFound
}

func (c *CI) Status(r deploy.Run) (deploy.RunState, error) {
	repo, err := c.getRepo()
	if err != nil {
		return 0, err
	}

	out, err := c.runner.Run(c.rootDir, "", "gh", "run", "view", r.ID, "--repo", repo, "--json", "status,conclusion")
	if err != nil {
		return 0, err
	}

	var data struct {
		Status     string `json:"status"`
		Conclusion string `json:"conclusion"`
	}
	if err := json.Unmarshal(out, &data); err != nil {
		return 0, err
	}

	switch data.Status {
	case "queued", "waiting", "pending":
		return deploy.RunQueued, nil
	case "in_progress":
		return deploy.RunInProgress, nil
	case "completed":
		if data.Conclusion == "success" {
			return deploy.RunSucceeded, nil
		}
		return deploy.RunFailed, nil
	}

	return 0, ErrUnknownRunStatus
}
