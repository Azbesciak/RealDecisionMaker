package weighted_sum

import (
	"fmt"
	. "github.com/Azbesciak/RealDecisionMaker/lib/model"
)

type WeightedSumPreferenceFunc struct {
}

type weightedSumParams struct {
	weightedCriteria *[]WeightedCriterion
}

func (w *WeightedSumPreferenceFunc) ParseParams(dm *DecisionMaker) interface{} {
	weights := ExtractWeights(dm)
	weightedCriteria := make([]WeightedCriterion, len(dm.Criteria))
	for i, c := range dm.Criteria {
		value, ok := weights[c.Id]
		if !ok {
			panic(fmt.Errorf("value for criterion '%s' not found in weights '%v'", c.Id, weights))
		}
		weightedCriteria[i] = WeightedCriterion{
			Criterion: c,
			Weight:    value,
		}
	}
	return weightedSumParams{weightedCriteria: &weightedCriteria}
}

func (w *WeightedSumPreferenceFunc) Identifier() string {
	return "weightedSum"
}

func (w *WeightedSumPreferenceFunc) MethodParameters() interface{} {
	return WeightsParamOnly()
}

func (w *WeightedSumPreferenceFunc) Evaluate(dm *DecisionMaker) *AlternativesRanking {
	params := w.ParseParams(dm).(weightedSumParams)
	prefFunc := func(alternative *AlternativeWithCriteria) *AlternativeResult {
		return WeightedSum(*alternative, *params.weightedCriteria)
	}
	return Rank(dm, prefFunc)
}

func WeightedSum(alternative AlternativeWithCriteria, criteria []WeightedCriterion) *AlternativeResult {
	var total Weight = 0
	for _, criterion := range criteria {
		total += alternative.CriterionValue(&criterion.Criterion)
	}
	return &AlternativeResult{Alternative: alternative, Value: total}
}
