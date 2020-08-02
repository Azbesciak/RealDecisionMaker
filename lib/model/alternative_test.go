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

func TestAlternativesShuffle(t *testing.T) {
	alternatives := []AlternativeWithCriteria{{Id: "1"}, {Id: "2"}, {Id: "3"}}
	current := -1.0
	generator := func() float64 {
		current++
		return current / float64(len(alternatives))
	}
	actual := ShuffleAlternatives(&alternatives, generator)
	expected := []string{"2", "3", "1"}
	validateAlternativesOrder(t, actual, expected)

	current = 3.0
	generator = func() float64 {
		current--
		return current / float64(len(alternatives))
	}
	actual = ShuffleAlternatives(&alternatives, generator)
	expected = []string{"3", "1", "2"}
	validateAlternativesOrder(t, actual, expected)
}

func validateAlternativesOrder(t *testing.T, actual *[]AlternativeWithCriteria, expected []string) {
	if len(*actual) != len(expected) {
		t.Errorf("different len after sort, expected %d, got %d", len(expected), len(*actual))
	} else {
		for i, a := range *actual {
			exp := expected[i]
			if exp != a.Id {
				t.Errorf("different alternative, index %d, expected %s, got %s", i, exp, a.Id)
			}
		}
	}
}

func TestRemoveAlternative(t *testing.T) {
	checkRemoval(t, []string{"1", "2", "3"}, "3", []string{"1", "2"})
	checkRemoval(t, []string{"1", "2", "3"}, "4", []string{"1", "2", "3"})
	checkRemoval(t, []string{"1", "2", "3"}, "1", []string{"2", "3"})
	checkRemoval(t, []string{"1", "2", "3"}, "2", []string{"1", "3"})
}

func checkRemoval(t *testing.T, ids []string, toRemove string, expected []string) {
	all := make([]AlternativeWithCriteria, len(ids))
	for i, id := range ids {
		all[i] = AlternativeWithCriteria{Id: id}
	}
	actual := RemoveAlternative(all, AlternativeWithCriteria{Id: toRemove})
	validateAlternativesOrder(t, &actual, expected)
}
