package github

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func writeFile(t *testing.T, contents []byte) string {
	t.Helper()
	tmpDir := t.TempDir()
	payloadPath := filepath.Join(tmpDir, "event.json")
	require.NoError(t, os.WriteFile(payloadPath, contents, 0o600))
	return payloadPath
}

func TestNewGHContext_WithEnvironmentVariables(t *testing.T) {
	t.Setenv("GITHUB_REPOSITORY", "some-owner/some/repo/here")
	t.Setenv("GITHUB_EVENT_NAME", "pull_request")
	t.Setenv("GITHUB_SHA", "abc123def456")
	t.Setenv("GITHUB_REF", "refs/heads/main")

	ctx, err := newGHContext()
	require.NoError(t, err)

	require.Equal(t, "some-owner/some/repo/here", ctx.Repository)
	require.Equal(t, "some-owner", ctx.Owner)
	require.Equal(t, "some/repo/here", ctx.Repo)
	require.Equal(t, "pull_request", ctx.EventName)
	require.Equal(t, "abc123def456", ctx.SHA)
	require.Equal(t, "refs/heads/main", ctx.Ref)
}

func TestNewGHContext_WithMissingRepository(t *testing.T) {
	ctx, err := newGHContext()
	require.NoError(t, err)

	require.Empty(t, ctx.Repository)
	require.Empty(t, ctx.Owner)
	require.Empty(t, ctx.Repo)
}

func TestNewGHContext_WithPayload(t *testing.T) {
	payload := Payload{
		PullRequest: PullRequest{
			Title:  "Test PR",
			Body:   "[JIRA-123] Fix bug",
			Number: 42,
		},
	}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	t.Setenv("GITHUB_EVENT_PATH", writeFile(t, payloadBytes))
	t.Setenv("GITHUB_REPOSITORY", "owner/repo")

	ctx, err := newGHContext()
	require.NoError(t, err)

	require.Equal(t, "Test PR", ctx.Payload.PullRequest.Title)
	require.Equal(t, "[JIRA-123] Fix bug", ctx.Payload.PullRequest.Body)
	require.Equal(t, 42, ctx.Payload.PullRequest.Number)
}

func TestNewGHContext_WithInvalidPayloadFile(t *testing.T) {
	t.Setenv("GITHUB_EVENT_PATH", "/nonexistent/path/event.json")
	t.Setenv("GITHUB_REPOSITORY", "owner/repo")

	_, err := newGHContext()
	require.Error(t, err)
}

func TestNewGHContext_WithMalformedPayload(t *testing.T) {
	t.Setenv("GITHUB_EVENT_PATH", writeFile(t, []byte("invalid json")))
	t.Setenv("GITHUB_REPOSITORY", "owner/repo")

	_, err := newGHContext()
	require.Error(t, err)
}

func TestPRNumber_FromPullRequest(t *testing.T) {
	ctx := &GHContext{
		Payload: Payload{
			PullRequest: PullRequest{Number: 123},
			Issue:       Issue{Number: 456},
		},
	}

	require.Equal(t, 123, ctx.PRNumber())
}

func TestPRNumber_FromIssue(t *testing.T) {
	ctx := &GHContext{
		Payload: Payload{
			PullRequest: PullRequest{Number: 0},
			Issue:       Issue{Number: 456},
		},
	}

	require.Equal(t, 456, ctx.PRNumber())
}

func TestPRNumber_BothZero(t *testing.T) {
	ctx := &GHContext{
		Payload: Payload{
			PullRequest: PullRequest{Number: 0},
			Issue:       Issue{Number: 0},
		},
	}

	require.Equal(t, 0, ctx.PRNumber())
}

func TestLoadPayload_Success(t *testing.T) {
	payload := Payload{
		PullRequest: PullRequest{
			Title:  "Feature request",
			Body:   "Description here",
			Number: 99,
		},
		Number: 99,
	}
	payloadBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	ctx := &GHContext{EventPath: writeFile(t, payloadBytes)}
	require.NoError(t, ctx.loadPayload())

	require.Equal(t, "Feature request", ctx.Payload.PullRequest.Title)
	require.Equal(t, 99, ctx.Payload.Number)
}

func TestRepositoryParsing_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		repo      string
		wantRepo  string
		wantOwner string
	}{
		{"Standard", "owner/repo", "repo", "owner"},
		{"With sub-dirs", "owner/repo/subdir", "repo/subdir", "owner"},
		{"Empty", "", "", ""},
		{"Only slash", "/", "", ""},
		{"No slash", "noslash", "", "noslash"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("GITHUB_REPOSITORY", tt.repo)

			ctx, err := newGHContext()
			require.NoError(t, err)

			require.Equal(t, tt.wantOwner, ctx.Owner)
			require.Equal(t, tt.wantRepo, ctx.Repo)
		})
	}
}
