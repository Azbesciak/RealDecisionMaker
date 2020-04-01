package choquet

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"sort"
)

const methodName = "choquetIntegral"

type ChoquetIntegralPreferenceFunc struct {
}

func (c *ChoquetIntegralPreferenceFunc) Identifier() string {
	return methodName
}

func (c *ChoquetIntegralPreferenceFunc) MethodParameters() interface{} {
	return model.WeightsParamOnly()
}

func (c *ChoquetIntegralPreferenceFunc) Evaluate(dmp *model.DecisionMakingParams) *model.AlternativesRanking {
	params := dmp.MethodParameters.(choquetParams)
	prefFunc := func(alternative *model.AlternativeWithCriteria) *model.AlternativeResult {
		return choquetIntegral(alternative, params.weights)
	}
	return model.Rank(dmp, prefFunc)
}

func ChoquetIntegral(
	alternative model.AlternativeWithCriteria,
	criteria model.Criteria,
	weights model.Weights,
) *model.AlternativeResult {
	resultWeights := parse(&criteria, &weights)
	return choquetIntegral(&alternative, resultWeights)
}

func choquetIntegral(
	alternative *model.AlternativeWithCriteria,
	weights *model.Weights,
) *model.AlternativeResult {
	sortedCriteria := prepareCriteriaInAscendingOrder(alternative)
	result, _ := computeTotalWeight(sortedCriteria, weights)
	return &model.AlternativeResult{Alternative: *alternative, Value: result}
}

type weightComponent struct {
	criteria   []string
	valueAdded model.Weight
}

func computeTotalWeight(sortedCriteria *criteriaWeights, weights *model.Weights) (model.Weight, []weightComponent) {
	var result model.Weight = 0
	var previousWeight model.Weight = 0
	var components []weightComponent
	totalElements := len(*sortedCriteria)
	for i := 0; i < totalElements; {
		commonWeightCriteria := make([]string, totalElements-i)
		for x := 0; x < totalElements-i; x++ {
			commonWeightCriteria[x] = (*sortedCriteria)[i+x].criterion
		}
		var current = (*sortedCriteria)[i]
		var j int
		for j = i + 1; j < totalElements; j++ {
			var nextValue = (*sortedCriteria)[j]
			if !utils.FloatsAreEqual(current.weight, nextValue.weight, 0.00001) {
				break
			}
		}
		criteriaUnionWeight := getWeightForCriteriaUnion(&commonWeightCriteria, weights)
		valueAdded := criteriaUnionWeight * (current.weight - previousWeight)
		result += valueAdded
		components = append(components, weightComponent{commonWeightCriteria, valueAdded})
		previousWeight = current.weight
		i = j
	}
	return result, components
}

func prepareCriteriaInAscendingOrder(alternative *model.AlternativeWithCriteria) *criteriaWeights {
	var sorted = make(criteriaWeights, len(alternative.Criteria))
	i := 0
	for k, v := range alternative.Criteria {
		sorted[i] = criterionWeight{k, v}
		i++
	}
	sort.Sort(&sorted)
	return &sorted
}

type criterionWeight struct {
	criterion string
	weight    model.Weight
}

type criteriaWeights []criterionWeight

func (c *criteriaWeights) Len() int {
	return len(*c)
}
func (c *criteriaWeights) Less(i, j int) bool {
	return (*c)[i].weight < (*c)[j].weight
}

func (c *criteriaWeights) Swap(i, j int) {
	(*c)[i], (*c)[j] = (*c)[j], (*c)[i]
}
