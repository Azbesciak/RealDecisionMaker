package weighted_sum

import (
	"fmt"
	. "github.com/Azbesciak/RealDecisionMaker/lib/model"
)

const methodName = "weightedSum"

type WeightedSumPreferenceFunc struct {
}

type weightedSumParams struct {
	weightedCriteria *WeightedCriteria
}

func (p *weightedSumParams) Criterion(criterion string) WeightedCriterion {
	for _, c := range *p.weightedCriteria {
		if c.Id == criterion {
			return c
		}
	}
	panic(fmt.Errorf("criterion '%s' not found in weights %v", criterion, *p.weightedCriteria))
}

func (w *WeightedSumPreferenceFunc) ParseParams(dm *DecisionMaker) interface{} {
	weights := ExtractWeights(dm)
	weightedCriteria := dm.Criteria.ZipWithWeights(&weights)
	return weightedSumParams{weightedCriteria: weightedCriteria}
}

func (w *WeightedSumPreferenceFunc) Identifier() string {
	return methodName
}

func (w *WeightedSumPreferenceFunc) MethodParameters() interface{} {
	return WeightsParamOnly()
}

func (w *WeightedSumPreferenceFunc) Evaluate(dmp *DecisionMakingParams) *AlternativesRanking {
	params := dmp.MethodParameters.(weightedSumParams)
	prefFunc := func(alternative *AlternativeWithCriteria) *AlternativeResult {
		return WeightedSum(*alternative, *params.weightedCriteria)
	}
	return Rank(dmp, prefFunc)
}

func WeightedSum(alternative AlternativeWithCriteria, criteria WeightedCriteria) *AlternativeResult {
	var total Weight = 0
	for _, criterion := range criteria {
		total += alternative.CriterionValue(&criterion.Criterion)
	}
	return ValueAlternativeResult(&alternative, total)
}
