package main

import (
	"fmt"
	"log"

	"github.com/tastybug/github-pr-enforcer/internal/enforcer"
)

func main() {
	ghUser := "tastybug"
	ghRepo := "github-pr-enforcer"
	ghPullNo := "1"

	result := enforcer.IsValid(ghUser, ghRepo, ghPullNo, &enforcer.RuleConfig{})
	if result {
		fmt.Println("Is valid.")
	} else {
		log.Fatalln("Invalid!")
	}
}
