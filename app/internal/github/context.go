package github

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// inspired by https://github.com/actions/toolkit/blob/main/packages/github/src/context.ts

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
	Payload   Payload
	EventName string
	EventPath string
	SHA       string
	Ref       string
}

func newGHContext() (*GHContext, error) {
	c := &GHContext{
		EventPath: os.Getenv("GITHUB_EVENT_PATH"),
		EventName: os.Getenv("GITHUB_EVENT_NAME"),
		SHA:       os.Getenv("GITHUB_SHA"),
		Ref:       os.Getenv("GITHUB_REF"),
	}

	if c.EventPath != "" {
		c.loadPayload()
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

func (c *GHContext) Issue() (string, string, int) {
	owner, repo := c.Repo()

	number := 0
	if c.Payload.Issue.Number > 0 {
		number = c.Payload.Issue.Number
	} else if c.Payload.PullRequest.Number > 0 {
		number = c.Payload.PullRequest.Number
	}
	return owner, repo, number
}

func (c *GHContext) Repo() (string, string) {
	full := os.Getenv("GITHUB_REPOSITORY")
	nuggets := strings.SplitN(full, "/", 2)
	return nuggets[0], nuggets[1]
}
