package model

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/utils"
	"reflect"
	"strings"
)

type DecisionMaker struct {
	PreferenceFunction string
	State              DecisionMakerState
	KnownAlternatives  []AlternativeWithCriteria
	ChoseToMake        []Alternative
	Criteria
	MethodParameters
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
	Focus   int
	Fatigue int
}

type DecisionMakerChoice struct {
	Result AlternativesRanking
	State  DecisionMakerState
}

type PreferenceFunction interface {
	utils.Identifiable
	Evaluate(dm *DecisionMaker) *AlternativesRanking
}

type PreferenceFunctions struct {
	Functions []PreferenceFunction
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
	preferenceFunMap := utils.AsMap(preferenceFunctions)
	fun, ok := (*preferenceFunMap)[dm.PreferenceFunction]
	if !ok {
		keys := reflect.ValueOf(preferenceFunctions).MapKeys()
		panic(fmt.Errorf("preference function '%s' not found, available are '%s'", dm.PreferenceFunction, keys))
	}
	res := fun.(PreferenceFunction).Evaluate(dm)
	return &DecisionMakerChoice{*res, dm.State}
}

func IsStringBlank(str *string) bool {
	return len(strings.TrimSpace(*str)) == 0
}
