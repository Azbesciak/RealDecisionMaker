package criteria_ordering

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type CriteriaOrderingResolver interface {
	utils.Identifiable
	OrderCriteria(
		params *model.DecisionMakingParams,
		props *model.BiasProps,
		listener *model.BiasListener,
	) *model.Criteria
}
