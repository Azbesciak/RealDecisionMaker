package electreIII

import (
	"fmt"
	. "github.com/Azbesciak/RealDecisionMaker/lib/model"
	. "github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type ElectreIIIPreferenceFunc struct {
}

func (e ElectreIIIPreferenceFunc) Identifier() string {
	return "electreIII"
}

func (e ElectreIIIPreferenceFunc) Evaluate(dm *DecisionMaker) *AlternativesRanking {
	criteria := *dm.AlternativesToConsider()
	eleIIICriteria := extractElectreIIICriteria(dm)
	distillationFunc := getDistillationFunc(dm)
	return ElectreIII(criteria, dm.Criteria, eleIIICriteria, distillationFunc)
}

func extractElectreIIICriteria(dm *DecisionMaker) *ElectreCriteria {
	potentialEleCriteria, ok := dm.MethodParameters["electreCriteria"]
	if !ok {
		panic(fmt.Errorf("criteria for electre not found in methodParameters: %v", dm.MethodParameters))
	}
	electreCriteria := make(ElectreCriteria)
	DecodeToStruct(potentialEleCriteria, &electreCriteria)
	for _, criterion := range dm.Criteria {
		electreCriterion, cOk := electreCriteria[criterion.Id]
		if !cOk {
			panic(fmt.Errorf("criterion '%s' not found in electre criteria: %v", criterion.Id, electreCriteria))
		}
		validateParameters(&criterion, &electreCriterion)
	}
	return &electreCriteria
}

func validateParameters(criterion *Criterion, crit *ElectreCriterion) {
	if crit.K <= 0 {
		panic(fmt.Errorf("electre criterion's weight must be positive, got %v for %s", crit.K, criterion.Id))
	}
	lastWeight := 0.0
	lastWeight = requireBValueAtLeast(&crit.Q, lastWeight, criterion.Id, "Q")
	lastWeight = requireBValueAtLeast(&crit.P, lastWeight, criterion.Id, "P")
	requireBValueAtLeast(&crit.V, lastWeight, criterion.Id, "V")
}

