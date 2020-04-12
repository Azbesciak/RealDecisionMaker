package owa

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type OwaBiasListener struct {
}

func (h *OwaBiasListener) Identifier() string {
	return methodName
}

func (h *OwaBiasListener) Merge(params model.MethodParameters, addition model.MethodParameters) model.MethodParameters {
	oldParams := params.(owaParams)
	newParams := addition.(owaParams)
	return *oldParams.withWeight(newParams.minWeight())
}

func (h *OwaBiasListener) OnCriterionAdded(
	criterion *model.Criterion,
	previousRankedCriteria *model.Criteria,
	params model.MethodParameters,
	generator utils.ValueGenerator,
) model.MethodParameters {
	owaPar := params.(owaParams)
	newWeight := generator() * owaPar.minWeight()
	return owaParams{weights: &[]model.Weight{newWeight}}
}

func (h *OwaBiasListener) OnCriteriaRemoved(
	removedCriteria *model.Criteria,
	leftCriteria *model.Criteria,
	params model.MethodParameters,
) model.MethodParameters {
	owaPar := params.(owaParams)
	return *owaPar.withoutNWorstWeights(len(*removedCriteria))
}

func (h *OwaBiasListener) RankCriteriaAscending(params *model.DecisionMakingParams) *model.Criteria {
	weights := *model.PrepareCumulatedWeightsMap(params, model.WeightIdentity)
	return params.Criteria.SortByWeights(weights)
}
