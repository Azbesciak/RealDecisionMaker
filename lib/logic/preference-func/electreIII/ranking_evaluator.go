package electreIII

import (
	. "github.com/Azbesciak/RealDecisionMaker/lib/model"
)

//go:generate easytags $GOFILE json:camel
type ElectreIIIEvaluation struct {
	AscendingIndex  int `json:"ascendingIndex"`
	DescendingIndex int `json:"descendingIndex"`
}

func EvaluateRanking(ascending, descending *[]int, alternatives *[]AlternativeWithCriteria) *AlternativesRanking {
	ranking := make(AlternativesRanking, 0)
	for ia, alt1Asc := range *ascending {
		alt1Desc := (*descending)[ia]
		betterOrSameAs := make([]Alternative, 0)

		for ib, alt2Asc := range *ascending {
			if ia == ib {
				continue
			}
			alt2Desc := (*descending)[ib]
			if alt1Asc <= alt2Asc && alt1Desc <= alt2Desc {
				betterOrSameAs = append(betterOrSameAs, (*alternatives)[ib].Id)
			}
		}
		ranking = append(ranking, AlternativesRankEntry{
			AlternativeResult: AlternativeResult{Alternative: (*alternatives)[ia], Evaluation: ElectreIIIEvaluation{
				AscendingIndex:  alt1Asc,
				DescendingIndex: alt1Desc,
			}},
			BetterThanOrSameAs: betterOrSameAs,
		})
	}
	return &ranking
}
