package github

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// inspired by https://github.com/actions/toolkit/blob/main/packages/github/src/context.ts
// payload ref: https://docs.github.com/en/graphql/reference/objects#pullrequest

type PullRequest struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Number int    `json:"number"`
}

type Issue struct {
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
	Number int    `json:"number"`
}

type Payload struct {
	PullRequest PullRequest `json:"pull_request"`
	Issue       Issue       `json:"issue"`
	Number      int         `json:"number"`
}

type GHContext struct {
	Payload    Payload
	EventName  string
	EventPath  string
	SHA        string
	Ref        string
	Repository string
	Owner      string
	Repo       string
}

func newGHContext() (*GHContext, error) {
	fullRepo := os.Getenv("GITHUB_REPOSITORY") // example: drmaples/pr-lint-jira
	owner, repo, _ := strings.Cut(fullRepo, "/")

	c := &GHContext{
		EventPath:  os.Getenv("GITHUB_EVENT_PATH"),
		EventName:  os.Getenv("GITHUB_EVENT_NAME"),
		SHA:        os.Getenv("GITHUB_SHA"),
		Ref:        os.Getenv("GITHUB_REF"),
		Repository: fullRepo,
		Owner:      owner,
		Repo:       repo,
	}

	if c.EventPath != "" {
		if err := c.loadPayload(); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *GHContext) loadPayload() error {
	content, err := os.ReadFile(c.EventPath)
	if err != nil {
		return fmt.Errorf("problem reading github event path: %w", err)
	}
	if err := json.Unmarshal(content, &c.Payload); err != nil {
		return fmt.Errorf("problem with payload unmarshal: %w", err)
	}
	return nil
}

func (c *GHContext) PRNumber() int {
	// depending on the event type this can move around
	if c.Payload.PullRequest.Number > 0 {
		return c.Payload.PullRequest.Number
	}
	return c.Payload.Issue.Number
}
