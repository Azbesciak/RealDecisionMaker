package electreIII

import (
	. "github.com/Azbesciak/RealDecisionMaker/lib/model"
	. "github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"testing"
)

func electreIIIRankingEntry(id string, ascending, descending int) AlternativeResult {
	return AlternativeResult{
		Alternative: AlternativeWithCriteria{Id: id},
		Evaluation: ElectreIIIEvaluation{
			AscendingIndex:  ascending,
			DescendingIndex: descending,
		},
	}
}

func electreIIIResults(ids *[]string, ascending, descending *[]int) AlternativeResults {
	result := make(AlternativeResults, len(*ids))
	for i, id := range *ids {
		result[i] = electreIIIRankingEntry(id, (*ascending)[i], (*descending)[i])
	}
	return result
}

func TestEvaluateRanking_PowerStation(t *testing.T) {
	ascending := &[]int{2, 3, 1, 3, 3}
	descending := &[]int{3, 4, 2, 5, 1}
	ids := &[]string{"ITA", "BEL", "GER", "AUT", "FRA"}
	alts := electreIIIResults(ids, ascending, descending)
	altsMap := AlternativesResultToMap(&alts)
	ranking := EvaluateRanking(ascending, descending, ExtractAlternativesFromResults(&alts))
	expectedRanking := AlternativesRanking{
		DummyRankingEntry(altsMap, "ITA", "BEL", "AUT"),
		DummyRankingEntry(altsMap, "BEL", "AUT"),
		DummyRankingEntry(altsMap, "GER", "ITA", "BEL", "AUT"),
		DummyRankingEntry(altsMap, "AUT"),
		DummyRankingEntry(altsMap, "FRA", "BEL", "AUT"),
	}
	CompareRankings(&expectedRanking, ranking, t)
}

func TestEvaluateRanking_Example(t *testing.T) {
	ascending := &[]int{1, 2, 3, 3, 3}
	descending := &[]int{1, 1, 3, 2, 3}
	ids := &[]string{"1", "2", "3", "4", "5"}
	alts := electreIIIResults(ids, ascending, descending)
	altsMap := AlternativesResultToMap(&alts)
	ranking := EvaluateRanking(ascending, descending, ExtractAlternativesFromResults(&alts))
	expectedRanking := AlternativesRanking{
		DummyRankingEntry(altsMap, "1", "2", "3", "4", "5"),
		DummyRankingEntry(altsMap, "2", "3", "4", "5"),
		DummyRankingEntry(altsMap, "3", "5"),
		DummyRankingEntry(altsMap, "4", "3", "5"),
		DummyRankingEntry(altsMap, "5", "3"),
	}
	CompareRankings(&expectedRanking, ranking, t)
}
