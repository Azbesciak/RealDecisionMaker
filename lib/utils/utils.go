package utils

import (
	"github.com/mitchellh/mapstructure"
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

func ContainsInts(slice *[]int, value *int) bool {
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
			t.Errorf("expected error with message '%s'", expectedMessage)
		} else if ErrorDiffers(e, expectedMessage) {
			t.Errorf("invalid error message:\nactual   '%s'\nexpected '%s'", e, expectedMessage)
		}
	}
}

func IsPositive(value float64) bool {
	return value > 0
}

func DecodeToStruct(src, target interface{}) {
	e := mapstructure.Decode(src, target)
	if e != nil {
		panic(e)
	}
}
