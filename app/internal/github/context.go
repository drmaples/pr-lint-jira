package github

// modeled after https://github.com/actions/toolkit/blob/main/packages/github/src/context.ts

type PullRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Issue struct {
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
	Number string `json:"number"` // issue number is same as PR number
}

type Payload struct {
	PullRequest PullRequest `json:"pull_request"`
	Issue       Issue       `json:"issue"`
}

type GHContext struct {
	Payload   Payload
	EventName string
	EventPath string
	SHA       string
	Ref       string
}
