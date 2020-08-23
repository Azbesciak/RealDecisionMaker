package criteria_omission

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/model/criteria-ordering"
	"github.com/Azbesciak/RealDecisionMaker/lib/model/criteria-splitting"
)

//go:generate easytags $GOFILE json:camel

type CriteriaOmission struct {
	omissionResolvers []criteria_ordering.CriteriaOrderingResolver
}

func NewCriteriaOmission(omissionResolvers []criteria_ordering.CriteriaOrderingResolver) *CriteriaOmission {
	if len(omissionResolvers) == 0 {
		panic("no criteria omission order resolvers")
	}
	return &CriteriaOmission{omissionResolvers: omissionResolvers}
}

// + CriteriaSplitCondition
// + CriteriaOrdering
type CriteriaOmissionParams struct {
}

type CriteriaOmissionResult struct {
	OmittedCriteria model.Criteria `json:"omittedCriteria"`
}

const BiasName = "criteriaOmission"

func (c *CriteriaOmission) Identifier() string {
	return BiasName
}

func (c *CriteriaOmission) Apply(
	_, current *model.DecisionMakingParams,
	props *model.BiasProps,
	listener *model.BiasListener,
) *model.BiasedResult {
	parsedProps, splitting := parseProps(props)
	resolver := criteria_ordering.FetchOrderingResolver(&c.omissionResolvers, parsedProps)
	sortedCriteria := resolver.OrderCriteria(current, props, listener)
	resParams, omitted := omitCriteria(sortedCriteria, splitting, current, listener)
	return &model.BiasedResult{
		DMP:   resParams,
		Props: CriteriaOmissionResult{OmittedCriteria: *omitted},
	}
}

func parseProps(props *model.BiasProps) (*criteria_ordering.CriteriaOrdering, *criteria_splitting.CriteriaSplitCondition) {
	ordering := criteria_ordering.Parse(props)
	splittingProps := criteria_splitting.Parse(props)
	return ordering, splittingProps
}
