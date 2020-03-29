package criteria_omission

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"sort"
	"testing"
)

func TestCriteriaOmission_splitCriteria(t *testing.T) {
	criteria := &model.Criteria{
		model.Criterion{Id: "1"},
		model.Criterion{Id: "2"},
		model.Criterion{Id: "3"},
		model.Criterion{Id: "4"},
		model.Criterion{Id: "5"},
		model.Criterion{Id: "6"},
	}
	validateOmission(t, criteria, 0, []string{}, []string{"1", "2", "3", "4", "5", "6"})
	validateOmission(t, criteria, 1, []string{"1", "2", "3", "4", "5", "6"}, []string{})
	validateOmission(t, criteria, 0.5, []string{"1", "2", "3"}, []string{"4", "5", "6"})
	validateOmission(t, criteria, 0.25, []string{"1"}, []string{"2", "3", "4", "5", "6"})
	validateOmission(t, criteria, 0.34, []string{"1", "2"}, []string{"3", "4", "5", "6"})
}

func validateOmission(t *testing.T, criteria *model.Criteria, ratio float64, omitted []string, kept []string) {
	division := splitCriteriaToOmit(ratio, criteria)
	actualOmittedLen := len(*division.omitted)
	actualKeptLen := len(*division.kept)

	if actualOmittedLen+actualKeptLen != len(*criteria) {
		t.Errorf("sum of kept (%d) and omitted (%d) criteria is not equal to total len (%d)", actualKeptLen, actualOmittedLen, len(*criteria))
	}
	checkCount(t, "omit", omitted, division.omitted)
	checkCount(t, "keep", kept, division.kept)
}

func checkCount(t *testing.T, typ string, expected []string, actual *model.Criteria) {
	expectedLen := len(expected)
	actualLen := len(*actual)
	if actualLen != expectedLen {
		t.Errorf("expected %d criteria to %s, but got %d", expectedLen, typ, actualLen)
		return
	}
	for i, expectedId := range expected {
		actualId := (*actual).Get(i).Identifier()
		if actualId != expectedId {
			t.Errorf("expected '%s' at %d in %s criteria, got '%s'", expectedId, i, typ, actualId)
		}
	}
}

func TestCriteriaOmission_Apply(t *testing.T) {
	omission := CriteriaOmission{
		newCriterionValueScalar: 1,
		generatorSource: func(seed int64) utils.ValueGenerator {
			counter := 3 // because 0.5 chance for added criterion
			return func() float64 {
				counter = (counter + 1) % 10
				return float64(counter) / 10
			}
		},
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
	listener := model.HeuristicListener(&dummyHeuListener{})
	m := model.HeuristicProps(map[string]interface{}{
		"addCriterionProbability": 0.5,
		"omittedCriteriaRatio":    0.4,
		"randomSeed":              0,
	})
	result := omission.Apply(&model.DecisionMakingParams{
		NotConsideredAlternatives: notConsidered,
		ConsideredAlternatives:    considered,
		Criteria:                  criteria,
		MethodParameters: dummyMethodParameters{
			criteria: []string{"1", "2", "3"},
		},
	}, &m, &listener)
	checkProps(t, result.Props, CriteriaOmissionResult{
		OmittedCriteria: model.Criteria{criteria[0]},
		AddedCriteria: []AddedCriterion{
			{
				Type:                model.Gain,
				AlternativesValues:  model.Weights{"a": 0.5, "b": 0.6, "x": 0.7, "y": 0.8},
				MethodParameters:    dummyMethodParameters{criteria: []string{addedCriterionName()}},
				CriterionValueRange: utils.ValueRange{Min: 0, Max: 1},
			},
		},
	})
}

func checkProps(t *testing.T, actual model.HeuristicProps, expected CriteriaOmissionResult) {
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
	actPar, ok := act.MethodParameters.(dummyMethodParameters)
	if !ok {
		t.Errorf("expected instance of dummyMethodParameters at criterion %d, got %v", i, act.MethodParameters)
		return
	}
	expPar := exp.MethodParameters.(dummyMethodParameters)
	if len(actPar.criteria) != len(expPar.criteria) {
		t.Errorf("expected %d criteria, got %d", len(expPar.criteria), len(actPar.criteria))
		return
	}
	for i, expCr := range expPar.criteria {
		actCr := actPar.criteria[i]
		if actCr != expCr {
			t.Errorf("criterion at %d are not equal, expected %s, got %s", i, expCr, actCr)
		}
	}
}

type dummyMethodParameters struct {
	criteria []string
}

type dummyHeuListener struct {
}

func (d *dummyHeuListener) Merge(params model.MethodParameters, addition model.MethodParameters) model.MethodParameters {
	prevParams := params.(dummyMethodParameters)
	addedParams := addition.(dummyMethodParameters)
	return dummyMethodParameters{criteria: append(prevParams.criteria, addedParams.criteria...)}
}

func (d *dummyHeuListener) Identifier() string {
	panic("should not call identifier in test")
}

func (d *dummyHeuListener) OnCriterionAdded(
	criterion *model.Criterion,
	previousRankedCriteria *model.Criteria,
	params model.MethodParameters,
	generator utils.ValueGenerator,
) model.AddedCriterionParams {
	return dummyMethodParameters{criteria: []string{criterion.Id}}
}

func (d *dummyHeuListener) OnCriteriaRemoved(removedCriteria *model.Criteria, leftCriteria *model.Criteria, params model.MethodParameters) model.MethodParameters {
	return dummyMethodParameters{criteria: *leftCriteria.Names()}
}

func (d *dummyHeuListener) RankCriteriaAscending(params *model.DecisionMakingParams) *model.Criteria {
	criteria := params.Criteria.ShallowCopy()
	sort.Slice(*criteria, func(i, j int) bool {
		return (*criteria)[i].Id < (*criteria)[j].Id
	})
	return criteria
}
