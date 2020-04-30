package satisfaction

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

//go:generate easytags $GOFILE json:camel

type idealCoefficientSatisfactionLevels struct {
	Coefficient          float64 `json:"coefficient"`
	MaxValue             float64 `json:"maxValue"`
	MinValue             float64 `json:"minValue"`
	currentValue         float64
	criteria             model.Criteria
	criteriaValuesRanges []utils.ValueRange
	manager              CoefficientManager
}

type CoefficientManager interface {
	Validate(params *idealCoefficientSatisfactionLevels)
	UpdateValue(current, coefficient float64) float64
	InitialValue(params *idealCoefficientSatisfactionLevels) float64
}

func (s *idealCoefficientSatisfactionLevels) Initialize(dmp *model.DecisionMakingParams) {
	s.manager.Validate(s)
	s.criteria = dmp.Criteria
	s.criteriaValuesRanges = make([]utils.ValueRange, len(dmp.Criteria))
	alternatives := dmp.AllAlternatives()
	for i, c := range s.criteria {
		s.criteriaValuesRanges[i] = *model.CriteriaValuesRange(&alternatives, &c)
	}
	s.currentValue = s.manager.InitialValue(s)
}

func (s *idealCoefficientSatisfactionLevels) HasNext() bool {
	return s.currentValue > s.MinValue
}

func (s *idealCoefficientSatisfactionLevels) Next() model.Weights {
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
	name               string
	coefficientManager CoefficientManager
}

func (s *IdealCoefficientSatisfactionLevelsSource) Name() string {
	return s.name
}

func (s *IdealCoefficientSatisfactionLevelsSource) BlankParams() SatisfactionLevels {
	return &idealCoefficientSatisfactionLevels{
		manager: s.coefficientManager,
	}
}
