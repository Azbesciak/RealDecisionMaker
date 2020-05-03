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
	addedCriterion := bias.OnCriterionAdded(
		criterion,
		&model.Criteria{{Id: "1", Type: model.Gain}, {Id: "5", Type: model.Cost}},
		owaParams{weights: &[]model.Weight{1, 3}},
		generator(123),
	)
	expected := model.SingleWeight(criterion, 0.5)
	if utils.Differs(expected, addedCriterion) {
		t.Errorf("wrong value for owa bias add result, expected %v, got %v", expected, addedCriterion)
	}
}

func TestOwaCriterionMerging(t *testing.T) {
	criterion := &model.Criterion{Id: "3", Type: model.Cost}
	added := model.SingleWeight(criterion, 0.5)
	expected := owaParams{weights: &[]model.Weight{0.5, 1, 3}}
	actual := bias.Merge(owaParams{weights: &[]model.Weight{1, 3}}, added)
	if p, ok := actual.(owaParams); !ok {
		t.Errorf("invalid result type from owa, expected owaParams")
	} else if utils.Differs(expected.weights, p.weights) {
		t.Errorf("wrong value for owa bias merge result, expected %v, got %v", expected, actual)
	}
}

func TestOwaBiasListener_OnCriteriaRemoved(t *testing.T) {
	removed := owaTestCriteria[:1]
	left := owaTestCriteria[1:]
	expected := []model.Weight{2, 3}
	original := []model.Weight{1, 2, 3}
	testRemoval(t, removed, left, original, expected)
}

func TestOwaBiasListener_OnCriteriaRemoved_OnlyLeftConsidered(t *testing.T) {
	removed := owaTestCriteria
	left := owaTestCriteria
	expected := []model.Weight{1, 2, 3}
	original := []model.Weight{1, 2, 3}
	testRemoval(t, removed, left, original, expected)
}

func TestOwaBiasListener_OnCriteriaRemoved_nothingLeft(t *testing.T) {
	removed := owaTestCriteria
	left := model.Criteria{}
	expected := make([]model.Weight, 0)
	original := []model.Weight{1, 2, 3}
	testRemoval(t, removed, left, original, expected)
}

func testRemoval(t *testing.T, removed model.Criteria, left model.Criteria, original []model.Weight, expected []model.Weight) {
	actual := bias.OnCriteriaRemoved(&removed, &left, owaParams{weights: &original})
	if a, ok := actual.(owaParams); !ok {
		t.Errorf("expected owaParams result type")
	} else if utils.Differs(expected, *a.weights) {
		t.Errorf("expected owa weights %v, got %v", expected, *a.weights)
	}
}
