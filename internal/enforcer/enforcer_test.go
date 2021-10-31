package enforcer

import (
	"testing"
)

func TestEmptyRulesMeansEverythingIsValid(t *testing.T) {
	testdata := []struct {
		pr    *InternalPullRequest
		valid bool
	}{
		{&InternalPullRequest{}, true},
		{aPrWithLabels([]string{"bug"}), true},
	}

	for _, data := range testdata {
		expected := data.valid
		if viol, actual := IsValidPr(data.pr, &RuleConfig{}); expected != actual {
			t.Errorf("%+v validity was expected as %t, but was %t. Violations: %v", data.pr, expected, actual, viol)
		}
	}
}

func TestBannedLabels(t *testing.T) {
	testdata := []struct {
		pr    *InternalPullRequest
		valid bool
	}{
		{aPrWithLabels([]string{"ok", "banned"}), false},
		{aPrWithLabels([]string{"bug", "ok"}), true},
		{aPrWithLabels([]string{"banned", "banned1", "bug"}), false},
		{aPrWithLabels([]string{"banned", "BANNED1"}), false}, // labels are case insensitive
	}
	rules := NewRules([]string{"banned", "banned1", "banned2"}, []string{"bug"})

	for _, data := range testdata {
		expected := data.valid
		if viol, actual := IsValidPr(data.pr, rules); expected != actual {
			t.Errorf("%+v validity was expected as %t, but was %t. violations: %v", data.pr, expected, actual, viol)
		}
	}
}

func TestAnyOfLabels(t *testing.T) {
	testdata := []struct {
		pr    *InternalPullRequest
		valid bool
	}{
		{aPrWithLabels([]string{"bog"}), false},
		{aPrWithLabels([]string{"bug"}), true},
		{aPrWithLabels([]string{"bug", "feature"}), true},
		{aPrWithLabels([]string{"FEATURE"}), true}, // labels are case insensitive
	}
	rules := NewRules([]string{}, []string{"bug", "feature"})

	for _, data := range testdata {
		expected := data.valid
		if viol, actual := IsValidPr(data.pr, rules); expected != actual {
			t.Errorf("%+v validity was expected as %t, but was %t. Violations: %v", data.pr, expected, actual, viol)
		}
	}
}

func aPrWithLabels(labels []string) *InternalPullRequest {
	pr := InternalPullRequest{}
	for _, l := range labels {
		pr.Labels = append(pr.Labels, InternalLabel{Name: l})
	}
	return &pr
}
