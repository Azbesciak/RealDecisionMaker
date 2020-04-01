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
