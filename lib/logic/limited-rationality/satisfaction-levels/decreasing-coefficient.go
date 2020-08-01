package satisfaction_levels

import (
	"fmt"
	"math"
)

type DecreasingCoefficientManager struct {
	updateCoefficient func(current, coefficient float64) float64
}

func (d *DecreasingCoefficientManager) Validate(params *IdealCoefficientSatisfactionLevels) {
	if params.Coefficient <= 0 || params.Coefficient >= 1 {
		panic(fmt.Errorf("satisfaction coefficient degradation level must be in range (0,1), got %f", params.Coefficient))
	}
	if params.MinValue <= 0 {
		panic(fmt.Errorf("minimum satisfaction coefficient level must be positive value, got %f", params.MinValue))
	}
	if params.MaxValue > 1 {
		panic(fmt.Errorf("max satisfaction coefficient level cannot be greater than 1, got %f", params.MaxValue))
	}
}

func (d *DecreasingCoefficientManager) UpdateValue(current, coefficient float64) float64 {
	return d.updateCoefficient(current, coefficient)
}

func (d *DecreasingCoefficientManager) InitialValue(params *IdealCoefficientSatisfactionLevels) float64 {
	return params.MaxValue
}

func (d *DecreasingCoefficientManager) HasNext(params *IdealCoefficientSatisfactionLevels) bool {
	return params.currentValue > params.MinValue
}

const IdealDecreasingMul = "idealMultipliedCoefficient"

var IdealDecreasingMulCoefficientSatisfaction = IdealCoefficientSatisfactionLevelsSource{
	id: IdealDecreasingMul,
	coefficientManager: &DecreasingCoefficientManager{
		updateCoefficient: func(current, coefficient float64) float64 {
			return current * coefficient
		},
	},
}

const IdealSubtractive = "idealSubtractiveCoefficient"

var IdealSubtrCoefficientSatisfaction = IdealCoefficientSatisfactionLevelsSource{
	id: IdealSubtractive,
	coefficientManager: &DecreasingCoefficientManager{
		updateCoefficient: func(current, coefficient float64) float64 {
			return math.Max(current-coefficient, 0)
		},
	},
}
