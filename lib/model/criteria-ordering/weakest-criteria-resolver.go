package criteria_ordering

import "github.com/Azbesciak/RealDecisionMaker/lib/model"

const WeakestCriteriaFirst = "weakest"

type WeakestCriteriaOrderingResolver struct {
}

func (w *WeakestCriteriaOrderingResolver) Identifier() string {
	return WeakestCriteriaFirst
}

func (w *WeakestCriteriaOrderingResolver) OrderCriteria(
	params *model.DecisionMakingParams,
	_ *model.BiasProps,
	listener *model.BiasListener,
) *model.Criteria {
	return (*listener).RankCriteriaAscending(params).Criteria()
}
