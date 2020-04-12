package owa

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"testing"
)

var owaTestCriteria = model.Criteria{{Id: "1"}, {Id: "2"}, {Id: "3"}}
var bias = &OwaBiasListener{}

func dmParams(alternativesWeights []model.Weights) *model.DecisionMakingParams {
	return &model.DecisionMakingParams{
		ConsideredAlternatives: testUtils.WrapAlternatives(alternativesWeights),
		Criteria:               owaTestCriteria,
	}
}

func TestOwaBiasListener_RankCriteriaAscending(t *testing.T) {
	testUtils.TestBiasRanking(0, t, bias, dmParams([]model.Weights{}), []string{"1", "2", "3"})
	testUtils.TestBiasRanking(1, t, bias, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}}), []string{"1", "2", "3"})
	testUtils.TestBiasRanking(2, t, bias, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 1, "3": 3}}), []string{"1", "2", "3"})
	testUtils.TestBiasRanking(3, t, bias, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 1, "3": 0}}), []string{"1", "2", "3"})
	testUtils.TestBiasRanking(4, t, bias, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 0, "3": 0}}), []string{"2", "1", "3"})
	testUtils.TestBiasRanking(5, t, bias, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 0, "3": -2}}), []string{"3", "2", "1"})
}
