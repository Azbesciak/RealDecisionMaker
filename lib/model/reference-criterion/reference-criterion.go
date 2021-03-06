package reference_criterion

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type ReferenceCriterionProvider interface {
	Provide(rankedCriteria *model.WeightedCriteria) *model.Criterion
}

type ReferenceCriterionFactory interface {
	utils.Identifiable
	// need to return pointer - it will be filled later
	NewProvider() ReferenceCriterionProvider
}

type ReferenceCriteriaManager struct {
	factories []ReferenceCriterionFactory
}

func NewReferenceCriteriaManager(factories []ReferenceCriterionFactory) *ReferenceCriteriaManager {
	return &ReferenceCriteriaManager{factories: factories}
}

type referenceParamsType struct {
	ReferenceCriterionType string `json:"referenceCriterionType"`
}

func (m *ReferenceCriteriaManager) ForParams(params *interface{}) ReferenceCriterionProvider {
	if len(m.factories) == 0 {
		panic(fmt.Errorf("no ReferenceCriterionFactory has been declared"))
	}
	referenceType := m.fetchFactoryTypeFromParams(params)
	factory := m.factory(&referenceType)
	provider := factory.NewProvider()
	utils.DecodeToStruct(*params, provider)
	return provider
}

func (m *ReferenceCriteriaManager) fetchFactoryTypeFromParams(params *interface{}) referenceParamsType {
	referenceType := referenceParamsType{}
	utils.DecodeToStruct(*params, &referenceType)
	if len(referenceType.ReferenceCriterionType) == 0 {
		referenceType.ReferenceCriterionType = m.factories[0].Identifier()
	}
	return referenceType
}

func (m *ReferenceCriteriaManager) factory(param *referenceParamsType) ReferenceCriterionFactory {
	for _, f := range m.factories {
		if f.Identifier() == param.ReferenceCriterionType {
			return f
		}
	}
	names := m.extractFactoriesNames()
	panic(fmt.Errorf("no reference criterion factory found for '%s' in %v", param.ReferenceCriterionType, names))
}

func (m *ReferenceCriteriaManager) extractFactoriesNames() []string {
	names := make([]string, len(m.factories))
	for i, f := range m.factories {
		names[i] = f.Identifier()
	}
	return names
}

func FindCriterionInRange(rankedCriteria *model.WeightedCriteria, expectedCumulatedWeight float64) *model.Criterion {
	currentWeight := 0.0
	for _, c := range *rankedCriteria {
		currentWeight += c.Weight
		if currentWeight >= expectedCumulatedWeight {
			return &c.Criterion
		}
	}
	return &(*rankedCriteria)[len(*rankedCriteria)-1].Criterion
}
