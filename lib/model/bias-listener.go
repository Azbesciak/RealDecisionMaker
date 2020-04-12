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

type BiasListener interface {
	utils.Identifiable
	CriterionAdder
	CriterionRemover
	CriteriaRanker
	MethodParametersMerger
}

type BiasListeners struct {
	Listeners []BiasListener
}

func (pf *BiasListeners) Get(index int) utils.Identifiable {
	return pf.Listeners[index]
}

func (pf *BiasListeners) Len() int {
	return len(pf.Listeners)
}

func (pf *BiasListeners) Fetch(listenerName string) *BiasListener {
	preferenceFunMap := utils.AsMap(pf)
	fun, ok := (*preferenceFunMap)[listenerName]
	if !ok {
		var keys []string
		for _, k := range pf.Listeners {
			keys = append(keys, k.Identifier())
		}
		panic(fmt.Errorf("bias listener for '%s' not found, available are '%s'", listenerName, keys))
	}
	listener := fun.(BiasListener)
	return &listener
}

func PrepareCumulatedWeightsMap(
	params *DecisionMakingParams,
	mapper func(criterion string, value Weight) Weight,
) *Weights {
	weights := make(Weights, len(params.Criteria))
	for _, a := range params.ConsideredAlternatives {
		for crit, v := range a.Criteria {
			w, ok := weights[crit]
			if !ok {
				weights[crit] = mapper(crit, v)
			} else {
				weights[crit] = w + mapper(crit, v)
			}
		}
	}
	return &weights
}

func WeightIdentity(criterion string, value Weight) Weight {
	return value
}
