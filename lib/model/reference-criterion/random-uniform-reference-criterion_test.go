package reference_criterion

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

func TestRandomUniformReferenceCriterionProvider(t *testing.T) {
	criteria := model.WeightedCriteria{
		{Criterion: model.Criterion{Id: "1"}, Weight: 10},
		{Criterion: model.Criterion{Id: "2"}, Weight: 20},
		{Criterion: model.Criterion{Id: "3"}, Weight: 30},
	}
	checkRandomUniformCriterionCorrect(t, 0, criteria, "1")
	checkRandomUniformCriterionCorrect(t, 10, criteria, "1")
	checkRandomUniformCriterionCorrect(t, 30, criteria, "1")
	checkRandomUniformCriterionCorrect(t, 40, criteria, "2")
	checkRandomUniformCriterionCorrect(t, 60, criteria, "2")
	checkRandomUniformCriterionCorrect(t, 67, criteria, "3")
	checkRandomUniformCriterionCorrect(t, 85, criteria, "3")
}

func checkRandomUniformCriterionCorrect(t *testing.T, seed int64, criteria model.WeightedCriteria, expectedId string) {
	manager := RandomUniformReferenceCriterionManager{
		RandomFactory: func(seed int64) utils.ValueGenerator {
			return func() float64 {
				return float64(seed) / 100.0
			}
		},
	}
	provider := manager.NewProvider().(*RandomUniformReferenceCriterionProvider)
	provider.NewCriterionRandomSeed = seed
	actual := provider.Provide(&criteria)
	if actual.Id != expectedId {
		t.Errorf("expected random uniform criterion for seed %v as '%s', got '%s'", seed, expectedId, actual.Id)
	}
}
