package satisfaction

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/limited-rationality/satisfaction-levels"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

var _satisfaction = Satisfaction{
	generator: func(seed int64) utils.ValueGenerator {
		return func() float64 {
			return 1
		}
	},
	functions: []satisfaction_levels.SatisfactionLevelsSource{
		&satisfaction_levels.IdealSubtrCoefficientSatisfaction,
		&satisfaction_levels.IdealDecreasingMulCoefficientSatisfaction,
	},
}

var criteria = testUtils.GenerateCriteria(3)

var a1 = model.AlternativeWithCriteria{
	Id:       "a",
	Criteria: model.Weights{"1": 3, "2": 4, "3": 5},
}

var a2 = model.AlternativeWithCriteria{
	Id:       "b",
	Criteria: model.Weights{"1": 6, "2": 5, "3": 4},
}

var a3 = model.AlternativeWithCriteria{
	Id:       "c",
	Criteria: model.Weights{"1": 3, "2": 6, "3": 3},
}

var a4 = model.AlternativeWithCriteria{
	Id:       "d",
	Criteria: model.Weights{"1": 10, "2": 10, "3": 10}, //not considered but as a reference for ideal
}

var cur = model.AlternativeWithCriteria{
	Id:       "cur",
	Criteria: model.Weights{"1": 0, "2": 0, "3": 0}, //not considered but as a reference for ideal
}

func TestSatisfaction_Evaluate(t *testing.T) {
	dm := model.DecisionMakingParams{
		ConsideredAlternatives:    []model.AlternativeWithCriteria{a1, a2, a3},
		NotConsideredAlternatives: []model.AlternativeWithCriteria{a4, cur},
		Criteria:                  criteria,
		MethodParameters: SatisfactionParameters{
			Function: satisfaction_levels.IdealSubtrCoefficientSatisfaction.Identifier(),
			Params: map[string]interface{}{
				"Coefficient": 0.1,
				"MinValue":    0.1,
				"MaxValue":    1.0,
			},
			RandomSeed:    0,
			CurrentChoice: "cur",
		},
	}
	actual := _satisfaction.Evaluate(&dm)
	expected := model.AlternativesRanking{
		{
			AlternativeResult: model.AlternativeResult{
				Alternative: a2,
				Evaluation: SatisfactionEvaluation{
					SatisfiedThresholds: model.Weights{
						"1": 3, "2": 3, "3": 3,
					},
					ThresholdsIndex: 7,
				},
			},
			BetterThanOrSameAs: []string{a1.Id, a3.Id, cur.Id},
		},
		{
			AlternativeResult: model.AlternativeResult{
				Alternative: a1,
				Evaluation: SatisfactionEvaluation{
					SatisfiedThresholds: model.Weights{
						"1": 2, "2": 2, "3": 2,
					},
					ThresholdsIndex: 8,
				},
			},
			BetterThanOrSameAs: []string{a3.Id, cur.Id},
		},
		{
			AlternativeResult: model.AlternativeResult{
				Alternative: a3,
				Evaluation: SatisfactionEvaluation{
					SatisfiedThresholds: model.Weights{
						"1": 2, "2": 2, "3": 2,
					},
					ThresholdsIndex: 8,
				},
			},
			BetterThanOrSameAs: []string{cur.Id},
		},
		{
			AlternativeResult: model.AlternativeResult{
				Alternative: cur,
				Evaluation: SatisfactionEvaluation{
					SatisfiedThresholds: model.Weights{
						"1": 0, "2": 0, "3": 0,
					},
					ThresholdsIndex: 10,
				},
			},
			BetterThanOrSameAs: []string{},
		},
	}
	testUtils.CompareRankings(&expected, actual, t)
}

