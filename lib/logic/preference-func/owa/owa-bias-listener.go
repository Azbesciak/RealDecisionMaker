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
	return *oldParams.merge(&newParams)
}

func (h *OwaBiasListener) OnCriterionAdded(
	criterion *model.Criterion,
	referenceCriterion *model.Criterion,
	params model.MethodParameters,
	generator utils.ValueGenerator,
) model.MethodParameters {
	owaPar := params.(owaParams)
	referenceCriterionValue := owaPar.find(referenceCriterion)
	newWeight := generator() * referenceCriterionValue.Weight
	return model.SingleWeight(criterion, newWeight)
}

func (h *OwaBiasListener) OnCriteriaRemoved(
	leftCriteria *model.Criteria,
	params model.MethodParameters,
) model.MethodParameters {
	owaPar := params.(owaParams)
	newWeights := make(model.WeightedCriteria, len(*leftCriteria))
	for i, c := range *leftCriteria {
		newWeights[i] = *owaPar.find(&c)
	}
	return owaParams{Weights: &newWeights}
}

func (h *OwaBiasListener) RankCriteriaAscending(params *model.DecisionMakingParams) *model.WeightedCriteria {
	weights := *model.PrepareCumulatedWeightsMap(params, model.WeightIdentity)
	return params.Criteria.SortByWeights(weights)
}
