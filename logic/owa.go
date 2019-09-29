package logic

import (
	"fmt"
	"sort"
)

func OWA(alternative AlternativeWithCriteria, weights []Weight) AlternativeResult {
	alternativeCriteria := len(alternative.Criteria)
	if alternativeCriteria != len(weights) {
		panic(fmt.Errorf("criteria and weights must have the same length, got %d and %d", alternativeCriteria, len(weights)))
	}
	tmpWeights := make([]Weight, len(weights))
	copy(tmpWeights, weights)
	sort.Float64s(tmpWeights)
	tmpCriteria := make([]Weight, alternativeCriteria)
	i := 0
	for _, c := range alternative.Criteria {
		tmpCriteria[i] = c
		i += 1
	}
	sort.Float64s(tmpCriteria)
	var total Weight = 0
	for i = 0; i < alternativeCriteria; i++ {
		total += tmpCriteria[i] * tmpWeights[i]
	}
	return AlternativeResult{Alternative: alternative, Value: total}
}
