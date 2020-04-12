package electreIII

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type ElectreIIIBiasLIstener struct {
}

func (e *ElectreIIIBiasLIstener) Identifier() string {
	return methodName
}

func (e *ElectreIIIBiasLIstener) Merge(params model.MethodParameters, addition model.MethodParameters) model.MethodParameters {
	oldEleParams := params.(electreIIIParams)
	newEleParams := addition.(electreIIIParams)
	newCriteria := make(ElectreCriteria, len(*oldEleParams.criteria)+len(*newEleParams.criteria))
	for c, v := range *oldEleParams.criteria {
		newCriteria[c] = v
	}
	for c, v := range *newEleParams.criteria {
		_, ok := newCriteria[c]
		if ok {
			panic(fmt.Errorf("criterion '%s' already exist in params in electre criteria, merge %v with %v", c, *oldEleParams.criteria, *newEleParams.criteria))
		}
		newCriteria[c] = v
	}
	return electreIIIParams{criteria: &newCriteria, distillationFun: oldEleParams.distillationFun}
}

func (e *ElectreIIIBiasLIstener) OnCriterionAdded(
	criterion *model.Criterion,
	previousRankedCriteria *model.Criteria,
	params model.MethodParameters,
	generator utils.ValueGenerator,
) model.AddedCriterionParams {
	eleParams := params.(electreIIIParams)
	weakestCriterion := (*eleParams.criteria)[previousRankedCriteria.First().Id]
	return electreIIIParams{criteria: &ElectreCriteria{
		criterion.Id: ElectreCriterion{
			K: generator() * weakestCriterion.K,
			Q: weakestCriterion.Q,
			P: weakestCriterion.P,
			V: weakestCriterion.V,
		},
	}}
}

func (e *ElectreIIIBiasLIstener) OnCriteriaRemoved(
	removedCriteria *model.Criteria,
	leftCriteria *model.Criteria,
	params model.MethodParameters,
) model.MethodParameters {
	eleParams := params.(electreIIIParams)
	resCriteria := make(ElectreCriteria, len(*leftCriteria))
	for _, c := range *leftCriteria {
		criterion, ok := (*eleParams.criteria)[c.Id]
		if !ok {
			panic(fmt.Errorf("criterion '%s' not found in electre criteria %v", c.Id, *eleParams.criteria))
		}
		resCriteria[c.Id] = criterion
	}
	return electreIIIParams{criteria: &resCriteria, distillationFun: eleParams.distillationFun}
}

func (e *ElectreIIIBiasLIstener) RankCriteriaAscending(params *model.DecisionMakingParams) *model.Criteria {
	eleParams := params.MethodParameters.(electreIIIParams)
	weights := make(model.Weights, len(*eleParams.criteria))
	for k, v := range *eleParams.criteria {
		weights[k] = v.K
	}
	return params.Criteria.SortByWeights(weights)
}
