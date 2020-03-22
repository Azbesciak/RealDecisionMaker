package model

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"github.com/alecthomas/jsonschema"
)

type PreferenceFunctions struct {
	Functions []PreferenceFunction `json:"functions"`
}

type PreferenceFunction interface {
	utils.Identifiable
	MethodParameters() interface{}
	ParseParams(dm *DecisionMaker) interface{}
	Evaluate(dmp *DecisionMakingParams) *AlternativesRanking
}

func (pf *PreferenceFunctions) Get(index int) utils.Identifiable {
	return pf.Functions[index]
}

func (pf *PreferenceFunctions) Len() int {
	return len(pf.Functions)
}

func (pf *PreferenceFunctions) Fetch(function string) *PreferenceFunction {
	preferenceFunMap := utils.AsMap(pf)
	fun, ok := (*preferenceFunMap)[function]
	if !ok {
		var keys []string
		for _, k := range pf.Functions {
			keys = append(keys, k.Identifier())
		}
		panic(fmt.Errorf("preference function '%s' not found, available are '%s'", function, keys))
	}
	preferenceFunction := fun.(PreferenceFunction)
	return &preferenceFunction
}

func (pf *PreferenceFunctions) FetchParameters() *map[string]interface{} {
	var functionsParameters = make(map[string]interface{}, pf.Len())
	reflector := jsonschema.Reflector{ExpandedStruct: true}
	for _, f := range pf.Functions {
		functionsParameters[f.Identifier()] = reflector.Reflect(f.MethodParameters())
	}
	return &functionsParameters
}
