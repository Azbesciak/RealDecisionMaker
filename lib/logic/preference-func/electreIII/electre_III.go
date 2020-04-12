package electreIII

import (
	"fmt"
	. "github.com/Azbesciak/RealDecisionMaker/lib/model"
)

const methodName = "electreIII"

type ElectreIIIPreferenceFunc struct {
}

type ElectreIIIInputParams struct {
	Criteria        ElectreCriteria          `json:"criteria"`
	DistillationFun LinearFunctionParameters `json:"distillationFun"`
}

func (e *ElectreIIIPreferenceFunc) Identifier() string {
	return methodName
}

func (e *ElectreIIIPreferenceFunc) MethodParameters() interface{} {
	return ElectreIIIInputParams{}
}

func (e *ElectreIIIPreferenceFunc) Evaluate(dmp *DecisionMakingParams) *AlternativesRanking {
	params := dmp.MethodParameters.(electreIIIParams)
	return ElectreIII(dmp.ConsideredAlternatives, dmp.Criteria, params.criteria, params.distillationFun)
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
