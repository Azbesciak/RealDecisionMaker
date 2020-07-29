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
		{Name: "b"},
		{Name: "a"},
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
		{Name: "b", Disabled: true},
		{Name: "c"},
		{Name: "a"},
	}
	expected := &BiasesParams{
		{Name: "c"},
		{Name: "a"},
	}
	result := ChooseBiases(availableBiases, params)
	checkResult(t, result, expected)
}

func TestBiasNotFound(t *testing.T) {
	availableBiases := AsBiasesMap(&Biases{
		&MockBias{"b"},
	})
	params := &BiasesParams{
		{Name: "a"},
	}
	defer utils.ExpectError(t, "bias 'a' not found, available are '[b]'")()
	ChooseBiases(availableBiases, params)
}

func TestSkippingDisabledNotExistingBiases(t *testing.T) {
	availableBiases := AsBiasesMap(&Biases{
		&MockBias{"b"},
	})
	params := &BiasesParams{
		{Name: "a", Disabled: true},
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
		{Name: "c", Props: MockProps{Id: "c"}},
		{Name: "a", Props: MockProps{Id: "a"}},
	}
	expected := &BiasesParams{
		{Name: "c", Props: MockProps{Id: "c"}},
		{Name: "a", Props: MockProps{Id: "a"}},
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
		actual := (*result)[i]
		identifier := (*actual.Bias).Identifier()
		if identifier != v.Name {
			t.Errorf("Invalid bias at %d: %s, expected %s", i, v.Name, identifier)
		}
		if v.Props == nil && actual.Props.Props == nil {
			continue
		} else if v.Props == nil {
			t.Errorf("result props for %d (%s) is not nil", i, identifier)
		} else if actual.Props.Props == nil {
			t.Errorf("result props for %d (%s) is nil", i, identifier)
		} else {
			ac := actual.Props.Props.(MockProps)
			ex := v.Props.(MockProps)
			if ac.Id != ex.Id {
				t.Errorf("expected mock props id == %s, got %s for prop %d (%s)", ex.Id, ac.Id, i, identifier)
			}
		}

	}
}
