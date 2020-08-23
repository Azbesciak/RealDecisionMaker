package criteria_ordering

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
)

const StrongestByProbabilityCriteriaFirst = "strongestByProbability"

type StrongestByProbabilityCriteriaOrderingResolver struct {
	WeakestByProbability *WeakestByProbabilityCriteriaOrderingResolver
}

func (s *StrongestByProbabilityCriteriaOrderingResolver) Identifier() string {
	return StrongestByProbabilityCriteriaFirst
}

func (s *StrongestByProbabilityCriteriaOrderingResolver) OrderCriteria(
	params *model.DecisionMakingParams,
	props *model.BiasProps,
	listener *model.BiasListener,
) *model.Criteria {
	criteria := s.WeakestByProbability.OrderCriteria(params, props, listener)
	criteriaCount := len(*criteria)
	result := make(model.Criteria, criteriaCount)
	for i, c := range *criteria {
		result[criteriaCount-i-1] = c
	}
	return &result
}
