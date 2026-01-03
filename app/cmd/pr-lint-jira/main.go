package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/drmaples/pr-lint-jira/app/github"
)

const (
	defaultErrorMessage   = "missing ticket reference in both PR title AND body"
	envVarPrefix          = "INPUT_"
	defaultTitleBodyRegex = `\[([a-zA-Z]{2,}-\d+)\]`
	defaultNoTicket       = "[no-ticket]"

	tokenKey         = "token"
	makePRCommentKey = "make_pr_comment"
	titleRegexKey    = "title_regex"
	bodyRegexKey     = "body_regex"
	noTicketKey      = "no_ticket"
	logLevelKey      = "log_level"
)

func getInputValue(key string) string {
	envVarKey := strings.ToUpper(fmt.Sprintf("%s%s", envVarPrefix, key))
	return strings.TrimSpace(os.Getenv(envVarKey))
}

func getRegex(reString string, defaultRe string) (*regexp.Regexp, error) {
	if reString == "" {
		reString = defaultRe
	}
	return regexp.Compile(reString)
}

func run(ctx context.Context) error {
	debug := getInputValue(logLevelKey) == "debug"

	if debug {
		fmt.Println("passed env vars:")
		for _, e := range sort.StringSlice(os.Environ()) {
			fmt.Println("\t", e)
		}
	}

	token := getInputValue(tokenKey)

	gh := github.NewGithub(token)
	ghCtx, err := gh.Context()
	if err != nil {
		return err
	}
	if debug {
		fmt.Println("parsed github event context:")
		fmt.Printf("%#v\n", ghCtx)
	}

	makePRComment, err := strconv.ParseBool(strings.ToLower(getInputValue(makePRCommentKey)))
	if err != nil {
		return fmt.Errorf("problem getting %s input: %w", makePRCommentKey, err)
	}

	titleRegex, err := getRegex(getInputValue(titleRegexKey), defaultTitleBodyRegex)
	if err != nil {
		return fmt.Errorf("problem parsing %s regex: %w", titleRegexKey, err)
	}

	bodyRegex, err := getRegex(getInputValue(bodyRegexKey), defaultTitleBodyRegex)
	if err != nil {
		return fmt.Errorf("problem parsing %s regex: %w", bodyRegexKey, err)
	}

	noTicket := getInputValue(noTicketKey)
	if noTicket == "" {
		noTicket = defaultNoTicket
	}

	fmt.Printf("input make_pr_comment : %#v\n", makePRComment)
	fmt.Printf("input title_regex     : %q\n", titleRegex)
	fmt.Printf("input body_regex      : %q\n", bodyRegex)
	fmt.Printf("input no_ticket       : %#v\n", noTicket)

	prTitle := ghCtx.Payload.PullRequest.Title
	prBody := ghCtx.Payload.PullRequest.Body

	foundTixInTitle := titleRegex.FindString(prTitle) != ""
	foundTixInBody := bodyRegex.FindString(prBody) != ""
	foundNoTixInTitle := strings.Contains(prTitle, noTicket)
	foundNoTixInBody := strings.Contains(prBody, noTicket)

	fmt.Println("found ticket in title    :", foundTixInTitle)
	fmt.Println("found ticket in body     :", foundTixInBody)
	fmt.Println("found no_ticket in title :", foundNoTixInTitle)
	fmt.Println("found no_ticket in body  :", foundNoTixInBody)

	if ghCtx.EventName != "pull_request" {
		fmt.Println("success, event is not a pull request")
		return nil
	}
	if foundTixInTitle && foundTixInBody {
		fmt.Println("success, found tix in both title and body")
		return nil
	}
	if foundNoTixInTitle && foundNoTixInBody {
		fmt.Printf("success, found %q in both title and body\n", noTicket)
		return nil
	}

	if makePRComment {
		if err := gh.CreateComment(ctx, ghCtx.Owner, ghCtx.Repo, ghCtx.PRNumber(), defaultErrorMessage); err != nil {
			return fmt.Errorf("problem creating PR comment: %w", err)
		}
	}
	return fmt.Errorf("missing ticket reference in both PR title AND body, use %q for no associated ticket", noTicket)
}

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
