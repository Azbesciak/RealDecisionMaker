package owa

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"sort"
)

type OWAPreferenceFunc struct {
}

func (O *OWAPreferenceFunc) Identifier() string {
	return "owa"
}

func (O *OWAPreferenceFunc) MethodParameters() interface{} {
	return model.WeightsParamOnly()
}

func (O *OWAPreferenceFunc) Evaluate(dm *model.DecisionMaker) *model.AlternativesRanking {
	originalWeights := model.ExtractWeights(dm)
	var weights = make([]model.Weight, len(originalWeights))
	i := 0
	for _, v := range originalWeights {
		weights[i] = v
		i++
	}
	prefFunc := func(alternative *model.AlternativeWithCriteria) *model.AlternativeResult {
		return OWA(*alternative, weights)
	}
	return model.Rank(dm, prefFunc)
}

func OWA(alternative model.AlternativeWithCriteria, weights []model.Weight) *model.AlternativeResult {
	alternativeCriteria := len(alternative.Criteria)
	if alternativeCriteria != len(weights) {
		panic(fmt.Errorf("criteria and weights must have the same length, got %d and %d", alternativeCriteria, len(weights)))
	}
	tmpWeights := make([]model.Weight, len(weights))
	copy(tmpWeights, weights)
	sort.Float64s(tmpWeights)
	tmpCriteria := make([]model.Weight, alternativeCriteria)
	i := 0
	for _, c := range alternative.Criteria {
		tmpCriteria[i] = c
		i += 1
	}
	sort.Float64s(tmpCriteria)
	var total model.Weight = 0
	for i = 0; i < alternativeCriteria; i++ {
		total += tmpCriteria[i] * tmpWeights[i]
	}
	return &model.AlternativeResult{Alternative: alternative, Value: total}
}
