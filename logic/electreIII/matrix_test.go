package electreIII

import (
	"reflect"
	"testing"
)

var testMat = Matrix{
	Size: 5,
	Data: []float64{
		1, 0, 1, 0.8, 1,
		0, 1, 0, 0.9, 0.67,
		0.6, 0, 1, 0.6, 0.8,
		0.25, 0.8, 0.67, 1, 0.85,
		0.67, 0, 0.8, 0.8, 1,
	},
}

func TestMatrix_Without(t *testing.T) {
	without1And3 := testMat.Without(&[]int{1, 3})
	expected := &Matrix{
		Size: 3,
		Data: []float64{
			1, 1, 1,
			0.6, 1, 0.8,
			0.67, 0.8, 1,
		},
	}
	validateMatrices(t, expected, without1And3, "cut")
}

func TestMatrix_Slice(t *testing.T) {
	slice1And4 := testMat.Slice(&[]int{1, 4})
	expected := &Matrix{
		Size: 2,
		Data: []float64{
			1, 0.67,
			0, 1,
		},
	}
	validateMatrices(t, expected, slice1And4, "slice")
}

func validateMatrices(t *testing.T, expected, actual *Matrix, action string) {
	if !reflect.DeepEqual(*expected, *actual) {
		t.Errorf("expected matrix after %s to be %v, got %v", action, *expected, *actual)
	}
}
