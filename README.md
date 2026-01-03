# pr-lint-jira

Lint PR to ensure a ticket is found in both the title and the body

## Usage

In `.github/workflows`, create a `lint_pr.yml` file with the appropriate config options.

If you do not have a ticket, then add `[no-ticket]` or whatever is configured in the `no_ticket`.

## Example

```yml
name: PR lint

on:
  pull_request:
    types: ["opened", "edited", "reopened", "synchronize"]

# required permissions IF make_pr_comment is set to true.
# by default these extra permissions are not needed
permissions:
  issues: write
  pull-requests: write

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
