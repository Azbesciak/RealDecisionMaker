package logic

import (
	"../model"
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

func ContainsByIdentity(slice *[]model.Identifiable, value *string) bool {
	for _, v := range *slice {
		if *value == v.Identifier() {
			return true
		}
	}
	return false
}

func ContainsAll(slice []model.Identifiable, values *[]string) bool {
	for _, v := range *values {
		if !ContainsByIdentity(&slice, &v) {
			return false
		}
	}
	return true
}

type ToIdentityConverter func(interface{}) model.Identifiable

func ToIdentifiable(objects IdentifiableIterable) []model.Identifiable {
	var total = objects.Len()
	var interfaceSlice = make([]model.Identifiable, total)
	for i := 0; i < total; i++ {
		interfaceSlice[i] = objects.Get(i)
	}
	return interfaceSlice
}

type IdentifiableIterable interface {
	Get(index int) model.Identifiable
	Len() int
}

func FloatsAreEqual(expected model.Weight, actual model.Weight, epsilon model.Weight) bool {
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
