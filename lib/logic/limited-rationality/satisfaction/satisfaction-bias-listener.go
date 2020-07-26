package satisfaction

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/limited-rationality/satisfaction-levels"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type SatisfactionBiasListener struct {
	satisfactionLevelsUpdateListeners satisfaction_levels.SatisfactionLevelsUpdateListeners
}

func (a *SatisfactionBiasListener) Identifier() string {
	return methodName
}

//go:generate easytags $GOFILE json:camel
type satisfactionAddedCriterion struct {
	Params interface{} `json:"params,omitempty"`
}

func (a *SatisfactionBiasListener) OnCriterionAdded(
	criterion *model.Criterion,
	referenceCriterion *model.Criterion,
	params model.MethodParameters,
	generator utils.ValueGenerator,
) model.AddedCriterionParams {
	pParams := params.(SatisfactionParameters)
	listener, methodParams := a.getMethodParams(pParams)
	addedParams := listener.OnCriterionAdded(criterion, referenceCriterion, methodParams, generator)
	return satisfactionAddedCriterion{addedParams}
}

func (a *SatisfactionBiasListener) getMethodParams(pParams SatisfactionParameters) (satisfaction_levels.SatisfactionLevelsUpdateListener, satisfaction_levels.SatisfactionLevels) {
	return a.satisfactionLevelsUpdateListeners.Get(pParams.Function, pParams.Params)
}

func (a *SatisfactionBiasListener) OnCriteriaRemoved(
	leftCriteria *model.Criteria,
	params model.MethodParameters,
) model.MethodParameters {
	pParams := params.(SatisfactionParameters)
	listener, methodParams := a.getMethodParams(pParams)
	afterRemoveParams := listener.OnCriteriaRemoved(leftCriteria, methodParams)
	return pParams.with(afterRemoveParams)
}

func (a *SatisfactionBiasListener) RankCriteriaAscending(params *model.DecisionMakingParams) *model.WeightedCriteria {
	weights := model.PrepareCumulatedWeightsMap(params, model.WeightIdentity)
	return params.Criteria.SortByWeights(*weights)
}

func (a *SatisfactionBiasListener) Merge(params model.MethodParameters, addition model.MethodParameters) model.MethodParameters {
	pParams := params.(SatisfactionParameters)
	aParams := addition.(satisfactionAddedCriterion)
	listener, methodParams := a.getMethodParams(pParams)
	return pParams.with(listener.Merge(methodParams, aParams.Params))
}
