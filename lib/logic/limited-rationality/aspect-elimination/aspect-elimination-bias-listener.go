package aspect_elimination

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/limited-rationality/satisfaction-levels"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type AspectEliminationBiasListener struct {
	satisfactionLevelsUpdateListeners satisfaction_levels.SatisfactionLevelsUpdateListeners
}

func NewAspectEliminationBiasListener(
	satisfactionLevelsUpdateListeners satisfaction_levels.SatisfactionLevelsUpdateListeners,
) *AspectEliminationBiasListener {
	return &AspectEliminationBiasListener{
		satisfactionLevelsUpdateListeners: satisfactionLevelsUpdateListeners,
	}
}

func (a *AspectEliminationBiasListener) Identifier() string {
	return methodName
}

//go:generate easytags $GOFILE json:camel
type aspectEliminationAddedCriterion struct {
	Weights model.Weights `json:"weights"`
	Params  interface{}   `json:"params,omitempty"`
}

func (a *AspectEliminationBiasListener) OnCriterionAdded(
	criterion *model.Criterion,
	referenceCriterion *model.Criterion,
	params model.MethodParameters,
	generator utils.ValueGenerator,
) model.AddedCriterionParams {
	pParams := params.(AspectEliminationHeuristicParams)
	newValue := model.NewCriterionValue(&pParams.Weights, referenceCriterion, &generator)
	listener, methodParams := a.getMethodParams(pParams)
	addedParams := listener.OnCriterionAdded(criterion, referenceCriterion, methodParams, generator)
	return aspectEliminationAddedCriterion{
		Weights: model.Weights{criterion.Id: newValue},
		Params:  addedParams,
	}
}

func (a *AspectEliminationBiasListener) getMethodParams(pParams AspectEliminationHeuristicParams) (satisfaction_levels.SatisfactionLevelsUpdateListener, satisfaction_levels.SatisfactionLevels) {
	return a.satisfactionLevelsUpdateListeners.Get(pParams.Function, pParams.Params)
}

func (a *AspectEliminationBiasListener) OnCriteriaRemoved(
	leftCriteria *model.Criteria,
	params model.MethodParameters,
) model.MethodParameters {
	pParams := params.(AspectEliminationHeuristicParams)
	listener, methodParams := a.getMethodParams(pParams)
	afterRemoveParams := listener.OnCriteriaRemoved(leftCriteria, methodParams)
	return pParams.with(afterRemoveParams, pParams.Weights.PreserveOnly(leftCriteria))
}

func (a *AspectEliminationBiasListener) RankCriteriaAscending(params *model.DecisionMakingParams) *model.WeightedCriteria {
	wParams := params.MethodParameters.(AspectEliminationHeuristicParams)
	return params.Criteria.SortByWeights(wParams.Weights)
}

func (a *AspectEliminationBiasListener) Merge(params model.MethodParameters, addition model.MethodParameters) model.MethodParameters {
	pParams := params.(AspectEliminationHeuristicParams)
	aParams := addition.(aspectEliminationAddedCriterion)
	listener, methodParams := a.getMethodParams(pParams)
	return pParams.with(listener.Merge(methodParams, aParams.Params), pParams.Weights.Merge(&aParams.Weights))
}
