package logic

import (
	"../model"
	"testing"
)

func TestChoquetIntegral(t *testing.T) {
	result := ChoquetIntegral(model.AlternativeWithCriteria{
		Alternative: model.Alternative{Id: "lama"},
		Criteria:    model.Weights{"1": 4, "2": 6, "3": 8},
	}, model.Criteria{
		model.Criterion{Id: "1", Type: model.Gain},
		model.Criterion{Id: "2", Type: model.Gain},
		model.Criterion{Id: "3", Type: model.Gain},
	}, model.Weights{
		"1,2,3": 1,
		"1,2":   0.6,
		"1,3":   0.5,
		"2,3":   0.7,
		"1":     0.1,
		"2":     0.5,
		"3":     0.4,
	})
	expectedWeight := 4*1 + (6-4)*0.7 + (8-6)*0.4
	if !FloatsAreEqual(expectedWeight, result.Value, 0.0001) {
		t.Errorf("invalid Choquet result, expected %f, got %f", expectedWeight, result.Value)
	}
}

func TestChoquetIntegral_MissingWeight(t *testing.T) {
	defer ExpectError(t, "weight for criteria union '1,2,3' not found")()
	ChoquetIntegral(model.AlternativeWithCriteria{
		Alternative: model.Alternative{Id: "lama"},
		Criteria:    model.Weights{"1": 4, "2": 6, "3": 8},
	}, model.Criteria{
		model.Criterion{Id: "1", Type: model.Gain},
		model.Criterion{Id: "2", Type: model.Gain},
		model.Criterion{Id: "3", Type: model.Gain},
	}, model.Weights{
		"1": 0.1,
		"2": 0.5,
		"3": 0.4,
	})
	t.Errorf("should fail due to the missing weight par")
}
