package model

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

func TestAlternativeWithCriteria_CriterionValueFound(t *testing.T) {
	criteria := Criteria{
		{Id: "1", Type: Gain},
		{Id: "2", Type: Gain},
		{Id: "3", Type: Gain},
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
		{Id: "1", Type: Gain},
		{Id: "2", Type: Gain},
		{Id: "3", Type: Gain},
	}

	weights := Weights{"4": 1, "5": 2, "6": 3}
	namedWeights := []namedWeight{{"4", 1}, {"5", 2}, {"6", 3}}
	for _, c := range criteria {
		func() {
			//noinspection GoDeferInLoop
			defer utils.ExpectError(t, fmt.Sprintf("weight for criterion '%s' not found in criteria %v", c.Id, namedWeights))()
			criteria.FindWeight(&weights, &c)
		}()
	}
}

func TestCriteriaValuesRangeConsidered(t *testing.T) {
	criteria := Criteria{
		{Id: "1", Type: Gain, ValuesRange: &utils.ValueRange{Min: 2, Max: 10}},
		{Id: "2", Type: Gain},
	}
	alternatives := []AlternativeWithCriteria{{
		Id:       "a",
		Criteria: Weights{"1": 1, "2": 4},
	}, {
		Id:       "b",
		Criteria: Weights{"1": 5, "2": 5},
	}}
	firstRange := CriteriaValuesRange(&alternatives, &criteria[0])
	expRange := criteria[0].ValuesRange
	if utils.Differs(firstRange, expRange) {
		t.Errorf("expected range from criteria for criterion %s (%v), got %v", criteria[0].Id, expRange, firstRange)
	}
	secondRange := CriteriaValuesRange(&alternatives, &criteria[1])
	expRange = &utils.ValueRange{Min: 4, Max: 5}
	if utils.Differs(secondRange, expRange) {
		t.Errorf("expected range from criteria for criterion %s (%v), got %v", criteria[1].Id, expRange, secondRange)
	}

}
