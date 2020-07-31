package choquet

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

var testCriteria = model.Criteria{
	model.Criterion{Id: "1", Type: model.Gain},
	model.Criterion{Id: "2", Type: model.Gain},
	model.Criterion{Id: "3", Type: model.Gain},
}

func TestChoquetIntegral(t *testing.T) {
	result := ChoquetIntegral(model.AlternativeWithCriteria{
		Id:       "lama",
		Criteria: model.Weights{"1": 4, "2": 6, "3": 8},
	},
		testCriteria,
		model.Weights{
			"1,2,3": 1,
			"1,2":   0.6,
			"1,3":   0.5,
			"2,3":   0.7,
			"1":     0.1,
			"2":     0.5,
			"3":     0.4,
		})
	expectedWeight := 4*1 + (6-4)*0.7 + (8-6)*0.4
	if !utils.FloatsAreEqual(expectedWeight, result.Value(), 0.0001) {
		t.Errorf("invalid Choquet result, expected %f, got %f", expectedWeight, result.Value())
	}
}

func TestChoquetIntegral_MissingWeight(t *testing.T) {
	defer utils.ExpectError(t, "weight for criteria union '1,2' not found")()
	ChoquetIntegral(model.AlternativeWithCriteria{
		Id:       "lama",
		Criteria: model.Weights{"1": 4, "2": 6, "3": 8},
	},
		testCriteria,
		model.Weights{
			"1": 0.1,
			"2": 0.5,
			"3": 0.4,
		})
	t.Errorf("should fail due to the missing weight par")
}

func TestCriteriaPresenceValidation(t *testing.T) {
	checkError(t, model.Weights{"1": 0.1, "2": 0.5, "1,2": 0.4}, []string{"1", "2"}, "")
	checkError(t, model.Weights{"1": 0.1, "2": 0.5, "2,1": 0.4}, []string{"1", "2"}, "weight for criteria union '1,2' not found")
	checkError(t, model.Weights{"1": 0.1, "1,2": 0.4}, []string{"1", "2"}, "weight for criteria union '2' not found")
	checkError(t, model.Weights{"1": 0.1, "2": 0.4}, []string{"1", "2"}, "weight for criteria union '1,2' not found")
	checkError(t, model.Weights{"1": 0.1, "1,2": 0.4}, []string{"1"}, "")
	checkError(t, model.Weights{"2": 0.1}, []string{"1"}, "weight for criteria union '1' not found")
}

func TestCriteriaAligning(t *testing.T) {
	actual := remapWeights(&model.Weights{"1": 1, "2": 0.5, "3": 0.75, "2,1": 0.6, "1,3": 0.7, "3,2": 0.1, "1,3,2": 0.9})
	expected := &model.Weights{"1": 1, "2": 0.5, "3": 0.75, "1,2": 0.6, "1,3": 0.7, "2,3": 0.1, "1,2,3": 0.9}
	if utils.Differs(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestDuplicatedCriteria(t *testing.T) {
	defer utils.ExpectError(t, "values for criteria 1,2 are redeclared")()
	remapWeights(&model.Weights{"1": 1, "2": 0.5, "1,2": 0.4, "2,1": 0.4})
}

func checkError(t *testing.T, weights model.Weights, criteria []string, error string) {
	_criteria := make(model.Criteria, len(criteria))
	for i, c := range criteria {
		_criteria[i] = model.Criterion{Id: c, Type: model.Gain}
	}
	if len(error) != 0 {
		defer utils.ExpectError(t, error)()
	}
	validateAllWeightsAvailable(&weights, &_criteria)
}
