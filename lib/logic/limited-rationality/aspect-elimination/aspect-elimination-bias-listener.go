package aspect_elimination

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/limited-rationality/satisfaction-levels"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type AspectEliminationBiasListener struct {
	satisfactionLevelsUpdateListeners satisfaction_levels.SatisfactionLevelsUpdateListeners
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
	previousRankedCriteria *model.Criteria,
	params model.MethodParameters,
	generator utils.ValueGenerator,
) model.AddedCriterionParams {
	pParams := params.(AspectEliminationHeuristicParams)
	newValue := model.NewCriterionValue(&pParams.Weights, previousRankedCriteria, &generator)
	listener, methodParams := a.getMethodParams(pParams)
	addedParams := listener.OnCriterionAdded(criterion, previousRankedCriteria, methodParams, generator)
	return aspectEliminationAddedCriterion{
		Weights: model.Weights{criterion.Id: newValue},
		Params:  addedParams,
	}
}

func (a *AspectEliminationBiasListener) getMethodParams(pParams AspectEliminationHeuristicParams) (satisfaction_levels.SatisfactionLevelsUpdateListener, satisfaction_levels.SatisfactionLevels) {
	listener := *a.satisfactionLevelsUpdateListeners.Fetch(pParams.Function)
	methodParams := listener.BlankParams()
	utils.DecodeToStruct(pParams.Params, &methodParams)
	return listener, methodParams
}

func (a *AspectEliminationBiasListener) OnCriteriaRemoved(
	leftCriteria *model.Criteria,
	params model.MethodParameters,
) model.MethodParameters {
	pParams := params.(AspectEliminationHeuristicParams)
	listener, methodParams := a.getMethodParams(pParams)
	afterRemoveParams := listener.OnCriteriaRemoved(leftCriteria, methodParams)
	return AspectEliminationHeuristicParams{
		Function: pParams.Function,
		Params:   afterRemoveParams,
		Seed:     pParams.Seed,
		Weights:  *pParams.Weights.PreserveOnly(leftCriteria),
	}
}

func (a *AspectEliminationBiasListener) RankCriteriaAscending(params *model.DecisionMakingParams) *model.Criteria {
	wParams := params.MethodParameters.(AspectEliminationHeuristicParams)
	return params.Criteria.SortByWeights(wParams.Weights)
}

func (a *AspectEliminationBiasListener) Merge(params model.MethodParameters, addition model.MethodParameters) model.MethodParameters {
	pParams := params.(AspectEliminationHeuristicParams)
	aParams := addition.(aspectEliminationAddedCriterion)
	listener, methodParams := a.getMethodParams(pParams)
	return AspectEliminationHeuristicParams{
		Function: pParams.Function,
		Params:   listener.Merge(methodParams, aParams.Params),
		Seed:     pParams.Seed,
		Weights:  *pParams.Weights.Merge(&aParams.Weights),
	}
}
