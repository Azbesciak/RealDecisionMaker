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
	newCriteria := make(ElectreCriteria, len(*oldEleParams.Criteria)+len(*newEleParams.Criteria))
	for c, v := range *oldEleParams.Criteria {
		newCriteria[c] = v
	}
	for c, v := range *newEleParams.Criteria {
		_, ok := newCriteria[c]
		if ok {
			panic(fmt.Errorf("criterion '%s' already exist in params in electre Criteria, merge %v with %v", c, *oldEleParams.Criteria, *newEleParams.Criteria))
		}
		newCriteria[c] = v
	}
	return electreIIIParams{Criteria: &newCriteria, DistillationFun: oldEleParams.DistillationFun}
}

func (e *ElectreIIIBiasLIstener) OnCriterionAdded(
	criterion *model.Criterion,
	previousRankedCriteria *model.Criteria,
	params model.MethodParameters,
	generator utils.ValueGenerator,
) model.AddedCriterionParams {
	eleParams := params.(electreIIIParams)
	weakestCriterion := (*eleParams.Criteria)[previousRankedCriteria.First().Id]
	return electreIIIParams{Criteria: &ElectreCriteria{
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
		criterion, ok := (*eleParams.Criteria)[c.Id]
		if !ok {
			panic(fmt.Errorf("criterion '%s' not found in electre Criteria %v", c.Id, *eleParams.Criteria))
		}
		resCriteria[c.Id] = criterion
	}
	return electreIIIParams{Criteria: &resCriteria, DistillationFun: eleParams.DistillationFun}
}

func (e *ElectreIIIBiasLIstener) RankCriteriaAscending(params *model.DecisionMakingParams) *model.Criteria {
	eleParams := params.MethodParameters.(electreIIIParams)
	weights := make(model.Weights, len(*eleParams.Criteria))
	for k, v := range *eleParams.Criteria {
		weights[k] = v.K
	}
	return params.Criteria.SortByWeights(weights)
}
