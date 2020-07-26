package reference_criterion

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"math"
)

type RandomWeightedReferenceCriterionProvider struct {
	NewCriterionRandomSeed int64 `json:"newCriterionRandomSeed"`
	generator              utils.SeededValueGenerator
}

func (i *RandomWeightedReferenceCriterionProvider) Provide(rankedCriteria *model.WeightedCriteria) *model.Criterion {
	minVal := math.MaxFloat64
	for _, c := range *rankedCriteria {
		if c.Weight < minVal {
			minVal = c.Weight
		}
	}
	mappedWeights := make(model.WeightedCriteria, len(*rankedCriteria))
	total := 0.0
	for i, c := range *rankedCriteria {
		weight := minVal / c.Weight
		total += weight
		mappedWeights[i] = model.WeightedCriterion{
			Criterion: c.Criterion,
			Weight:    weight,
		}
	}
	generator := i.generator(i.NewCriterionRandomSeed)
	expectedWeight := generator() * total
	return FindCriterionInRange(&mappedWeights, expectedWeight)
}

type RandomWeightedReferenceCriterionManager struct {
	RandomFactory utils.SeededValueGenerator
}

const RandomWeightedReferenceCriterion = "randomWeighted"

func (i *RandomWeightedReferenceCriterionManager) Identifier() string {
	return RandomWeightedReferenceCriterion
}

func (i *RandomWeightedReferenceCriterionManager) NewProvider() ReferenceCriterionProvider {
	return &RandomWeightedReferenceCriterionProvider{
		generator: i.RandomFactory,
	}
}
