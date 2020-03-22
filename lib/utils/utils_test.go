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
	CheckValueRange(t, *res, expMin, expMax)
}
