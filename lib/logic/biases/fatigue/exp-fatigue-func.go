package fatigue

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

//go:generate easytags $GOFILE json:camel

type ExpFatigueParams struct {
	Alpha       float64 `json:"alpha"`
	Multiplier  float64 `json:"multiplier"`
	QueryNumber int64   `json:"queryNumber"`
}

type ExponentialFromZeroFatigue struct {
}

const FatExpFromZero = utils.ExpFromZeroFunctionName

func (e *ExponentialFromZeroFatigue) Name() string {
	return FatExpFromZero
}

func (e *ExponentialFromZeroFatigue) BlankParams() interface{} {
	return &ExpFatigueParams{}
}

func (e *ExponentialFromZeroFatigue) Evaluate(params interface{}) float64 {
	p := params.(*ExpFatigueParams)
	function := utils.ExpFromZeroFunction{
		Alpha:      p.Alpha,
		Multiplier: p.Multiplier,
	}
	return function.Evaluate(float64(p.QueryNumber))
}
