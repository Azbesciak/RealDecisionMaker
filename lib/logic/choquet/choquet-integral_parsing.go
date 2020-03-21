package choquet

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"sort"
	"strings"
)

type choquetParams struct {
	weights *model.Weights
}

func (c *ChoquetIntegralPreferenceFunc) ParseParams(dm *model.DecisionMaker) interface{} {
	weights := model.ExtractWeights(dm)
	parsedWeights := parse(&dm.Criteria, &weights)
	return &choquetParams{weights: parsedWeights}
}

func parse(criteria *model.Criteria, weights *model.Weights) *model.Weights {
	validateAllCriteriaAreGain(criteria)
	return prepareWeights(weights, criteria)
}

func validateAllCriteriaAreGain(criteria *model.Criteria) {
	for _, c := range *criteria {
		if c.Type != model.Gain {
			panic(fmt.Errorf("%s: only Gain criteria acceptable for Choquet integral", c.Id))
		}
	}
}

func prepareWeights(weights *model.Weights, criteria *model.Criteria) *model.Weights {
	resultWeights := make(model.Weights, len(*weights))
	for k, v := range *weights {
		splittedValues := strings.Split(k, ",")
		identifiable := utils.ToIdentifiable(criteria)
		if !utils.ContainsAll(identifiable, &splittedValues) {
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
