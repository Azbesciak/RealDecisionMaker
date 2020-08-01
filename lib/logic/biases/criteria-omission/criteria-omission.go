package criteria_omission

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

//go:generate easytags $GOFILE json:camel

type CriteriaOmission struct {
}

type CriteriaOmissionParams struct {
	OmittedCriteriaRatio float64 `json:"omittedCriteriaRatio"`
}

type CriteriaOmissionResult struct {
	OmittedCriteria model.Criteria `json:"omittedCriteria"`
}

const BiasName = "criteriaOmission"

func (c *CriteriaOmission) Identifier() string {
	return BiasName
}

func (c *CriteriaOmission) Apply(
	original, current *model.DecisionMakingParams,
	props *model.BiasProps,
	listener *model.BiasListener,
) *model.BiasedResult {
	parsedProps := *parseProps(props)
	if parsedProps.OmittedCriteriaRatio == 0 {
		return &model.BiasedResult{DMP: current, Props: CriteriaOmissionResult{}}
	}
	paramsWithSortedCriteria := paramsWithSortedCriteria(original, listener)
	resParams, omitted := omitCriteria(&parsedProps, paramsWithSortedCriteria, listener)
	return &model.BiasedResult{
		DMP:   resParams,
		Props: CriteriaOmissionResult{OmittedCriteria: *omitted},
	}
}

func paramsWithSortedCriteria(
	params *model.DecisionMakingParams,
	listener *model.BiasListener,
) *model.DecisionMakingParams {
	return &model.DecisionMakingParams{
		NotConsideredAlternatives: params.NotConsideredAlternatives,
		ConsideredAlternatives:    params.ConsideredAlternatives,
		Criteria:                  *(*listener).RankCriteriaAscending(params).Criteria(),
		MethodParameters:          params.MethodParameters,
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
