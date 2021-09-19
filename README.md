# pr-lint-jira
Lint PR to ensure a JIRA ticket is found in both the title and the body

## Usage

In `.github/workflows`, create a `lint_pr.yml` file with the appropriate config options.

## Example

```yml
name: PR lint

on:
  pull_request:
    types: ['opened', 'edited', 'reopened', 'synchronize']

jobs:
  title:
    name: pr-lint-jira
    runs-on: ubuntu-latest
    steps:
      - uses: drmaples/pr-lint-jira@v1
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
```

## Inputs

| Name | Required | Default | Description |
| --- | --- | --- | --- |
| token | ✅ | | github access token |
| quiet | | `true` | if enabled, create a PR comment on failure |
| titleRegexInput | | `^\[([A-Z]{2,})(-)(\d+)\]` | regex that matches JIRA ticket in PR title. defaults to starting with JIRA ticket in square brackets. |
| bodyRegexInput | | `\[([A-Z]{2,})(-)(\d+)\]` | regex that matches JIRA ticket in PR body. defaults to JIRA ticket in square brackets being anywhere in PR body |
| noTicketInput | | `[no-ticket]` | your PR has no associated JIRA ticket. use this string in both the PR title and body to pass linting |
