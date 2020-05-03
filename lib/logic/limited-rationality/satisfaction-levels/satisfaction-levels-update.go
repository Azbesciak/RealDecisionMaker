package satisfaction_levels

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"sort"
)

type ParamsAddition = interface{}

type CriterionAdder interface {
	OnCriterionAdded(
		criterion *model.Criterion,
		previousRankedCriteria *model.Criteria,
		params SatisfactionLevels,
		generator utils.ValueGenerator,
	) ParamsAddition
}

type SatisfactionParamsMerger interface {
	Merge(params SatisfactionLevels, addition ParamsAddition) SatisfactionLevels
}

type CriterionRemover interface {
	OnCriteriaRemoved(leftCriteria *model.Criteria, params SatisfactionLevels) SatisfactionLevels
}

type SatisfactionLevelsUpdateListener interface {
	SatisfactionParamsSource
	CriterionAdder
	CriterionRemover
	SatisfactionParamsMerger
}
type ListenersMap = map[string]SatisfactionLevelsUpdateListener

type SatisfactionLevelsUpdateListeners struct {
	Listeners ListenersMap
}

func (sl *SatisfactionLevelsUpdateListeners) Fetch(listenerName string) *SatisfactionLevelsUpdateListener {
	fun, ok := (sl.Listeners)[listenerName]
	if !ok {
		keys := make([]string, 0)
		for k := range sl.Listeners {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		panic(fmt.Errorf("satisfaction levels update listener for '%s' not found, available are '%s'", listenerName, keys))
	}
	listener := fun.(SatisfactionLevelsUpdateListener)
	return &listener
}

func (sl *SatisfactionLevelsUpdateListeners) Get(listenerName string, params interface{}) (SatisfactionLevelsUpdateListener, SatisfactionLevels) {
	listener := *sl.Fetch(listenerName)
	methodParams := listener.BlankParams()
	utils.DecodeToStruct(params, &methodParams)
	return listener, methodParams
}
