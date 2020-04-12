package fatigue

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

const testEps = 1e-6

func TestFatigue_Apply(t *testing.T) {
	mixing := Fatigue{
		valueGeneratorSource: testUtils.CyclicRandomGenerator(0, 10),
		signGeneratorSource:  testUtils.CyclicRandomGenerator(0, 2),
		functions:            []FatigueFunction{&ConstFatigueFunction{}},
	}
	notConsidered := []model.AlternativeWithCriteria{
		{Id: "x", Criteria: model.Weights{"1": 1, "2": 2, "3": 3}},
		{Id: "y", Criteria: model.Weights{"1": 6, "2": 20, "3": 4}},
	}
	considered := []model.AlternativeWithCriteria{
		{Id: "a", Criteria: model.Weights{"1": 7, "2": 3, "3": 1}},
		{Id: "b", Criteria: model.Weights{"1": 6, "2": 5, "3": 10}},
	}
	criteria := model.Criteria{
		{Id: "1", Type: model.Gain},
		{Id: "2", Type: model.Gain},
		{Id: "3", Type: model.Gain},
	}
	listener := model.BiasListener(&testUtils.DummyBiasListener{})
	m := model.BiasProps(map[string]interface{}{
		"function": FatConstFunc,
		"params": map[string]interface{}{
			"value": 0.5,
		},
		"randomSeed": 123,
	})
	original := &model.DecisionMakingParams{
		NotConsideredAlternatives: notConsidered,
		ConsideredAlternatives:    considered,
		Criteria:                  criteria,
		MethodParameters:          nil,
	}
	result := mixing.Apply(original, original, &m, &listener)
	expected := FatigueResult{
		EffectiveFatigueRatio: .5,
		ConsideredAlternatives: []model.AlternativeWithCriteria{
			{Id: "a", Criteria: model.Weights{"1": 6.65, "2": 3.3, "3": 0.85}},
			{Id: "b", Criteria: model.Weights{"1": 7.2, "2": 3.75, "3": 13}},
		},
		NotConsideredAlternatives: []model.AlternativeWithCriteria{
			{Id: "x", Criteria: model.Weights{"1": 0.65, "2": 2.8, "3": 1.65}},
			{Id: "y", Criteria: model.Weights{"1": 6, "2": 19, "3": 4.4}},
		},
	}
	checkProps(t, expected, result.Props)
}

func checkProps(t *testing.T, expected FatigueResult, actual model.BiasProps) {
	v, ok := actual.(FatigueResult)
	if !ok {
		t.Errorf("expected instance of FatigueResult")
		return
	}
	if !utils.FloatsAreEqual(expected.EffectiveFatigueRatio, v.EffectiveFatigueRatio, testEps) {
		t.Errorf("expected fatigue: %v, got %v", expected.EffectiveFatigueRatio, v.EffectiveFatigueRatio)
	}
	checkAlternatives(t, "considered", expected.ConsideredAlternatives, v.ConsideredAlternatives)
	checkAlternatives(t, "notConsidered", expected.NotConsideredAlternatives, v.NotConsideredAlternatives)
}

func checkAlternatives(t *testing.T, name string, expected, actual []model.AlternativeWithCriteria) {
	if len(expected) != len(actual) {
		t.Errorf("%s: expected %d alternatives , got %d; %v vs %v", name, len(expected), len(actual), expected, actual)
		return
	}
	for i, e := range expected {
		a := actual[i]
		if a.Id != e.Id {
			t.Errorf("invalid alternative at %d, expected '%s', got '%s'", i, e.Id, a.Id)
		} else {
			testUtils.ValidateWeights(t, name, e.Criteria, a.Criteria)
		}
	}
}
