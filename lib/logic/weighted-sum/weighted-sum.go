package weighted_sum

import (
	"fmt"
	. "github.com/Azbesciak/RealDecisionMaker/lib/model"
)

type WeightedSumPreferenceFunc struct {
}

func (w *WeightedSumPreferenceFunc) Identifier() string {
	return "weightedSum"
}

func (w *WeightedSumPreferenceFunc) Evaluate(dm *DecisionMaker) *AlternativesRanking {
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
	prefFunc := func(alternative *AlternativeWithCriteria) *AlternativeResult {
		return WeightedSum(*alternative, weightedCriteria)
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
