package enforcer

import (
	"fmt"
	"log"
)

type RuleConfig struct {
	BannedLabels map[string]bool
	AnyOfThis    map[string]bool
}

type violations []string

func IsValid(ghUser, ghRepo, ghPullNo string, rules *RuleConfig) bool {

	if prPtr, err := fetchPr(ghUser, ghRepo, ghPullNo); err != nil {
		log.Printf("Problem fetching PR: %s", err.Error())
		return false
	} else {
		return IsValidPr(prPtr, rules)
	}
}

func IsValidPr(pr *PullRequest, rules *RuleConfig) bool {
	report := violations{}
	for _, label := range pr.Labels {
		if rules.containsBannedLabel(label.Name) {
			report = append(report, fmt.Sprintf("%s is not allowed", label.Name))
		}
	}

	return len(report) == 0

}

func NewRules(bannedLabels []string, anyOfTheseLabels []string) *RuleConfig {
	config := new(RuleConfig)
	for _, banned := range bannedLabels {
		config.BannedLabels[banned] = true
	}
	for _, anyOfThis := range anyOfTheseLabels {
		config.AnyOfThis[anyOfThis] = true
	}
	return config
}

func (c *RuleConfig) containsBannedLabel(label string) bool {
	return c.BannedLabels[label]
}
