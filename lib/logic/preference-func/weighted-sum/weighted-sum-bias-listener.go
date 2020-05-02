package weighted_sum

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type WeightedSumBiasListener struct {
}

func (w *WeightedSumBiasListener) Identifier() string {
	return methodName
}

type WeightedSumAddedCriterion struct {
	Weights model.Weights `json:"weights"`
	weights []model.WeightedCriterion
}

func (w *WeightedSumBiasListener) Merge(params model.MethodParameters, addition model.MethodParameters) model.MethodParameters {
	oldWeights := *params.(weightedSumParams).weightedCriteria
	addedWeights := addition.(WeightedSumAddedCriterion).weights
	merged := append(oldWeights, addedWeights...)
	return weightedSumParams{weightedCriteria: &merged}
}

func (w *WeightedSumBiasListener) OnCriterionAdded(
	criterion *model.Criterion,
	previousRankedCriteria *model.Criteria,
	params model.MethodParameters,
	generator utils.ValueGenerator,
) model.AddedCriterionParams {
	wParams := params.(weightedSumParams)
	leastImportantParam := wParams.Criterion(previousRankedCriteria.First().Identifier())
	newCriterionWeight := generator() * leastImportantParam.Weight
	newCriterion := []model.WeightedCriterion{{
		Criterion: *criterion,
		Weight:    newCriterionWeight,
	}}
	return WeightedSumAddedCriterion{
		Weights: model.Weights{criterion.Id: newCriterionWeight},
		weights: newCriterion,
	}
}

func (w *WeightedSumBiasListener) OnCriteriaRemoved(removedCriteria *model.Criteria, leftCriteria *model.Criteria, params model.MethodParameters) model.MethodParameters {
	wParams := params.(weightedSumParams)
	result := make([]model.WeightedCriterion, len(*leftCriteria))
	for i, c := range *leftCriteria {
		result[i] = wParams.Criterion(c.Id)
	}
	return weightedSumParams{weightedCriteria: &result}
}

func (w *WeightedSumBiasListener) RankCriteriaAscending(params *model.DecisionMakingParams) *model.Criteria {
	wParams := params.MethodParameters.(weightedSumParams)
	weights := model.PrepareCumulatedWeightsMap(params, func(criterion string, value model.Weight) model.Weight {
		cryt := wParams.Criterion(criterion)
		return cryt.Weight * value
	})
	return params.Criteria.SortByWeights(*weights)
}
