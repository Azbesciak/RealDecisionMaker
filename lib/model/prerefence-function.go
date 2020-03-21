package model

import "github.com/Azbesciak/RealDecisionMaker/lib/utils"

type PreferenceFunctions struct {
	Functions []PreferenceFunction `json:"functions"`
}

type PreferenceFunction interface {
	utils.Identifiable
	MethodParameters() interface{}
	ParseParams(dm *DecisionMaker) interface{}
	Evaluate(dm *DecisionMaker) *AlternativesRanking
}

func (pf PreferenceFunctions) Get(index int) utils.Identifiable {
	return pf.Functions[index]
}

func (pf PreferenceFunctions) Len() int {
	return len(pf.Functions)
}
