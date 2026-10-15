package github

import (
	"bytes"
	"os"
	"path/filepath"

	"webtyp.com/deploy"
	"webtyp.com/lang"
)

func (c *CI) Check(p deploy.Plan, job deploy.Job) ([]deploy.Missing, error) {
	var missing []deploy.Missing

	// G1: gh --version
	_, err := c.runner.Run(c.rootDir, "", "gh", "--version")
	if err != nil {
		return []deploy.Missing{{
			Item:  textGhMissing,
			Why:   textGhMissingWhy,
			Steps: []lang.Text{textGhInstall},
			Fix:   deploy.FixManual,
		}}, nil
	}

	// G2: gh auth status
	_, err = c.runner.Run(c.rootDir, "", "gh", "auth", "status")
	if err != nil {
		return []deploy.Missing{{
			Item:  textGhAuthMissing,
			Why:   textGhAuthMissingWhy,
			Steps: []lang.Text{textGhAuthFix},
			Fix:   deploy.FixManual,
		}}, nil
	}

	// G3: repo
	_, err = c.getRepo()
	if err != nil {
		missing = append(missing, deploy.Missing{
			Item:  textGhRemoteMissing,
			Why:   textGhRemoteMissingWhy,
			Steps: []lang.Text{textGhRemoteFix},
			Fix:   deploy.FixManual,
		})
	}

	rendered := render(p, job)

	// G4: read rootDir/.github/workflows/deploy.yml
	path := filepath.Join(c.rootDir, ".github", "workflows", "deploy.yml")
	content, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(content, rendered) {
		missing = append(missing, deploy.Missing{
			Item:  textGhWorkflowMissing,
			Why:   textGhWorkflowMissingWhy,
			Steps: []lang.Text{textGhWorkflowFix},
			Fix:   deploy.FixAutomatic,
		})
	}

	// G5: git show origin/main:.github/workflows/deploy.yml
	gitOut, gitErr := c.runner.Run(c.rootDir, "", "git", "show", "origin/main:.github/workflows/deploy.yml")
	if gitErr != nil || !bytes.Equal(gitOut, rendered) {
		missing = append(missing, deploy.Missing{
			Item:  textGhWorkflowOutdated,
			Why:   textGhWorkflowOutdatedWhy,
			Steps: []lang.Text{textGhWorkflowCommitFix},
			Fix:   deploy.FixManual,
		})
	}

	return missing, nil
}
