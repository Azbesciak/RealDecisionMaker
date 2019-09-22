package logic

import "testing"

func TestOWAResult(t *testing.T) {
	criteria := AlternativeWithCriteria{
		Alternative: Alternative{"test"},
		Criteria: Weights{
			"1": 30,
			"2": 10,
			"3": 5,
			"4": 60,
		},
	}
	result := OWA(criteria, []Weight{1, 2, 3, 4})
	var expected Weight = 4*60 + 30*3 + 10*2 + 5*1
	if result.Value != expected {
		t.Errorf("Expected %f, got %f", expected, result.Value)
	}
	result = OWA(criteria, []Weight{6, 4, 3, 4})
	expected = 6*60 + 30*4 + 10*4 + 5*3
	if result.Value != expected {
		t.Errorf("Expected %f, got %f", expected, result.Value)
	}
}

func TestOWAInvalidArguments(t *testing.T) {
	criteria := AlternativeWithCriteria{
		Alternative: Alternative{"test"},
		Criteria: Weights{
			"1": 30,
			"2": 10,
		},
	}
	defer func() {
		if err := recover(); err == nil {
			t.Error("should throw error")
		} else if err != "criteria and weights must have the same length, got 2 and 1" {
			t.Error("Invalid error message", err)
		}
	}()
	OWA(criteria, []Weight{1})
}
