package electreIII

import (
	. "github.com/Azbesciak/RealDecisionMaker/model"
	. "github.com/Azbesciak/RealDecisionMaker/testUtils"
	"testing"
)

func TestEvaluateRanking_PowerStation(t *testing.T) {
	ascending := &[]int{2, 3, 1, 3, 3}
	descending := &[]int{3, 4, 2, 5, 1}
	alts := AlternativeResults{
		DummyAlternative("ITA", 0),
		DummyAlternative("BEL", 0),
		DummyAlternative("GER", 0),
		DummyAlternative("AUT", 0),
		DummyAlternative("FRA", 0),
	}
	altsMap := AlternativesResultToMap(&alts)
	ranking := EvaluateRanking(ascending, descending, ExtractAlternativesFromResults(&alts))
	expectedRanking := AlternativesRanking{
		DummyRankingEntry(altsMap, "ITA", "BEL", "AUT"),
		DummyRankingEntry(altsMap, "BEL", "AUT"),
		DummyRankingEntry(altsMap, "GER", "ITA", "BEL", "AUT"),
		DummyRankingEntry(altsMap, "AUT"),
		DummyRankingEntry(altsMap, "FRA", "BEL", "AUT"),
	}
	CompareSize(&expectedRanking, ranking, t)
	CompareRankings(&expectedRanking, ranking, t)
}

func TestEvaluateRanking_Example(t *testing.T) {
	ascending := &[]int{1, 2, 3, 3, 3}
	descending := &[]int{1, 1, 3, 2, 3}
	alts := AlternativeResults{
		DummyAlternative("1", 0),
		DummyAlternative("2", 0),
		DummyAlternative("3", 0),
		DummyAlternative("4", 0),
		DummyAlternative("5", 0),
	}
	altsMap := AlternativesResultToMap(&alts)
	ranking := EvaluateRanking(ascending, descending, ExtractAlternativesFromResults(&alts))
	expectedRanking := AlternativesRanking{
		DummyRankingEntry(altsMap, "1", "2", "3", "4", "5"),
		DummyRankingEntry(altsMap, "2", "3", "4", "5"),
		DummyRankingEntry(altsMap, "3", "5"),
		DummyRankingEntry(altsMap, "4", "3", "5"),
		DummyRankingEntry(altsMap, "5", "3"),
	}
	CompareSize(&expectedRanking, ranking, t)
	CompareRankings(&expectedRanking, ranking, t)
}
