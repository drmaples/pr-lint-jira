package github

import (
	"fmt"
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
	return newGHContext()
}

func (a *api) CreateComment(comment string) error {
	fmt.Println(comment)
	return nil
}
