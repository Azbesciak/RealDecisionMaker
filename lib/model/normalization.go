package model

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"math"
)

func ValuesRangeWithGroundZero(alternatives *[]AlternativeWithCriteria, criterion *Criterion) *utils.ValueRange {
	valRange := CriteriaValuesRange(alternatives, criterion)
	minAbs := math.Abs(valRange.Min)
	maxAbs := math.Abs(valRange.Max)
	return &utils.ValueRange{
		Min: 0,
		Max: math.Max(math.Max(minAbs, maxAbs), valRange.Diff()),
	}
}

func RescaleCriterion(c *Criterion, alternatives *[]AlternativeWithCriteria, target *utils.ValueRange) Weights {
	currentRange := CriteriaValuesRange(alternatives, c)
	values := make(Weights, len(*alternatives))
	targetDif := target.Diff()
	currentDif := currentRange.Diff()
	scale := 0.0
	if currentDif != 0 {
		scale = targetDif / currentDif
	}
	for _, a := range *alternatives {
		value := a.CriterionRawValue(c)
		if c.Type == Cost {
			value = (currentRange.Max-value)*scale + target.Min
		} else {
			value = (value-currentRange.Min)*scale + target.Min
		}
		values[a.Id] = value
	}
	return values
}