func TestSatisfaction_ParseParams(t *testing.T) {
	expected := SatisfactionParameters{
		Function: satisfaction_levels.IdealDecreasingMulCoefficientSatisfaction.Identifier(),
		Params: utils.Map{
			"minValue":    0.01,
			"maxValue":    1,
			"coefficient": 0.2,
		},
		RandomSeed:    159,
		CurrentChoice: "a",
	}
	actual := _satisfaction.ParseParams(&model.DecisionMaker{
		PreferenceFunction: methodName,
		KnownAlternatives: []model.AlternativeWithCriteria{{
			Id:       "a",
			Criteria: model.Weights{"1": 1, "2": 2},
		}, {
			Id:       "b",
			Criteria: model.Weights{"1": 2, "2": 1},
		}},
		ChoseToMake: model.Alternatives{"a", "b"},
		Criteria:    testUtils.GenerateCriteria(2),
		MethodParameters: utils.Map{
			"function": "idealMultipliedCoefficient",
			"params": utils.Map{
				"minValue":    0.01,
				"maxValue":    1,
				"coefficient": 0.2,
			},
			"randomSeed":    159,
			"currentChoice": "a",
		},
	})

	if _, ok := actual.(SatisfactionParameters); !ok {
		t.Errorf("expected Satisfaction parameters")
	} else if utils.Differs(expected, actual) {
		t.Errorf("different than expected: %v vs %v", expected, actual)
	}
}

func TestAlternativesIterationWithThresholds(t *testing.T) {
	satisfaction := Satisfaction{
		generator: func(seed int64) utils.ValueGenerator {
			return func() float64 {
				return 1
			}
		},
		functions: []satisfaction_levels.SatisfactionLevelsSource{
			&satisfaction_levels.DecreasingThresholds,
		},
	}
	considered := []model.AlternativeWithCriteria{
		{Id: "Avengers", Criteria: model.Weights{"popularity": 7, "length": 14, "averageMark": 7.7}},
		{Id: "Dunkirk", Criteria: model.Weights{"popularity": 3, "length": 11, "averageMark": 7.4}},
		{Id: "Knives Out", Criteria: model.Weights{"popularity": 6, "length": 11, "averageMark": 7.8}},
	}

	notConsidered := []model.AlternativeWithCriteria{{
		Id: "Best",
		Criteria: model.Weights{
			"popularity":  10,
			"length":      1,
			"averageMark": 10,
		},
	}, {
		Id: "Worst",
		Criteria: model.Weights{
			"popularity":  1,
			"length":      19,
			"averageMark": 2,
		},
	}}

	criteria := model.Criteria{
		{Id: "popularity", Type: "gain"},
		{Id: "length", Type: "cost"},
		{Id: "averageMark", Type: "gain"},
	}

	thresholds := []model.Weights{
		{"popularity": 10, "length": 6, "averageMark": 8},
		{"popularity": 8, "length": 10, "averageMark": 7.7},
		{"popularity": 6, "length": 11, "averageMark": 7.5},
		{"popularity": 6, "length": 20, "averageMark": 7},
		{"popularity": 3, "length": 20, "averageMark": 7},
	}
	actual := satisfaction.Evaluate(&model.DecisionMakingParams{
		ConsideredAlternatives:    considered,
		NotConsideredAlternatives: notConsidered,
		Criteria:                  criteria,
		MethodParameters: SatisfactionParameters{
			Function: satisfaction_levels.DecreasingThresholds.Identifier(),
			Params: utils.Map{
				"thresholds": thresholds,
			},
			RandomSeed: 1,
		},
	})
	expected := model.AlternativesRanking{
		{
			AlternativeResult: model.AlternativeResult{
				Alternative: considered[2],
				Evaluation: SatisfactionEvaluation{
					SatisfiedThresholds: thresholds[2],
					ThresholdsIndex:     2,
				},
			},
			BetterThanOrSameAs: []string{considered[0].Id, considered[1].Id},
		},
		{
			AlternativeResult: model.AlternativeResult{
				Alternative: considered[0],
				Evaluation: SatisfactionEvaluation{
					SatisfiedThresholds: thresholds[3],
					ThresholdsIndex:     3,
				},
			},
			BetterThanOrSameAs: []string{considered[1].Id},
		},
		{
			AlternativeResult: model.AlternativeResult{
				Alternative: considered[1],
				Evaluation: SatisfactionEvaluation{
					SatisfiedThresholds: thresholds[4],
					ThresholdsIndex:     4,
				},
			},
			BetterThanOrSameAs: []string{},
		},
	}
	testUtils.CompareRankings(&expected, actual, t)
}
