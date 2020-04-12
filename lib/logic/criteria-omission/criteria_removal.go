package criteria_omission

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"math"
)

func omitCriteria(
	parsedProps *CriteriaOmissionParams,
	current *model.DecisionMakingParams,
	listener *model.BiasListener,
) (*model.DecisionMakingParams, *model.Criteria) {
	omissionPartition := splitCriteriaToOmit(parsedProps.OmittedCriteriaRatio, &current.Criteria)
	resultMethodParameters := (*listener).OnCriteriaRemoved(omissionPartition.omitted, omissionPartition.kept, current.MethodParameters)
	consideredAlternatives := model.PreserveCriteriaForAlternatives(&current.ConsideredAlternatives, omissionPartition.kept)
	notConsideredAlternatives := model.PreserveCriteriaForAlternatives(&current.NotConsideredAlternatives, omissionPartition.kept)
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
