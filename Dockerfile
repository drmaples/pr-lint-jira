# keep the golang version in sync with the mise.toml file
FROM golang:1.25.5 AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN update-ca-certificates
RUN useradd -u 10001 scratchuser

WORKDIR /code

# vendor dir may not exist, thus the [r] trick.
# https://stackoverflow.com/questions/31528384/conditional-copy-add-in-dockerfile
# COPY vendo[r] vendor

COPY go.mod go.sum ./
COPY app app

RUN go build -o=/go/bin ./app/cmd/...

# ----------------------------------------
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/pr-lint-jira .

# dont run as root in a container:
# https://docs.docker.com/develop/develop-images/dockerfile_best-practices/#user
COPY --from=builder /etc/passwd /etc/passwd
USER scratchuser

CMD ["./pr-lint-jira"]
