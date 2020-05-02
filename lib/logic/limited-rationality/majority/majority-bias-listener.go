package majority

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type MajorityBiasListener struct {
}

func (m *MajorityBiasListener) Identifier() string {
	return methodName
}

func (m *MajorityBiasListener) OnCriterionAdded(
	criterion *model.Criterion,
	previousRankedCriteria *model.Criteria,
	params model.MethodParameters,
	generator utils.ValueGenerator,
) model.AddedCriterionParams {
	wParams := params.(MajorityHeuristicParams)
	leastImportantParam := wParams.Weights[previousRankedCriteria.First().Identifier()]
	newWeight := generator() * leastImportantParam
	return model.WeightedCriterion{
		Criterion: *criterion,
		Weight:    newWeight,
	}
}

func (m *MajorityBiasListener) OnCriteriaRemoved(removedCriteria *model.Criteria, leftCriteria *model.Criteria, params model.MethodParameters) model.MethodParameters {
	wParams := params.(MajorityHeuristicParams)
	leftWeights := wParams.Weights.PreserveOnly(leftCriteria)
	return MajorityHeuristicParams{
		Weights: *leftWeights,
		Current: wParams.Current,
		Seed:    wParams.Seed,
	}
}

func (m *MajorityBiasListener) RankCriteriaAscending(params *model.DecisionMakingParams) *model.Criteria {
	wParams := params.MethodParameters.(MajorityHeuristicParams)
	return params.Criteria.SortByWeights(wParams.Weights)
}

func (m *MajorityBiasListener) Merge(params model.MethodParameters, addition model.MethodParameters) model.MethodParameters {
	oldParams := params.(MajorityHeuristicParams)
	newParams := addition.(model.WeightedCriterion)
	return MajorityHeuristicParams{
		Weights: *oldParams.Weights.Merge(newParams.AsWeights()),
		Current: oldParams.Current,
		Seed:    oldParams.Seed,
	}
}
