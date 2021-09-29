package enforcer

import (
	"testing"
)

func TestEmptyRulesMeansEverythingIsValid(t *testing.T) {
	testdata := []*PullRequest{
		&PullRequest{},
		aPrWithLabels([]string{"bug"}),
	}

	for _, data := range testdata {
		if IsValidPr(data, &RuleConfig{}) != true {
			t.Errorf("PR %+v unexpectedly invalid", *data)
		}
	}
}

func aPrWithLabels(labels []string) *PullRequest {
	pr := PullRequest{}
	for _, l := range labels {
		pr.Labels = append(pr.Labels, Label{Name: l})
	}
	return &pr
}
