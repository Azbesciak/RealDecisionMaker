package fatigue

import "math"

//go:generate easytags $GOFILE json:camel

type ExpFatigueParams struct {
	Alpha       float64 `json:"alpha"`
	Multiplier  float64 `json:"multiplier"`
	QueryNumber int64   `json:"queryNumber"`
}

type ExponentialFromZeroFatigue struct {
}

const FatExpFromZero = "expFromZero"

func (e *ExponentialFromZeroFatigue) Name() string {
	return FatExpFromZero
}

func (e *ExponentialFromZeroFatigue) BlankParams() interface{} {
	return &ExpFatigueParams{}
}

func (e *ExponentialFromZeroFatigue) Evaluate(params interface{}) float64 {
	p := params.(*ExpFatigueParams)
	return p.Multiplier*math.Exp(p.Alpha*float64(p.QueryNumber)) - p.Multiplier
}
