package service

import (
	"fmt"
	"github.com/tastybug/github-pr-enforcer/internal/enforcer/domain"
	"strings"
)

func IsValidPr(pr *domain.PullRequest, rules *domain.RuleConfig) (domain.Violations, bool) {
	report := domain.Violations{}
	for _, label := range pr.Labels {
		l := strings.ToLower(label.Name)
		if rules.ContainsBannedLabel(l) {
			report = append(report, fmt.Sprintf("%s is on the blacklist: %v", label.Name, rules.BannedAsList()))
		}
	}
	if !rules.ContainsAnyRequiredLabel(pr) {
		report = append(report, fmt.Sprintf("does not contained a required label out of: %v", rules.AnyOfThisAsList()))
	}

	return report, len(report) == 0
}

func DefaultRules() *domain.RuleConfig {
	return domain.CreateRuleConfig(
		[]string{"wip", "do-not-merge"},
		[]string{"bug", "feature", "enabler", "rework"},
	)
}
