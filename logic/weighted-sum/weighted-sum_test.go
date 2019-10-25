package weighted_sum

import (
	"../../model"
	"../../utils"
	"testing"
)

func TestValidFunc(t *testing.T) {
	criteria := []model.WeightedCriterion{{
		model.Criterion{Id: "Cost", Type: model.Cost}, 1,
	}, {
		model.Criterion{Id: "Color", Type: model.Gain}, 2,
	}}
	alternative := model.AlternativeWithCriteria{
		Id: "Ferrari",
		Criteria: model.Weights{
			"Cost":  200,
			"Color": 10,
		},
	}
	res := WeightedSum(alternative, criteria)
	if !utils.FloatsAreEqual(-190.0, res.Value, 0.01) {
		t.Error("invalid result", res.Value)
	}
}

func TestMissingCriterion(t *testing.T) {
	criteria := []model.WeightedCriterion{{model.Criterion{Id: "Color", Type: model.Gain}, 1}}
	var alternative = model.AlternativeWithCriteria{
		Id: "Ferrari",
		Criteria: model.Weights{
			"Cost": 200,
		},
	}
	defer utils.ExpectError(t, "alternative 'Ferrari' does not have value for criterion 'Color'")()
	WeightedSum(alternative, criteria)
}
