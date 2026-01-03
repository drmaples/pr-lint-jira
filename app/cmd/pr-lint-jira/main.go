package main

import (
	"fmt"
	"os"

	"github.com/drmaples/pr-lint-jira/app/internal"
)

func main() {
	if err := internal.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
