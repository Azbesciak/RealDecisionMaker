package criteria_bounding

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

const epsilon = 1e-6

func TestFromParams(t *testing.T) {
	checkParsing(t, utils.Map{}, CriteriaBounding{
		AllowedValuesRangeScaling: -1,
		DisallowNegativeValues:    false,
	})
	checkParsing(t, utils.Map{"allowedValuesRangeScaling": 2}, CriteriaBounding{
		AllowedValuesRangeScaling: 2,
		DisallowNegativeValues:    false,
	})
	checkParsing(t, utils.Map{"disallowNegativeValues": true}, CriteriaBounding{
		AllowedValuesRangeScaling: -1,
		DisallowNegativeValues:    true,
	})
}

func checkParsing(t *testing.T, data interface{}, expected CriteriaBounding) {
	actual := FromParams(&data)
	if !utils.FloatsAreEqual(expected.AllowedValuesRangeScaling, actual.AllowedValuesRangeScaling, epsilon) {
		t.Errorf("different valuesRangeScaling, expected %f, got %f", expected.AllowedValuesRangeScaling, actual.AllowedValuesRangeScaling)
	}
	if expected.DisallowNegativeValues != actual.DisallowNegativeValues {
		t.Errorf("different disallowNegativeValues, expected %v, got %v", expected.DisallowNegativeValues, actual.DisallowNegativeValues)
	}
}

func TestCriteriaBounding_BoundValueNoLimit(t *testing.T) {
	bounding := CriteriaBounding{
		AllowedValuesRangeScaling: -1,
		DisallowNegativeValues:    false,
	}
	validateBounding(t, bounding, 1, 1, 2, 1)
	validateBounding(t, bounding, 1, 2, 2, 1)
	validateBounding(t, bounding, 4, 2, 2, 4)
	validateBounding(t, bounding, -1, 2, 2, -1)
}

func TestCriteriaBounding_BoundValueBelowZeroDisallowed(t *testing.T) {
	bounding := CriteriaBounding{
		AllowedValuesRangeScaling: -1,
		DisallowNegativeValues:    true,
	}
	validateBounding(t, bounding, 1, 1, 2, 1)
	validateBounding(t, bounding, 1, 2, 2, 1)
	validateBounding(t, bounding, 4, 2, 2, 4)
	validateBounding(t, bounding, -1, 2, 2, 0)
	validateBounding(t, bounding, -1, -2, 2, 0)
}

func TestCriteriaBounding_BoundValueBelowZeroDisallowedRegardlessRange(t *testing.T) {
	bounding := CriteriaBounding{
		AllowedValuesRangeScaling: 1,
		DisallowNegativeValues:    true,
	}
	validateBounding(t, bounding, 1, 1, 2, 1)
	validateBounding(t, bounding, 1, 2, 2, 2)
	validateBounding(t, bounding, 4, 1, 2, 2)
	validateBounding(t, bounding, 4, 1, 3, 3)
	validateBounding(t, bounding, -1, 2, 2, 2)
	validateBounding(t, bounding, -1, -2, 2, 0)
}

func TestCriteriaBounding_BoundValueWithRangeIncreased(t *testing.T) {
	bounding := CriteriaBounding{
		AllowedValuesRangeScaling: 2,
		DisallowNegativeValues:    true,
	}
	validateBounding(t, bounding, 10, 5, 15, 10)
	validateBounding(t, bounding, 1, 10, 20, 5)
	validateBounding(t, bounding, 30, 10, 20, 25)
	validateBounding(t, bounding, 4, 1, 2, 2.5)
	validateBounding(t, bounding, -10, 1, 3, 0)
	validateBounding(t, bounding, -1, -10, 10, 0)
}

func validateBounding(t *testing.T, bounding CriteriaBounding, value, min, max, expected float64) {
	valueRange := &utils.ValueRange{
		Min: min,
		Max: max,
	}
	actual := bounding.BoundValue(value, valueRange)
	if !utils.FloatsAreEqual(actual, expected, epsilon) {
		t.Errorf("different value after bounding, expected %f, got %f", expected, actual)
	}
	withRangeResult := bounding.WithRange(valueRange).BoundValue(value)
	if !utils.FloatsAreEqual(actual, withRangeResult, epsilon) {
		t.Errorf("different results after WithRange, expected %v, got %v", actual, withRangeResult)
	}
}
