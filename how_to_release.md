how to release:

1. click `Draft a new release` https://github.com/drmaples/pr-lint-jira/releases
1. choose a new tag with bumped version, target latest commit, publish release
1. set v1 tag to latest release
   ```
   git tag -f v1
   git push origin v1 --force
   ```
