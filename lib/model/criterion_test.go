package model

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

func TestAlternativeWithCriteria_CriterionValueFound(t *testing.T) {
	criteria := Criteria{
		{"1", Gain},
		{"2", Gain},
		{"3", Gain},
	}
	weights := Weights{"1": 1, "2": 2, "3": 3}
	for _, c := range criteria {
		v := criteria.FindWeight(&weights, &c)
		weight := weights[c.Id]
		if !utils.FloatsAreEqual(weight, v, 1e-6) {
			t.Errorf("Expected weight %f, got %f for criterion %s", weight, v, c.Id)
		}
	}
}

func TestAlternativeWithCriteria_CriterionValueMissing(t *testing.T) {
	criteria := Criteria{
		{"1", Gain},
		{"2", Gain},
		{"3", Gain},
	}

	weights := Weights{"4": 1, "5": 2, "6": 3}
	namedWeights := []namedWeight{{"4", 1}, {"5", 2}, {"6", 3}}
	for _, c := range criteria {
		//noinspection GoDeferInLoop
		defer utils.ExpectError(t, fmt.Sprintf("weight for criterion '%s' not found in criteria %v", c.Id, namedWeights))()
		criteria.FindWeight(&weights, &c)
	}
}

func TestWeights_Merge(t *testing.T) {
	oldWeights := Weights{"1": 1, "2": 2, "3": 3}
	newWeights := Weights{"4": 4, "5": 5}
	actual := *oldWeights.Merge(&newWeights)
	expected := Weights{"1": 1, "2": 2, "3": 3, "4": 4, "5": 5}
	if utils.Differs(actual, expected) {
		t.Errorf("wrong weights after merge, expected %v, got %v", expected, actual)
	}
}

func TestWeights_MergeOverlapError(t *testing.T) {
	oldWeights := Weights{"1": 1, "2": 2, "3": 3}
	newWeights := Weights{"2": 2, "4": 4}
	defer utils.ExpectError(t, "criterion '2' from [{1 1} {2 2} {3 3}] already exists in [{2 2} {4 4}]")()
	oldWeights.Merge(&newWeights)
}

func TestWeights_PreserveOnly(t *testing.T) {
	weights := Weights{"1": 1, "2": 2, "3": 3, "4": 4}
	actual := *weights.PreserveOnly(&Criteria{{Id: "1"}, {Id: "3"}})
	expected := Weights{"1": 1, "3": 3}
	if utils.Differs(expected, actual) {
		t.Errorf("different values after preserve, expected %v, got %v", expected, actual)
	}
}

func TestWeights_PreserveOnlyError(t *testing.T) {
	weights := Weights{"1": 1, "2": 2, "3": 3, "4": 4}
	defer utils.ExpectError(t, "criterion 5 not found in [{1 1} {2 2} {3 3} {4 4}]")()
	weights.PreserveOnly(&Criteria{{Id: "1"}, {Id: "5"}})
}
