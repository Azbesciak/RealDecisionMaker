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
	omission := CriteriaOmission{}
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
	m := model.BiasProps(utils.Map{"omittedCriteriaRatio": 0.4})
	original := &model.DecisionMakingParams{
		NotConsideredAlternatives: notConsidered,
		ConsideredAlternatives:    considered,
		Criteria:                  criteria,
		MethodParameters: testUtils.DummyMethodParameters{
			Criteria: []string{"1", "2", "3"},
		},
	}
	result := omission.Apply(original, original, &m, &listener)
	checkOmissionResult(t, result.Props, CriteriaOmissionResult{OmittedCriteria: model.Criteria{criteria[0]}})
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

func checkOmissionResult(t *testing.T, actual model.BiasProps, expected CriteriaOmissionResult) {
	r, ok := actual.(CriteriaOmissionResult)
	if !ok {
		t.Errorf("expected instance of CriteriaOmissionResult")
		return
	}
	if len(r.OmittedCriteria) != len(expected.OmittedCriteria) {
		t.Errorf("expected %d ommited criteria, got %d", len(expected.OmittedCriteria), len(r.OmittedCriteria))
		return
	}
	for i, exp := range expected.OmittedCriteria {
		act := r.OmittedCriteria[i]
		if act.Id != exp.Id {
			t.Errorf("expected '%s' criterion ommited, got '%s' at index '%d'", exp.Id, act.Id, i)
		}
	}
}
