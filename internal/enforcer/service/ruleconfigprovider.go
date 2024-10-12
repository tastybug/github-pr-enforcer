package service

import (
	"encoding/json"
	"fmt"
	"github.com/tastybug/github-pr-enforcer/internal/enforcer/domain"
	"os"
	"strings"
)

const MUST_HAVE_ONE_OF = "MUST_HAVE_ONE_OF"
const MUST_HAVE_NONE_OF = "MUST_HAVE_NONE_OF"

func GetRules(customRules *domain.RuleConfig) *domain.RuleConfig {
	if customRules != nil {
		fmt.Printf("Using custom rules: %+v\n", customRules)
		return customRules
	} else if rulesFromEnv := GetRulesFromEnv(); rulesFromEnv != nil {
		fmt.Printf("Using env based rules: %+v\n", rulesFromEnv)
		return rulesFromEnv
	} else {
		fallbackRules := getFallbackRules()
		fmt.Printf("Using fallback rules: %+v\n", fallbackRules)
		return fallbackRules
	}
}

func getFallbackRules() *domain.RuleConfig {
	return &domain.RuleConfig{
		[]string{"wip", "do-not-merge"},
		[]string{"bug", "feature", "enabler", "rework"},
	}
}

func GetRulesFromJsonString(rulesJson string) (*domain.RuleConfig, error) {
	var ruleConfig domain.RuleConfig
	if err := json.NewDecoder(strings.NewReader(rulesJson)).Decode(&ruleConfig); err != nil {
		return nil, fmt.Errorf("Query param rule set broken: %s\n", err)
	} else {
		return &ruleConfig, nil
	}
}

func GetRulesFromEnv() *domain.RuleConfig {
	_, existsA := os.LookupEnv(MUST_HAVE_ONE_OF)
	_, existsB := os.LookupEnv(MUST_HAVE_NONE_OF)

	if !(existsA || existsB) {
		return nil
	}

	var anyOfTheseLabels []string
	var bannedLabels []string
	if val, exists := os.LookupEnv(MUST_HAVE_ONE_OF); exists {
		anyOfTheseLabels = strings.Split(val, ",")
	}
	if val, exists := os.LookupEnv(MUST_HAVE_NONE_OF); exists {
		bannedLabels = strings.Split(val, ",")
	}

	return &domain.RuleConfig{bannedLabels, anyOfTheseLabels}
}
