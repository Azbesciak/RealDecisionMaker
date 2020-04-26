package limited_rationality

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type HeuristicParams interface {
	CurrentChoice() string
	RandomSeed() int64
}

func GetAlternativesSearchOrder(
	dm *model.DecisionMakingParams,
	params HeuristicParams,
	generator utils.ValueGenerator,
) (model.AlternativeWithCriteria, []model.AlternativeWithCriteria) {
	if len(params.CurrentChoice()) > 0 {
		allAlternatives := dm.AllAlternatives()
		choice := model.FetchAlternative(&allAlternatives, params.CurrentChoice())
		leftAlternatives := model.RemoveAlternative(dm.ConsideredAlternatives, choice)
		return choice, *model.ShuffleAlternatives(&leftAlternatives, generator)
	} else {
		alternatives := *model.ShuffleAlternatives(&dm.ConsideredAlternatives, generator)
		return alternatives[0], alternatives[1:]
	}
}
