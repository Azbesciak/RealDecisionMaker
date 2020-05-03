package model

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

func TestWeights_Merge(t *testing.T) {
	oldWeights := Weights{"1": 1, "2": 2, "3": 3}
	newWeights := Weights{"4": 4, "5": 5}
	actual := *oldWeights.Merge(&newWeights)
	expected := Weights{"1": 1, "2": 2, "3": 3, "4": 4, "5": 5}
	if utils.Differs(actual, expected) {
		t.Errorf("wrong weights after merge, expected %v, got %v", expected, actual)
	}
}

func TestWeights_MergeOverlapError(t *testing.T) {
	oldWeights := Weights{"1": 1, "2": 2, "3": 3}
	newWeights := Weights{"2": 2, "4": 4}
	defer utils.ExpectError(t, "criterion '2' from [{1 1} {2 2} {3 3}] already exists in [{2 2} {4 4}]")()
	oldWeights.Merge(&newWeights)
}

func TestWeights_PreserveOnly(t *testing.T) {
	weights := Weights{"1": 1, "2": 2, "3": 3, "4": 4}
	actual := *weights.PreserveOnly(&Criteria{{Id: "1"}, {Id: "3"}})
	expected := Weights{"1": 1, "3": 3}
	if utils.Differs(expected, actual) {
		t.Errorf("different values after preserve, expected %v, got %v", expected, actual)
	}
}

func TestWeights_PreserveOnlyError(t *testing.T) {
	weights := Weights{"1": 1, "2": 2, "3": 3, "4": 4}
	defer utils.ExpectError(t, "criterion 5 not found in [{1 1} {2 2} {3 3} {4 4}]")()
	weights.PreserveOnly(&Criteria{{Id: "1"}, {Id: "5"}})
}
