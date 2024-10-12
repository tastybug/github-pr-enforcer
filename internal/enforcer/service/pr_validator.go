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
			report = append(report, fmt.Sprintf("FAIL: label %s prevents merge.", label.Name))
		}
	}
	if !containsOneRequiredLabel(rules, pr) {
		report = append(report, fmt.Sprintf("FAIL: required label missing: %v", rules.OneOfLabels))
	}

	if containsTooManyRequiredLabels(rules, pr) {
		report = append(report, fmt.Sprintf("FAIL: only one of : %v", rules.OneOfLabels))
	}

	return report, len(report) == 0
}

func containsTooManyRequiredLabels(rules *domain.RuleConfig, pr *domain.PullRequest) bool {
	var matches = 0
	for _, label := range pr.Labels {
		normalizedLabelToLookFor := strings.ToLower(label.Name)
		if slices.Contains(rules.OneOfLabels, normalizedLabelToLookFor) {
			matches++
		}
	}
	return matches > 1
}

func containsBannedLabel(c *domain.RuleConfig, label string) bool {
	return slices.Contains(c.NoneOfLabels, label)
}

func containsOneRequiredLabel(c *domain.RuleConfig, pr *domain.PullRequest) bool {
	if len(c.OneOfLabels) == 0 {
		return true
	}
	for _, label := range pr.Labels {
		normalizedLabelToLookFor := strings.ToLower(label.Name)
		if slices.Contains(c.OneOfLabels, normalizedLabelToLookFor) {
			return true
		}
	}
	return false
}
