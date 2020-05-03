package choquet

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type ChoquetIntegralBiasListener struct {
}

func (c *ChoquetIntegralBiasListener) Identifier() string {
	return methodName
}

func (c *ChoquetIntegralBiasListener) OnCriterionAdded(
	criterion *model.Criterion,
	previousRankedCriteria *model.Criteria,
	params model.MethodParameters,
	generator utils.ValueGenerator,
) model.AddedCriterionParams {
	oldWeights := params.(choquetParams).weights
	newCriteria := previousRankedCriteria.Add(criterion)
	newWeightsKeys := PowerSet(*newCriteria.Names())
	newWeights := make(model.Weights, len(*newWeightsKeys))
	for _, k := range *newWeightsKeys {
		cKey := criterionKey(&k)
		weight, ok := (*oldWeights)[cKey]
		if ok {
			newWeights[cKey] = weight
			continue
		}
		originalKeyCriteriaWithoutNewOne := utils.RemoveSingleStringOccurrence(k, criterion.Id)
		if len(originalKeyCriteriaWithoutNewOne) == 0 {
			newWeights[criterion.Id] = generator()
			continue
		}
		newWeights[cKey] = getWeightForCriteriaUnion(&originalKeyCriteriaWithoutNewOne, oldWeights)
	}
	return choquetParams{weights: &newWeights}
}

func (c *ChoquetIntegralBiasListener) OnCriteriaRemoved(
	leftCriteria *model.Criteria,
	params model.MethodParameters,
) model.MethodParameters {
	cParams := params.(choquetParams)
	expectedSet := *PowerSet(*leftCriteria.Names())
	filteredWeights := make(model.Weights, len(expectedSet))
	for _, criteria := range expectedSet {
		key := criterionKey(&criteria)
		filteredWeights[key] = cParams.weights.Fetch(key)
	}
	return choquetParams{weights: &filteredWeights}
}

func (c *ChoquetIntegralBiasListener) RankCriteriaAscending(params *model.DecisionMakingParams) *model.Criteria {
	criteriaWeights := decomposeWeights(params)
	return params.Criteria.SortByWeights(*criteriaWeights)
}

func decomposeWeights(params *model.DecisionMakingParams) *model.Weights {
	combinedWeights := *params.MethodParameters.(choquetParams).weights
	weights := make(model.Weights, len(params.Criteria))
	for _, c := range params.Criteria {
		weights[c.Id] = 0
	}
	for _, a := range params.ConsideredAlternatives {
		sortedCriteria := prepareCriteriaInAscendingOrder(&a)
		_, w := computeTotalWeight(sortedCriteria, &combinedWeights)
		for _, criteriaValues := range w {
			for _, c := range criteriaValues.criteria {
				w, ok := weights[c]
				if !ok {
					panic(fmt.Errorf("criterion '%s' not found in criteria %v", c, params.Criteria))
				} else {
					weights[c] = w + criteriaValues.valueAdded
				}
			}
		}
	}
	return &weights
}

func (c *ChoquetIntegralBiasListener) Merge(params model.MethodParameters, addition model.MethodParameters) model.MethodParameters {
	oldWeights := *params.(choquetParams).weights
	addedWeights := *addition.(choquetParams).weights
	return choquetParams{weights: oldWeights.Merge(&addedWeights)}
}
