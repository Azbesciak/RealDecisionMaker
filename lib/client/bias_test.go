package client

import (
	. "github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

type MockBias struct {
	Name string
}

func (m *MockBias) Identifier() string {
	return m.Name
}
func (m *MockBias) Apply(original, current *DecisionMakingParams, props *BiasProps, listener *BiasListener) *BiasedResult {
	panic("should not be called")
}

func TestBiasesSelection(t *testing.T) {
	availableBiases := AsBiasesMap(&Biases{
		&MockBias{"a"},
		&MockBias{"c"},
		&MockBias{"b"},
	})
	params := &BiasesParams{
		BiasParams{Name: "b"},
		BiasParams{Name: "a"},
	}
	result := ChooseBiases(availableBiases, params)
	checkResult(t, result, params)
}

func TestBiasesFiltering(t *testing.T) {
	availableBiases := AsBiasesMap(&Biases{
		&MockBias{"a"},
		&MockBias{"c"},
		&MockBias{"b"},
	})
	params := &BiasesParams{
		BiasParams{Name: "b", Disabled: true},
		BiasParams{Name: "c"},
		BiasParams{Name: "a"},
	}
	expected := &BiasesParams{
		BiasParams{Name: "c"},
		BiasParams{Name: "a"},
	}
	result := ChooseBiases(availableBiases, params)
	checkResult(t, result, expected)
}

func TestBiasNotFound(t *testing.T) {
	availableBiases := AsBiasesMap(&Biases{
		&MockBias{"b"},
	})
	params := &BiasesParams{
		BiasParams{Name: "a"},
	}
	defer utils.ExpectError(t, "bias 'a' not found, available are '[b]'")()
	ChooseBiases(availableBiases, params)
}

func TestSkippingDisabledNotExistingBiases(t *testing.T) {
	availableBiases := AsBiasesMap(&Biases{
		&MockBias{"b"},
	})
	params := &BiasesParams{
		BiasParams{Name: "a", Disabled: true},
	}
	result := ChooseBiases(availableBiases, params)
	checkResult(t, result, &BiasesParams{})
}

func checkResult(t *testing.T, result *BiasesWithProps, expected *BiasesParams) {
	actualLen := len(*expected)
	if len(*result) != actualLen {
		t.Errorf("expected %d results, get %v", actualLen, *result)
	}
	for i, v := range *expected {
		identifier := (*(*result)[i].Bias).Identifier()
		if identifier != v.Name {
			t.Errorf("Invalid bias at %d: %s, expected %s", i, v.Name, identifier)
		}
	}
}
