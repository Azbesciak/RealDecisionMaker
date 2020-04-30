package satisfaction

import (
	"fmt"
	"math"
)

type IncreasingCoefficientManager struct {
	updateCoefficient func(current, coefficient float64) float64
}

func (i *IncreasingCoefficientManager) Validate(params *idealCoefficientSatisfactionLevels) {
	if params.Coefficient <= 0 || params.Coefficient >= 1 {
		panic(fmt.Errorf("satisfaction Coefficient degradation level must be in range (0,1), got %f", params.Coefficient))
	}
	if params.MinValue < 0 {
		panic(fmt.Errorf("minimum satisfaction Coefficient level must be positive value, got %f", params.MinValue))
	}
	if params.MaxValue >= 1 {
		panic(fmt.Errorf("max satisfaction Coefficient level cannot be greater than 1, got %f", params.MinValue))
	}
}

func (i *IncreasingCoefficientManager) UpdateValue(current, coefficient float64) float64 {
	return current * (1 + coefficient)
}

func (i *IncreasingCoefficientManager) InitialValue(params *idealCoefficientSatisfactionLevels) float64 {
	return params.MinValue
}

var IdealIncreasingMulCoefficientSatisfaction = IdealCoefficientSatisfactionLevelsSource{
	name: "idealMulCoefficient",
	coefficientManager: &IncreasingCoefficientManager{
		updateCoefficient: func(current, coefficient float64) float64 {
			return current * coefficient
		},
	},
}

var IdealAdditiveCoefficientSatisfaction = IdealCoefficientSatisfactionLevelsSource{
	name: "idealSubtCoefficient",
	coefficientManager: &DecreasingCoefficientManager{
		updateCoefficient: func(current, coefficient float64) float64 {
			return math.Max(current+coefficient, 1)
		},
	},
}
