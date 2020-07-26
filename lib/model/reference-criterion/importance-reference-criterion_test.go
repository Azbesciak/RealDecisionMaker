package reference_criterion

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"testing"
)

func TestImportanceRatioReferenceCriterionProvider(t *testing.T) {
	criteria := &model.WeightedCriteria{
		{Criterion: model.Criterion{Id: "2"}, Weight: 5},
		{Criterion: model.Criterion{Id: "1"}, Weight: 10},
		{Criterion: model.Criterion{Id: "3"}, Weight: 15},
	}
	checkImportanceProvision(t, -1.0, criteria, "2")
	checkImportanceProvision(t, 0, criteria, "2")
	checkImportanceProvision(t, 0.16, criteria, "2")
	checkImportanceProvision(t, 0.3, criteria, "1")
	checkImportanceProvision(t, 0.5, criteria, "1")
	checkImportanceProvision(t, 0.51, criteria, "3")
	checkImportanceProvision(t, 0.7, criteria, "3")
	checkImportanceProvision(t, 1, criteria, "3")
	checkImportanceProvision(t, 2, criteria, "3")
}

func checkImportanceProvision(t *testing.T, importance float64, criteria *model.WeightedCriteria, expectedId string) {
	manager := ImportanceRatioReferenceCriterionManager{}
	provider := manager.NewProvider().(*ImportanceRatioReferenceCriterionProvider)
	provider.NewCriterionImportance = importance
	actual := provider.Provide(criteria)
	if actual.Id != expectedId {
		t.Errorf("expected criterion by importance %f - '%s', got '%s'", importance, expectedId, actual.Id)
	}

}
