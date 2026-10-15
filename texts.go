package github

import "webtyp.com/lang"

const (
	textGhMissing              lang.Text = "GitHub CLI"
	textGhMissingWhy           lang.Text = "webtyp uses it to reach GitHub with your own session"
	textGhInstall              lang.Text = "Install it from https://cli.github.com and open the project again"
	textGhAuthMissing          lang.Text = "GitHub CLI auth"
	textGhAuthMissingWhy       lang.Text = "webtyp needs your GitHub session to configure the deployment"
	textGhAuthFix              lang.Text = "Run 'gh auth login' in a terminal, follow the device code, and open the project again"
	textGhRemoteMissing        lang.Text = "GitHub repository"
	textGhRemoteMissingWhy     lang.Text = "The deploy runs on GitHub, but this folder has no remote"
	textGhRemoteFix            lang.Text = "Run 'gh repo create' in a terminal and open the project again"
	textGhWorkflowMissing      lang.Text = "Deploy workflow"
	textGhWorkflowMissingWhy   lang.Text = "GitHub Actions needs a workflow file to run the deploy"
	textGhWorkflowFix          lang.Text = "webtyp will write .github/workflows/deploy.yml"
	textGhWorkflowOutdated     lang.Text = "Deploy workflow outdated"
	textGhWorkflowOutdatedWhy  lang.Text = "The deploy config changed, but the workflow file on GitHub is old"
	textGhWorkflowCommitFix    lang.Text = "Commit and push .github/workflows/deploy.yml to main"
)
