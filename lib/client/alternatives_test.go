package client

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	. "github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"testing"
)

func TestAlternativeResults_Ranking(t *testing.T) {
	alts := model.AlternativeResults{
		DummyAlternative("1", 5),
		DummyAlternative("2", 3),
		DummyAlternative("3", 5),
		DummyAlternative("7", 2),
		DummyAlternative("4", 1),
		DummyAlternative("0", 3),
	}
	ranking := alts.Ranking()
	altsMap := AlternativesResultToMap(&alts)
	expectedRanking := model.AlternativesRanking{
		DummyRankingEntry(altsMap, "1", "3", "0", "2"),
		DummyRankingEntry(altsMap, "3", "1", "0", "2"),
		DummyRankingEntry(altsMap, "0", "2", "7"),
		DummyRankingEntry(altsMap, "2", "0", "7"),
		DummyRankingEntry(altsMap, "7", "4"),
		DummyRankingEntry(altsMap, "4"),
	}
	CompareRankings(&expectedRanking, ranking, t)
}
