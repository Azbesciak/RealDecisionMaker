package criteria_omission

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/model/criteria-ordering"
)

func omitCriteria(
	omissionOrderCriteria *model.Criteria,
	parsedProps *CriteriaOmissionParams,
	current *model.DecisionMakingParams,
	listener *model.BiasListener,
) (*model.DecisionMakingParams, *model.Criteria) {
	omissionPartition := criteria_ordering.SplitCriteriaByOrdering(parsedProps.OmittedCriteriaRatio, omissionOrderCriteria)
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
