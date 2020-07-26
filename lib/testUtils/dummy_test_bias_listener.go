package testUtils

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"sort"
)

type DummyMethodParameters struct {
	Criteria []string
}

type DummyBiasListener struct {
}

func (d *DummyBiasListener) Merge(params model.MethodParameters, addition model.MethodParameters) model.MethodParameters {
	prevParams := params.(DummyMethodParameters)
	addedParams := addition.(DummyMethodParameters)
	return DummyMethodParameters{Criteria: append(prevParams.Criteria, addedParams.Criteria...)}
}

func (d *DummyBiasListener) Identifier() string {
	panic("should not call identifier in test")
}

func (d *DummyBiasListener) OnCriterionAdded(
	criterion *model.Criterion,
	referenceCriterion *model.Criterion,
	params model.MethodParameters,
	generator utils.ValueGenerator,
) model.AddedCriterionParams {
	return DummyMethodParameters{Criteria: []string{criterion.Id}}
}

func (d *DummyBiasListener) OnCriteriaRemoved(leftCriteria *model.Criteria, params model.MethodParameters) model.MethodParameters {
	return DummyMethodParameters{Criteria: *leftCriteria.Names()}
}

func (d *DummyBiasListener) RankCriteriaAscending(params *model.DecisionMakingParams) *model.WeightedCriteria {
	criteria := params.Criteria.ShallowCopy()
	sort.Slice(*criteria, func(i, j int) bool {
		return (*criteria)[i].Id < (*criteria)[j].Id
	})
	result := make(model.WeightedCriteria, len(*criteria))
	for i, c := range *criteria {
		result[i] = model.WeightedCriterion{
			Criterion: c,
			Weight:    float64(i + 1),
		}
	}
	return &result
}
