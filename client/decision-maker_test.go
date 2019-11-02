package client

import (
	. "github.com/Azbesciak/RealDecisionMaker/logic/electreIII"
	. "github.com/Azbesciak/RealDecisionMaker/logic/owa"
	. "github.com/Azbesciak/RealDecisionMaker/logic/weighted-sum"
	. "github.com/Azbesciak/RealDecisionMaker/model"
	"reflect"
	"testing"
)

func TestDecisionMaker_MakeDecision_WeightedSum(t *testing.T) {
	dm := DecisionMaker{
		PreferenceFunction: "weightedSum",
		State:              DecisionMakerState{},
		KnownAlternatives: []AlternativeWithCriteria{
			{"1", Weights{"1": 10, "2": 20}},
			{"2", Weights{"1": 20, "2": 10}},
		},
		ChoseToMake:      []Alternative{"1", "2"},
		Criteria:         Criteria{Criterion{"1", "gain"}, Criterion{"2", "cost"}},
		MethodParameters: map[string]interface{}{"weights": Weights{"1": 1, "2": 2}},
	}
	weightedSum := &WeightedSumPreferenceFunc{}
	owa := &OWAPreferenceFunc{}
	funcs := PreferenceFunctions{Functions: []PreferenceFunction{weightedSum, owa}}
	ranking := dm.MakeDecision(funcs)
	expectedRanking := weightedSum.Evaluate(&dm)

	if !reflect.DeepEqual(ranking.Result, *expectedRanking) {
		t.Errorf("Invalid result, expected '%v', got '%v'", ranking.Result, *expectedRanking)
	}
}

func TestDecisionMaker_MakeDecision_ElectreIII(t *testing.T) {
	dm := DecisionMaker{
		PreferenceFunction: "electreIII",
		State:              DecisionMakerState{},
		KnownAlternatives: []AlternativeWithCriteria{
			{"1", Weights{"1": 10, "2": 20}},
			{"2", Weights{"1": 20, "2": 10}},
		},
		ChoseToMake: []Alternative{"1", "2"},
		Criteria:    Criteria{Criterion{"1", "gain"}, Criterion{"2", "cost"}},
		MethodParameters: map[string]interface{}{
			"electreCriteria": ElectreCriteria{
				"1": {
					K: 1,
					Q: LinearFunctionParameters{A: 0, B: 2},
					P: LinearFunctionParameters{A: 1, B: 4},
					V: LinearFunctionParameters{A: 0, B: 10},
				},
				"2": {
					K: 2,
					Q: LinearFunctionParameters{A: 0, B: 10},
					P: LinearFunctionParameters{A: 0, B: 12},
					V: LinearFunctionParameters{A: 0, B: 14},
				},
			},
			"electreDistillation": LinearFunctionParameters{
				A: -0.1,
				B: 0.5,
			}},
	}
	weightedSum := &WeightedSumPreferenceFunc{}
	ele := &ElectreIIIPreferenceFunc{}
	funcs := PreferenceFunctions{Functions: []PreferenceFunction{weightedSum, ele}}
	ranking := dm.MakeDecision(funcs)
	expectedRanking := ele.Evaluate(&dm)

	if !reflect.DeepEqual(ranking.Result, *expectedRanking) {
		t.Errorf("Invalid result, expected '%v', got '%v'", ranking.Result, *expectedRanking)
	}
}
