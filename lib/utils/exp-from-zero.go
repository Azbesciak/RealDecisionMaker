package utils

import (
	"fmt"
	"math"
)

const ExpFromZeroFunctionName = "expFromZero"

type ExpFromZeroFunction struct {
	Alpha      float64 `json:"alpha"`
	Multiplier float64 `json:"multiplier"`
}

func (e *ExpFromZeroFunction) String() string {
	return fmt.Sprintf("alpha:%v, multiplier:%v", e.Alpha, e.Multiplier)
}

func (e *ExpFromZeroFunction) Evaluate(value float64) float64 {
	return e.Multiplier*math.Exp(e.Alpha*value) - e.Multiplier
}
