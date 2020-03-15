package model

import (
	"fmt"
	"strings"
)

//go:generate easytags $GOFILE json:camel

type DecisionMaker struct {
	PreferenceFunction string                    `json:"preferenceFunction"`
	Heuristics         HeuristicsParams          `json:"heuristics"`
	KnownAlternatives  []AlternativeWithCriteria `json:"knownAlternatives"`
	ChoseToMake        []Alternative             `json:"choseToMake"`
	Criteria           Criteria                  `json:"criteria"`
	MethodParameters   MethodParameters          `json:"methodParameters"`
}

type MethodParameters = map[string]interface{}

func (dm *DecisionMaker) Alternative(id string) AlternativeWithCriteria {
	for _, a := range dm.KnownAlternatives {
		if a.Id == id {
			return a
		}
	}
	panic(fmt.Errorf("alternative '%s' is unknown", id))
}

func (dm *DecisionMaker) AlternativesToConsider() *[]AlternativeWithCriteria {
	results := make([]AlternativeWithCriteria, len(dm.ChoseToMake))
	for i, r := range dm.ChoseToMake {
		results[i] = dm.Alternative(r)
	}
	return &results
}

type DecisionMakerChoice struct {
	Result     AlternativesRanking `json:"result"`
	Heuristics HeuristicsParams    `json:"heuristics"`
}

func (dm *DecisionMaker) MakeDecision(preferenceFunctions PreferenceFunctions, availableHeuristics *HeuristicsMap) *DecisionMakerChoice {
	if IsStringBlank(&dm.PreferenceFunction) {
		panic(fmt.Errorf("preference function must not be empty"))
	}
	fun := FetchPreferenceFunction(preferenceFunctions, dm.PreferenceFunction)
	preferenceFunction := (*fun).(PreferenceFunction)
	chosenHeuristics := ChooseHeuristics(availableHeuristics, &dm.Heuristics)
	processedDm, heuristicProps := dm.processHeuristics(chosenHeuristics)
	res := preferenceFunction.Evaluate(processedDm)
	return &DecisionMakerChoice{*res, *heuristicProps}
}

func (dm *DecisionMaker) processHeuristics(heuristics *HeuristicsWithProps) (*DecisionMaker, *HeuristicsParams) {
	result := make(HeuristicsParams, len(*heuristics))
	tempDM := dm
	for i, h := range *heuristics {
		res := (*h.Heuristic).Apply(tempDM, &h.Props.Props)
		tempDM = res.Dm
		result[i] = *UpdateHeuristicProps(h.Props, res.Props)
	}
	return tempDM, &result
}

func IsStringBlank(str *string) bool {
	return len(strings.TrimSpace(*str)) == 0
}
