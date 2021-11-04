package enforcer

import (
	"fmt"
	"strings"
)

type RuleConfig struct {
	bannedLabels map[string]bool
	anyOfThis    map[string]bool
}

type Violations []string

type InternalLabel struct {
	Name        string
	Description string
}

type InternalPullRequest struct {
	RepoName string
	Number   int
	Labels   []InternalLabel
}

func IsValidPr(pr *InternalPullRequest, rules *RuleConfig) (Violations, bool) {
	report := Violations{}
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

func DefaultRules() *RuleConfig {
	return NewRules(
		[]string{"wip", "do-not-merge"},
		[]string{"bug", "feature", "enabler", "rework"},
	)
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

func (c *RuleConfig) containsAnyRequiredLabel(pr *InternalPullRequest) bool {
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

func (v Violations) String() string {
	return strings.Join(v, `, `)
}

func (pr *InternalPullRequest) UID() string {
	return fmt.Sprintf(`%s:%d`, pr.RepoName, pr.Number)
}
