# webtyp/github
<img src="docs/img/badges.svg">

Implements `deploy.CI` for GitHub Actions.

## Config

Use this CI in `config/deploy.go`:

```go
CI: github.Actions
```

## Requirements

This integration uses the `gh` CLI (the developer's own session) to interact with GitHub. The `gh` CLI must be installed and authenticated via `gh auth login`.

## Security

Secret values are passed to `gh` on **stdin only**, never as an argument or environment variable, and never logged.
