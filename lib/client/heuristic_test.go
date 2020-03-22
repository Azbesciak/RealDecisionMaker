package client

import (
	. "github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

type MockHeuristic struct {
	Name string
}

func (m *MockHeuristic) Identifier() string {
	return m.Name
}
func (m *MockHeuristic) Apply(dm *DecisionMakingParams, props *HeuristicProps, listener *HeuristicListener) *HeuristicResult {
	panic("mock")
}

func TestHeuristicsSelection(t *testing.T) {
	availableHeuristics := AsHeuristicsMap(&Heuristics{
		&MockHeuristic{"a"},
		&MockHeuristic{"c"},
		&MockHeuristic{"b"},
	})
	params := &HeuristicsParams{
		HeuristicParams{Name: "b"},
		HeuristicParams{Name: "a"},
	}
	result := ChooseHeuristics(availableHeuristics, params)
	checkResult(t, result, params)
}

func TestHeuristicsFiltering(t *testing.T) {
	availableHeuristics := AsHeuristicsMap(&Heuristics{
		&MockHeuristic{"a"},
		&MockHeuristic{"c"},
		&MockHeuristic{"b"},
	})
	params := &HeuristicsParams{
		HeuristicParams{Name: "b", Disabled: true},
		HeuristicParams{Name: "c"},
		HeuristicParams{Name: "a"},
	}
	expected := &HeuristicsParams{
		HeuristicParams{Name: "c"},
		HeuristicParams{Name: "a"},
	}
	result := ChooseHeuristics(availableHeuristics, params)
	checkResult(t, result, expected)
}

func TestHeuristicNotFound(t *testing.T) {
	availableHeuristics := AsHeuristicsMap(&Heuristics{
		&MockHeuristic{"b"},
	})
	params := &HeuristicsParams{
		HeuristicParams{Name: "a"},
	}
	defer utils.ExpectError(t, "heuristic 'a' not found, available are '[b]'")()
	ChooseHeuristics(availableHeuristics, params)
}

func TestSkippingDisabledNotExistingHeuristic(t *testing.T) {
	availableHeuristics := AsHeuristicsMap(&Heuristics{
		&MockHeuristic{"b"},
	})
	params := &HeuristicsParams{
		HeuristicParams{Name: "a", Disabled: true},
	}
	result := ChooseHeuristics(availableHeuristics, params)
	checkResult(t, result, &HeuristicsParams{})
}

func checkResult(t *testing.T, result *HeuristicsWithProps, expected *HeuristicsParams) {
	actualLen := len(*expected)
	if len(*result) != actualLen {
		t.Errorf("expected %d results, get %v", actualLen, *result)
	}
	for i, v := range *expected {
		identifier := (*(*result)[i].Heuristic).Identifier()
		if identifier != v.Name {
			t.Errorf("Invalid heuristic at %d: %s, expected %s", i, v.Name, identifier)
		}
	}
}
