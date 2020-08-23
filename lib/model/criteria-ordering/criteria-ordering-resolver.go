package criteria_ordering

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"math"
)

type CriteriaOrderingResolver interface {
	utils.Identifiable
	OrderCriteria(
		params *model.DecisionMakingParams,
		props *model.BiasProps,
		listener *model.BiasListener,
	) *model.Criteria
}

func FetchOrderingResolver(resolvers *[]CriteriaOrderingResolver, resolver string) CriteriaOrderingResolver {
	if len(resolver) == 0 {
		return (*resolvers)[0]
	}
	for _, r := range *resolvers {
		if r.Identifier() == resolver {
			return r
		}
	}
	names := make([]string, len(*resolvers))
	for i, r := range *resolvers {
		names[i] = r.Identifier()
	}
	panic(fmt.Errorf("omission order resolver '%s' not found in %v", resolver, names))
}

func SplitCriteriaByOrdering(ratio float64, sortedCriteria *model.Criteria) *CriteriaPartition {
	criteriaCount := len(*sortedCriteria)
	pivot := int(math.Floor(float64(criteriaCount) * ratio))
	left := (*sortedCriteria)[0:pivot]
	right := (*sortedCriteria)[pivot:]
	return &CriteriaPartition{
		Left:  &right,
		Right: &left,
	}
}

type CriteriaPartition struct {
	Left  *model.Criteria
	Right *model.Criteria
}
