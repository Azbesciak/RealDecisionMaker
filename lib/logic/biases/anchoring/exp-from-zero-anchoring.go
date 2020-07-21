package anchoring

import "github.com/Azbesciak/RealDecisionMaker/lib/utils"

type ExpFromZeroAnchoringEvaluator struct {
}

const ExpFromZeroFunctionName = utils.ExpFromZeroFunctionName

func (e *ExpFromZeroAnchoringEvaluator) Identifier() string {
	return ExpFromZeroFunctionName
}

func (e *ExpFromZeroAnchoringEvaluator) BlankParams() FunctionParams {
	return &utils.ExpFromZeroFunction{}
}

func (e *ExpFromZeroAnchoringEvaluator) Evaluate(params FunctionParams, difference float64) float64 {
	p := params.(*utils.ExpFromZeroFunction)
	return p.Evaluate(difference)
}
