package logic

import (
	"../model"
	"fmt"
)

type WeightedSumPreferenceFunc struct {
}

func (w *WeightedSumPreferenceFunc) Identifier() string {
	return "weightedSum"
}

func (w *WeightedSumPreferenceFunc) Evaluate(dm *model.DecisionMaker) *model.AlternativesRanking {
	weightedCriteria := make([]model.WeightedCriterion, len(dm.Criteria))
	for i, c := range dm.Criteria {
		value, ok := dm.Weights[c.Id]
		if !ok {
			panic(fmt.Errorf("value for criterion '%s' not found in weights '%v'", c.Id, dm.Weights))
		}
		weightedCriteria[i] = model.WeightedCriterion{
			Criterion: c,
			Weight:    value,
		}
	}
	prefFunc := func(alternative *model.AlternativeWithCriteria) *model.AlternativeResult {
		return WeightedSum(*alternative, weightedCriteria)
	}
	return model.Rank(dm, prefFunc)
}

func WeightedSum(alternative model.AlternativeWithCriteria, criteria []model.WeightedCriterion) *model.AlternativeResult {
	var total model.Weight = 0
	for _, criterion := range criteria {
		var value, ok = alternative.Criteria[criterion.Id]
		if !ok {
			panic(fmt.Errorf("criterion '%s' not found in criteria", criterion.Id))
		}
		total += value * model.Weight(criterion.Multiplier())
	}
	return &model.AlternativeResult{Alternative: alternative, Value: total}
}
