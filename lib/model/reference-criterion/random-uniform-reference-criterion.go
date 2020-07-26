package reference_criterion

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"math"
)

type RandomUniformReferenceCriterionProvider struct {
	NewCriterionRandomSeed int64 `json:"newCriterionRandomSeed"`
	generator              utils.SeededValueGenerator
}

func (i *RandomUniformReferenceCriterionProvider) Provide(rankedCriteria *model.WeightedCriteria) *model.Criterion {
	criteriaCount := len(*rankedCriteria)
	generator := i.generator(i.NewCriterionRandomSeed)
	expectedIndex := int(math.Floor(generator() * float64(criteriaCount)))
	return &(*rankedCriteria)[expectedIndex].Criterion
}

type RandomUniformReferenceCriterionManager struct {
	RandomFactory utils.SeededValueGenerator
}

const RandomUniformReferenceCriterion = "randomUniform"

func (i *RandomUniformReferenceCriterionManager) Identifier() string {
	return RandomUniformReferenceCriterion
}

func (i *RandomUniformReferenceCriterionManager) NewProvider() ReferenceCriterionProvider {
	return &RandomUniformReferenceCriterionProvider{
		generator: i.RandomFactory,
	}
}
