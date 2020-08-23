package criteria_ordering

import (
	"fmt"
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
