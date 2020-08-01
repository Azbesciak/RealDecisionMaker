package satisfaction_levels

import (
	"fmt"
	"math"
)

type IncreasingCoefficientManager struct {
	updateCoefficient func(current, coefficient float64) float64
}

func (i *IncreasingCoefficientManager) Validate(params *IdealCoefficientSatisfactionLevels) {
	if params.Coefficient <= 0 || params.Coefficient >= 1 {
		panic(fmt.Errorf("satisfaction coefficient degradation level must be in range (0,1), got %f", params.Coefficient))
	}
	if params.MinValue < 0 {
		panic(fmt.Errorf("minimum satisfaction coefficient level must be non-negative value, got %f", params.MinValue))
	}
	if params.MaxValue >= 1 {
		panic(fmt.Errorf("max satisfaction coefficient level must be lower than 1, got %f", params.MaxValue))
	}
}

func (i *IncreasingCoefficientManager) UpdateValue(current, coefficient float64) float64 {
	return i.updateCoefficient(current, coefficient)
}

func (i *IncreasingCoefficientManager) InitialValue(params *IdealCoefficientSatisfactionLevels) float64 {
	return params.MinValue
}

func (i *IncreasingCoefficientManager) HasNext(params *IdealCoefficientSatisfactionLevels) bool {
	return params.currentValue < params.MaxValue
}

const IdealIncreasingMul = "idealMultipliedCoefficient"

var IdealIncreasingMulCoefficientSatisfaction = IdealCoefficientSatisfactionLevelsSource{
	id: IdealIncreasingMul,
	coefficientManager: &IncreasingCoefficientManager{
		updateCoefficient: func(current, coefficient float64) float64 {
			return math.Min((1+current)*(1+coefficient)-1, 1)
		},
	},
}

const IdealAdditive = "idealAdditiveCoefficient"

var IdealAdditiveCoefficientSatisfaction = IdealCoefficientSatisfactionLevelsSource{
	id: IdealAdditive,
	coefficientManager: &IncreasingCoefficientManager{
		updateCoefficient: func(current, coefficient float64) float64 {
			return math.Min(current+coefficient, 1)
		},
	},
}
