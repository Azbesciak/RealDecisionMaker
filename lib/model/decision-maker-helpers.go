package model

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"github.com/alecthomas/jsonschema"
)

type AlternativeWeightFunction func(alternative *AlternativeWithCriteria) *AlternativeResult

func Rank(dmp *DecisionMakingParams, pref AlternativeWeightFunction) *AlternativesRanking {
	results := make(AlternativeResults, len(dmp.ConsideredAlternatives))
	for i, alternative := range dmp.ConsideredAlternatives {
		results[i] = *pref(&alternative)
	}
	return results.Ranking()
}

const WeightsParam = "weights"

type WeightType struct {
	Weights `json:"weights"`
}

func WeightsParamOnly() interface{} {
	return WeightType{}
}

func ExtractWeights(dm *DecisionMaker) Weights {
	weights, ok := dm.MethodParameters[WeightsParam]
	if !ok {
		panic(fmt.Errorf("weights not found"))
	}
	weightsParsed := make(Weights)
	utils.DecodeToStruct(weights, &weightsParsed)
	return weightsParsed
}

func FetchPreferenceFunction(preferenceFunctions PreferenceFunctions, function string) *utils.Identifiable {
	preferenceFunMap := utils.AsMap(preferenceFunctions)
	fun, ok := (*preferenceFunMap)[function]
	if !ok {
		var keys []string
		for _, k := range preferenceFunctions.Functions {
			keys = append(keys, k.Identifier())
		}
		panic(fmt.Errorf("preference function '%s' not found, available are '%s'", function, keys))
	}
	return &fun
}

func FetchPreferenceFunctionsParameters(functions PreferenceFunctions) *map[string]interface{} {
	var functionsParameters = make(map[string]interface{}, functions.Len())
	reflector := jsonschema.Reflector{ExpandedStruct: true}
	for _, f := range functions.Functions {
		functionsParameters[f.Identifier()] = reflector.Reflect(f.MethodParameters())
	}
	return &functionsParameters
}
