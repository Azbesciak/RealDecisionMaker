package model

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"strings"
)

//go:generate easytags $GOFILE json:camel

type DecisionMaker struct {
	PreferenceFunction string                    `json:"preferenceFunction"`
	State              DecisionMakerState        `json:"state"`
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

type DecisionMakerState struct {
	Focus   int `json:"focus"`
	Fatigue int `json:"fatigue"`
}

type DecisionMakerChoice struct {
	Result AlternativesRanking `json:"result"`
	State  DecisionMakerState  `json:"state"`
}

type PreferenceFunction interface {
	utils.Identifiable
	MethodParameters() *map[string]interface{}
	Evaluate(dm *DecisionMaker) *AlternativesRanking
}

type PreferenceFunctions struct {
	Functions []PreferenceFunction `json:"functions"`
}

func (pf PreferenceFunctions) Get(index int) utils.Identifiable {
	return pf.Functions[index]
}

func (pf PreferenceFunctions) Len() int {
	return len(pf.Functions)
}

func (dm *DecisionMaker) MakeDecision(preferenceFunctions PreferenceFunctions) *DecisionMakerChoice {
	if IsStringBlank(&dm.PreferenceFunction) {
		panic(fmt.Errorf("preference function must not be empty"))
	}
	fun := FetchPreferenceFunction(preferenceFunctions, dm.PreferenceFunction)
	res := (*fun).(PreferenceFunction).Evaluate(dm)
	return &DecisionMakerChoice{*res, dm.State}
}

func IsStringBlank(str *string) bool {
	return len(strings.TrimSpace(*str)) == 0
}
