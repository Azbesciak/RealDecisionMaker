package limited_rationality

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type HeuristicParams interface {
	GetCurrentChoice() string
	GetRandomSeed() int64
}

func GetAlternativesSearchOrder(
	dm *model.DecisionMakingParams,
	params HeuristicParams,
	generator utils.ValueGenerator,
) (model.AlternativeWithCriteria, []model.AlternativeWithCriteria) {
	if len(params.GetCurrentChoice()) > 0 {
		allAlternatives := dm.AllAlternatives()
		choice := model.FetchAlternative(&allAlternatives, params.GetCurrentChoice())
		leftAlternatives := model.RemoveAlternative(dm.ConsideredAlternatives, choice)
		return choice, *model.ShuffleAlternatives(&leftAlternatives, generator)
	} else {
		alternatives := *model.ShuffleAlternatives(&dm.ConsideredAlternatives, generator)
		return alternatives[0], alternatives[1:]
	}
}

func PrepareSequentialRanking(result model.AlternativeResults, resultIds []model.Alternative) model.AlternativesRanking {
	resultsCount := len(result)
	ranking := make(model.AlternativesRanking, resultsCount)
	for i, r := range result {
		ranking[i] = model.AlternativesRankEntry{
			AlternativeResult:  r,
			BetterThanOrSameAs: resultIds[i+1:],
		}
	}
	return ranking
}
