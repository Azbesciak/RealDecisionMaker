package testUtils

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"strconv"
	"testing"
)

func TestHeuristicRanking(index int, t *testing.T, heu model.HeuristicListener, params *model.DecisionMakingParams, expectedOrder []string) {
	ranked := heu.RankCriteriaAscending(params)
	if ranked.Len() != len(expectedOrder) {
		t.Errorf("lengths differ, expected %d, got %d", len(expectedOrder), ranked.Len())
		return
	}
	for i, expectedId := range expectedOrder {
		actualId := ranked.Get(i).Identifier()
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
