package anchoring

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type LinearAnchoringEvaluator struct {
}

const LinearFunctionName = utils.LinearFunctionName

func (i *LinearAnchoringEvaluator) Identifier() string {
	return LinearFunctionName
}

func (i *LinearAnchoringEvaluator) BlankParams() FunctionParams {
	return &utils.LinearFunctionParameters{}
}

func (i *LinearAnchoringEvaluator) Evaluate(params FunctionParams, difference float64) float64 {
	p := params.(*utils.LinearFunctionParameters)
	res, _ := p.Evaluate(difference)
	return res
}
