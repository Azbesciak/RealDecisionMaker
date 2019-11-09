package electreIII

import (
	. "github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

var power = "power"
var safety = "safety"
var cost = "cost"
var powerCriterion = &Criterion{Id: power, Type: Gain}
var safetyCriterion = &Criterion{Id: safety, Type: Gain}
var costCriterion = &Criterion{Id: cost, Type: Cost}
var criteria = &Criteria{*powerCriterion, *safetyCriterion, *costCriterion}
var electreIIICriteria = &map[string]ElectreCriterion{
	power: {
		K: 3,
		Q: LinearFunctionParameters{B: 4},
		P: LinearFunctionParameters{B: 12},
		V: LinearFunctionParameters{B: 28},
	},
	safety: {
		K: 3,
		Q: LinearFunctionParameters{B: 1},
		P: LinearFunctionParameters{B: 2},
		V: LinearFunctionParameters{B: 8},
	},
	cost: {
		K: 4,
		Q: LinearFunctionParameters{B: 100},
		P: LinearFunctionParameters{B: 200},
		V: LinearFunctionParameters{B: 600},
	},
}
var fra = &AlternativeWithCriteria{
	Id:       "FRA",
	Criteria: Weights{power: 98, safety: 6, cost: 800},
}
var ita = &AlternativeWithCriteria{
	Id:       "ITA",
	Criteria: Weights{power: 90, safety: 4, cost: 600},
}

func TestElectreIIIPreferenceFunc_Evaluate(t *testing.T) {
	checkEvaluate(LinearFunctionParameters{B: 10}, 10, 10.0, true, t)
	checkEvaluate(LinearFunctionParameters{A: 10}, 10, 100.0, true, t)
	checkEvaluate(LinearFunctionParameters{A: 0, B: 0}, 10, 0.0, false, t)
	checkEvaluate(LinearFunctionParameters{A: 14, B: 5}, 10, 145.0, true, t)
}

func checkEvaluate(f LinearFunctionParameters, value, expectedValue float64, expectedOk bool, t *testing.T) {
	res, ok := f.evaluate(value)
	if ok != expectedOk {
		t.Errorf("expected ok to be '%v', got '%v'", expectedOk, ok)
	}
	if ok && res != expectedValue {
		t.Errorf("invalid value from evaluate for %v, expected %f, got %f", &f, expectedValue, res)
	}
}

func TestElectreIIIEvaluatePair(t *testing.T) {
	assertPair(t, fra, ita, powerCriterion, electreIIICriteria, &ElectreResult{C: 1, D: 0})
	assertPair(t, fra, ita, safetyCriterion, electreIIICriteria, &ElectreResult{C: 1, D: 0})
	assertPair(t, fra, ita, costCriterion, electreIIICriteria, &ElectreResult{C: 0, D: 0})
	assertPair(t, ita, fra, powerCriterion, electreIIICriteria, &ElectreResult{C: 0.5, D: 0})
	assertPair(t, ita, fra, safetyCriterion, electreIIICriteria, &ElectreResult{C: 0, D: 0})
	assertPair(t, ita, fra, costCriterion, electreIIICriteria, &ElectreResult{C: 1, D: 0})
}

func TestElectreIIIFinalResult(t *testing.T) {
	assertElectreIIIPairCredibilityPerCriterion(t, fra, ita, criteria, electreIIICriteria, &ElectreResult{C: 0.6, D: 0.6})
	assertElectreIIIPairCredibilityPerCriterion(t, ita, fra, criteria, electreIIICriteria, &ElectreResult{C: 0.55, D: 0.55})
}

func assertElectreIIIPairCredibilityPerCriterion(
	t *testing.T,
	a1, a2 *AlternativeWithCriteria,
	criteria *Criteria,
	electreCriteria *ElectreCriteria,
	expected *ElectreResult,
) {
	result := electreIIICredibility(a1, a2, criteria, electreCriteria)
	noneCriterion := &Criterion{Id: "none", Type: Cost}
	validate(t, expected.C, result.C, a1, a2, noneCriterion, "C")
	validate(t, expected.D, result.D, a1, a2, noneCriterion, "D")
}

func assertPair(
	t *testing.T,
	a1, a2 *AlternativeWithCriteria,
	criterion *Criterion,
	criteria *ElectreCriteria,
	expected *ElectreResult,
) {
	result := evaluatePair(a1, a2, criterion, criteria)
	validate(t, expected.C, result.result.C, a1, a2, criterion, "C")
	validate(t, expected.D, result.result.D, a1, a2, criterion, "D")
}

func validate(t *testing.T, expected, result float64, a1, a2 *AlternativeWithCriteria, criterion *Criterion, typ string) {
	if !utils.FloatsAreEqual(expected, result, 1e-4) {
		t.Errorf(
			"alterative 1: '%s', 2: '%s', criterion '%s': expected %s to be %v, got %v",
			a1.Id, a2.Id, criterion.Id, typ, expected, result,
		)
	}
}
