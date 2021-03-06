package criteria_concealment

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/model/reference-criterion"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

func TestCriteriaOmission_Apply(t *testing.T) {
	omission := CriteriaConcealment{
		generatorSource: testUtils.CyclicRandomGenerator(3, 10),
		referenceCriterionManager: *reference_criterion.NewReferenceCriteriaManager(
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
		"randomSeed": 0,
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
	criterionName := newConcealedCriterionName(&model.Criteria{})
	checkProps(t, result.Props, CriteriaConcealmentResult{
		AddedCriteria: []AddedCriterion{
			{
				Id:                 criterionName,
				Type:               model.Gain,
				AlternativesValues: model.Weights{"a": 0.4, "b": 0.5, "x": 0.6, "y": 0.7},
				MethodParameters:   testUtils.DummyMethodParameters{Criteria: []string{criterionName}},
				ValuesRange:        utils.ValueRange{Min: 0, Max: 1},
			},
		},
	})
}

func TestCriteriaConcealment_Apply_multiple(t *testing.T) {
	omission := CriteriaConcealment{
		generatorSource: testUtils.CyclicRandomGenerator(3, 10),
		referenceCriterionManager: *reference_criterion.NewReferenceCriteriaManager(
			[]reference_criterion.ReferenceCriterionFactory{&reference_criterion.ImportanceRatioReferenceCriterionManager{}},
		),
	}
	var notConsidered []model.AlternativeWithCriteria
	considered := []model.AlternativeWithCriteria{
		{Id: "a", Criteria: model.Weights{"1": 0, "2": 3}},
		{Id: "b", Criteria: model.Weights{"1": 3, "2": 5}},
	}
	criteria := testUtils.GenerateCriteria(2)
	listener := model.BiasListener(&testUtils.DummyBiasListener{})
	m := model.BiasProps(utils.Map{
		"randomSeed": 0,
	})
	original := &model.DecisionMakingParams{
		NotConsideredAlternatives: notConsidered,
		ConsideredAlternatives:    considered,
		Criteria:                  criteria,
		MethodParameters: testUtils.DummyMethodParameters{
			Criteria: []string{"1", "2"},
		},
	}
	result := omission.Apply(original, original, &m, &listener)
	result2 := omission.Apply(original, result.DMP, &m, &listener)
	criterionName := newConcealedCriterionName(&model.Criteria{{Id: baseConcealedCriterionName}})
	checkProps(t, result2.Props, CriteriaConcealmentResult{
		AddedCriteria: []AddedCriterion{{
			Id:                 criterionName,
			Type:               model.Gain,
			AlternativesValues: model.Weights{"a": 1.2, "b": 1.5},
			MethodParameters: testUtils.DummyMethodParameters{Criteria: []string{
				criterionName},
			},
			ValuesRange: utils.ValueRange{Min: 0, Max: 3},
		}},
	})
}

func TestCriteriaConcealment_Apply_strongest_criterion(t *testing.T) {
	omission := CriteriaConcealment{
		generatorSource: testUtils.CyclicRandomGenerator(3, 10),
		referenceCriterionManager: *reference_criterion.NewReferenceCriteriaManager(
			[]reference_criterion.ReferenceCriterionFactory{&reference_criterion.ImportanceRatioReferenceCriterionManager{}},
		),
	}
	var notConsidered []model.AlternativeWithCriteria
	considered := []model.AlternativeWithCriteria{
		{Id: "a", Criteria: model.Weights{"1": 0, "2": 3}},
		{Id: "b", Criteria: model.Weights{"1": 3, "2": 5}},
	}
	criteria := testUtils.GenerateCriteria(2)
	listener := model.BiasListener(&testUtils.DummyBiasListener{})
	m := model.BiasProps(utils.Map{
		"randomSeed":             0,
		"newCriterionImportance": 1,
	})
	original := &model.DecisionMakingParams{
		NotConsideredAlternatives: notConsidered,
		ConsideredAlternatives:    considered,
		Criteria:                  criteria,
		MethodParameters: testUtils.DummyMethodParameters{
			Criteria: []string{"1", "2"},
		},
	}
	result := omission.Apply(original, original, &m, &listener)
	criterionName := newConcealedCriterionName(&model.Criteria{})
	checkProps(t, result.Props, CriteriaConcealmentResult{
		AddedCriteria: []AddedCriterion{{
			Id:                 criterionName,
			Type:               model.Gain,
			AlternativesValues: model.Weights{"a": 3.8, "b": 4},
			MethodParameters: testUtils.DummyMethodParameters{Criteria: []string{
				criterionName},
			},
			ValuesRange: utils.ValueRange{Min: 3, Max: 5},
		}},
	})
}

func TestCriteriaConcealment_Apply_strongest_criterionWithScaling(t *testing.T) {
	omission := CriteriaConcealment{
		generatorSource: testUtils.CyclicRandomGenerator(3, 10),
		referenceCriterionManager: *reference_criterion.NewReferenceCriteriaManager(
			[]reference_criterion.ReferenceCriterionFactory{&reference_criterion.ImportanceRatioReferenceCriterionManager{}},
		),
	}
	var notConsidered []model.AlternativeWithCriteria
	considered := []model.AlternativeWithCriteria{
		{Id: "a", Criteria: model.Weights{"1": 0, "2": 3}},
		{Id: "b", Criteria: model.Weights{"1": 3, "2": 5}},
	}
	criteria := testUtils.GenerateCriteria(2)
	listener := model.BiasListener(&testUtils.DummyBiasListener{})
	m := model.BiasProps(utils.Map{
		"randomSeed":             0,
		"newCriterionImportance": 1,
		"newCriterionScaling":    2,
	})
	original := &model.DecisionMakingParams{
		NotConsideredAlternatives: notConsidered,
		ConsideredAlternatives:    considered,
		Criteria:                  criteria,
		MethodParameters: testUtils.DummyMethodParameters{
			Criteria: []string{"1", "2"},
		},
	}
	result := omission.Apply(original, original, &m, &listener)
	criterionName := newConcealedCriterionName(&model.Criteria{})
	checkProps(t, result.Props, CriteriaConcealmentResult{
		AddedCriteria: []AddedCriterion{{
			Id:                 criterionName,
			Type:               model.Gain,
			AlternativesValues: model.Weights{"a": 3.6, "b": 4},
			MethodParameters: testUtils.DummyMethodParameters{Criteria: []string{
				criterionName},
			},
			ValuesRange: utils.ValueRange{Min: 2, Max: 6},
		}},
	})
}

func TestCriteriaConcealment_Apply_weakest_withBounding(t *testing.T) {
	omission := CriteriaConcealment{
		generatorSource: testUtils.CyclicRandomGenerator(3, 10),
		referenceCriterionManager: *reference_criterion.NewReferenceCriteriaManager(
			[]reference_criterion.ReferenceCriterionFactory{&reference_criterion.ImportanceRatioReferenceCriterionManager{}},
		),
	}
	var notConsidered []model.AlternativeWithCriteria
	considered := []model.AlternativeWithCriteria{
		{Id: "a", Criteria: model.Weights{"1": -20, "2": 3}},
		{Id: "b", Criteria: model.Weights{"1": 30, "2": 5}},
	}
	criteria := testUtils.GenerateCriteria(2)
	listener := model.BiasListener(&testUtils.DummyBiasListener{})
	m := model.BiasProps(utils.Map{
		"randomSeed":                0,
		"newCriterionImportance":    0,
		"newCriterionScaling":       2,
		"disallowNegativeValues":    false,
		"allowedValuesRangeScaling": 0.15,
	})
	original := &model.DecisionMakingParams{
		NotConsideredAlternatives: notConsidered,
		ConsideredAlternatives:    considered,
		Criteria:                  criteria,
		MethodParameters: testUtils.DummyMethodParameters{
			Criteria: []string{"1", "2"},
		},
	}
	result := omission.Apply(original, original, &m, &listener)
	criterionName := newConcealedCriterionName(&model.Criteria{})
	checkProps(t, result.Props, CriteriaConcealmentResult{
		AddedCriteria: []AddedCriterion{{
			Id:                 criterionName,
			Type:               model.Gain,
			AlternativesValues: model.Weights{"a": -2.5, "b": 5},
			MethodParameters: testUtils.DummyMethodParameters{Criteria: []string{
				criterionName},
			},
			ValuesRange: utils.ValueRange{Min: -45, Max: 55},
		}},
	})
}

func checkProps(t *testing.T, actual model.BiasProps, expected CriteriaConcealmentResult) {
	r, ok := actual.(CriteriaConcealmentResult)
	if !ok {
		t.Errorf("expected instance of CriteriaConcealmentResult")
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

func validateMethodAddedCriterion(t *testing.T, r CriteriaConcealmentResult, i int, exp AddedCriterion) {
	act := r.AddedCriteria[i]
	checkMethodParameters(act, t, i, exp)
	if act.Type != exp.Type {
		t.Errorf("wrong added criterion type, expected %s, got %s", exp.Type, act.Type)
	}
	utils.CheckValueRange(t, act.ValuesRange, exp.ValuesRange.Min, exp.ValuesRange.Max)
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
