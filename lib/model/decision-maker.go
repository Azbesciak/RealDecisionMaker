package model

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"strings"
)

//go:generate easytags $GOFILE json:camel

type DecisionMaker struct {
	PreferenceFunction string                    `json:"preferenceFunction"`
	Biases             BiasesParams              `json:"biases"`
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

func (p *DecisionMakingParams) AllAlternatives() []AlternativeWithCriteria {
	notConsider := p.NotConsideredAlternatives
	if notConsider == nil {
		notConsider = make([]AlternativeWithCriteria, 0)
	}
	toConsider := p.ConsideredAlternatives
	if toConsider == nil {
		toConsider = make([]AlternativeWithCriteria, 0)
	}
	return append(toConsider, notConsider...)
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
	Result AlternativesRanking `json:"result"`
	Biases BiasesParams        `json:"biases"`
}

func (dm *DecisionMaker) MakeDecision(
	preferenceFunctions PreferenceFunctions,
	biasListeners BiasListeners,
	availableBiases *BiasMap,
) *DecisionMakerChoice {
	if IsStringBlank(&dm.PreferenceFunction) {
		panic(fmt.Errorf("preference function must not be empty"))
	}
	dm.validateCriteria()
	dm.validateAlternatives()
	preferenceFunction := preferenceFunctions.Fetch(dm.PreferenceFunction)
	params := dm.prepareParams(preferenceFunction)
	chosenBiases := ChooseBiases(availableBiases, &dm.Biases)
	processedParams, biasesProps := dm.processBiases(chosenBiases, params, &biasListeners)
	res := (*preferenceFunction).Evaluate(processedParams)
	return &DecisionMakerChoice{*res, *biasesProps}
}

func (dm *DecisionMaker) validateCriteria() {
	criteriaSet := make(map[string]bool)
	for i, c := range dm.Criteria {
		if _, ok := criteriaSet[c.Id]; ok {
			panic(fmt.Errorf("criterion '%s' [index %d] is not unique", c.Id, i))
		}
		criteriaSet[c.Id] = true
	}
}

func (dm *DecisionMaker) validateAlternatives() {
	for i, a := range dm.KnownAlternatives {
		for _, c := range dm.Criteria {
			_, ok := a.Criteria[c.Id]
			if !ok {
				panic(fmt.Errorf("value of criterion '%s' not found for alternative %d '%s'", c.Id, i, a.Id))
			}
		}
	}
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

func (dm *DecisionMaker) processBiases(
	biases *BiasesWithProps,
	params *DecisionMakingParams,
	listeners *BiasListeners,
) (*DecisionMakingParams, *BiasesParams) {
	biasesToProcessCount := len(*biases)
	result := make(BiasesParams, biasesToProcessCount)
	current := params
	if biasesToProcessCount == 0 {
		return current, &result
	}
	listener := listeners.Fetch(dm.PreferenceFunction)
	for i, h := range *biases {
		res := (*h.Bias).Apply(params, current, &h.Props.Props, listener)
		current = res.DMP
		result[i] = *UpdateBiasesProps(h.Props, res.Props)
	}
	return current, &result
}

func IsStringBlank(str *string) bool {
	return len(strings.TrimSpace(*str)) == 0
}
