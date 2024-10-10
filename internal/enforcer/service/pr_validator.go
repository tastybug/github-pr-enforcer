package service

import (
	"fmt"
	"github.com/tastybug/github-pr-enforcer/internal/enforcer/domain"
	"strings"
)

func ValidatePr(pr *domain.PullRequest, rules *domain.RuleConfig) (domain.Violations, bool) {
	report := domain.Violations{}
	for _, label := range pr.Labels {
		l := strings.ToLower(label.Name)
		if containsBannedLabel(rules, l) {
			report = append(report, fmt.Sprintf("%s is on the blacklist: %v", label.Name, rules.BannedAsList()))
		}
	}
	if !containsAnyRequiredLabel(rules, pr) {
		report = append(report, fmt.Sprintf("does not contained a required label out of: %v", rules.AnyOfThisAsList()))
	}

	return report, len(report) == 0
}

func containsBannedLabel(c *domain.RuleConfig, label string) bool {
	return c.MustNotHaveOneOf[label]
}

func containsAnyRequiredLabel(c *domain.RuleConfig, pr *domain.PullRequest) bool {
	if len(c.NeedsOneOf) == 0 {
		return true
	}
	for _, label := range pr.Labels {
		l := strings.ToLower(label.Name)
		if c.NeedsOneOf[l] {
			return true
		}
	}
	return false
}

func DefaultRules() *domain.RuleConfig {
	return domain.CreateRuleConfig(
		[]string{"wip", "do-not-merge"},
		[]string{"bug", "feature", "enabler", "rework"},
	)
}
