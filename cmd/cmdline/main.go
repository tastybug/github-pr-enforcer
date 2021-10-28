package main

import (
	"fmt"
	"log"

	"github.com/tastybug/github-pr-enforcer/internal/enforcer"
)

func main() {
	// TODO: do input validation here
	repoFullName := "tastybug/github-pr-enforcer"
	ghPullNo := "1"

	_, ok := enforcer.ValidatePullRequest(repoFullName, ghPullNo, &enforcer.RuleConfig{})
	if ok {
		fmt.Println("Is valid.")
	} else {
		log.Fatalln("Invalid!")
	}
}
