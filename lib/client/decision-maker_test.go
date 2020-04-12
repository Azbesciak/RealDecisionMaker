package client

import (
	. "github.com/Azbesciak/RealDecisionMaker/lib/logic/electreIII"
	. "github.com/Azbesciak/RealDecisionMaker/lib/logic/owa"
	. "github.com/Azbesciak/RealDecisionMaker/lib/logic/weighted-sum"
	. "github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"reflect"
	"testing"
)

func TestDecisionMaker_MakeDecision_WeightedSum(t *testing.T) {
	dm := DecisionMaker{
		PreferenceFunction: "weightedSum",
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
	ranking := dm.MakeDecision(funcs, BiasListeners{}, &BiasMap{})
	expectedRanking := weightedSum.Evaluate(&DecisionMakingParams{
		NotConsideredAlternatives: nil,
		ConsideredAlternatives:    dm.KnownAlternatives,
		Criteria:                  dm.Criteria,
		MethodParameters:          weightedSum.ParseParams(&dm),
	})

	if !reflect.DeepEqual(ranking.Result, *expectedRanking) {
		t.Errorf("Invalid result, expected '%v', got '%v'", ranking.Result, *expectedRanking)
	}
}

func TestDecisionMaker_MakeDecision_ElectreIII(t *testing.T) {
	dm := DecisionMaker{
		PreferenceFunction: "electreIII",
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
	ranking := dm.MakeDecision(funcs, BiasListeners{}, &BiasMap{})
	expectedRanking := ele.Evaluate(&DecisionMakingParams{
		NotConsideredAlternatives: nil,
		ConsideredAlternatives:    dm.KnownAlternatives,
		Criteria:                  dm.Criteria,
		MethodParameters:          ele.ParseParams(&dm),
	})

	if !reflect.DeepEqual(ranking.Result, *expectedRanking) {
		t.Errorf("Invalid result, expected '%v', got '%v'", ranking.Result, *expectedRanking)
	}
}

func TestFetchPreferenceFunc(t *testing.T) {
	functions := PreferenceFunctions{Functions: []PreferenceFunction{
		&ElectreIIIPreferenceFunc{},
		&OWAPreferenceFunc{},
	}}
	defer utils.ExpectError(t, "preference function 'xyz' not found, available are '[electreIII owa]'")()
	functions.Fetch("xyz")
}
