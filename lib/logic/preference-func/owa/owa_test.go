package owa

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

func TestOWAResult(t *testing.T) {
	criteria := model.AlternativeWithCriteria{
		Id: "test",
		Criteria: model.Weights{
			"1": 30,
			"2": 10,
			"3": 5,
			"4": 60,
		},
	}
	result := OWA(criteria, model.WeightedCriteria{
		{Criterion: model.Criterion{Id: "1"}, Weight: 1},
		{Criterion: model.Criterion{Id: "2"}, Weight: 2},
		{Criterion: model.Criterion{Id: "3"}, Weight: 3},
		{Criterion: model.Criterion{Id: "4"}, Weight: 4},
	})
	var expected model.Weight = 4*60 + 30*3 + 10*2 + 5*1
	if result.Value() != expected {
		t.Errorf("Expected %f, got %f", expected, result.Value())
	}
	result = OWA(criteria, model.WeightedCriteria{
		{Criterion: model.Criterion{Id: "1"}, Weight: 6},
		{Criterion: model.Criterion{Id: "2"}, Weight: 4},
		{Criterion: model.Criterion{Id: "3"}, Weight: 3},
		{Criterion: model.Criterion{Id: "4"}, Weight: 4},
	})
	expected = 6*60 + 30*4 + 10*4 + 5*3
	if result.Value() != expected {
		t.Errorf("Expected %f, got %f", expected, result.Value())
	}
}

func TestOWAInvalidArguments(t *testing.T) {
	criteria := model.AlternativeWithCriteria{
		Id: "test",
		Criteria: model.Weights{
			"1": 30,
			"2": 10,
		},
	}
	defer utils.ExpectError(t, "criteria and Weights must have the same length, got 2 and 1")()
	OWA(criteria, model.WeightedCriteria{
		{Criterion: model.Criterion{Id: "1"}, Weight: 1},
	})
}
