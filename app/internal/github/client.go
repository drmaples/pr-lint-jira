package github

import (
	"encoding/json"
	"fmt"
	"os"
)

type Iface interface {
	Context() (*GHContext, error)
	CreateComment(comment string) error
}

type api struct {
	Token string
	// client
}

func NewGithub(token string) Iface {
	return &api{
		Token: token,
	}
}

func (a *api) Context() (*GHContext, error) {
	gh := &GHContext{
		EventPath: os.Getenv("GITHUB_EVENT_PATH"),
		EventName: os.Getenv("GITHUB_EVENT_NAME"),
		SHA:       os.Getenv("GITHUB_SHA"),
		Ref:       os.Getenv("GITHUB_REF"),
	}

	if gh.EventPath != "" {
		content, err := os.ReadFile(gh.EventPath)
		if err != nil {
			return nil, fmt.Errorf("problem reading github event path: %w", err)
		}

		if err := json.Unmarshal(content, &gh.Payload); err != nil {
			return nil, fmt.Errorf("problem with payload unmarshal: %w", err)
		}
	}

	return gh, nil
}

func (a *api) CreateComment(comment string) error {
	fmt.Println(comment)
	return nil
}
