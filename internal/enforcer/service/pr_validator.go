package service

import (
	"fmt"
	"github.com/tastybug/github-pr-enforcer/internal/enforcer/domain"
	"slices"
	"strings"
)

func ValidatePr(pr *domain.PullRequest, rules *domain.RuleConfig) (domain.Violations, bool) {
	report := domain.Violations{}
	for _, label := range pr.Labels {
		l := strings.ToLower(label.Name)
		if containsBannedLabel(rules, l) {
			report = append(report, fmt.Sprintf("%s is on the blacklist: %v", label.Name, rules.BannedLabels))
		}
	}
	if !containsAnyRequiredLabel(rules, pr) {
		report = append(report, fmt.Sprintf("does not contained a required label out of: %v", rules.AnyOfTheseLabels))
	}

	return report, len(report) == 0
}

func containsBannedLabel(c *domain.RuleConfig, label string) bool {
	return slices.Contains(c.BannedLabels, label)
}

func containsAnyRequiredLabel(c *domain.RuleConfig, pr *domain.PullRequest) bool {
	if len(c.AnyOfTheseLabels) == 0 {
		return true
	}
	for _, label := range pr.Labels {
		normalizedLabelToLookFor := strings.ToLower(label.Name)
		if slices.Contains(c.AnyOfTheseLabels, normalizedLabelToLookFor) {
			return true
		}
	}
	return false
}
