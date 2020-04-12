package criteria_mixing

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"testing"
)

const testMixing = 0.25

func mixValues(v1, v2 float64) float64 {
	return v1*testMixing + v2*(1-testMixing)
}

func TestCriteriaMixing_Apply(t *testing.T) {
	mixing := CriteriaMixing{
		generatorSource: testUtils.CyclicRandomGenerator(0, 10),
	}
	notConsidered := []model.AlternativeWithCriteria{
		{Id: "x", Criteria: model.Weights{"1": 1, "2": 2, "3": 3}},
		{Id: "y", Criteria: model.Weights{"1": 0, "2": 1, "3": 4}},
	}
	considered := []model.AlternativeWithCriteria{
		{Id: "a", Criteria: model.Weights{"1": 0, "2": 3, "3": 1}},
		{Id: "b", Criteria: model.Weights{"1": 0, "2": 5, "3": 0}},
	}
	criteria := model.Criteria{
		{Id: "1", Type: model.Gain},
		{Id: "2", Type: model.Gain},
		{Id: "3", Type: model.Gain},
	}
	listener := model.BiasListener(&testUtils.DummyBiasListener{})
	m := model.BiasProps(map[string]interface{}{
		"mixingRatio": testMixing,
		"randomSeed":  0,
	})
	original := &model.DecisionMakingParams{
		NotConsideredAlternatives: notConsidered,
		ConsideredAlternatives:    considered,
		Criteria:                  criteria,
		MethodParameters: testUtils.DummyMethodParameters{
			Criteria: []string{"1", "2", "3"},
		},
	}
	result := mixing.Apply(original, original, &m, &listener)
	expected := MixedCriterion{
		Component1: CriterionComponent{
			Criterion:    "1",
			Type:         model.Gain,
			ScaledValues: model.Weights{"x": 1, "y": 0, "a": 0, "b": 0},
		},
		Component2: CriterionComponent{
			Criterion:    "2",
			Type:         model.Gain,
			ScaledValues: model.Weights{"x": 0.25, "y": 0, "a": 0.5, "b": 1},
		},
		NewCriterion: CriterionComponent{
			Criterion:    "__1+2__",
			Type:         model.Gain,
			ScaledValues: model.Weights{"x": mixValues(1, 0.25), "y": 0, "a": mixValues(0, 0.5), "b": mixValues(0, 1)},
		},
		Params: []string{"1", "2", "3", "__1+2__"},
	}
	checkProps(t, expected, result.Props)
}

func checkProps(t *testing.T, expected MixedCriterion, actual model.BiasProps) {
	v, ok := actual.(MixedCriterion)
	if !ok {
		t.Errorf("expected instance of MixedCriterion")
		return
	}
	checkCriterion(t, "cmp1", expected.Component1, v.Component1)
	checkCriterion(t, "cmp2", expected.Component2, v.Component2)
	checkCriterion(t, "mix", expected.NewCriterion, v.NewCriterion)
}

func checkCriterion(t *testing.T, name string, expected, actual CriterionComponent) {
	if expected.Criterion != actual.Criterion {
		t.Errorf("%s: expected name '%s', got '%s'", name, expected.Criterion, actual.Criterion)
	}
	if expected.Type != actual.Type {
		t.Errorf("%s: expected type '%s', got '%s'", name, expected.Type, actual.Type)
	}
	testUtils.ValidateWeights(t, name, expected.ScaledValues, actual.ScaledValues)
}
