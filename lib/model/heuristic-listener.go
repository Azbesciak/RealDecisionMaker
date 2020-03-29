package model

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type AddedCriterionParams = MethodParameters

type CriterionAdder interface {
	OnCriterionAdded(
		criterion *Criterion,
		previousRankedCriteria *Criteria,
		params MethodParameters,
		generator utils.ValueGenerator,
	) AddedCriterionParams
}

type MethodParametersMerger interface {
	Merge(params MethodParameters, addition MethodParameters) MethodParameters
}

type CriterionRemover interface {
	OnCriteriaRemoved(removedCriteria *Criteria, leftCriteria *Criteria, params MethodParameters) MethodParameters
}

type CriteriaRanker interface {
	RankCriteriaAscending(params *DecisionMakingParams) *Criteria
}

type HeuristicListener interface {
	utils.Identifiable
	CriterionAdder
	CriterionRemover
	CriteriaRanker
	MethodParametersMerger
}

type HeuristicListeners struct {
	Listeners []HeuristicListener
}

func (pf *HeuristicListeners) Get(index int) utils.Identifiable {
	return pf.Listeners[index]
}

func (pf *HeuristicListeners) Len() int {
	return len(pf.Listeners)
}

func (pf *HeuristicListeners) Fetch(listenerName string) *HeuristicListener {
	preferenceFunMap := utils.AsMap(pf)
	fun, ok := (*preferenceFunMap)[listenerName]
	if !ok {
		var keys []string
		for _, k := range pf.Listeners {
			keys = append(keys, k.Identifier())
		}
		panic(fmt.Errorf("heuristic listener for '%s' not found, available are '%s'", listenerName, keys))
	}
	listener := fun.(HeuristicListener)
	return &listener
}
