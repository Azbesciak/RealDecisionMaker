package criteria_omission

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

func TestCriteriaOmission_splitCriteria(t *testing.T) {
	criteria := &model.Criteria{
		{Id: "1"},
		{Id: "2"},
		{Id: "3"},
		{Id: "4"},
		{Id: "5"},
		{Id: "6"},
	}
	validateOmission(t, criteria, 0, []string{}, []string{"1", "2", "3", "4", "5", "6"})
	validateOmission(t, criteria, 1, []string{"1", "2", "3", "4", "5", "6"}, []string{})
	validateOmission(t, criteria, 0.5, []string{"1", "2", "3"}, []string{"4", "5", "6"})
	validateOmission(t, criteria, 0.25, []string{"1"}, []string{"2", "3", "4", "5", "6"})
	validateOmission(t, criteria, 0.34, []string{"1", "2"}, []string{"3", "4", "5", "6"})
}

func TestCriteriaOmission_Apply(t *testing.T) {
	omission := CriteriaOmission{
		newCriterionValueScalar: 1,
		generatorSource:         testUtils.CyclicRandomGenerator(3, 10),
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
	m := model.BiasProps(map[string]interface{}{
		"addCriterionProbability": 0.5,
		"omittedCriteriaRatio":    0.4,
		"randomSeed":              0,
	})
	original := &model.DecisionMakingParams{
		NotConsideredAlternatives: notConsidered,
		ConsideredAlternatives:    considered,
		Criteria:                  criteria,
		MethodParameters: testUtils.DummyMethodParameters{
			Criteria: []string{"1", "2", "3"},
		},
	}
	result := omission.Apply(original, original, &m, &listener)
	checkProps(t, result.Props, CriteriaOmissionResult{
		OmittedCriteria: model.Criteria{criteria[0]},
		AddedCriteria: []AddedCriterion{
			{
				Type:                model.Gain,
				AlternativesValues:  model.Weights{"a": 0.5, "b": 0.6, "x": 0.7, "y": 0.8},
				MethodParameters:    testUtils.DummyMethodParameters{Criteria: []string{addedCriterionName()}},
				CriterionValueRange: utils.ValueRange{Min: 0, Max: 1},
			},
		},
	})
}

func validateOmission(t *testing.T, criteria *model.Criteria, ratio float64, omitted []string, kept []string) {
	division := splitCriteriaToOmit(ratio, criteria)
	actualOmittedLen := len(*division.omitted)
	actualKeptLen := len(*division.kept)

	if actualOmittedLen+actualKeptLen != len(*criteria) {
		t.Errorf("sum of kept (%d) and omitted (%d) criteria is not equal to total len (%d)", actualKeptLen, actualOmittedLen, len(*criteria))
	}
	testUtils.CheckCount(t, "omit", omitted, division.omitted)
	testUtils.CheckCount(t, "keep", kept, division.kept)
}

func checkProps(t *testing.T, actual model.BiasProps, expected CriteriaOmissionResult) {
	r, ok := actual.(CriteriaOmissionResult)
	if !ok {
		t.Errorf("expected instance of CriteriaOmissionResult")
		return
	}
	if len(r.AddedCriteria) != len(expected.AddedCriteria) {
		t.Errorf("expected %d added criteria, got %d", len(expected.AddedCriteria), len(r.AddedCriteria))
		return
	}
	for i, exp := range expected.AddedCriteria {
		validateMethodAddedCriterion(t, r, i, exp)
	}
}

func validateMethodAddedCriterion(t *testing.T, r CriteriaOmissionResult, i int, exp AddedCriterion) {
	act := r.AddedCriteria[i]
	checkMethodParameters(act, t, i, exp)
	if act.Type != exp.Type {
		t.Errorf("wrong added criterion type, expected %s, got %s", exp.Type, act.Type)
	}
	utils.CheckValueRange(t, act.CriterionValueRange, 0, 1)
	checkAlternatives(exp, act, t)
}

func checkAlternatives(exp AddedCriterion, act AddedCriterion, t *testing.T) {
	for ek, ev := range exp.AlternativesValues {
		av, ok := act.AlternativesValues[ek]
		if !ok {
			t.Errorf("alternative '%s' not found in values %v, expected %v", ek, act.AlternativesValues, exp.AlternativesValues)
		} else if !utils.FloatsAreEqual(av, ev, 1e-6) {
			t.Errorf("expected %f for alternative '%s', got %f", ev, ek, av)
		}
	}
}

func checkMethodParameters(act AddedCriterion, t *testing.T, i int, exp AddedCriterion) {
	actPar, ok := act.MethodParameters.(testUtils.DummyMethodParameters)
	if !ok {
		t.Errorf("expected instance of testUtils.DummyMethodParameters at criterion %d, got %v", i, act.MethodParameters)
		return
	}
	expPar := exp.MethodParameters.(testUtils.DummyMethodParameters)
	if len(actPar.Criteria) != len(expPar.Criteria) {
		t.Errorf("expected %d criteria, got %d", len(expPar.Criteria), len(actPar.Criteria))
		return
	}
	for i, expCr := range expPar.Criteria {
		actCr := actPar.Criteria[i]
		if actCr != expCr {
			t.Errorf("criterion at %d are not equal, expected %s, got %s", i, expCr, actCr)
		}
	}
}
