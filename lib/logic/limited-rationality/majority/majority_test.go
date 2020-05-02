package majority

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

var _majority = Majority{
	generator: func(seed int64) utils.ValueGenerator {
		return func() float64 {
			return 1
		}
	},
}

var criteria = testUtils.GenerateCriteria(3)
var criteriaValues = []model.WeightedCriterion{{
	Criterion: criteria[0],
	Weight:    1,
}, {
	Criterion: criteria[1],
	Weight:    2,
}, {
	Criterion: criteria[2],
	Weight:    1,
}}
var a1 = model.AlternativeWithCriteria{
	Id:       "a",
	Criteria: model.Weights{"1": 2, "2": 4, "3": 3},
}

var a2 = model.AlternativeWithCriteria{
	Id:       "b",
	Criteria: model.Weights{"1": 0, "2": 5, "3": 1},
}

var a3 = model.AlternativeWithCriteria{
	Id:       "c",
	Criteria: model.Weights{"1": 0, "2": 4, "3": 1},
}

var a4 = model.AlternativeWithCriteria{
	Id:       "d",
	Criteria: model.Weights{"1": 4, "2": 6, "3": 2},
}

var cur = model.AlternativeWithCriteria{
	Id:       "cur",
	Criteria: model.Weights{"1": 0, "2": 3, "3": 1},
}

func TestMajority_Evaluate(t *testing.T) {
	dm := model.DecisionMakingParams{
		ConsideredAlternatives:    []model.AlternativeWithCriteria{a1, a2, a3},
		NotConsideredAlternatives: []model.AlternativeWithCriteria{a4, cur},
		Criteria:                  criteria,
		MethodParameters: MajorityHeuristicParams{
			Weights: map[string]float64{
				"1": 1, "2": 2, "3": 1,
			},
			Current: "cur",
			Seed:    0,
		},
	}
	actual := _majority.Evaluate(&dm)
	expected := model.AlternativesRanking{
		{
			AlternativeResult: model.AlternativeResult{
				Alternative: a1,
				Evaluation: MajorityEvaluation{
					Value:                    2,
					ComparedWith:             "",
					ComparedAlternativeValue: 0,
				},
			},
			BetterThanOrSameAs: []string{cur.Id, a3.Id, a2.Id},
		},
		{
			AlternativeResult: model.AlternativeResult{
				Alternative: a2,
				Evaluation: MajorityEvaluation{
					Value:                    2,
					ComparedWith:             "a",
					ComparedAlternativeValue: 2,
				},
			},
			BetterThanOrSameAs: []string{cur.Id, a3.Id, a1.Id},
		},
		{
			AlternativeResult: model.AlternativeResult{
				Alternative: a3,
				Evaluation: MajorityEvaluation{
					Value:                    0,
					ComparedWith:             "a",
					ComparedAlternativeValue: 2,
				},
			},
			BetterThanOrSameAs: []string{cur.Id},
		},
		{
			AlternativeResult: model.AlternativeResult{
				Alternative: cur,
				Evaluation: MajorityEvaluation{
					Value:                    0,
					ComparedWith:             "a",
					ComparedAlternativeValue: 4,
				},
			},
			BetterThanOrSameAs: []string{},
		},
	}
	testUtils.CompareRankings(&expected, actual, t)
}

func TestMajority_compare(t *testing.T) {
	s1, s2 := compare(&criteriaValues, &a1, &a2)
	if !utils.FloatsAreEqual(s1, s2, eps) {
		t.Errorf("expected same values, got %f and %f", s1, s2)
	} else if !utils.FloatsAreEqual(2, s1, eps) {
		t.Errorf("expected same values, got %f, expected 2", s1)
	}
}

func TestMajority_ParseParams(t *testing.T) {
	expected := MajorityHeuristicParams{
		Weights: model.Weights{
			"1": 111, "2": 2, "3": 0,
		},
		Current: "1234",
		Seed:    951,
	}
	actual := _majority.ParseParams(&model.DecisionMaker{
		PreferenceFunction: methodName,
		KnownAlternatives: []model.AlternativeWithCriteria{{
			Id:       "1111",
			Criteria: model.Weights{"1": 3, "2": 7, "3": 10},
		}, {
			Id:       "1234",
			Criteria: model.Weights{"1": 101, "2": 18, "3": 6},
		}},
		ChoseToMake: model.Alternatives{"1111", "1234"},
		Criteria:    testUtils.GenerateCriteria(3),
		MethodParameters: utils.Map{
			"weights": utils.Map{"1": 111, "2": 2, "3": 0},
			"current": "1234",
			"seed":    951,
		},
	})
	if _, ok := actual.(MajorityHeuristicParams); !ok {
		t.Errorf("expected MajorityHeuristicParams")
	} else if testUtils.Differs(expected, actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
