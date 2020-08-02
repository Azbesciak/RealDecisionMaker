package aspect_elimination

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/limited-rationality/satisfaction-levels"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

var testAspElim = AspectEliminationHeuristic{
	functions: []satisfaction_levels.SatisfactionLevelsSource{
		&satisfaction_levels.IdealAdditiveCoefficientSatisfaction,
		&satisfaction_levels.IdealIncreasingMulCoefficientSatisfaction,
		&satisfaction_levels.ThresholdSatisfactionLevelsSource{},
	},
	generator: testUtils.CyclicRandomGenerator(0, 1),
}

func TestAspectEliminationHeuristic_Evaluate(t *testing.T) {
	alternatives := []model.AlternativeWithCriteria{{
		Id:       "b",
		Criteria: model.Weights{"1": 2, "2": 3, "3": 4},
	}, {
		Id:       "c",
		Criteria: model.Weights{"1": 40, "2": 1, "3": 60},
	}, {
		Id:       "d",
		Criteria: model.Weights{"1": 20, "2": 2, "3": 3},
	}}
	params := model.DecisionMakingParams{
		NotConsideredAlternatives: []model.AlternativeWithCriteria{{
			Id:       "a",
			Criteria: model.Weights{"1": 1, "2": 2, "3": 3},
		}},
		ConsideredAlternatives: alternatives,
		Criteria:               testUtils.GenerateCriteria(3),
		MethodParameters: AspectEliminationHeuristicParams{
			Function: "thresholds",
			Params: map[string]interface{}{
				"thresholds": []map[string]float64{{
					"1": 2, "2": 1.5, "3": 3.5,
				}},
			},
			RandomSeed: 0,
			Weights: model.Weights{
				"2": 3, "3": 2, "1": 1,
			},
		},
	}
	actual := testAspElim.Evaluate(&params)
	expected := model.AlternativesRanking{{
		AlternativeResult: model.AlternativeResult{
			Alternative: alternatives[0],
			Evaluation: AspectEliminationEvaluation{
				NotSatisfiedThreshold: model.Weights{},
				ThresholdsIndex:       1,
			},
		},
		BetterThanOrSameAs: []string{alternatives[2].Id}, // order matters
	}, {
		AlternativeResult: model.AlternativeResult{
			Alternative: alternatives[2],
			Evaluation: AspectEliminationEvaluation{
				NotSatisfiedThreshold: model.Weights{"3": 3.5},
				ThresholdsIndex:       0,
			},
		},
		BetterThanOrSameAs: []string{alternatives[1].Id},
	}, {
		AlternativeResult: model.AlternativeResult{
			Alternative: alternatives[1],
			Evaluation: AspectEliminationEvaluation{
				NotSatisfiedThreshold: model.Weights{"2": 1.5},
				ThresholdsIndex:       0,
			},
		},
		BetterThanOrSameAs: []string{},
	}}
	testUtils.CompareRankings(&expected, actual, t)
}

func TestAspectEliminationHeuristic_ParseParams(t *testing.T) {
	criteria := testUtils.GenerateCriteria(2)
	dm := model.DecisionMaker{
		KnownAlternatives: nil,
		ChoseToMake:       nil,
		Criteria:          criteria,
		MethodParameters: map[string]interface{}{
			"function": "thresholds",
			"params": map[string]interface{}{
				"thresholds": []map[string]interface{}{{"1": 1, "2": 3}, {"1": 3, "2": 4}},
			},
			"randomSeed": 123,
			"weights":    map[string]interface{}{"1": 2, "2": 2.5},
		},
	}
	actual := testAspElim.ParseParams(&dm)
	expected := AspectEliminationHeuristicParams{
		Function: "thresholds",
		Params: map[string]interface{}{
			"thresholds": []map[string]interface{}{{"1": 1, "2": 3}, {"1": 3, "2": 4}},
		},
		RandomSeed: 123,
		Weights:    model.Weights{"1": 2, "2": 2.5},
	}
	compareParams(t, expected, actual)
}

func compareParams(t *testing.T, expected AspectEliminationHeuristicParams, actual interface{}) {
	parsed, ok := actual.(AspectEliminationHeuristicParams)
	if !ok {
		t.Errorf("expected AspectEliminationHeuristicParams as params")
		return
	}
	if utils.Differs(expected.Weights, parsed.Weights) {
		t.Errorf("different weights, expected %v, got %v", expected.Weights, parsed.Weights)
	}
	if utils.Differs(expected.Params, parsed.Params) {
		t.Errorf("different params, expected %v, got %v", expected.Params, parsed.Params)
	}
	if expected.RandomSeed != parsed.RandomSeed {
		t.Errorf("different seed, expected %v, got %v", expected.RandomSeed, parsed.RandomSeed)
	}
}
