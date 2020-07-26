package reference_criterion

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"math"
	"testing"
)

func TestRandomWeightedReferenceCriterionProvider(t *testing.T) {
	criteria := model.WeightedCriteria{
		{Criterion: model.Criterion{Id: "1"}, Weight: 10},
		{Criterion: model.Criterion{Id: "2"}, Weight: 20},
		{Criterion: model.Criterion{Id: "3"}, Weight: 30},
	}
	checkRandomCriterionCorrect(t, 0, criteria, model.Criterion{Id: "1"})
	checkRandomCriterionCorrect(t, 10, criteria, model.Criterion{Id: "1"})
	checkRandomCriterionCorrect(t, 20, criteria, model.Criterion{Id: "1"})
	checkRandomCriterionCorrect(t, 50, criteria, model.Criterion{Id: "1"})
	criterion2MinRequirement := int64(math.Floor(100.0*(6.0/11.0))) + 1
	checkRandomCriterionCorrect(t, criterion2MinRequirement, criteria, model.Criterion{Id: "2"})
	checkRandomCriterionCorrect(t, 80, criteria, model.Criterion{Id: "2"})
	checkRandomCriterionCorrect(t, 85, criteria, model.Criterion{Id: "3"})
}

func checkRandomCriterionCorrect(t *testing.T, seed int64, criteria model.WeightedCriteria, expected model.Criterion) {
	manager := RandomWeightedReferenceCriterionManager{
		RandomFactory: func(seed int64) utils.ValueGenerator {
			return func() float64 {
				return float64(seed) / 100.0
			}
		},
	}
	provider := manager.NewProvider().(*RandomWeightedReferenceCriterionProvider)
	provider.NewCriterionRandomSeed = seed
	actual := provider.Provide(&criteria)
	if actual.Id != expected.Id {
		t.Errorf("expected random criterion for seed %v as '%s', got '%s'", seed, expected.Id, actual.Id)
	}
}
