package owa

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"testing"
)

var owaTestCriteria = model.Criteria{{Id: "1"}, {Id: "2"}, {Id: "3"}}
var heu = &OwaHeuristic{}

func dmParams(alternativesWeights []model.Weights) *model.DecisionMakingParams {
	return &model.DecisionMakingParams{
		ConsideredAlternatives: testUtils.WrapAlternatives(alternativesWeights),
		Criteria:               owaTestCriteria,
	}
}

func TestOwaHeuristic_RankCriteriaAscending(t *testing.T) {
	testUtils.TestHeuristicRanking(t, heu, dmParams([]model.Weights{}), []string{"1", "2", "3"})
	testUtils.TestHeuristicRanking(t, heu, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}}), []string{"1", "2", "3"})
	testUtils.TestHeuristicRanking(t, heu, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 1, "3": 3}}), []string{"1", "2", "3"})
	testUtils.TestHeuristicRanking(t, heu, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 1, "3": 0}}), []string{"1", "2", "3"})
	testUtils.TestHeuristicRanking(t, heu, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 0, "3": 0}}), []string{"2", "1", "3"})
	testUtils.TestHeuristicRanking(t, heu, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 0, "3": -2}}), []string{"3", "2", "1"})
}
