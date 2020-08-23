package criteria_omission

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/model/criteria-splitting"
)

func omitCriteria(
	omissionOrderCriteria *model.Criteria,
	parsedProps *criteria_splitting.CriteriaSplitCondition,
	current *model.DecisionMakingParams,
	listener *model.BiasListener,
) (*model.DecisionMakingParams, *model.Criteria) {
	omissionPartition := parsedProps.SplitCriteriaByOrdering(omissionOrderCriteria)
	resultMethodParameters := (*listener).OnCriteriaRemoved(omissionPartition.Right, current.MethodParameters)
	consideredAlternatives := model.PreserveCriteriaForAlternatives(&current.ConsideredAlternatives, omissionPartition.Right)
	notConsideredAlternatives := model.PreserveCriteriaForAlternatives(&current.NotConsideredAlternatives, omissionPartition.Right)
	return &model.DecisionMakingParams{
		NotConsideredAlternatives: *notConsideredAlternatives,
		ConsideredAlternatives:    *consideredAlternatives,
		Criteria:                  *omissionPartition.Right,
		MethodParameters:          resultMethodParameters,
	}, omissionPartition.Left
}
