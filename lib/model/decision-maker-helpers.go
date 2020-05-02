package model

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
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
	Weights Weights `json:"weights"`
}

func SingleWeight(criterion *Criterion, value Weight) WeightType {
	return WeightType{Weights: Weights{criterion.Id: value}}
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
