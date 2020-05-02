package owa

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"sort"
)

type OWAPreferenceFunc struct {
}

type owaParams struct {
	weights *[]model.Weight
}

const methodName = "owa"

func (O *OWAPreferenceFunc) ParseParams(dm *model.DecisionMaker) interface{} {
	originalWeights := model.ExtractWeights(dm)
	weightsCount := len(originalWeights)
	if weightsCount != len(dm.Criteria) {
		panic(fmt.Errorf("weights count (%d) not equal to criteria count (%d) for OWA", weightsCount, len(dm.Criteria)))
	}
	weights := toArray(originalWeights)
	_sortWeightsMutate(&weights)
	return owaParams{weights: &weights}
}

func toArray(weights model.Weights) []model.Weight {
	result := make([]model.Weight, len(weights))
	i := 0
	for _, v := range weights {
		result[i] = v
		i++
	}
	return result
}

func (O *OWAPreferenceFunc) Identifier() string {
	return methodName
}

func (O *OWAPreferenceFunc) MethodParameters() interface{} {
	return model.WeightsParamOnly()
}

func (O *OWAPreferenceFunc) Evaluate(dmp *model.DecisionMakingParams) *model.AlternativesRanking {
	weights := dmp.MethodParameters.(owaParams)
	prefFunc := func(alternative *model.AlternativeWithCriteria) *model.AlternativeResult {
		return OWA(*alternative, *weights.weights)
	}
	return model.Rank(dmp, prefFunc)
}

func OWA(alternative model.AlternativeWithCriteria, weights []model.Weight) *model.AlternativeResult {
	sortedWeights := sortWeights(&weights)
	return owa(&alternative, sortedWeights)
}

func owa(alternative *model.AlternativeWithCriteria, sortedWeights *[]model.Weight) *model.AlternativeResult {
	validateSameCriteriaAndWeightsCount(alternative, sortedWeights)
	sortedAlternativeCriteriaWeights := sortAlternativeCriteriaWeights(alternative)
	total := calculateTotalAlternativeValue(sortedWeights, sortedAlternativeCriteriaWeights)
	return model.ValueAlternativeResult(alternative, total)
}

func calculateTotalAlternativeValue(sortedWeights *[]model.Weight, sortedCriteriaWeights *[]model.Weight) model.Weight {
	var total model.Weight = 0
	for i := range *sortedWeights {
		total += (*sortedCriteriaWeights)[i] * (*sortedWeights)[i]
	}
	return total
}

func sortAlternativeCriteriaWeights(alternative *model.AlternativeWithCriteria) *[]model.Weight {
	tmpCriteria := make([]model.Weight, len(alternative.Criteria))
	i := 0
	for _, c := range alternative.Criteria {
		tmpCriteria[i] = c
		i += 1
	}
	_sortWeightsMutate(&tmpCriteria)
	return &tmpCriteria
}

func sortWeights(weights *[]model.Weight) *[]model.Weight {
	tmpWeights := make([]model.Weight, len(*weights))
	copy(tmpWeights, *weights)
	_sortWeightsMutate(&tmpWeights)
	return &tmpWeights
}

func (o *owaParams) withoutNWorstWeights(n int) *owaParams {
	res := (*o.weights)[n:]
	return &owaParams{weights: &res}
}

func (o *owaParams) withWeights(weights model.WeightType) *owaParams {
	newWeights := toArray(weights.Weights)
	result := append(newWeights, *o.weights...)
	_sortWeightsMutate(&result)
	return &owaParams{weights: &result}
}

func (o *owaParams) minWeight() float64 {
	return (*o.weights)[0]
}

func _sortWeightsMutate(weights *[]model.Weight) {
	sort.Float64s(*weights)
}

func validateSameCriteriaAndWeightsCount(alternative *model.AlternativeWithCriteria, weights *[]model.Weight) {
	alternativeCriteriaCount := len(alternative.Criteria)
	weightsCount := len(*weights)
	if alternativeCriteriaCount != weightsCount {
		panic(fmt.Errorf("criteria and weights must have the same length, got %d and %d", alternativeCriteriaCount, weightsCount))
	}
}
