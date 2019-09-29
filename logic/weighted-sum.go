package logic

import (
	"../model"
	"fmt"
)

func WeightedSum(alternative model.AlternativeWithCriteria, criteria []model.WeightedCriterion) model.AlternativeResult {
	var total model.Weight = 0
	for _, criterion := range criteria {
		var value, ok = alternative.Criteria[criterion.Id]
		if !ok {
			panic(fmt.Errorf("criterion '%s' not found in criteria", criterion.Id))
		}
		total += value * model.Weight(criterion.Multiplier())
	}
	return model.AlternativeResult{Alternative: alternative, Value: total}
}
