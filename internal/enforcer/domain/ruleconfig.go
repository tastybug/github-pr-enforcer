package domain

type RuleConfig struct {
	NoneOfLabels []string `json:"banned"`
	OneOfLabels  []string `json:"needs-one-of"`
}
