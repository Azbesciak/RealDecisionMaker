package criteria_mixing

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	reference_criterion "github.com/Azbesciak/RealDecisionMaker/lib/model/reference-criterion"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

const testMixing = 0.25

func mixValues(v1, v2 float64) float64 {
	return v1*testMixing + v2*(1-testMixing)
}

func TestCriteriaMixing_Apply(t *testing.T) {
	mixing := CriteriaMixing{
		generatorSource: testUtils.CyclicRandomGenerator(0, 10),
		referenceCriteriaManager: *reference_criterion.NewReferenceCriteriaManager(
			[]reference_criterion.ReferenceCriterionFactory{&reference_criterion.ImportanceRatioReferenceCriterionManager{}},
		),
	}
	notConsidered := []model.AlternativeWithCriteria{
		{Id: "x", Criteria: model.Weights{"1": 1, "2": 2, "3": 3}},
		{Id: "y", Criteria: model.Weights{"1": 0, "2": 1, "3": 4}},
	}
	considered := []model.AlternativeWithCriteria{
		{Id: "a", Criteria: model.Weights{"1": 0, "2": 3, "3": 1}},
		{Id: "b", Criteria: model.Weights{"1": 0, "2": 5, "3": 0}},
	}
	criteria := testUtils.GenerateCriteria(3)
	listener := model.BiasListener(&testUtils.DummyBiasListener{})
	m := model.BiasProps(utils.Map{
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
			Id:           "1",
			Type:         model.Gain,
			ScaledValues: model.Weights{"x": 1, "y": 0, "a": 0, "b": 0},
		},
		Component2: CriterionComponent{
			Id:           "2",
			Type:         model.Gain,
			ScaledValues: model.Weights{"x": 0.25, "y": 0, "a": 0.5, "b": 1},
		},
		NewCriterion: CriterionComponent{
			Id:           "__1+2__",
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
	if expected.Id != actual.Id {
		t.Errorf("%s: expected name '%s', got '%s'", name, expected.Id, actual.Id)
	}
	if expected.Type != actual.Type {
		t.Errorf("%s: expected type '%s', got '%s'", name, expected.Type, actual.Type)
	}
	testUtils.ValidateWeights(t, name, expected.ScaledValues, actual.ScaledValues)
}
