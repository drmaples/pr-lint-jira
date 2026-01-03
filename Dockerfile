# keep the golang version in sync with the mise.toml file
FROM golang:1.25.5 AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN update-ca-certificates

WORKDIR /code

# vendor dir may not exist, thus the [r] trick.
# https://stackoverflow.com/questions/31528384/conditional-copy-add-in-dockerfile
# COPY vendo[r] vendor

COPY go.mod go.sum ./
COPY app app

RUN go build -o=/go/bin ./app/cmd/...

########################
# NOTE:
#   https://docs.github.com/en/actions/reference/workflows-and-actions/dockerfile-support
#   must run as root user
#   do not use workdir, GHA sets it
#   use absolute path for the cmd or entrypoint since GHA modifies workdir
########################
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/pr-lint-jira .

CMD ["/pr-lint-jira"]
