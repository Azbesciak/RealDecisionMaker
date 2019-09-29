package logic

import (
	"../model"
	"fmt"
	"sort"
	"strings"
)

func ChoquetIntegral(alternative model.AlternativeWithCriteria, criteria model.Criteria, weights model.Weights) model.AlternativeResult {
	validateAllCriteriaAreGain(&criteria)
	resultWeights := prepareWeights(&weights, &criteria)
	sortedCriteria := prepareCriteriaInAscendingOrder(alternative)
	result := computeTotalWeight(sortedCriteria, resultWeights)
	return model.AlternativeResult{alternative, result}
}

func validateAllCriteriaAreGain(criteria *model.Criteria) {
	for _, c := range *criteria {
		if c.Type != model.Gain {
			panic(fmt.Errorf("%s: only Gain criteria acceptable for Choquet integral", c.Id))
		}
	}
}

func computeTotalWeight(sortedCriteria *criteriaWeights, weights *model.Weights) model.Weight {
	var result model.Weight = 0
	var previousWeight model.Weight = 0
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
			if !FloatsAreEqual(current.weight, nextValue.weight, 0.00001) {
				break
			}
		}
		sort.Strings(commonWeightCriteria)
		weightKey := strings.Join(commonWeightCriteria, ",")
		criteriaUnionWeight, ok := (*weights)[weightKey]

		if !ok {
			panic(fmt.Errorf("weight for criteria union '%s' not found", weightKey))
		}
		result += criteriaUnionWeight * (current.weight - previousWeight)
		previousWeight = current.weight
		i = j
	}
	return result
}

func prepareCriteriaInAscendingOrder(alternative model.AlternativeWithCriteria) *criteriaWeights {
	var sorted = make(criteriaWeights, len(alternative.Criteria))
	i := 0
	for k, v := range alternative.Criteria {
		sorted[i] = criterionWeight{k, v}
		i++
	}
	sort.Sort(&sorted)
	return &sorted
}

func prepareWeights(weights *model.Weights, criteria *model.Criteria) *model.Weights {
	resultWeights := make(model.Weights, len(*weights))
	for k, v := range *weights {
		splittedValues := strings.Split(k, ",")
		identificable := ToIdentifiable(criteria)
		if !ContainsAll(identificable, &splittedValues) {
			panic(fmt.Errorf("%s: not all weights are present in criteria %s", k, *criteria))
		}
		if v < 0 || v > 1 {
			panic(fmt.Errorf("%s: weight must be in range [0,1], got %f", k, v))
		}
		sort.Strings(splittedValues)
		resultWeights[strings.Join(splittedValues, ",")] = v
	}
	return &resultWeights
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
