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
	referenceCriterion *model.Criterion,
	params model.MethodParameters,
	generator utils.ValueGenerator,
) model.AddedCriterionParams {
	wParams := params.(MajorityHeuristicParams)
	newWeight := model.NewCriterionValue(&wParams.Weights, referenceCriterion, &generator)
	return model.SingleWeight(criterion, newWeight)
}

func (m *MajorityBiasListener) OnCriteriaRemoved(leftCriteria *model.Criteria, params model.MethodParameters) model.MethodParameters {
	wParams := params.(MajorityHeuristicParams)
	leftWeights := wParams.Weights.PreserveOnly(leftCriteria)
	return MajorityHeuristicParams{
		Weights:                    *leftWeights,
		CurrentChoice:              wParams.CurrentChoice,
		RandomSeed:                 wParams.RandomSeed,
		RandomAlternativesOrdering: wParams.RandomAlternativesOrdering,
	}
}

func (m *MajorityBiasListener) RankCriteriaAscending(params *model.DecisionMakingParams) *model.WeightedCriteria {
	wParams := params.MethodParameters.(MajorityHeuristicParams)
	return params.Criteria.SortByWeights(wParams.Weights)
}

func (m *MajorityBiasListener) Merge(params model.MethodParameters, addition model.MethodParameters) model.MethodParameters {
	oldParams := params.(MajorityHeuristicParams)
	newParams := addition.(model.WeightType)
	return MajorityHeuristicParams{
		Weights:                    *oldParams.Weights.Merge(&newParams.Weights),
		CurrentChoice:              oldParams.CurrentChoice,
		RandomSeed:                 oldParams.RandomSeed,
		RandomAlternativesOrdering: oldParams.RandomAlternativesOrdering,
	}
}
