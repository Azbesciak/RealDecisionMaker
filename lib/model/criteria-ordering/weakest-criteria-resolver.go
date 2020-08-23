package criteria_ordering

import "github.com/Azbesciak/RealDecisionMaker/lib/model"

const WeakestCriteriaFirst = "weakest"

type WeakestCriteriaOmissionResolver struct {
}

func (w *WeakestCriteriaOmissionResolver) Identifier() string {
	return WeakestCriteriaFirst
}

func (w *WeakestCriteriaOmissionResolver) OrderCriteria(
	params *model.DecisionMakingParams,
	_ *model.BiasProps,
	listener *model.BiasListener,
) *model.Criteria {
	return (*listener).RankCriteriaAscending(params).Criteria()
}
