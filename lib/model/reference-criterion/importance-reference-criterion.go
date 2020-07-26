package reference_criterion

import "github.com/Azbesciak/RealDecisionMaker/lib/model"

type ImportanceRatioReferenceCriterionProvider struct {
	NewCriterionImportance float64 `json:"newCriterionImportance"`
}

func (i *ImportanceRatioReferenceCriterionProvider) Provide(rankedCriteria *model.WeightedCriteria) *model.Criterion {
	total := 0.0
	for _, c := range *rankedCriteria {
		total += c.Weight
	}
	expectedWeight := i.NewCriterionImportance * total
	return FindCriterionInRange(rankedCriteria, expectedWeight)
}

type ImportanceRatioReferenceCriterionManager struct {
}

const ImportanceRatioReferenceCriterion = "importanceRatio"

func (i *ImportanceRatioReferenceCriterionManager) Identifier() string {
	return ImportanceRatioReferenceCriterion
}

func (i *ImportanceRatioReferenceCriterionManager) NewProvider() ReferenceCriterionProvider {
	return &ImportanceRatioReferenceCriterionProvider{}
}
