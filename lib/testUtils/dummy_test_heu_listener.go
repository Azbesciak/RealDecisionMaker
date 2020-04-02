package testUtils

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"sort"
)

type DummyMethodParameters struct {
	Criteria []string
}

type DummyHeuListener struct {
}

func (d *DummyHeuListener) Merge(params model.MethodParameters, addition model.MethodParameters) model.MethodParameters {
	prevParams := params.(DummyMethodParameters)
	addedParams := addition.(DummyMethodParameters)
	return DummyMethodParameters{Criteria: append(prevParams.Criteria, addedParams.Criteria...)}
}

func (d *DummyHeuListener) Identifier() string {
	panic("should not call identifier in test")
}

func (d *DummyHeuListener) OnCriterionAdded(
	criterion *model.Criterion,
	previousRankedCriteria *model.Criteria,
	params model.MethodParameters,
	generator utils.ValueGenerator,
) model.AddedCriterionParams {
	return DummyMethodParameters{Criteria: []string{criterion.Id}}
}

func (d *DummyHeuListener) OnCriteriaRemoved(removedCriteria *model.Criteria, leftCriteria *model.Criteria, params model.MethodParameters) model.MethodParameters {
	return DummyMethodParameters{Criteria: *leftCriteria.Names()}
}

func (d *DummyHeuListener) RankCriteriaAscending(params *model.DecisionMakingParams) *model.Criteria {
	criteria := params.Criteria.ShallowCopy()
	sort.Slice(*criteria, func(i, j int) bool {
		return (*criteria)[i].Id < (*criteria)[j].Id
	})
	return criteria
}
