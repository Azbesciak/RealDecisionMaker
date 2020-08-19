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

var omission = NewCriteriaOmission([]OmissionResolver{
	&WeakestCriteriaOmissionResolver{},
	&StrongestCriteriaOmissionResolver{},
	&RandomCriteriaOmissionResolver{
		Generator: func(seed int64) utils.ValueGenerator {
			maxVal := float64(len(criteria))
			counter := -1
			return func() float64 {
				counter++
				if seed == 0 {
					return float64(counter) / maxVal
				} else {
					return (maxVal - float64(counter) - 1) / maxVal
				}
			}
		},
	},
	&WeakestByProbabilityCriteriaOmissionResolver{
		Generator: func(seed int64) utils.ValueGenerator {
			maxVal := float64(len(criteria))
			counter := -1
			return func() float64 {
				counter++
				if counter > len(criteria) {
					counter = 0
				}
				if seed == 0 {
					return float64(counter) / maxVal
				} else if seed == 1 {
					return 1
				} else {
					actual := counter - 1
					if actual < 0 {
						actual = len(criteria)
					}
					return float64(actual) / maxVal
				}
			}
		},
	},
})
var notConsidered = []model.AlternativeWithCriteria{
	{Id: "x", Criteria: model.Weights{"1": 1, "2": 2, "3": 3}},
	{Id: "y", Criteria: model.Weights{"1": 0, "2": 1, "3": 4}},
}
var considered = []model.AlternativeWithCriteria{
	{Id: "a", Criteria: model.Weights{"1": 0, "2": 3, "3": 1}},
	{Id: "b", Criteria: model.Weights{"1": 0, "2": 5, "3": 0}},
}
var criteria = testUtils.GenerateCriteria(3)
var listener = model.BiasListener(&testUtils.DummyBiasListener{})
var original = &model.DecisionMakingParams{
	NotConsideredAlternatives: notConsidered,
	ConsideredAlternatives:    considered,
	Criteria:                  criteria,
	MethodParameters: testUtils.DummyMethodParameters{
		Criteria: []string{"1", "2", "3"},
	},
}

func TestCriteriaOmission_ApplyWeakestAsDefault(t *testing.T) {
	m := model.BiasProps(utils.Map{"omittedCriteriaRatio": 0.4})
	result := omission.Apply(original, original, &m, &listener)
	checkOmissionResult(t, result.Props, CriteriaOmissionResult{OmittedCriteria: model.Criteria{criteria[0]}})
}

func TestCriteriaOmission_ApplyWeakest(t *testing.T) {
	m := model.BiasProps(utils.Map{"omittedCriteriaRatio": 0.4, "omissionOrder": "weakest"})
	result := omission.Apply(original, original, &m, &listener)
	checkOmissionResult(t, result.Props, CriteriaOmissionResult{OmittedCriteria: model.Criteria{criteria[0]}})
}

func TestCriteriaOmission_ApplyStrongest(t *testing.T) {
	m := model.BiasProps(utils.Map{"omittedCriteriaRatio": 0.4, "omissionOrder": "strongest"})
	result := omission.Apply(original, original, &m, &listener)
	checkOmissionResult(t, result.Props, CriteriaOmissionResult{OmittedCriteria: model.Criteria{criteria[2]}})
}

func TestCriteriaOmission_ApplyRandom(t *testing.T) {
	m := model.BiasProps(utils.Map{"omittedCriteriaRatio": 0.4, "omissionOrder": "random"})
	result := omission.Apply(original, original, &m, &listener)
	checkOmissionResult(t, result.Props, CriteriaOmissionResult{OmittedCriteria: model.Criteria{criteria[1]}})
}

func TestCriteriaOmission_ApplyWeakestRandom(t *testing.T) {
	m := model.BiasProps(utils.Map{"omittedCriteriaRatio": 0.4, "omissionOrder": "weakestByProbability", "randomSeed": 0})
	result := omission.Apply(original, original, &m, &listener)
	checkOmissionResult(t, result.Props, CriteriaOmissionResult{OmittedCriteria: model.Criteria{criteria[0]}})
}

func TestCriteriaOmission_ApplyWeakestRandomDesc(t *testing.T) {
	m := model.BiasProps(utils.Map{"omittedCriteriaRatio": 0.4, "omissionOrder": "weakestByProbability", "randomSeed": 1})
	result := omission.Apply(original, original, &m, &listener)
	checkOmissionResult(t, result.Props, CriteriaOmissionResult{OmittedCriteria: model.Criteria{criteria[2]}})
}

func TestCriteriaOmission_ApplyWeakestRandomTwo(t *testing.T) {
	m := model.BiasProps(utils.Map{"omittedCriteriaRatio": 0.7, "omissionOrder": "weakestByProbability", "randomSeed": 2})
	result := omission.Apply(original, original, &m, &listener)
	checkOmissionResult(t, result.Props, CriteriaOmissionResult{OmittedCriteria: model.Criteria{criteria[2], criteria[0]}})
}

func TestCriteriaOmission_ApplyRandomDesc(t *testing.T) {
	m := model.BiasProps(utils.Map{"omittedCriteriaRatio": 0.4, "omissionOrder": "random", "randomSeed": 1})
	result := omission.Apply(original, original, &m, &listener)
	checkOmissionResult(t, result.Props, CriteriaOmissionResult{OmittedCriteria: model.Criteria{criteria[2]}})
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
