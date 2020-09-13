package anchoring

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
)

type IdealReferenceAlternativeEvaluator struct {
}

const IdealReferenceAltEvaluator = "ideal"

func (i *IdealReferenceAlternativeEvaluator) Identifier() string {
	return IdealReferenceAltEvaluator
}

func (i *IdealReferenceAlternativeEvaluator) BlankParams() FunctionParams {
	return i
}

type valueWithCoefficient struct {
	value       float64
	coefficient float64
}

func (i *IdealReferenceAlternativeEvaluator) Evaluate(
	params FunctionParams,
	alternatives *[]AnchoringAlternativeWithCriteria,
	criteria *model.Criteria,
) []model.AlternativeWithCriteria {
	return findBest(alternatives, criteria, IdealReferenceAltEvaluator, func(c *model.Criterion, a, b valueWithCoefficient) bool {
		return canNewBeBetter(a, b) && isBetter(c, a, b)
	})
}

const NadirReferenceAltEvaluator = "nadir"

type NadirReferenceAlternativeEvaluator struct {
}

func (i *NadirReferenceAlternativeEvaluator) Identifier() string {
	return NadirReferenceAltEvaluator
}

func (i *NadirReferenceAlternativeEvaluator) BlankParams() FunctionParams {
	return i
}

// we already assume that a & b coefficients are != 0
func isBetter(criterion *model.Criterion, a, b valueWithCoefficient) bool {
	if criterion.IsGain() {
		aVal, bVal := a.value*a.coefficient, b.value*b.coefficient
		if aVal == bVal {
			return a.value <= b.value
		} else {
			return aVal < bVal
		}
	} else {
		// division replacement
		aVal, bVal := a.value*b.coefficient, b.value*a.coefficient
		if aVal == bVal {
			return a.value >= b.value
		} else {
			return aVal > bVal
		}
	}
}

func canNewBeBetter(a, b valueWithCoefficient) bool {
	if b.coefficient == 0 && a.coefficient == 0 {
		return true
	} else if a.coefficient == 0 {
		return true
	} else if b.coefficient == 0 {
		return false
	}
	return true
}

func (i *NadirReferenceAlternativeEvaluator) Evaluate(
	params FunctionParams,
	alternatives *[]AnchoringAlternativeWithCriteria,
	criteria *model.Criteria,
) []model.AlternativeWithCriteria {
	return findBest(alternatives, criteria, NadirReferenceAltEvaluator, func(c *model.Criterion, a, b valueWithCoefficient) bool {
		return canNewBeBetter(a, b) && !isBetter(c, a, b)
	})
}

func findBest(
	alternatives *[]AnchoringAlternativeWithCriteria,
	criteria *model.Criteria,
	name string,
	isBetter func(c *model.Criterion, a, b valueWithCoefficient) bool,
) []model.AlternativeWithCriteria {
	best := prepareCriteriaWithCoefficients(alternatives, criteria)
	findBestCriteriaValues(alternatives, criteria, best, isBetter)
	result := extractCriteriaValues(best)
	return []model.AlternativeWithCriteria{{
		Id:       name,
		Criteria: result,
	}}
}

func extractCriteriaValues(best *map[string]valueWithCoefficient) model.Weights {
	result := make(model.Weights, len(*best))
	for criterion, value := range *best {
		result[criterion] = value.value
	}
	return result
}

func findBestCriteriaValues(
	alternatives *[]AnchoringAlternativeWithCriteria,
	criteria *model.Criteria,
	best *map[string]valueWithCoefficient,
	isBetter func(c *model.Criterion, a valueWithCoefficient, b valueWithCoefficient) bool,
) {
	for i := 1; i < len(*alternatives); i++ {
		alt := (*alternatives)[i]
		for _, c := range *criteria {
			criterionValue := alt.Alternative.CriterionRawValue(&c)
			if oldValue, ok := (*best)[c.Id]; !ok {
				panic(fmt.Errorf("criterion '%s' not found in criteria %v", c.Id, *criteria.Names()))
			} else {
				newValue := valueWithCoefficient{
					value:       criterionValue,
					coefficient: alt.Coefficient,
				}
				if isBetter(&c, oldValue, newValue) {
					(*best)[c.Id] = newValue
				}
			}
		}
	}
}

func prepareCriteriaWithCoefficients(
	alternatives *[]AnchoringAlternativeWithCriteria,
	criteria *model.Criteria,
) *map[string]valueWithCoefficient {
	best := make(map[string]valueWithCoefficient, len(*criteria))
	firstAlternative := (*alternatives)[0]
	for criterionId, criterionValue := range firstAlternative.Alternative.Criteria {
		best[criterionId] = valueWithCoefficient{
			value:       criterionValue,
			coefficient: firstAlternative.Coefficient,
		}
	}
	return &best
}
