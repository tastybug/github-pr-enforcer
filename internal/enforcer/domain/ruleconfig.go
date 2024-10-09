package domain

import "strings"

type RuleConfig struct {
	NeedsNoneOf map[string]bool
	NeedsOneOf  map[string]bool
}

func (c *RuleConfig) ContainsBannedLabel(label string) bool {
	return c.NeedsNoneOf[label]
}

func (c *RuleConfig) ContainsAnyRequiredLabel(pr *PullRequest) bool {
	if len(c.NeedsOneOf) == 0 {
		return true
	}
	matchesAnyLabel := false
	for _, label := range pr.Labels {
		l := strings.ToLower(label.Name)
		matchesAnyLabel = matchesAnyLabel || c.NeedsOneOf[l]
	}
	return matchesAnyLabel
}

func CreateRuleConfig(bannedLabels []string, anyOfTheseLabels []string) *RuleConfig {
	config := emptyRuleConfig()

	for _, banned := range bannedLabels {
		config.NeedsNoneOf[strings.ToLower(banned)] = true
	}
	for _, anyOfThis := range anyOfTheseLabels {
		config.NeedsOneOf[strings.ToLower(anyOfThis)] = true
	}
	return &config
}

func emptyRuleConfig() RuleConfig {
	return RuleConfig{
		make(map[string]bool),
		make(map[string]bool),
	}
}

func (c *RuleConfig) BannedAsList() []string {
	l := make([]string, 0)
	for key, _ := range c.NeedsNoneOf {
		l = append(l, key)
	}
	return l
}

func (c *RuleConfig) AnyOfThisAsList() []string {
	l := make([]string, 0)
	for key, _ := range c.NeedsOneOf {
		l = append(l, key)
	}
	return l
}
