package model

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"strings"
)

//go:generate easytags $GOFILE json:camel

type DecisionMaker struct {
	PreferenceFunction string                    `json:"preferenceFunction"`
	Heuristics         HeuristicsParams          `json:"heuristics"`
	KnownAlternatives  []AlternativeWithCriteria `json:"knownAlternatives"`
	ChoseToMake        []Alternative             `json:"choseToMake"`
	Criteria           Criteria                  `json:"criteria"`
	MethodParameters   RawMethodParameters       `json:"methodParameters"`
}

type DecisionMakingParams struct {
	NotConsideredAlternatives []AlternativeWithCriteria
	ConsideredAlternatives    []AlternativeWithCriteria
	Criteria                  Criteria
	MethodParameters          interface{}
}

type RawMethodParameters = map[string]interface{}
type MethodParameters = interface{}

func (dm *DecisionMaker) Alternative(id Alternative) AlternativeWithCriteria {
	return FetchAlternative(&dm.KnownAlternatives, id)
}

func UpdateAlternatives(old *[]AlternativeWithCriteria, newOnes *[]AlternativeWithCriteria) *[]AlternativeWithCriteria {
	res := make([]AlternativeWithCriteria, len(*old))
	for i, a := range *old {
		res[i] = FetchAlternative(newOnes, a.Id)
	}
	return &res
}

func FetchAlternative(a *[]AlternativeWithCriteria, id Alternative) AlternativeWithCriteria {
	for _, a := range *a {
		if a.Id == id {
			return a
		}
	}
	panic(fmt.Errorf("alternative '%s' is unknown", id))
}

func (dm *DecisionMaker) AlternativesToConsider() *[]AlternativeWithCriteria {
	return FetchAlternatives(&dm.KnownAlternatives, &dm.ChoseToMake)
}

func FetchAlternatives(a *[]AlternativeWithCriteria, ids *[]Alternative) *[]AlternativeWithCriteria {
	results := make([]AlternativeWithCriteria, len(*ids))
	for i, id := range *ids {
		results[i] = FetchAlternative(a, id)
	}
	return &results
}

type DecisionMakerChoice struct {
	Result     AlternativesRanking `json:"result"`
	Heuristics HeuristicsParams    `json:"heuristics"`
}

func (dm *DecisionMaker) MakeDecision(
	preferenceFunctions PreferenceFunctions,
	heuristicListeners HeuristicListeners,
	availableHeuristics *HeuristicsMap,
) *DecisionMakerChoice {
	if IsStringBlank(&dm.PreferenceFunction) {
		panic(fmt.Errorf("preference function must not be empty"))
	}
	preferenceFunction := preferenceFunctions.Fetch(dm.PreferenceFunction)
	params := dm.prepareParams(preferenceFunction)
	chosenHeuristics := ChooseHeuristics(availableHeuristics, &dm.Heuristics)
	processedParams, heuristicProps := dm.processHeuristics(chosenHeuristics, params, &heuristicListeners)
	res := (*preferenceFunction).Evaluate(processedParams)
	return &DecisionMakerChoice{*res, *heuristicProps}
}

func (dm *DecisionMaker) prepareParams(preferenceFunction *PreferenceFunction) *DecisionMakingParams {
	return &DecisionMakingParams{
		NotConsideredAlternatives: *dm.NotConsideredAlternatives(),
		ConsideredAlternatives:    *dm.AlternativesToConsider(),
		Criteria:                  dm.Criteria,
		MethodParameters:          (*preferenceFunction).ParseParams(dm),
	}
}

func (dm *DecisionMaker) NotConsideredAlternatives() *[]AlternativeWithCriteria {
	var result []AlternativeWithCriteria
	for _, a := range dm.KnownAlternatives {
		if !utils.ContainsString(&dm.ChoseToMake, &a.Id) {
			result = append(result, a)
		}
	}
	return &result
}

func (dm *DecisionMaker) processHeuristics(
	heuristics *HeuristicsWithProps,
	params *DecisionMakingParams,
	listeners *HeuristicListeners,
) (*DecisionMakingParams, *HeuristicsParams) {
	heuristicsToProcessCount := len(*heuristics)
	result := make(HeuristicsParams, heuristicsToProcessCount)
	tempDM := params
	if heuristicsToProcessCount == 0 {
		return tempDM, &result
	}
	listener := listeners.Fetch(dm.PreferenceFunction)
	for i, h := range *heuristics {
		res := (*h.Heuristic).Apply(tempDM, &h.Props.Props, listener)
		tempDM = res.DMP
		result[i] = *UpdateHeuristicProps(h.Props, res.Props)
	}
	return tempDM, &result
}

func IsStringBlank(str *string) bool {
	return len(strings.TrimSpace(*str)) == 0
}
