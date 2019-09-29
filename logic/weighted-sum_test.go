package logic

import (
	"../model"
	"testing"
)

func TestValidFunc(t *testing.T) {
	criteria := []model.WeightedCriterion{{
		model.Criterion{Id: "Cost", Type: model.Cost}, 1,
	}, {
		model.Criterion{Id: "Color", Type: model.Gain}, 2,
	}}
	alternative := model.AlternativeWithCriteria{
		Alternative: model.Alternative{Id: "Ferrari"},
		Criteria: model.Weights{
			"Cost":  200,
			"Color": 10,
		},
	}
	res := WeightedSum(alternative, criteria)
	if !FloatsAreEqual(-190.0, res.Value, 0.01) {
		t.Error("invalid result", res.Value)
	}
}

func TestMissingCriterion(t *testing.T) {
	criteria := []model.WeightedCriterion{{model.Criterion{Id: "Color", Type: model.Gain}, 1}}
	var alternative = model.AlternativeWithCriteria{
		Alternative: model.Alternative{Id: "Ferrari"},
		Criteria: model.Weights{
			"Cost": 200,
		},
	}
	defer ExpectError(t, "criterion 'Color' not found in criteria")()
	WeightedSum(alternative, criteria)
}
