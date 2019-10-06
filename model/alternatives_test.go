package model

import (
	"reflect"
	"testing"
)

func TestAlternativeResults_Ranking(t *testing.T) {
	alts := AlternativeResults{
		dummyAlternative("1", 5),
		dummyAlternative("2", 3),
		dummyAlternative("3", 5),
		dummyAlternative("7", 2),
		dummyAlternative("4", 1),
		dummyAlternative("0", 3),
	}
	ranking := alts.Ranking()
	altsMap := alts.toMap()
	expectedRanking := AlternativesRanking{
		dummyRankingEntry(altsMap, "1", "3", "0", "2"),
		dummyRankingEntry(altsMap, "3", "1", "0", "2"),
		dummyRankingEntry(altsMap, "0", "2", "7"),
		dummyRankingEntry(altsMap, "2", "0", "7"),
		dummyRankingEntry(altsMap, "7", "4"),
		dummyRankingEntry(altsMap, "4"),
	}
	compareSize(&expectedRanking, ranking, t)
	compareRankings(&expectedRanking, ranking, t)
}

func compareSize(expectedRanking, ranking *AlternativesRanking, t *testing.T) {
	expectedRankingSize := len(*expectedRanking)
	receivedRankingSize := len(*ranking)
	if receivedRankingSize != expectedRankingSize {
		t.Errorf("Expected ranking of length %d , got %d", expectedRankingSize, receivedRankingSize)
	}
}

func compareRankings(expected, received *AlternativesRanking, t *testing.T) {
	for i, e := range *expected {
		rec := (*received)[i]
		if e.Alternative.Id != rec.Alternative.Id {
			t.Errorf("Expected id of '%s' at position %d, got '%s'", e.Alternative.Id, i, rec.Alternative.Id)
		}
		if !reflect.DeepEqual(e.BetterThanOrSameAs, rec.BetterThanOrSameAs) {
			t.Errorf(
				"Invalid Preferrence of id '%s' at position %d, expected '%s', got '%s'",
				e.Alternative.Id, i, e.BetterThanOrSameAs, rec.BetterThanOrSameAs,
			)
		}
	}
}

type AltsMap = *map[string]AlternativeResult

func (a *AlternativeResults) toMap() AltsMap {
	altsById := map[string]AlternativeResult{}
	for _, alt := range *a {
		altsById[alt.Alternative.Id] = alt
	}
	return &altsById
}

func dummyRankingEntry(alts AltsMap, thisAlt string, betterThanOrSameAs ...Alternative) AlternativesRankEntry {
	var betThOrSamAs = make(Alternatives, len(betterThanOrSameAs))
	for i, a := range betterThanOrSameAs {
		betThOrSamAs[i] = a
	}
	alt := (*alts)[thisAlt]
	return AlternativesRankEntry{
		AlternativeResult:  alt,
		BetterThanOrSameAs: betThOrSamAs,
	}
}

func dummyAlternative(id string, value Weight) AlternativeResult {
	return AlternativeResult{
		Alternative: AlternativeWithCriteria{
			Id:       id,
			Criteria: nil,
		},
		Value: value,
	}
}
