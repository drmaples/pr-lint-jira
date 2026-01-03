package github

import (
	"fmt"
	"os"
)

type Iface interface {
	Context() GHContext
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

func (a *api) Context() GHContext {
	gh := GHContext{
		EventPath: os.Getenv("GITHUB_EVENT_PATH"),
		EventName: os.Getenv("GITHUB_EVENT_NAME"),
		SHA:       os.Getenv("GITHUB_SHA"),
		Ref:       os.Getenv("GITHUB_REF"),
	}

	if gh.EventPath != "" {
		// do nothing
		fmt.Println("empty path")
	}

	return gh
}

func (a *api) CreateComment(comment string) error {
	fmt.Println(comment)
	return nil
}
