package choquet

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

var wsTestCriteria = testUtils.GenerateCriteria(3)
var bias = &ChoquetIntegralBiasListener{}

type cw struct {
	_1, _2, _3, _12, _13, _23, _123 model.Weight
}

func (w cw) toMap() *model.Weights {
	return &model.Weights{
		"1":     w._1,
		"2":     w._2,
		"3":     w._3,
		"1,2":   w._12,
		"1,3":   w._13,
		"2,3":   w._23,
		"1,2,3": w._123,
	}
}

type Alts = []model.Weights

func dm(alternativesWeights Alts, combinationsWeights cw) *model.DecisionMakingParams {
	alternatives := testUtils.WrapAlternatives(alternativesWeights)
	return &model.DecisionMakingParams{
		ConsideredAlternatives: alternatives,
		Criteria:               wsTestCriteria,
		MethodParameters:       choquetParams{weights: combinationsWeights.toMap()},
	}
}

func TestChoquetIntegralBiasListener_RankCriteriaAscending(t *testing.T) {
	testUtils.TestBiasRanking(0, t, bias, dm(
		Alts{},
		cw{_1: 0, _2: 0, _3: 0, _12: 0, _13: 0, _23: 0, _123: 0},
	), []string{"1", "2", "3"})
	testUtils.TestBiasRanking(1, t, bias, dm(
		Alts{{"1": 1, "2": 2, "3": 3}},
		cw{_1: 1, _2: 1, _3: 1, _12: 1, _13: 1, _23: 1, _123: 1},
	), []string{"1", "2", "3"})
	testUtils.TestBiasRanking(2, t, bias, dm(
		Alts{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 1, "3": 3}},
		cw{_1: 0.5, _2: 0.5, _3: 0.5, _12: 0.5, _13: 0.2, _23: 0.7, _123: 1},
	), []string{"1", "2", "3"})
	testUtils.TestBiasRanking(3, t, bias, dm(
		Alts{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 1, "3": 0}},
		cw{_1: 0.5, _2: .6, _3: 1, _12: 1, _13: 0.7, _23: .7, _123: 0},
	), []string{"1", "2", "3"})
	testUtils.TestBiasRanking(4, t, bias, dm(
		Alts{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 0, "3": 0}},
		cw{_1: .45, _2: 0, _3: 1, _12: 0, _13: 0, _23: 0, _123: 0},
	), []string{"2", "1", "3"})
	testUtils.TestBiasRanking(5, t, bias, dm(
		Alts{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 0, "3": -2}},
		cw{_1: 1, _2: 1, _3: 0, _12: 1, _13: 0, _23: .5, _123: 1},
	), []string{"3", "2", "1"})
}

func TestChoquetIntegralBias_RankCriteriaAscending_EqualWeights(t *testing.T) {
	testUtils.TestBiasRanking(0, t, bias, dm(
		Alts{},
		cw{_1: 0, _2: 0, _3: 0, _12: 0, _13: 0, _23: 0, _123: 0},
	), []string{"1", "2", "3"})
	testUtils.TestBiasRanking(1, t, bias, dm(
		Alts{{"1": 1, "2": 2, "3": 3}},
		cw{_1: .5, _2: .5, _3: .5, _12: .5, _13: .5, _23: .5, _123: .5},
	), []string{"1", "2", "3"})
	testUtils.TestBiasRanking(2, t, bias, dm(
		Alts{{"1": 10, "2": 6, "3": 3}, {"1": 2, "2": 4, "3": 3}},
		cw{_1: .25, _2: .25, _3: .25, _12: .25, _13: .25, _23: .25, _123: .25},
	), []string{"3", "2", "1"})
	testUtils.TestBiasRanking(3, t, bias, dm(
		Alts{{"1": 1, "2": 7, "3": 2}, {"1": 2, "2": 1, "3": 0}},
		cw{_1: .1, _2: .1, _3: .1, _12: .1, _13: .1, _23: .1, _123: .1},
	), []string{"3", "1", "2"})
	testUtils.TestBiasRanking(4, t, bias, dm(
		Alts{{"1": 1, "2": 2, "3": 3}, {"1": 2.1, "2": 0, "3": 0}},
		cw{_1: .4, _2: .4, _3: .4, _12: .4, _13: .4, _23: .4, _123: .4},
	), []string{"2", "3", "1"})
	testUtils.TestBiasRanking(5, t, bias, dm(
		Alts{{"1": 1, "2": 1000, "3": 3}, {"1": 2, "2": 600, "3": -2}},
		cw{_1: 1, _2: 1, _3: 1, _12: 1, _13: 1, _23: 1, _123: 1},
	), []string{"3", "1", "2"})
}

func TestChoquetIntegralBiasListener_OnCriteriaRemoved_criteriaOrder(t *testing.T) {
	original := *cw{_1: 1, _2: 2, _3: 3, _12: 12, _13: 13, _23: 23, _123: 123}.toMap()
	expected := model.Weights{"2": 2, "3": 3, "2,3": 23}
	left := model.Criteria{wsTestCriteria[2], wsTestCriteria[1]}
	testChoquetRemoval(t, original, expected, left)
}

func TestChoquetIntegralBiasListener_OnCriteriaRemoved(t *testing.T) {
	original := *cw{_1: 1, _2: 2, _3: 3, _12: 12, _13: 13, _23: 23, _123: 123}.toMap()
	expected := model.Weights{"2": 2, "3": 3, "2,3": 23}
	left := wsTestCriteria[1:]
	testChoquetRemoval(t, original, expected, left)
}

func TestChoquetIntegralBiasListener_OnCriteriaRemoved_more(t *testing.T) {
	original := *cw{_1: 1, _2: 2, _3: 3, _12: 12, _13: 13, _23: 23, _123: 123}.toMap()
	expected := model.Weights{"2": 2}
	left := model.Criteria{wsTestCriteria[1]}
	testChoquetRemoval(t, original, expected, left)
}

func testChoquetRemoval(t *testing.T, original, expected model.Weights, left model.Criteria) {
	actual := bias.OnCriteriaRemoved(&left, choquetParams{weights: &original})
	if a, ok := actual.(choquetParams); !ok {
		t.Errorf("Expected instance of choquetParams")
	} else if utils.Differs(expected, *a.weights) {
		t.Errorf("Different weights in choquet, expected %v, got %v", expected, *a.weights)
	}
}
