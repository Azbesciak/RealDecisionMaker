package criteria_omission

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

//go:generate easytags $GOFILE json:camel

type CriteriaOmission struct {
	omissionResolvers []OmissionResolver
}

func NewCriteriaOmission(omissionResolvers []OmissionResolver) *CriteriaOmission {
	if len(omissionResolvers) == 0 {
		panic("no criteria omission order resolvers")
	}
	return &CriteriaOmission{omissionResolvers: omissionResolvers}
}

type CriteriaOmissionParams struct {
	OmittedCriteriaRatio float64 `json:"omittedCriteriaRatio"`
	OmissionOrder        string  `json:"omissionOrder"`
}

type CriteriaOmissionResult struct {
	OmittedCriteria model.Criteria `json:"omittedCriteria"`
}

const BiasName = "criteriaOmission"

func (c *CriteriaOmission) Identifier() string {
	return BiasName
}

func (c *CriteriaOmission) omissionResolver(params *CriteriaOmissionParams) OmissionResolver {
	if len(params.OmissionOrder) == 0 {
		return c.omissionResolvers[0]
	}
	for _, r := range c.omissionResolvers {
		if r.Identifier() == params.OmissionOrder {
			return r
		}
	}
	names := make([]string, len(c.omissionResolvers))
	for i, r := range c.omissionResolvers {
		names[i] = r.Identifier()
	}
	panic(fmt.Errorf("omission order resolver '%s' not found in %v", c.omissionResolvers, names))
}

func (c *CriteriaOmission) Apply(
	_, current *model.DecisionMakingParams,
	props *model.BiasProps,
	listener *model.BiasListener,
) *model.BiasedResult {
	parsedProps := *parseProps(props)
	if parsedProps.OmittedCriteriaRatio == 0 {
		return &model.BiasedResult{DMP: current, Props: CriteriaOmissionResult{}}
	}
	resolver := c.omissionResolver(&parsedProps)
	sortedCriteria := resolver.CriteriaOmissionOrder(current, props, listener)
	resParams, omitted := omitCriteria(sortedCriteria, &parsedProps, current, listener)
	return &model.BiasedResult{
		DMP:   resParams,
		Props: CriteriaOmissionResult{OmittedCriteria: *omitted},
	}
}

func parseProps(props *model.BiasProps) *CriteriaOmissionParams {
	parsedProps := CriteriaOmissionParams{}
	utils.DecodeToStruct(*props, &parsedProps)
	parsedProps.validate()
	return &parsedProps
}

func (params *CriteriaOmissionParams) validate() {
	if !utils.IsProbability(params.OmittedCriteriaRatio) {
		panic(fmt.Errorf("'omittedCriteriaRatio' need to be in range [0,1], got %f", params.OmittedCriteriaRatio))
	}
}
