package domain

import "strings"

type RuleConfig struct {
	BannedLabels map[string]bool
	AnyOfThis    map[string]bool
}

func (c *RuleConfig) ContainsBannedLabel(label string) bool {
	return c.BannedLabels[label]
}

func (c *RuleConfig) ContainsAnyRequiredLabel(pr *PullRequest) bool {
	if len(c.AnyOfThis) == 0 {
		return true
	}
	matchesAnyLabel := false
	for _, label := range pr.Labels {
		l := strings.ToLower(label.Name)
		matchesAnyLabel = matchesAnyLabel || c.AnyOfThis[l]
	}
	return matchesAnyLabel
}

func (c *RuleConfig) BannedAsList() []string {
	l := make([]string, 0)
	for key, _ := range c.BannedLabels {
		l = append(l, key)
	}
	return l
}

func (c *RuleConfig) AnyOfThisAsList() []string {
	l := make([]string, 0)
	for key, _ := range c.AnyOfThis {
		l = append(l, key)
	}
	return l
}
