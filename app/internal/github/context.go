package github

// modeled after https://github.com/actions/toolkit/blob/main/packages/github/src/context.ts

type PullRequest struct {
	Title string
	Body  string
}

type Issue struct {
	Owner  string
	Repo   string
	Number string // issue number is same as PR number
}

type Payload struct {
	PullRequest PullRequest
	Issue       Issue
}

type GHContext struct {
	Payload   Payload
	EventName string
	EventPath string
	SHA       string
	Ref       string
}
