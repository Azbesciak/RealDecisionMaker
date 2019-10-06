package utils

import (
	"math"
	"testing"
)

func ContainsString(slice *[]string, value *string) bool {
	for _, v := range *slice {
		if v == *value {
			return true
		}
	}
	return false
}

func FloatsAreEqual(expected float64, actual float64, epsilon float64) bool {
	return math.Abs(expected-actual) <= epsilon
}

func ErrorDiffers(err interface{}, message string) bool {
	return err.(error).Error() != message
}

func ExpectError(t *testing.T, expectedMessage string) func() {
	return func() {
		if e := recover(); e == nil {
			t.Errorf("expected error")
		} else if ErrorDiffers(e, expectedMessage) {
			t.Errorf("invalid error message, got '%s'", e)
		}
	}
}
