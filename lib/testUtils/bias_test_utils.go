package testUtils

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"strconv"
	"testing"
)

func TestBiasRanking(index int, t *testing.T, bias model.BiasListener, params *model.DecisionMakingParams, expectedOrder []string) {
	ranked := bias.RankCriteriaAscending(params)
	if len(*ranked) != len(expectedOrder) {
		t.Errorf("lengths differ, expected %d, got %d", len(expectedOrder), len(*ranked))
		return
	}
	for i, expectedId := range expectedOrder {
		actualId := (*ranked)[i].Criterion.Identifier()
		if actualId != expectedId {
			t.Errorf("%d: expected '%s' at %d index, got '%s', alternatives: %v",
				index, expectedId, i, actualId, params.ConsideredAlternatives)
		}
	}
}

func WrapAlternatives(alternativesWeights []model.Weights) []model.AlternativeWithCriteria {
	alternatives := make([]model.AlternativeWithCriteria, len(alternativesWeights))
	for i, w := range alternativesWeights {
		alternatives[i] = model.AlternativeWithCriteria{
			Id:       strconv.Itoa(i),
			Criteria: w,
		}
	}
	return alternatives
}

func CheckCount(t *testing.T, typ string, expected []string, actual *model.Criteria) {
	expectedLen := len(expected)
	actualLen := len(*actual)
	if actualLen != expectedLen {
		t.Errorf("expected %d criteria to %s, but got %d", expectedLen, typ, actualLen)
		return
	}
	for i, expectedId := range expected {
		actualId := (*actual).Get(i).Identifier()
		if actualId != expectedId {
			t.Errorf("expected '%s' at %d in %s criteria, got '%s'", expectedId, i, typ, actualId)
		}
	}
}