func requireBValueAtLeast(f *LinearFunctionParameters, current float64, criterion, funcName string) float64 {
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

func getDistillationFunc(dm *DecisionMaker) *LinearFunctionParameters {
	params, ok := dm.MethodParameters["electreDistillation"]
	if !ok {
		return &DefaultDistillationFunc
	} else {
		parameters := LinearFunctionParameters{}
		DecodeToStruct(params, &parameters)
		return &parameters
	}
}

func ElectreIII(
	alternatives []AlternativeWithCriteria,
	criteria Criteria,
	electreCriteria *ElectreCriteria,
	distillationFun *LinearFunctionParameters,
) *AlternativesRanking {
	matrix := evaluateCredibilityMatrix(&alternatives, &criteria, electreCriteria)
	ascending := RankAscending(matrix, distillationFun)
	descending := RankDescending(matrix, distillationFun)
	return EvaluateRanking(ascending, descending, &alternatives)
}

func evaluateCredibilityMatrix(
	alternatives *[]AlternativeWithCriteria,
	criteria *Criteria,
	electreCriteria *ElectreCriteria,
) *AlternativesMatrix {
	alternativesNum := len(*alternatives)
	credibilityFlatMatrix := make([]float64, alternativesNum*alternativesNum)
	alternativesIds := make(Alternatives, alternativesNum)
	for i, a1 := range *alternatives {
		alternativesIds[i] = a1.Id
		for j, a2 := range *alternatives {
			credibilityFlatMatrix[i*alternativesNum+j] = evaluateAlternativesPair(i, j, &a1, &a2, criteria, electreCriteria)
		}
	}
	return &AlternativesMatrix{&alternativesIds, &Matrix{
		Size: alternativesNum,
		Data: credibilityFlatMatrix,
	}}
}

func evaluateAlternativesPair(i, j int, a1, a2 *AlternativeWithCriteria, criteria *Criteria, electreCriteria *ElectreCriteria) float64 {
	if i == j {
		return 1
	} else {
		eleRes := electreIIICredibility(a1, a2, criteria, electreCriteria)
		return eleRes.D
	}
}

func electreIIICredibility(
	a1, a2 *AlternativeWithCriteria,
	criteria *Criteria,
	criteriaThresholds *ElectreCriteria,
) *ElectreResult {
	electreRes := make([]*electreIIISingleResult, len(*criteria))
	for i, c := range *criteria {
		electreRes[i] = evaluatePair(a1, a2, &c, criteriaThresholds)
	}
	c := calculateTotalC(&electreRes)
	d := calculateCredibility(c, &electreRes)
	return &ElectreResult{C: c, D: d}
}

func calculateTotalC(results *[]*electreIIISingleResult) float64 {
	weightSum := 0.0
	totalC := 0.0
	for _, c := range *results {
		weightSum += c.criterion.K
		totalC += c.criterion.K * c.result.C
	}
	return totalC / weightSum
}

func calculateCredibility(C float64, results *[]*electreIIISingleResult) float64 {
	credibility := C
	for _, res := range *results {
		if res.result.D > C {
			credibility *= (1 - res.result.D) / (1 - C)
		}
	}
	return credibility
}

// TODO what with diff == 0?
func evaluatePair(
	a1, a2 *AlternativeWithCriteria,
	c *Criterion,
	criteriaThresholds *ElectreCriteria,
) *electreIIISingleResult {
	c1Val := a1.CriterionValue(c)
	c2Val := a2.CriterionValue(c)
	ths, foundThreshold := (*criteriaThresholds)[c.Id]
	if !foundThreshold {
		panic(fmt.Errorf("properties for criterion '%s' not found", c.Id))
	}
	return &electreIIISingleResult{
		criterion: &ths,
		result:    calculateElectreResult(c1Val, c2Val, c, &ths),
	}
}

func calculateElectreResult(c1Val, c2Val Weight, c *Criterion, ths *ElectreCriterion) *ElectreResult {
	if c1Val > c2Val {
		return &ElectreResult{C: 1}
	}
	originalFirstCriterionValue := c1Val * Weight(c.Multiplier())
	criteriaValueDifference := c2Val - c1Val
	q, qok := ths.Q.evaluate(originalFirstCriterionValue)
	if qok && q >= criteriaValueDifference {
		return &ElectreResult{C: 1}
	}
	p, pok := ths.P.evaluate(originalFirstCriterionValue)
	if pok && p >= criteriaValueDifference {
		return &ElectreResult{C: 1 - (criteriaValueDifference-q)/(p-q)}
	}
	v, vok := ths.V.evaluate(originalFirstCriterionValue)
	if vok && v >= criteriaValueDifference {
		return &ElectreResult{D: (criteriaValueDifference - p) / (v - p)}
	}
	if vok && v < criteriaValueDifference {
		return &ElectreResult{D: 1}
	}
	return &ElectreResult{}
}

type ElectreResult struct {
	C float64
	D float64
}

type ElectreCriteria = map[string]ElectreCriterion

type ElectreCriterion struct {
	K float64
	Q LinearFunctionParameters
	P LinearFunctionParameters
	V LinearFunctionParameters
}

type LinearFunctionParameters struct {
	A float64
	B float64
}

func (f *LinearFunctionParameters) String() string {
	return fmt.Sprintf("a:%v, b:%v", f.A, f.B)
}

func (f *LinearFunctionParameters) evaluate(value float64) (result float64, ok bool) {
	if f.A == 0 && f.B == 0 {
		return 0, false
	}
	return f.A*value + f.B, true
}

type AlternativesMatrix struct {
	Alternatives *Alternatives
	Values       *Matrix
}

type electreIIISingleResult struct {
	criterion *ElectreCriterion
	result    *ElectreResult
}
