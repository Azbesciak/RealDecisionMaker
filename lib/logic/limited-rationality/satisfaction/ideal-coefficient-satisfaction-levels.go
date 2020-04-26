package satisfaction

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"math"
)

//go:generate easytags $GOFILE json:camel

type idealCoefficientSatisfactionLevels struct {
	Coefficient          float64 `json:"coefficient"`
	MinValue             float64 `json:"minValue"`
	currentValue         float64
	criteria             model.Criteria
	criteriaValuesRanges []utils.ValueRange
	coefficientManager   CoefficientManager
}

type CoefficientManager = func(current, coefficient float64) float64

func (s *idealCoefficientSatisfactionLevels) Initialize(dmp *model.DecisionMakingParams) {
	if s.Coefficient <= 0 || s.Coefficient >= 1 {
		panic(fmt.Errorf("satisfaction Coefficient degradation level must be in range (0,1), got %f", s.Coefficient))
	}
	if s.MinValue <= 0 {
		panic(fmt.Errorf("minimum satisfaction Coefficient level must be positive value, got %f", s.MinValue))
	}
	s.criteria = dmp.Criteria
	s.criteriaValuesRanges = make([]utils.ValueRange, len(dmp.Criteria))
	alternatives := dmp.AllAlternatives()
	for i, c := range s.criteria {
		s.criteriaValuesRanges[i] = *model.CriteriaValuesRange(&alternatives, &c)
	}
	s.currentValue = 1
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
	s.currentValue = s.coefficientManager(s.currentValue, s.Coefficient)
	return weights
}

type IdealCoefficientSatisfactionLevelsSource struct {
	name               string
	coefficientManager CoefficientManager
}

var IdealMulCoefficientSatisfaction = IdealCoefficientSatisfactionLevelsSource{
	name: "idealMulCoefficient",
	coefficientManager: func(current, coefficient float64) float64 {
		return current * coefficient
	},
}

var IdealSubtrCoefficientSatisfaction = IdealCoefficientSatisfactionLevelsSource{
	name: "idealSubtCoefficient",
	coefficientManager: func(current, coefficient float64) float64 {
		return math.Max(current-coefficient, 0)
	},
}

func (s *IdealCoefficientSatisfactionLevelsSource) Name() string {
	return s.name
}

func (s *IdealCoefficientSatisfactionLevelsSource) BlankParams() SatisfactionLevels {
	return &idealCoefficientSatisfactionLevels{
		coefficientManager: s.coefficientManager,
	}
}
