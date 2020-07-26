package owa

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"sort"
)

type OWAPreferenceFunc struct {
}

type owaParams struct {
	Weights *model.WeightedCriteria `json:"weights"`
}

const methodName = "owa"

func (O *OWAPreferenceFunc) ParseParams(dm *model.DecisionMaker) interface{} {
	originalWeights := model.ExtractWeights(dm)
	weightsCount := len(originalWeights)
	if weightsCount != len(dm.Criteria) {
		panic(fmt.Errorf("Weights count (%d) not equal to criteria count (%d) for OWA", weightsCount, len(dm.Criteria)))
	}
	weights := toArray(&originalWeights, &dm.Criteria)
	_sortWeightsMutate(weights)
	return owaParams{Weights: weights}
}

func toArray(weights *model.Weights, criteria *model.Criteria) *model.WeightedCriteria {
	result := make(model.WeightedCriteria, len(*weights))
	for i, v := range *criteria {
		result[i] = model.WeightedCriterion{Criterion: v, Weight: weights.Fetch(v.Id)}
	}
	return &result
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
		return OWA(*alternative, *weights.Weights)
	}
	return model.Rank(dmp, prefFunc)
}

func OWA(alternative model.AlternativeWithCriteria, weights model.WeightedCriteria) *model.AlternativeResult {
	sortedWeights := sortWeights(&weights)
	return owa(&alternative, sortedWeights)
}

func owa(alternative *model.AlternativeWithCriteria, sortedWeights *model.WeightedCriteria) *model.AlternativeResult {
	validateSameCriteriaAndWeightsCount(alternative, sortedWeights)
	sortedAlternativeCriteriaWeights := sortAlternativeCriteriaWeights(alternative)
	total := calculateTotalAlternativeValue(sortedWeights, sortedAlternativeCriteriaWeights)
	return model.ValueAlternativeResult(alternative, total)
}

func calculateTotalAlternativeValue(sortedWeights *model.WeightedCriteria, sortedCriteriaWeights *[]model.Weight) model.Weight {
	var total model.Weight = 0
	for i := range *sortedWeights {
		total += (*sortedCriteriaWeights)[i] * (*sortedWeights)[i].Weight
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
	sort.Float64s(tmpCriteria)
	return &tmpCriteria
}

func sortWeights(weights *model.WeightedCriteria) *model.WeightedCriteria {
	tmpWeights := make(model.WeightedCriteria, len(*weights))
	copy(tmpWeights, *weights)
	_sortWeightsMutate(&tmpWeights)
	return &tmpWeights
}

func (o *owaParams) merge(other *owaParams) *owaParams {
	originalLen := len(*o.Weights)
	result := make(model.WeightedCriteria, originalLen+len(*other.Weights))
	validationCache := make(map[string]bool, originalLen+len(*other.Weights))
	addCriteria(o.Weights, &result, &validationCache, 0)
	addCriteria(other.Weights, &result, &validationCache, originalLen)
	_sortWeightsMutate(&result)
	return &owaParams{Weights: &result}
}

func (o *owaParams) find(criterion *model.Criterion) *model.WeightedCriterion {
	for _, c := range *o.Weights {
		if c.Criterion.Id == criterion.Id {
			return &c
		}
	}
	panic(fmt.Errorf("criterion '%s' not found in criteria %v", criterion.Id, *o.Weights))
}

func addCriteria(toAdd, result *model.WeightedCriteria, validationCache *map[string]bool, offset int) {
	for i, w := range *toAdd {
		if _, ok := (*validationCache)[w.Id]; ok {
			criterionAlreadyExist(&w.Criterion, result)
		}
		(*result)[i+offset] = w
		(*validationCache)[w.Id] = true
	}
}

func criterionAlreadyExist(w *model.Criterion, result *model.WeightedCriteria) {
	panic(fmt.Errorf("criterion '%s' already exist in result %v", w, *result))
}

func _sortWeightsMutate(weights *model.WeightedCriteria) {
	sort.SliceStable(*weights, func(i, j int) bool {
		return (*weights)[i].Weight < (*weights)[j].Weight
	})
}

func validateSameCriteriaAndWeightsCount(alternative *model.AlternativeWithCriteria, weights *model.WeightedCriteria) {
	alternativeCriteriaCount := len(alternative.Criteria)
	weightsCount := len(*weights)
	if alternativeCriteriaCount != weightsCount {
		panic(fmt.Errorf("criteria and Weights must have the same length, got %d and %d", alternativeCriteriaCount, weightsCount))
	}
}
