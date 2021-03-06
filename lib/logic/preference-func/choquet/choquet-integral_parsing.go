package choquet

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"sort"
	"strings"
)

type choquetParams struct {
	weights  *model.Weights
	criteria *model.Criteria
}

func (c *ChoquetIntegralPreferenceFunc) ParseParams(dm *model.DecisionMaker) interface{} {
	weights := model.ExtractWeights(dm)
	parsedWeights := parse(&dm.Criteria, &weights)
	return choquetParams{weights: parsedWeights, criteria: &dm.Criteria}
}

func parse(criteria *model.Criteria, weights *model.Weights) *model.Weights {
	validateAllCriteriaAreGain(criteria)
	newWeights := remapWeights(weights)
	validateAllWeightsAvailable(newWeights, criteria)
	return prepareWeights(newWeights, criteria)
}

func remapWeights(weights *model.Weights) *model.Weights {
	result := make(model.Weights, len(*weights))
	for k, v := range *weights {
		criteria := containedCriteria(k)
		validCriterionKey := criterionKey(&criteria)
		if _, ok := result[validCriterionKey]; ok {
			if len(criteria) == 1 {
				panic(fmt.Errorf("value for criterion %v is redeclared", validCriterionKey))
			} else {
				panic(fmt.Errorf("values for criteria %v are redeclared", validCriterionKey))
			}
		}
		result[validCriterionKey] = v
	}
	return &result
}

func validateAllCriteriaAreGain(criteria *model.Criteria) {
	for _, c := range *criteria {
		if c.Type != model.Gain {
			panic(fmt.Errorf("%s: only Gain criteria acceptable for Choquet integral", c.Id))
		}
	}
}

func validateAllWeightsAvailable(weights *model.Weights, criteria *model.Criteria) {
	criteriaNames := criteria.Names()
	requiredCriteriaCombinations := *PowerSet(*criteriaNames)
	for _, rcc := range requiredCriteriaCombinations {
		getWeightForCriteriaUnion(&rcc, weights)
	}
}

const criteriaSeparator = ","

func containedCriteria(key string) []string {
	return strings.Split(key, criteriaSeparator)
}

func prepareWeights(weights *model.Weights, criteria *model.Criteria) *model.Weights {
	resultWeights := make(model.Weights, len(*weights))
	for k, v := range *weights {
		splittedValues := containedCriteria(k)
		identifiable := utils.ToIdentifiable(criteria)
		if !utils.ContainsAll(identifiable, &splittedValues) {
			panic(fmt.Errorf("%s: not all weights are present in criteria %s", k, *criteria))
		}
		validateWeightValue(&splittedValues, v)
		resultWeights[criterionKey(&splittedValues)] = v
	}
	return &resultWeights
}

func validateWeightValue(criterionKey *[]string, v model.Weight) {
	if v < 0 || v > 1 {
		panic(fmt.Errorf("%s: weight must be in range [0,1], got %f", *criterionKey, v))
	}
}

func getWeightForCriteriaUnion(commonWeightCriteria *[]string, weights *model.Weights) model.Weight {
	weightKey := criterionKey(commonWeightCriteria)
	criteriaUnionWeight := getWeightForCombinedCriterion(weights, &weightKey)
	return criteriaUnionWeight
}

func getWeightForCombinedCriterion(weights *model.Weights, weightKey *string) model.Weight {
	criteriaUnionWeight, ok := (*weights)[*weightKey]
	if !ok {
		panic(fmt.Errorf("weight for criteria union '%s' not found", *weightKey))
	}
	return criteriaUnionWeight
}

func criterionKey(criteria *[]string) string {
	sort.Strings(*criteria)
	return strings.Join(*criteria, criteriaSeparator)
}
