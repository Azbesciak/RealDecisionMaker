package logic

import "fmt"

func WeightedSum(alternative AlternativeWithCriteria, criteria []WeightedCriterion) AlternativeResult {
	var total Weight = 0
	for _, criterion := range criteria {
		var value, ok = alternative.Criteria[criterion.Id]
		if !ok {
			panic(fmt.Sprintf("criterion '%s' not found in criteria", criterion.Id))
		}
		total += value * Weight(criterion.multiplier())
	}
	return AlternativeResult{Alternative: alternative, Value: total}
}
