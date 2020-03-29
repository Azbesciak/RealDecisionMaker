package weighted_sum

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"testing"
)

func TestWeightedSumHeuristic_RankCriteriaAscending_EqualWeights(t *testing.T) {
	testUtils.TestHeuristicRanking(t, heu, dmParams([]model.Weights{}, []model.Weight{}), []string{"1", "2", "3"})
	testUtils.TestHeuristicRanking(t, heu, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}}, []model.Weight{1, 1, 1}), []string{"1", "2", "3"})
	testUtils.TestHeuristicRanking(t, heu, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 1, "3": 3}}, []model.Weight{1, 1, 1}), []string{"1", "2", "3"})
	testUtils.TestHeuristicRanking(t, heu, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 1, "3": 0}}, []model.Weight{1, 1, 1}), []string{"1", "2", "3"})
	testUtils.TestHeuristicRanking(t, heu, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 0, "3": 0}}, []model.Weight{1, 1, 1}), []string{"2", "1", "3"})
	testUtils.TestHeuristicRanking(t, heu, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 0, "3": -2}}, []model.Weight{1, 1, 1}), []string{"3", "2", "1"})
}

func TestWeightedSumHeuristic_RankCriteriaAscending(t *testing.T) {
	testUtils.TestHeuristicRanking(t, heu, dmParams([]model.Weights{}, []model.Weight{}), []string{"1", "2", "3"})
	testUtils.TestHeuristicRanking(t, heu, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}}, []model.Weight{1, 1, 2}), []string{"1", "2", "3"})
	testUtils.TestHeuristicRanking(t, heu, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 1, "3": 3}}, []model.Weight{10, 0.5, 0.1}), []string{"3", "2", "1"})
	testUtils.TestHeuristicRanking(t, heu, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 1, "3": 0}}, []model.Weight{1, 1, .1}), []string{"3", "1", "2"})
	testUtils.TestHeuristicRanking(t, heu, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 0, "3": 0}}, []model.Weight{10, 1, 1}), []string{"2", "3", "1"})
	testUtils.TestHeuristicRanking(t, heu, dmParams([]model.Weights{{"1": 1, "2": 1000, "3": 3}, {"1": 2, "2": 600, "3": -2}}, []model.Weight{1, 0, 1}), []string{"2", "3", "1"})
}

var wsTestCriteria = model.Criteria{{Id: "1"}, {Id: "2"}, {Id: "3"}}
var heu = &WeightedSumHeuristic{}

func dmParams(alternativesWeights []model.Weights, criteriaWeights []model.Weight) *model.DecisionMakingParams {
	alternatives := testUtils.WrapAlternatives(alternativesWeights)
	criteria := assignCriteria(wsTestCriteria, criteriaWeights)
	return &model.DecisionMakingParams{
		ConsideredAlternatives: alternatives,
		Criteria:               wsTestCriteria,
		MethodParameters:       weightedSumParams{weightedCriteria: &criteria},
	}
}

func assignCriteria(criteria []model.Criterion, criteriaWeights []model.Weight) []model.WeightedCriterion {
	weightedCriteria := make([]model.WeightedCriterion, len(criteriaWeights))
	for i, c := range criteriaWeights {
		weightedCriteria[i] = model.WeightedCriterion{
			Criterion: criteria[i],
			Weight:    c,
		}
	}
	return weightedCriteria
}
