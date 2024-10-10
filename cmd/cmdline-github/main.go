package main

import (
	"fmt"
	"github.com/tastybug/github-pr-enforcer/internal/enforcer/domain"
	"github.com/tastybug/github-pr-enforcer/internal/enforcer/service"
	"log"
)

/*
Purpose: read arguments from command line and call port to process
*/
func main() {
	// TODO: do input validation here
	repoFullName := "tastybug/github-pr-enforcer"
	ghPullNo := 1

	_, ok := validatePullRequest(repoFullName, ghPullNo, &domain.RuleConfig{})

	if ok {
		fmt.Println("Is valid.")
	} else {
		log.Fatalln("Invalid!")
	}
}

func validatePullRequest(repoFullName string, ghPullNo int, rules *domain.RuleConfig) (domain.Violations, bool) {

	if prPtr, err := service.FetchPrViaFullName(repoFullName, ghPullNo); err != nil {
		log.Printf("Problem fetching PR: %s", err.Error())
		return domain.Violations{}, false
	} else {
		// convert external PR representation to internal one, then call IsValidPr...
		log.Printf("Implement converting %+v to internal format, then call IsValidPr", *prPtr)
		// return enforcer.IsValidPr(prPtr, rules)
		return nil, true
	}
}
