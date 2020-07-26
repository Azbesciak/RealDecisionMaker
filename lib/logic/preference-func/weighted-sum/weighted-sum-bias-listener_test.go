package weighted_sum

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

func TestWeightedSumBias_RankCriteriaAscending_EqualWeights(t *testing.T) {
	testUtils.TestBiasRanking(0, t, bias, dmParams([]model.Weights{}, []model.Weight{}), []string{"1", "2", "3"})
	testUtils.TestBiasRanking(1, t, bias, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}}, []model.Weight{1, 1, 1}), []string{"1", "2", "3"})
	testUtils.TestBiasRanking(2, t, bias, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 1, "3": 3}}, []model.Weight{1, 1, 1}), []string{"1", "2", "3"})
	testUtils.TestBiasRanking(3, t, bias, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 1, "3": 0}}, []model.Weight{1, 1, 1}), []string{"1", "2", "3"})
	testUtils.TestBiasRanking(4, t, bias, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 0, "3": 0}}, []model.Weight{1, 1, 1}), []string{"2", "1", "3"})
	testUtils.TestBiasRanking(5, t, bias, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 0, "3": -2}}, []model.Weight{1, 1, 1}), []string{"3", "2", "1"})
}

func TestWeightedSumBias_RankCriteriaAscending(t *testing.T) {
	testUtils.TestBiasRanking(0, t, bias, dmParams([]model.Weights{}, []model.Weight{}), []string{"1", "2", "3"})
	testUtils.TestBiasRanking(1, t, bias, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}}, []model.Weight{1, 1, 2}), []string{"1", "2", "3"})
	testUtils.TestBiasRanking(2, t, bias, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 1, "3": 3}}, []model.Weight{10, 0.5, 0.1}), []string{"3", "2", "1"})
	testUtils.TestBiasRanking(3, t, bias, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 1, "3": 0}}, []model.Weight{1, 1, .1}), []string{"3", "1", "2"})
	testUtils.TestBiasRanking(4, t, bias, dmParams([]model.Weights{{"1": 1, "2": 2, "3": 3}, {"1": 2, "2": 0, "3": 0}}, []model.Weight{10, 1, 1}), []string{"2", "3", "1"})
	testUtils.TestBiasRanking(5, t, bias, dmParams([]model.Weights{{"1": 1, "2": 1000, "3": 3}, {"1": 2, "2": 600, "3": -2}}, []model.Weight{1, 0, 1}), []string{"2", "3", "1"})
}

var wsTestCriteria = model.Criteria{{Id: "1"}, {Id: "2"}, {Id: "3"}}
var bias = &WeightedSumBiasListener{}

func dmParams(alternativesWeights []model.Weights, criteriaWeights []model.Weight) *model.DecisionMakingParams {
	alternatives := testUtils.WrapAlternatives(alternativesWeights)
	criteria := assignCriteria(wsTestCriteria, criteriaWeights)
	return &model.DecisionMakingParams{
		ConsideredAlternatives: alternatives,
		Criteria:               wsTestCriteria,
		MethodParameters:       weightedSumParams{weightedCriteria: &criteria},
	}
}

func assignCriteria(criteria []model.Criterion, criteriaWeights []model.Weight) model.WeightedCriteria {
	weightedCriteria := make(model.WeightedCriteria, len(criteriaWeights))
	for i, c := range criteriaWeights {
		weightedCriteria[i] = model.WeightedCriterion{
			Criterion: criteria[i],
			Weight:    c,
		}
	}
	return weightedCriteria
}

func TestWSCriterionAdding(t *testing.T) {
	generator := testUtils.CyclicRandomGenerator(0, 2)
	criterion := &model.Criterion{Id: "3", Type: model.Cost}
	c1 := model.Criterion{Id: "1", Type: model.Gain}
	c2 := model.Criterion{Id: "2", Type: model.Cost}
	addedCriterion := bias.OnCriterionAdded(
		criterion,
		&c1,
		weightedSumParams{
			weightedCriteria: &model.WeightedCriteria{{
				Criterion: c1,
				Weight:    12,
			}, {
				Criterion: c2,
				Weight:    5,
			}},
		},
		generator(123),
	)
	expected := WeightedSumAddedCriterion{
		Weights: model.Weights{criterion.Id: 6},
		weights: model.WeightedCriteria{{*criterion, 6}},
	}
	if utils.Differs(expected, addedCriterion) {
		t.Errorf("wrong value for owa bias add result, expected %v, got %v", expected, addedCriterion)
	}
}

func TestWSCriterionMerging(t *testing.T) {
	c1 := model.Criterion{Id: "1", Type: model.Gain}
	c2 := model.Criterion{Id: "2", Type: model.Cost}
	c1w := model.WeightedCriterion{Criterion: c1, Weight: 12}
	c2w := model.WeightedCriterion{Criterion: c2, Weight: 5}

	c3 := model.Criterion{Id: "3", Type: model.Cost}
	c3w := model.WeightedCriterion{Criterion: c3, Weight: 2.5}

	added := WeightedSumAddedCriterion{
		Weights: model.Weights{c3.Id: c3w.Weight},
		weights: model.WeightedCriteria{c3w},
	}
	current := weightedSumParams{weightedCriteria: &model.WeightedCriteria{c1w, c2w}}
	actual := bias.Merge(current, added)
	expected := weightedSumParams{weightedCriteria: &model.WeightedCriteria{c1w, c2w, c3w}}
	if wsParsed, ok := actual.(weightedSumParams); !ok {
		t.Errorf("invalid result type from weighted sum, expected weightedSumParams")
	} else if utils.Differs(expected.weightedCriteria, wsParsed.weightedCriteria) {
		t.Errorf("wrong value for weighted sum bias merge result, expected %v, got %v", expected, actual)
	}
}
