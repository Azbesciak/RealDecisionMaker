package client

import (
	. "../logic/owa"
	. "../logic/weighted-sum"
	. "../model"
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
