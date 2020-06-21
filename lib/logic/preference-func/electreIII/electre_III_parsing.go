package electreIII

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

//go:generate easytags $GOFILE json:camel
type electreIIIParams struct {
	Criteria        *ElectreCriteria                `json:"criteria"`
	DistillationFun *utils.LinearFunctionParameters `json:"distillationFun,omitempty"`
}

const distillationFun = "electreDistillation"
const criteria = "electreCriteria"

func (e *ElectreIIIPreferenceFunc) ParseParams(dm *model.DecisionMaker) interface{} {
	return electreIIIParams{
		Criteria:        extractElectreIIICriteria(dm),
		DistillationFun: getDistillationFunc(dm),
	}
}

func extractElectreIIICriteria(dm *model.DecisionMaker) *ElectreCriteria {
	potentialEleCriteria, ok := dm.MethodParameters[criteria]
	if !ok {
		panic(fmt.Errorf("Criteria for electre not found in methodParameters: %v", dm.MethodParameters))
	}
	electreCriteria := make(ElectreCriteria)
	utils.DecodeToStruct(potentialEleCriteria, &electreCriteria)
	for _, criterion := range dm.Criteria {
		electreCriterion, cOk := electreCriteria[criterion.Id]
		if !cOk {
			panic(fmt.Errorf("criterion '%s' not found in electre Criteria: %v", criterion.Id, electreCriteria))
		}
		validateParameters(&criterion, &electreCriterion)
	}
	return &electreCriteria
}

func validateParameters(criterion *model.Criterion, crit *ElectreCriterion) {
	if crit.K <= 0 {
		panic(fmt.Errorf("electre criterion's weight must be positive, got %v for %s", crit.K, criterion.Id))
	}
	lastWeight := 0.0
	lastWeight = requireBValueAtLeast(&crit.Q, lastWeight, criterion.Id, "Q")
	lastWeight = requireBValueAtLeast(&crit.P, lastWeight, criterion.Id, "P")
	requireBValueAtLeast(&crit.V, lastWeight, criterion.Id, "V")
}

func requireBValueAtLeast(f *utils.LinearFunctionParameters, current float64, criterion, funcName string) float64 {
	if f.A == 0 && f.B != 0 && f.B <= current {
		panic(fmt.Errorf(
			"b parameter of electre pref func %s %v for criterion %s must be greater than %f",
			funcName, f, criterion, current,
		))
	}
	if f.B > 0 {
		return f.B
	}
	return current
}

func getDistillationFunc(dm *model.DecisionMaker) *utils.LinearFunctionParameters {
	params, ok := dm.MethodParameters[distillationFun]
	if !ok {
		return &DefaultDistillationFunc
	} else {
		parameters := utils.LinearFunctionParameters{}
		utils.DecodeToStruct(params, &parameters)
		return &parameters
	}
}
