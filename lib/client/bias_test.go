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

type MockProps struct {
	Id string
}

func TestBiasesSelection(t *testing.T) {
	availableBiases := AsBiasesMap(&Biases{
		&MockBias{"a"},
		&MockBias{"c"},
		&MockBias{"b"},
	})
	params := &BiasesParams{
		biasProps("b"),
		biasProps("a"),
	}
	result := ChooseBiases(availableBiases, params)
	checkResult(t, result, params)
}

func biasProps(name string) interface{} {
	return BiasParams{
		Name:             name,
		ApplyProbability: 1,
	}
}

func disabledBiasProps(name string) interface{} {
	return BiasParams{
		Name:             name,
		ApplyProbability: 1,
		Disabled:         true,
	}
}

func TestBiasesFiltering(t *testing.T) {
	availableBiases := AsBiasesMap(&Biases{
		&MockBias{"a"},
		&MockBias{"c"},
		&MockBias{"b"},
	})
	params := &BiasesParams{
		disabledBiasProps("b"),
		biasProps("c"),
		biasProps("a"),
	}
	expected := &BiasesParams{
		biasProps("c"),
		biasProps("a"),
	}
	result := ChooseBiases(availableBiases, params)
	checkResult(t, result, expected)
}

func TestBiasNotFound(t *testing.T) {
	availableBiases := AsBiasesMap(&Biases{
		&MockBias{"b"},
	})
	params := &BiasesParams{
		biasProps("a"),
	}
	defer utils.ExpectError(t, "bias 'a' not found, available are '[b]'")()
	ChooseBiases(availableBiases, params)
}

func TestSkippingDisabledNotExistingBiases(t *testing.T) {
	availableBiases := AsBiasesMap(&Biases{
		&MockBias{"b"},
	})
	params := &BiasesParams{
		disabledBiasProps("a"),
	}
	result := ChooseBiases(availableBiases, params)
	checkResult(t, result, &BiasesParams{})
}

func TestValidPropsAssigned(t *testing.T) {
	availableBiases := AsBiasesMap(&Biases{
		&MockBias{"a"},
		&MockBias{"c"},
		&MockBias{"b"},
	})
	params := &BiasesParams{
		interface{}(BiasParams{
			Name:             "c",
			ApplyProbability: 1,
			Props:            MockProps{Id: "c"},
		}),
		interface{}(BiasParams{
			Name:             "a",
			ApplyProbability: 1,
			Props:            MockProps{Id: "a"},
		}),
	}
	expected := &BiasesParams{
		interface{}(BiasParams{
			Name:             "c",
			ApplyProbability: 1,
			Props:            MockProps{Id: "c"},
		}),
		interface{}(BiasParams{
			Name:             "a",
			ApplyProbability: 1,
			Props:            MockProps{Id: "a"},
		}),
	}
	result := ChooseBiases(availableBiases, params)
	checkResult(t, result, expected)
}

func checkResult(t *testing.T, result *BiasesWithProps, expected *BiasesParams) {
	actualLen := len(*expected)
	if len(*result) != actualLen {
		t.Errorf("expected %d results, get %v", actualLen, *result)
	}
	for i, v := range *expected {
		expectedBias := v.(BiasParams)
		actual := (*result)[i]
		identifier := (*actual.Bias).Identifier()
		if identifier != expectedBias.Name {
			t.Errorf("Invalid bias at %d: %s, expected %s", i, expectedBias.Name, identifier)
		}
		if expectedBias.Props == nil && actual.Props.Props == nil {
			continue
		} else if expectedBias.Props == nil {
			t.Errorf("result props for %d (%s) is not nil", i, identifier)
		} else if actual.Props.Props == nil {
			t.Errorf("result props for %d (%s) is nil", i, identifier)
		} else {
			ac := actual.Props.Props.(MockProps)
			ex := expectedBias.Props.(MockProps)
			if ac.Id != ex.Id {
				t.Errorf("expected mock props id == %s, got %s for prop %d (%s)", ex.Id, ac.Id, i, identifier)
			}
		}

	}
}
