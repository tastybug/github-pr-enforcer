package domain

import "strings"

type RuleConfig struct {
	MustNotHaveOneOf map[string]bool
	NeedsOneOf       map[string]bool
}

func CreateRuleConfig(bannedLabels []string, anyOfTheseLabels []string) *RuleConfig {
	config := EmptyRuleConfig()

	for _, banned := range bannedLabels {
		config.MustNotHaveOneOf[strings.ToLower(banned)] = true
	}
	for _, anyOfThis := range anyOfTheseLabels {
		config.NeedsOneOf[strings.ToLower(anyOfThis)] = true
	}
	return &config
}

func EmptyRuleConfig() RuleConfig {
	return RuleConfig{
		make(map[string]bool),
		make(map[string]bool),
	}
}

func (c *RuleConfig) BannedAsList() []string {
	l := make([]string, 0)
	for key, _ := range c.MustNotHaveOneOf {
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
