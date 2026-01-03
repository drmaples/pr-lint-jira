package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Iface interface {
	Context() (*GHContext, error)
	CreateComment(ctx context.Context, owner string, repo string, issueNumber int, comment string) error
}

type api struct {
	Token string
}

func NewGithub(token string) Iface {
	return &api{
		Token: token,
	}
}

func (a *api) Context() (*GHContext, error) {
	return newGHContext()
}

// https://docs.github.com/en/rest/issues/comments?apiVersion=2022-11-28#create-an-issue-comment
func (a *api) CreateComment(ctx context.Context, owner string, repo string, issueNumber int, comment string) error {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/issues/%d/comments",
		owner,
		repo,
		issueNumber,
	)

	body, err := json.Marshal(map[string]string{
		"body": comment,
	})
	if err != nil {
		return fmt.Errorf("problem with payload marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("problem creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+a.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "pr-lint-jira-action")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("problem sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf(
			"GitHub API error: status=%d body=%s",
			resp.StatusCode,
			string(respBody),
		)
	}

	return nil
}
