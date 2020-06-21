package electreIII

import (
	"fmt"
	. "github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

const methodName = "electreIII"

//go:generate easytags $GOFILE json:camel
type ElectreIIIPreferenceFunc struct {
}

type ElectreIIIInputParams struct {
	Criteria        ElectreCriteria                `json:"criteria"`
	DistillationFun utils.LinearFunctionParameters `json:"distillationFun,omitempty"`
}

func (e *ElectreIIIPreferenceFunc) Identifier() string {
	return methodName
}

func (e *ElectreIIIPreferenceFunc) MethodParameters() interface{} {
	return ElectreIIIInputParams{}
}

func (e *ElectreIIIPreferenceFunc) Evaluate(dmp *DecisionMakingParams) *AlternativesRanking {
	params := dmp.MethodParameters.(electreIIIParams)
	return ElectreIII(dmp.ConsideredAlternatives, dmp.Criteria, params.Criteria, params.DistillationFun)
}

func ElectreIII(
	alternatives []AlternativeWithCriteria,
	criteria Criteria,
	electreCriteria *ElectreCriteria,
	distillationFun *utils.LinearFunctionParameters,
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
	q, qok := ths.Q.Evaluate(originalFirstCriterionValue)
	if qok && q >= criteriaValueDifference {
		return &ElectreResult{C: 1}
	}
	p, pok := ths.P.Evaluate(originalFirstCriterionValue)
	if pok && p >= criteriaValueDifference {
		return &ElectreResult{C: 1 - (criteriaValueDifference-q)/(p-q)}
	}
	v, vok := ths.V.Evaluate(originalFirstCriterionValue)
	if vok && v >= criteriaValueDifference {
		return &ElectreResult{D: (criteriaValueDifference - p) / (v - p)}
	}
	if vok && v < criteriaValueDifference {
		return &ElectreResult{D: 1}
	}
	return &ElectreResult{}
}

type ElectreResult struct {
	C float64 `json:"c"`
	D float64 `json:"d"`
}

type ElectreCriteria = map[string]ElectreCriterion

type ElectreCriterion struct {
	K float64                        `json:"k"`
	Q utils.LinearFunctionParameters `json:"q"`
	P utils.LinearFunctionParameters `json:"p"`
	V utils.LinearFunctionParameters `json:"v"`
}

type AlternativesMatrix struct {
	Alternatives *Alternatives `json:"alternatives"`
	Values       *Matrix       `json:"values"`
}

type electreIIISingleResult struct {
	criterion *ElectreCriterion
	result    *ElectreResult
}
