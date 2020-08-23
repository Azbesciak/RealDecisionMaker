package criteria_ordering

import "github.com/Azbesciak/RealDecisionMaker/lib/model"

const StrongestCriteriaFirst = "strongest"

type StrongestCriteriaOrderingResolver struct {
}

func (w *StrongestCriteriaOrderingResolver) Identifier() string {
	return StrongestCriteriaFirst
}

func (w *StrongestCriteriaOrderingResolver) OrderCriteria(
	params *model.DecisionMakingParams,
	_ *model.BiasProps,
	listener *model.BiasListener,
) *model.Criteria {
	ascending := (*listener).RankCriteriaAscending(params).Criteria()
	totalCount := len(*ascending)
	descending := make(model.Criteria, totalCount)
	for i, a := range *ascending {
		descending[totalCount-i-1] = a
	}
	return &descending
}
