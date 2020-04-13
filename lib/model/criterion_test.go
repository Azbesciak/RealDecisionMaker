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
