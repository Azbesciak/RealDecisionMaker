package utils

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/mitchellh/mapstructure"
	"math"
	"math/rand"
	"reflect"
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

func RemoveSingleStringOccurrence(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
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

func CheckValueRange(t *testing.T, valueRange ValueRange, expMin, expMax float64) {
	validateValue(t, "min", expMin, valueRange.Min)
	validateValue(t, "max", expMax, valueRange.Max)
}

func validateValue(t *testing.T, name string, expected, actual float64) {
	if !FloatsAreEqual(actual, expected, 1e-6) {
		t.Errorf("%s value expected %f, got %f", name, expected, actual)
	}
}

func IsPositive(value float64) bool {
	return value > 0
}

func IsInBounds(value float64, lower float64, upper float64) bool {
	return value >= lower && value <= upper
}

func IsProbability(value float64) bool {
	return IsInBounds(value, 0, 1)
}

func DecodeToStruct(src, target interface{}) {
	e := mapstructure.Decode(src, target)
	if e != nil {
		panic(e)
	}
}

type ValueRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

func (r *ValueRange) String() string {
	return fmt.Sprintf("min: %.4f, max: %.4f", r.Min, r.Max)
}

func (r *ValueRange) Diff() float64 {
	return r.Max - r.Min
}

func (r *ValueRange) ScaleEqually(scale float64) *ValueRange {
	dif := r.Diff() / 2
	return &ValueRange{
		Min: r.Min + dif - dif*scale,
		Max: r.Max - dif + dif*scale,
	}
}

func NewValueRange() *ValueRange {
	return &ValueRange{
		Min: 0,
		Max: 0,
	}
}

type SeededValueGenerator = func(seed int64) ValueGenerator
type ValueGenerator = func() float64

func RandomGenerator(seed int64) *rand.Rand {
	source := rand.NewSource(seed)
	return rand.New(source)
}

func RandomBasedSeedValueGenerator(seed int64) ValueGenerator {
	gen := RandomGenerator(seed)
	return func() float64 {
		return gen.Float64()
	}
}

func NewValueInRangeGenerator(generator ValueGenerator, valueRange *ValueRange) ValueGenerator {
	dif := valueRange.Max - valueRange.Min
	return func() float64 {
		return generator()*dif + valueRange.Min
	}
}

type Map = map[string]interface{}
type Array = []interface{}

func Differs(a, b interface{}) bool {
	return !cmp.Equal(a, b, cmpopts.EquateApprox(0, 1e-8), cmp.Exporter(func(r reflect.Type) bool {
		return true
	}))
}
