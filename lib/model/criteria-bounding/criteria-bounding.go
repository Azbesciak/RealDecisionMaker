package criteria_bounding

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type CriteriaBounding struct {
	AllowedValuesRangeScaling float64 `json:"allowedValuesRangeScaling"`
	DisallowNegativeValues    bool    `json:"disallowNegativeValues"`
}

type CriteriaInRangeBounding struct {
	bounding   *CriteriaBounding
	valueRange *utils.ValueRange
}

func DefaultParams() *CriteriaBounding {
	return &CriteriaBounding{
		AllowedValuesRangeScaling: -1.0,
		DisallowNegativeValues:    false,
	}
}

func FromParams(params *interface{}) *CriteriaBounding {
	bounding := DefaultParams()
	utils.DecodeToStruct(*params, bounding)
	if bounding.AllowedValuesRangeScaling == 0 {
		panic(fmt.Errorf("allowedValuesRangeScaling cannot be 0"))
	}
	return bounding
}

func (b *CriteriaBounding) WithRange(valueRange *utils.ValueRange) *CriteriaInRangeBounding {
	var scaled *utils.ValueRange = nil
	if b.AllowedValuesRangeScaling > 0 {
		scaled = scaleRange(valueRange, b.AllowedValuesRangeScaling)
	}
	return &CriteriaInRangeBounding{
		bounding:   b,
		valueRange: scaled,
	}
}

func (b *CriteriaInRangeBounding) BoundValue(value float64) float64 {
	value = b.bounding.trimBelowZeroIfRequired(value)
	if b.valueRange == nil {
		return value
	}
	return boundValueInRange(value, b.valueRange)
}

func (b *CriteriaBounding) BoundValue(value float64, valueRange *utils.ValueRange) float64 {
	value = b.trimBelowZeroIfRequired(value)
	return boundValue(value, b.AllowedValuesRangeScaling, valueRange)
}

func (b *CriteriaBounding) trimBelowZeroIfRequired(value float64) float64 {
	if b.DisallowNegativeValues && value < 0 {
		value = 0
	}
	return value
}

func boundValue(value, scaling float64, valueRange *utils.ValueRange) float64 {
	if scaling > 0 {
		scaledRange := scaleRange(valueRange, scaling)
		return boundValueInRange(value, scaledRange)
	}
	return value
}

func boundValueInRange(value float64, scaledRange *utils.ValueRange) float64 {
	if value < scaledRange.Min {
		value = scaledRange.Min
	}
	if value > scaledRange.Max {
		value = scaledRange.Max
	}
	return value
}

func scaleRange(valueRange *utils.ValueRange, scaling float64) *utils.ValueRange {
	if scaling == 1 {
		return valueRange
	} else {
		return valueRange.ScaleEqually(scaling)
	}
}
