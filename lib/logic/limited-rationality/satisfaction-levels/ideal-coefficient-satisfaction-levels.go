package satisfaction_levels

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

//go:generate easytags $GOFILE json:camel

type IdealCoefficientSatisfactionLevels struct {
	Coefficient          float64 `json:"coefficient"`
	MaxValue             float64 `json:"maxValue"`
	MinValue             float64 `json:"minValue"`
	currentValue         float64
	criteria             model.Criteria
	criteriaValuesRanges []utils.ValueRange
	manager              CoefficientManager
}

type CoefficientManager interface {
	Validate(params *IdealCoefficientSatisfactionLevels)
	UpdateValue(current, coefficient float64) float64
	InitialValue(params *IdealCoefficientSatisfactionLevels) float64
}

func (s *IdealCoefficientSatisfactionLevels) Initialize(dmp *model.DecisionMakingParams) {
	s.manager.Validate(s)
	s.criteria = dmp.Criteria
	s.criteriaValuesRanges = make([]utils.ValueRange, len(dmp.Criteria))
	alternatives := dmp.AllAlternatives()
	for i, c := range s.criteria {
		s.criteriaValuesRanges[i] = *model.CriteriaValuesRange(&alternatives, &c)
	}
	s.currentValue = s.manager.InitialValue(s)
}

func (s *IdealCoefficientSatisfactionLevels) HasNext() bool {
	return s.currentValue > s.MinValue
}

func (s *IdealCoefficientSatisfactionLevels) Next() model.Weights {
	weights := make(model.Weights, len(s.criteria))
	for i, c := range s.criteria {
		valRange := s.criteriaValuesRanges[i]
		delta := valRange.Diff() * s.currentValue
		if c.Multiplier() > 0 {
			weights[c.Id] = valRange.Min + delta
		} else {
			weights[c.Id] = valRange.Max - delta
		}
	}
	s.currentValue = s.manager.UpdateValue(s.currentValue, s.Coefficient)
	return weights
}

type IdealCoefficientSatisfactionLevelsSource struct {
	id                 string
	coefficientManager CoefficientManager
}

func (s *IdealCoefficientSatisfactionLevelsSource) OnCriterionAdded(criterion *model.Criterion, previousRankedCriteria *model.Criteria, params SatisfactionLevels, generator utils.ValueGenerator) ParamsAddition {
	return nil
}

func (s *IdealCoefficientSatisfactionLevelsSource) OnCriteriaRemoved(leftCriteria *model.Criteria, params SatisfactionLevels) SatisfactionLevels {
	return params
}

func (s *IdealCoefficientSatisfactionLevelsSource) Merge(params SatisfactionLevels, addition ParamsAddition) SatisfactionLevels {
	return params
}

func (s *IdealCoefficientSatisfactionLevelsSource) Identifier() string {
	return s.id
}

func (s *IdealCoefficientSatisfactionLevelsSource) BlankParams() SatisfactionLevels {
	return &IdealCoefficientSatisfactionLevels{
		manager: s.coefficientManager,
	}
}
