package model

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/utils"
)

type AlternativeWeightFunction func(alternative *AlternativeWithCriteria) *AlternativeResult

func Rank(dm *DecisionMaker, pref AlternativeWeightFunction) *AlternativesRanking {
	results := make(AlternativeResults, len(dm.ChoseToMake))
	for i, r := range dm.ChoseToMake {
		var alternative = dm.Alternative(r)
		results[i] = *pref(&alternative)
	}
	return results.Ranking()
}

func ExtractWeights(dm *DecisionMaker) Weights {
	weights, ok := dm.MethodParameters["weights"]
	if !ok {
		panic(fmt.Errorf("weights not found"))
	}
	weightsParsed := make(Weights)
	utils.DecodeToStruct(weights, &weightsParsed)
	return weightsParsed
}
