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

var _normalRange = utils.ValueRange{
	Min: 0,
	Max: 1,
}

func RescaleCriterion(c *Criterion, alternatives *[]AlternativeWithCriteria, target *utils.ValueRange) Weights {
	currentRange := CriteriaValuesRange(alternatives, c)
	scaledCriterionValues := make(Weights, len(*alternatives))
	scale := GetScaleRatio(target, currentRange)
	for _, a := range *alternatives {
		scaledCriterionValues[a.Id] = scaleCriterion(c, a, currentRange, scale, target)
	}
	return scaledCriterionValues
}

func scaleCriterion(c *Criterion, a AlternativeWithCriteria, currentRange *utils.ValueRange, scale float64, target *utils.ValueRange) Weight {
	value := a.CriterionRawValue(c)
	if c.Type == Cost {
		return (currentRange.Max-value)*scale + target.Min
	} else {
		return (value-currentRange.Min)*scale + target.Min
	}
}

func GetNormalScaleRatio(currentRange *utils.ValueRange) float64 {
	return GetScaleRatio(&_normalRange, currentRange)
}

func GetScaleRatio(target *utils.ValueRange, currentRange *utils.ValueRange) float64 {
	targetDif := target.Diff()
	currentDif := currentRange.Diff()
	scale := 0.0
	if currentDif != 0 {
		scale = targetDif / currentDif
	}
	return scale
}
