package domain

type RuleConfig struct {
	BannedLabels     []string `json:"banned"`
	AnyOfTheseLabels []string `json:"needs-one-of"`
}
