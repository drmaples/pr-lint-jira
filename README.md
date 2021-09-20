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

see [action.yml](action.yml)
