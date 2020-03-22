package owa

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"math/rand"
)

type OwaHeuristic struct {
}

func (h *OwaHeuristic) Identifier() string {
	return methodName
}

func (h *OwaHeuristic) OnCriterionAdded(
	criterion *model.Criterion,
	previousRankedCriteria *model.Criteria,
	params model.MethodParameters,
	generator utils.ValueGenerator,
) *model.AddCriterionResult {
	owaPar := params.(owaParams)
	newWeight := rand.Float64() * owaPar.minWeight()
	return &model.AddCriterionResult{
		MethodParameters:     owaPar.withWeight(newWeight),
		AddedCriterionParams: owaParams{weights: &[]model.Weight{newWeight}},
	}
}

func (h *OwaHeuristic) OnCriteriaRemoved(
	removedCriteria *model.Criteria,
	leftCriteria *model.Criteria,
	params model.MethodParameters,
) model.MethodParameters {
	owaPar := params.(owaParams)
	return owaPar.withoutNWorstWeights(len(*removedCriteria))
}

func (h *OwaHeuristic) RankCriteriaAscending(params *model.DecisionMakingParams) *model.Criteria {
	weights := *prepareAccumulatedValuesMap(params)
	return params.Criteria.SortByWeights(weights)
}

func prepareAccumulatedValuesMap(params *model.DecisionMakingParams) *model.Weights {
	weights := make(model.Weights, len(params.Criteria))
	for _, a := range params.ConsideredAlternatives {
		for crit, v := range a.Criteria {
			w, ok := weights[crit]
			if !ok {
				weights[crit] = v
			} else {
				weights[crit] = w + v
			}
		}
	}
	return &weights
}
