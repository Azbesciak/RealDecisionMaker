package preference_reversal

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/model/criteria-ordering"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

var bias = PreferenceReversal{
	orderingResolvers: []criteria_ordering.CriteriaOrderingResolver{
		&criteria_ordering.WeakestCriteriaOrderingResolver{},
		&criteria_ordering.StrongestCriteriaOrderingResolver{},
	},
}
var params = model.DecisionMakingParams{
	NotConsideredAlternatives: []model.AlternativeWithCriteria{{
		Id:       "a",
		Criteria: model.Weights{"1": 5, "2": 1},
	}},
	ConsideredAlternatives: []model.AlternativeWithCriteria{{
		Id:       "b",
		Criteria: model.Weights{"1": 15, "2": 15},
	}, {
		Id:       "c",
		Criteria: model.Weights{"1": 10, "2": 21},
	}},
	Criteria: model.Criteria{{
		Id:   "1",
		Type: model.Gain,
		ValuesRange: &utils.ValueRange{
			Min: 10,
			Max: 20,
		},
	}, {
		Id:   "2",
		Type: model.Cost,
	}},
}

func TestPreferenceReversal_ApplyWeakest(t *testing.T) {
	expectedParams := model.DecisionMakingParams{
		NotConsideredAlternatives: []model.AlternativeWithCriteria{{
			Id:       "a",
			Criteria: model.Weights{"1": 25, "2": 21},
		}},
		ConsideredAlternatives: []model.AlternativeWithCriteria{{
			Id:       "b",
			Criteria: model.Weights{"1": 15, "2": 7},
		}, {
			Id:       "c",
			Criteria: model.Weights{"1": 20, "2": 1},
		}},
		Criteria: model.Criteria{{
			Id:   "1",
			Type: model.Gain,
			ValuesRange: &utils.ValueRange{
				Min: 10,
				Max: 20,
			},
		}, {
			Id:   "2",
			Type: model.Cost,
		}},
	}
	expectedResult := PreferenceReversalResult{
		ReversedPreferenceCriteria: []ReversedPreferenceCriterion{{
			Id:          "1",
			Type:        model.Gain,
			ValuesRange: utils.ValueRange{Min: 10, Max: 20},
			AlternativesValues: model.Weights{
				"a": 25, "b": 15, "c": 20,
			},
		}, {
			Id:          "2",
			Type:        model.Cost,
			ValuesRange: utils.ValueRange{Min: 1, Max: 21},
			AlternativesValues: model.Weights{
				"a": 21, "b": 7, "c": 1,
			},
		}},
	}
	var props model.BiasProps = utils.Map{
		"ordering": "weakest",
		"ratio":    1,
	}
	checkPreferenceReversal(t, bias, params, props, expectedParams, expectedResult)
}

func TestPreferenceReversal_ApplyWeakestJustOne(t *testing.T) {
	expectedParams := model.DecisionMakingParams{
		NotConsideredAlternatives: []model.AlternativeWithCriteria{{
			Id:       "a",
			Criteria: model.Weights{"1": 25, "2": 1},
		}},
		ConsideredAlternatives: []model.AlternativeWithCriteria{{
			Id:       "b",
			Criteria: model.Weights{"1": 15, "2": 15},
		}, {
			Id:       "c",
			Criteria: model.Weights{"1": 20, "2": 21},
		}},
		Criteria: model.Criteria{{
			Id:   "1",
			Type: model.Gain,
			ValuesRange: &utils.ValueRange{
				Min: 10,
				Max: 20,
			},
		}, {
			Id:   "2",
			Type: model.Cost,
		}},
	}
	expectedResult := PreferenceReversalResult{
		ReversedPreferenceCriteria: []ReversedPreferenceCriterion{{
			Id:          "1",
			Type:        model.Gain,
			ValuesRange: utils.ValueRange{Min: 10, Max: 20},
			AlternativesValues: model.Weights{
				"a": 25, "b": 15, "c": 20,
			},
		}},
	}
	var props model.BiasProps = utils.Map{
		"ordering": "weakest",
		"ratio":    0.5,
		"min":      1,
	}
	checkPreferenceReversal(t, bias, params, props, expectedParams, expectedResult)
}

func TestPreferenceReversal_ApplyStrongest(t *testing.T) {
	expectedParams := model.DecisionMakingParams{
		NotConsideredAlternatives: []model.AlternativeWithCriteria{{
			Id:       "a",
			Criteria: model.Weights{"1": 25, "2": 21},
		}},
		ConsideredAlternatives: []model.AlternativeWithCriteria{{
			Id:       "b",
			Criteria: model.Weights{"1": 15, "2": 7},
		}, {
			Id:       "c",
			Criteria: model.Weights{"1": 20, "2": 1},
		}},
		Criteria: model.Criteria{{
			Id:   "1",
			Type: model.Gain,
			ValuesRange: &utils.ValueRange{
				Min: 10,
				Max: 20,
			},
		}, {
			Id:   "2",
			Type: model.Cost,
		}},
	}
	expectedResult := PreferenceReversalResult{
		ReversedPreferenceCriteria: []ReversedPreferenceCriterion{{
			Id:          "2",
			Type:        model.Cost,
			ValuesRange: utils.ValueRange{Min: 1, Max: 21},
			AlternativesValues: model.Weights{
				"a": 21, "b": 7, "c": 1,
			},
		}, {
			Id:          "1",
			Type:        model.Gain,
			ValuesRange: utils.ValueRange{Min: 10, Max: 20},
			AlternativesValues: model.Weights{
				"a": 25, "b": 15, "c": 20,
			},
		}},
	}
	var props model.BiasProps = utils.Map{
		"ordering": "strongest",
		"ratio":    1,
	}
	checkPreferenceReversal(t, bias, params, props, expectedParams, expectedResult)
}

func checkPreferenceReversal(t *testing.T, bias PreferenceReversal, params model.DecisionMakingParams, props model.BiasProps, expectedParams model.DecisionMakingParams, expectedResult PreferenceReversalResult) {
	listener := model.BiasListener(&testUtils.DummyBiasListener{})
	actual := bias.Apply(&params, &params, &props, &listener)
	if utils.Differs(expectedParams.ConsideredAlternatives, actual.DMP.ConsideredAlternatives) {
		t.Errorf("different considered alternatives, expected %v, got %v", expectedParams.ConsideredAlternatives, actual.DMP.ConsideredAlternatives)
	}
	if utils.Differs(expectedParams.NotConsideredAlternatives, actual.DMP.NotConsideredAlternatives) {
		t.Errorf("different not considered alternatives, expected %v, got %v", expectedParams.NotConsideredAlternatives, actual.DMP.NotConsideredAlternatives)
	}
	if utils.Differs(expectedParams.Criteria, actual.DMP.Criteria) {
		t.Errorf("different criteria, expected %v, got %v", expectedParams.Criteria, actual.DMP.Criteria)
	}
	if utils.Differs(expectedResult, actual.Props) {
		t.Errorf("different bias result, expected %v, got %v", expectedResult, actual)
	}
}
