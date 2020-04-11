package criteria_omission

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

//go:generate easytags $GOFILE json:camel

type CriteriaOmission struct {
	generatorSource         utils.SeededValueGenerator
	newCriterionValueScalar float64
}

type CriteriaOmissionParams struct {
	AddCriterionProbability float64 `json:"addCriterionProbability"`
	OmittedCriteriaRatio    float64 `json:"omittedCriteriaRatio"`
	RandomSeed              int64   `json:"randomSeed"`
}

type CriteriaOmissionResult struct {
	OmittedCriteria model.Criteria   `json:"omittedCriteria"`
	AddedCriteria   []AddedCriterion `json:"addedCriterion"`
}

func (c *CriteriaOmission) Identifier() string {
	return "criteriaOmission"
}

func (c *CriteriaOmission) Apply(
	original, current *model.DecisionMakingParams,
	props *model.HeuristicProps,
	listener *model.HeuristicListener,
) *model.HeuristicResult {
	parsedProps := *parseProps(props)
	if parsedProps.OmittedCriteriaRatio == 0 && parsedProps.AddCriterionProbability == 0 {
		return &model.HeuristicResult{DMP: current, Props: CriteriaOmissionResult{}}
	}
	paramsWithSortedCriteria := paramsWithSortedCriteria(original, listener)
	resParams, omitted := omitCriteria(&parsedProps, paramsWithSortedCriteria, listener)
	resParams, addedCriterion := c.addCriterion(parsedProps, paramsWithSortedCriteria, resParams, listener)
	return &model.HeuristicResult{
		DMP: resParams,
		Props: CriteriaOmissionResult{
			OmittedCriteria: *omitted,
			AddedCriteria:   addedCriterion,
		},
	}
}

func paramsWithSortedCriteria(
	params *model.DecisionMakingParams,
	listener *model.HeuristicListener,
) *model.DecisionMakingParams {
	return &model.DecisionMakingParams{
		NotConsideredAlternatives: params.NotConsideredAlternatives,
		ConsideredAlternatives:    params.ConsideredAlternatives,
		Criteria:                  *((*listener).RankCriteriaAscending(params)),
		MethodParameters:          params.MethodParameters,
	}
}

func parseProps(props *model.HeuristicProps) *CriteriaOmissionParams {
	parsedProps := CriteriaOmissionParams{}
	utils.DecodeToStruct(*props, &parsedProps)
	parsedProps.validate()
	return &parsedProps
}

func (params *CriteriaOmissionParams) validate() {
	if !utils.IsProbability(params.AddCriterionProbability) {
		panic(fmt.Errorf("'addCriterionProbability' need to be in range [0,1], got %f", params.AddCriterionProbability))
	}
	if !utils.IsProbability(params.OmittedCriteriaRatio) {
		panic(fmt.Errorf("'omittedCriteriaRatio' need to be in range [0,1], got %f", params.AddCriterionProbability))
	}
}
