package owa

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"testing"
)

var owaTestCriteria = model.Criteria{{Id: "1"}, {Id: "2"}, {Id: "3"}}

func TestOwaHeuristic_RankCriteriaAscending(t *testing.T) {
	testRanking(t, []model.AlternativeWithCriteria{}, []string{"1", "2", "3"})
	testRanking(t, []model.AlternativeWithCriteria{
		{Id: "x", Criteria: model.Weights{"1": 1, "2": 2, "3": 3}},
	}, []string{"1", "2", "3"})
	testRanking(t, []model.AlternativeWithCriteria{
		{Id: "x", Criteria: model.Weights{"1": 1, "2": 2, "3": 3}},
		{Id: "x", Criteria: model.Weights{"1": 2, "2": 1, "3": 3}},
	}, []string{"1", "2", "3"})
	testRanking(t, []model.AlternativeWithCriteria{
		{Id: "x", Criteria: model.Weights{"1": 1, "2": 2, "3": 3}},
		{Id: "x", Criteria: model.Weights{"1": 2, "2": 1, "3": 0}},
	}, []string{"1", "2", "3"})
	testRanking(t, []model.AlternativeWithCriteria{
		{Id: "x", Criteria: model.Weights{"1": 1, "2": 2, "3": 3}},
		{Id: "x", Criteria: model.Weights{"1": 2, "2": 0, "3": 0}},
	}, []string{"2", "1", "3"})
	testRanking(t, []model.AlternativeWithCriteria{
		{Id: "x", Criteria: model.Weights{"1": 1, "2": 2, "3": 3}},
		{Id: "x", Criteria: model.Weights{"1": 2, "2": 0, "3": -2}},
	}, []string{"3", "2", "1"})
}

func testRanking(t *testing.T, alternatives []model.AlternativeWithCriteria, expectedOrder []string) {
	heu := OwaHeuristic{}
	ranked := heu.RankCriteriaAscending(&model.DecisionMakingParams{
		NotConsideredAlternatives: nil,
		ConsideredAlternatives:    alternatives,
		Criteria:                  owaTestCriteria,
		MethodParameters:          nil,
	})
	if ranked.Len() != len(expectedOrder) {
		t.Errorf("lengths differ, expected %d, got %d", len(expectedOrder), ranked.Len())
		return
	}
	for i, expectedId := range expectedOrder {
		actualId := ranked.Get(i).Identifier()
		if actualId != expectedId {
			t.Errorf("expected '%s' at %d index, got '%s', alternatives: %v", expectedId, i, actualId, alternatives)
		}
	}
}
