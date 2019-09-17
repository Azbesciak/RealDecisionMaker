package logic

import (
	"math"
	"testing"
)

func TestValidFunc(t *testing.T) {
	criteria := []WeightedCriterion{{
		Criterion{"Cost", Cost}, 1,
	}, {
		Criterion{"Color", Gain}, 2,
	}}
	alternative := AlternativeWithCriteria{
		Alternative{"Ferrari"},
		map[string]Weight{
			"Cost":  200,
			"Color": 10,
		},
	}
	res := WeightedSum(alternative, criteria)
	if !FloatsEquals(-190.0, float64(res.Value), 0.01) {
		t.Error("invalid result", res.Value)
	}
}

func TestMissingCriterion(t *testing.T) {
	criteria := []WeightedCriterion{{Criterion{"Color", Gain}, 1}}
	alternative := AlternativeWithCriteria{
		Alternative{"Ferrari"},
		map[string]Weight{
			"Cost": 200,
		},
	}
	defer func() {
		if err := recover(); err == nil {
			t.Error("should throw error")
		} else if err != "criterion 'Color' not found in criteria" {
			t.Error("Invalid error message", err)
		}
	}()
	WeightedSum(alternative, criteria)
}

func FloatsEquals(expected float64, actual Weight, epsilon Weight) bool {
	return math.Abs(expected-actual) <= epsilon
}
