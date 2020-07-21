package utils

import "fmt"

type LinearFunctionParameters struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

const LinearFunctionName = "linear"

func (f *LinearFunctionParameters) String() string {
	return fmt.Sprintf("a:%v, b:%v", f.A, f.B)
}

func (f *LinearFunctionParameters) Evaluate(value float64) (result float64, ok bool) {
	if f.A == 0 && f.B == 0 {
		return 0, false
	}
	return f.A*value + f.B, true
}
