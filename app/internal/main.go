package internal

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/drmaples/pr-lint-jira/app/internal/github"
)

const (
	envVarPrefix          = "INPUT_"
	defaultTitleBodyRegex = `\[([a-zA-Z]{2,}-\d+)\]`
	defaultNoTicket       = "[no-ticket]"

	tokenKey         = "token"
	makePRCommentKey = "make_pr_comment"
	titleRegexKey    = "title_regex"
	bodyRegexKey     = "body_regex"
	noTicketKey      = "no_ticket"
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

func Run() error {
	fmt.Println("========== env ==========")
	for _, e := range os.Environ() {
		fmt.Printf("%#v\n", e)
	}
	fmt.Println("========== env ==========")

	token := getInputValue(tokenKey)

	gh := github.NewGithub(token)
	gh.Context()

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

	fmt.Println("==========")
	if false {
		fmt.Println(token)
	}
	fmt.Printf("makePRComment: %#v\n", makePRComment)
	fmt.Printf("titleRegex:    %#v\n", titleRegex)
	fmt.Printf("bodyRegex:     %#v\n", bodyRegex)
	fmt.Printf("noTicket:      %#v\n", noTicket)
	fmt.Println("==========")

	return nil
}
