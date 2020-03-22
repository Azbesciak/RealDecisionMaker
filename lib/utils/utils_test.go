package utils

import "testing"

func TestValueRange_ScaleEqually(t *testing.T) {
	checkRange(t, 0, 10, 2, 8, 0.6)
	checkRange(t, 0, 5, 1, 4, 0.6)
	checkRange(t, 0, 5, 2.5, 2.5, 0.0)
	checkRange(t, 0, 0, 0, 0, 0.0)
	checkRange(t, 0, 10, -10, 20, 3)
	checkRange(t, 10, 0, 20, -10, 3)
}

func checkRange(t *testing.T, min, max, expMin, expMax, scale float64) {
	valueRange := ValueRange{Min: min, Max: max}
	res := valueRange.ScaleEqually(scale)
	validateValue(t, "min", expMin, res.Min)
	validateValue(t, "max", expMax, res.Max)
}

func validateValue(t *testing.T, name string, expected, actual float64) {
	if !FloatsAreEqual(actual, expected, 1e-6) {
		t.Errorf("%s value expected %f, got %f", name, expected, actual)
	}
}
