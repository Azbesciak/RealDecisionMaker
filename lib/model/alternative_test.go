package model

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

func TestPreserveCriteriaForAlternatives(t *testing.T) {
	testAlternatives := []AlternativeWithCriteria{
		{Id: "1", Criteria: Weights{"1": 10, "2": 20, "3": 20, "4": 5}},
		{Id: "2", Criteria: Weights{"1": 20, "2": 30, "3": 10, "4": 0}},
		{Id: "3", Criteria: Weights{"1": 5, "2": 45, "3": 2, "4": 8}},
	}
	checkAlternativeHasOnlyCriteria(t, testAlternatives, []string{"1", "4"})
	checkAlternativeHasOnlyCriteria(t, testAlternatives, []string{"2", "3"})
	checkAlternativeHasOnlyCriteria(t, testAlternatives, []string{"3"})
	defer utils.ExpectError(t, "alternative '1' does not have value for criterion '789'")()
	checkAlternativeHasOnlyCriteria(t, testAlternatives, []string{"789"})
}

func checkAlternativeHasOnlyCriteria(t *testing.T, alternatives []AlternativeWithCriteria, criteria []string) {
	crt := generateCriteria(criteria)
	result := *PreserveCriteriaForAlternatives(&alternatives, &crt)
	if len(result) != len(alternatives) {
		t.Errorf("invalid length of returned alternatives, expected %d, got %d", len(alternatives), len(result))
		return
	}
	for i, r := range result {
		originalAlternative := alternatives[i]
		validateCriteriaCount(t, r, criteria, originalAlternative)
		validateCriteriaPresence(t, criteria, r, originalAlternative)
	}
}

func generateCriteria(criteria []string) Criteria {
	crt := make(Criteria, len(criteria))
	for i, c := range criteria {
		crt[i] = Criterion{Id: c, Type: Gain}
	}
	return crt
}

func validateCriteriaCount(t *testing.T, r AlternativeWithCriteria, criteria []string, originalAlternative AlternativeWithCriteria) {
	if len(r.Criteria) != len(criteria) {
		t.Errorf(
			"expected %d criteria for alternative %v, got %d: %v",
			len(criteria), originalAlternative, len(r.Criteria), r.Criteria,
		)
	}
}

func validateCriteriaPresence(t *testing.T, criteria []string, r AlternativeWithCriteria, originalAlternative AlternativeWithCriteria) {
	for _, criterion := range criteria {
		value, ok := r.Criteria[criterion]
		if !ok {
			t.Errorf("criterion '%s' not found in criteria %v", criterion, r.Criteria)
		}
		originalValue := originalAlternative.Criteria[criterion]
		if value != originalValue {
			t.Errorf(
				"expected value %f for criterion %s and alternative %v, got %f",
				originalValue, criterion, originalAlternative, value,
			)
		}
	}
}
