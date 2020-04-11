package criteria_omission

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"math"
)

func omitCriteria(
	parsedProps *CriteriaOmissionParams,
	params *model.DecisionMakingParams,
	listener *model.HeuristicListener,
) (*model.DecisionMakingParams, *model.Criteria) {
	omissionPartition := splitCriteriaToOmit(parsedProps.OmittedCriteriaRatio, &params.Criteria)
	resultMethodParameters := (*listener).OnCriteriaRemoved(omissionPartition.omitted, omissionPartition.kept, params.MethodParameters)
	consideredAlternatives := model.PreserveCriteriaForAlternatives(&params.ConsideredAlternatives, omissionPartition.kept)
	notConsideredAlternatives := model.PreserveCriteriaForAlternatives(&params.NotConsideredAlternatives, omissionPartition.kept)
	return &model.DecisionMakingParams{
		NotConsideredAlternatives: *notConsideredAlternatives,
		ConsideredAlternatives:    *consideredAlternatives,
		Criteria:                  *omissionPartition.kept,
		MethodParameters:          resultMethodParameters,
	}, omissionPartition.omitted
}

func splitCriteriaToOmit(ratio float64, sortedCriteria *model.Criteria) *criteriaOmissionPartition {
	criteriaCount := len(*sortedCriteria)
	toOmitCount := int(math.Floor(float64(criteriaCount) * ratio))
	toOmit := (*sortedCriteria)[0:toOmitCount]
	toKeep := (*sortedCriteria)[toOmitCount:]
	return &criteriaOmissionPartition{
		kept:    &toKeep,
		omitted: &toOmit,
	}
}

type criteriaOmissionPartition struct {
	kept    *model.Criteria
	omitted *model.Criteria
}
