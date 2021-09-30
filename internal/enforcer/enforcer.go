package enforcer

import (
	"fmt"
	"log"
	"strings"
)

type RuleConfig struct {
	bannedLabels map[string]bool
	anyOfThis    map[string]bool
}

type violations []string

func IsValid(repoFullName, ghPullNo string, rules *RuleConfig) (violations, bool) {

	if prPtr, err := fetchPrViaFullName(repoFullName, ghPullNo); err != nil {
		log.Printf("Problem fetching PR: %s", err.Error())
		return violations{}, false
	} else {
		return IsValidPr(prPtr, rules)
	}
}

func IsValidPr(pr *PullRequest, rules *RuleConfig) (violations, bool) {
	report := violations{}
	for _, label := range pr.Labels {
		l := strings.ToLower(label.Name)
		if rules.containsBannedLabel(l) {
			report = append(report, fmt.Sprintf("%s is on the blacklist: %v", label.Name, rules.bannedAsList()))
		}
	}
	if !rules.containsAnyRequiredLabel(pr) {
		report = append(report, fmt.Sprintf("does not contained a required label out of: %v", rules.anyOfThisAsList()))
	}

	return report, len(report) == 0
}

func NewRules(bannedLabels []string, anyOfTheseLabels []string) *RuleConfig {
	config := RuleConfig{
		make(map[string]bool),
		make(map[string]bool),
	}
	for _, banned := range bannedLabels {
		config.bannedLabels[strings.ToLower(banned)] = true
	}
	for _, anyOfThis := range anyOfTheseLabels {
		config.anyOfThis[strings.ToLower(anyOfThis)] = true
	}
	return &config
}

func (c *RuleConfig) containsBannedLabel(label string) bool {
	return c.bannedLabels[label]
}

func (c *RuleConfig) containsAnyRequiredLabel(pr *PullRequest) bool {
	if len(c.anyOfThis) == 0 {
		return true
	}
	matchesAnyLabel := false
	for _, label := range pr.Labels {
		l := strings.ToLower(label.Name)
		matchesAnyLabel = matchesAnyLabel || c.anyOfThis[l]
	}
	return matchesAnyLabel
}

func (c *RuleConfig) bannedAsList() []string {
	l := make([]string, 0)
	for key, _ := range c.bannedLabels {
		l = append(l, key)
	}
	return l
}

func (c *RuleConfig) anyOfThisAsList() []string {
	l := make([]string, 0)
	for key, _ := range c.anyOfThis {
		l = append(l, key)
	}
	return l
}

func (v violations) String() string {
	return strings.Join(v, `, `)
}
