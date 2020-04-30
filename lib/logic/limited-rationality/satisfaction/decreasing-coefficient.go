package satisfaction

import (
	"fmt"
	"math"
)

type DecreasingCoefficientManager struct {
	updateCoefficient func(current, coefficient float64) float64
}

func (d *DecreasingCoefficientManager) Validate(params *idealCoefficientSatisfactionLevels) {
	if params.Coefficient <= 0 || params.Coefficient >= 1 {
		panic(fmt.Errorf("satisfaction Coefficient degradation level must be in range (0,1), got %f", params.Coefficient))
	}
	if params.MinValue <= 0 {
		panic(fmt.Errorf("minimum satisfaction Coefficient level must be positive value, got %f", params.MinValue))
	}
	if params.MaxValue > 1 {
		panic(fmt.Errorf("max satisfaction Coefficient level cannot be greater than 1, got %f", params.MinValue))
	}
}

func (d *DecreasingCoefficientManager) UpdateValue(current, coefficient float64) float64 {
	return d.updateCoefficient(current, coefficient)
}

func (d *DecreasingCoefficientManager) InitialValue(params *idealCoefficientSatisfactionLevels) float64 {
	return params.MaxValue
}

var IdealDecreasingMulCoefficientSatisfaction = IdealCoefficientSatisfactionLevelsSource{
	name: "idealMulCoefficient",
	coefficientManager: &DecreasingCoefficientManager{
		updateCoefficient: func(current, coefficient float64) float64 {
			return current * coefficient
		},
	},
}

var IdealSubtrCoefficientSatisfaction = IdealCoefficientSatisfactionLevelsSource{
	name: "idealSubtCoefficient",
	coefficientManager: &IncreasingCoefficientManager{
		updateCoefficient: func(current, coefficient float64) float64 {
			return math.Max(current-coefficient, 0)
		},
	},
}
