package weighted_sum

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type WeightedSumHeuristic struct {
}

func (w *WeightedSumHeuristic) Identifier() string {
	return methodName
}

func (w *WeightedSumHeuristic) Merge(params model.MethodParameters, addition model.MethodParameters) model.MethodParameters {
	oldWeights := *params.(weightedSumParams).weightedCriteria
	addedWeights := *addition.(weightedSumParams).weightedCriteria
	merged := append(oldWeights, addedWeights...)
	return weightedSumParams{weightedCriteria: &merged}
}

func (w *WeightedSumHeuristic) OnCriterionAdded(
	criterion *model.Criterion,
	previousRankedCriteria *model.Criteria,
	params model.MethodParameters,
	generator utils.ValueGenerator,
) model.AddedCriterionParams {
	wParams := params.(weightedSumParams)
	leastImportantParam := wParams.Criterion(previousRankedCriteria.First().Identifier())
	newCriterionWeight := []model.WeightedCriterion{{
		Criterion: *criterion,
		Weight:    generator() * leastImportantParam.Weight,
	}}
	return weightedSumParams{weightedCriteria: &newCriterionWeight}
}

func (w *WeightedSumHeuristic) OnCriteriaRemoved(removedCriteria *model.Criteria, leftCriteria *model.Criteria, params model.MethodParameters) model.MethodParameters {
	wParams := params.(weightedSumParams)
	result := make([]model.WeightedCriterion, len(*leftCriteria))
	for i, c := range *leftCriteria {
		result[i] = wParams.Criterion(c.Id)
	}
	return weightedSumParams{weightedCriteria: &result}
}

func (w *WeightedSumHeuristic) RankCriteriaAscending(params *model.DecisionMakingParams) *model.Criteria {
	wParams := params.MethodParameters.(weightedSumParams)
	weights := model.PrepareCumulatedWeightsMap(params, func(criterion string, value model.Weight) model.Weight {
		cryt := wParams.Criterion(criterion)
		return cryt.Weight * value
	})
	return params.Criteria.SortByWeights(*weights)
}
