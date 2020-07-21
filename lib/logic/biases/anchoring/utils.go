package anchoring

import "github.com/Azbesciak/RealDecisionMaker/lib/utils"

//go:generate easytags $GOFILE json:camel

type FunctionBase interface {
	utils.Identifiable
	// need to return pointer - it will be filled later
	BlankParams() FunctionParams
}

type FunctionParams = interface{}
type FunctionDefinition struct {
	Function string         `json:"function"`
	Params   FunctionParams `json:"params"`
}

func parseFuncParams(fun FunctionBase, parsedProps *FunctionDefinition) FunctionParams {
	funParams := fun.BlankParams()
	utils.DecodeToStruct(parsedProps.Params, funParams)
	return funParams
}
