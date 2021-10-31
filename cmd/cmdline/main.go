package main

import (
	"fmt"
	"log"

	"github.com/tastybug/github-pr-enforcer/internal/enforcer"
)

func main() {
	// TODO: do input validation here
	repoFullName := "tastybug/github-pr-enforcer"
	ghPullNo := 1

	_, ok := validatePullRequest(repoFullName, ghPullNo, &enforcer.RuleConfig{})

	if ok {
		fmt.Println("Is valid.")
	} else {
		log.Fatalln("Invalid!")
	}
}

func validatePullRequest(repoFullName string, ghPullNo int, rules *enforcer.RuleConfig) (enforcer.Violations, bool) {

	if prPtr, err := enforcer.FetchPrViaFullName(repoFullName, ghPullNo); err != nil {
		log.Printf("Problem fetching PR: %s", err.Error())
		return enforcer.Violations{}, false
	} else {
		// convert external PR representation to internal one, then call IsValidPr...
		log.Printf("Implement converting %+v to internal format, then call IsValidPr", *prPtr)
		// return enforcer.IsValidPr(prPtr, rules)
		return nil, true
	}
}
