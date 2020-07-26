package owa

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
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

func TestOWACriterionAdding(t *testing.T) {
	generator := testUtils.CyclicRandomGenerator(0, 2)
	criterion := &model.Criterion{Id: "3", Type: model.Cost}
	criteria := model.Criteria{{Id: "1", Type: model.Gain}, {Id: "5", Type: model.Cost}}
	addedCriterion := bias.OnCriterionAdded(
		criterion,
		&criteria[0],
		owaParams{Weights: &model.WeightedCriteria{
			{Criterion: criteria[0], Weight: 1},
			{Criterion: criteria[1], Weight: 3},
		}},
		generator(123),
	)
	expected := model.SingleWeight(criterion, 0.5)
	if utils.Differs(expected, addedCriterion) {
		t.Errorf("wrong value for owa bias add result, expected %v, got %v", expected, addedCriterion)
	}
}

func TestOwaCriterionMerging(t *testing.T) {
	added := model.WeightedCriteria{{Criterion: model.Criterion{Id: "3", Type: model.Cost}, Weight: 0.5}}
	current := model.WeightedCriteria{
		{Criterion: model.Criterion{Id: "1", Type: model.Gain}, Weight: 1},
		{Criterion: model.Criterion{Id: "2", Type: model.Gain}, Weight: 3},
	}
	expected := owaParams{Weights: &model.WeightedCriteria{added[0], current[0], current[1]}}
	actual := bias.Merge(owaParams{Weights: &current}, owaParams{Weights: &added})
	if p, ok := actual.(owaParams); !ok {
		t.Errorf("invalid result type from owa, expected owaParams")
	} else if utils.Differs(expected.Weights, p.Weights) {
		t.Errorf("wrong value for owa bias merge result, expected %v, got %v", expected, actual)
	}
}

func TestOwaBiasListener_OnCriteriaRemoved(t *testing.T) {
	left := owaTestCriteria[1:]
	expected := model.WeightedCriteria{
		{Criterion: owaTestCriteria[1], Weight: 2},
		{Criterion: owaTestCriteria[2], Weight: 3},
	}
	original := model.WeightedCriteria{
		{Criterion: owaTestCriteria[0], Weight: 1},
		{Criterion: owaTestCriteria[1], Weight: 2},
		{Criterion: owaTestCriteria[2], Weight: 3},
	}
	testRemoval(t, left, original, expected)
}

func TestOwaBiasListener_OnCriteriaRemoved_OnlyLeftConsidered(t *testing.T) {
	left := owaTestCriteria
	original := model.WeightedCriteria{
		{Criterion: owaTestCriteria[0], Weight: 1},
		{Criterion: owaTestCriteria[1], Weight: 2},
		{Criterion: owaTestCriteria[2], Weight: 3},
	}
	expected := original
	testRemoval(t, left, original, expected)
}

func TestOwaBiasListener_OnCriteriaRemoved_nothingLeft(t *testing.T) {
	left := model.Criteria{}
	original := model.WeightedCriteria{
		{Criterion: owaTestCriteria[0], Weight: 1},
		{Criterion: owaTestCriteria[1], Weight: 2},
		{Criterion: owaTestCriteria[2], Weight: 3},
	}
	expected := make(model.WeightedCriteria, 0)
	testRemoval(t, left, original, expected)
}

func testRemoval(t *testing.T, left model.Criteria, original, expected model.WeightedCriteria) {
	actual := bias.OnCriteriaRemoved(&left, owaParams{Weights: &original})
	if a, ok := actual.(owaParams); !ok {
		t.Errorf("expected owaParams result type")
	} else if utils.Differs(expected, *a.Weights) {
		t.Errorf("expected owa Weights %v, got %v", expected, *a.Weights)
	}
}
