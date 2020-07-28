package criteria_concealment

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	reference_criterion "github.com/Azbesciak/RealDecisionMaker/lib/model/reference-criterion"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

//go:generate easytags $GOFILE json:camel

type CriteriaConcealment struct {
	generatorSource           utils.SeededValueGenerator
	newCriterionValueScalar   float64
	referenceCriterionManager reference_criterion.ReferenceCriteriaManager
}

func NewCriteriaConcealment(
	generatorSource utils.SeededValueGenerator,
	newCriterionValueScalar float64,
	referenceCriterionManager reference_criterion.ReferenceCriteriaManager,
) *CriteriaConcealment {
	return &CriteriaConcealment{
		generatorSource:           generatorSource,
		newCriterionValueScalar:   newCriterionValueScalar,
		referenceCriterionManager: referenceCriterionManager,
	}
}

type CriteriaConcealmentParams struct {
	CriterionConcealmentProbability float64 `json:"criterionConcealmentProbability"`
	RandomSeed                      int64   `json:"randomSeed"`
}

type CriteriaConcealmentResult struct {
	AddedCriteria []AddedCriterion `json:"addedCriterion"`
}

const BiasName = "criteriaConcealment"

func (c *CriteriaConcealment) Identifier() string {
	return BiasName
}

func (c *CriteriaConcealment) Apply(
	original, current *model.DecisionMakingParams,
	props *model.BiasProps,
	listener *model.BiasListener,
) *model.BiasedResult {
	parsedProps := *parseProps(props)
	if parsedProps.CriterionConcealmentProbability == 0 {
		return &model.BiasedResult{DMP: current, Props: CriteriaConcealmentResult{}}
	}
	resParams, addedCriterion := c.addCriterion(props, parsedProps, original, current, listener)
	return &model.BiasedResult{
		DMP: resParams,
		Props: CriteriaConcealmentResult{
			AddedCriteria: addedCriterion,
		},
	}
}

func parseProps(props *model.BiasProps) *CriteriaConcealmentParams {
	parsedProps := CriteriaConcealmentParams{}
	utils.DecodeToStruct(*props, &parsedProps)
	parsedProps.validate()
	return &parsedProps
}

func (params *CriteriaConcealmentParams) validate() {
	if !utils.IsProbability(params.CriterionConcealmentProbability) {
		panic(fmt.Errorf("'criterionConcealmentProbability' need to be in range [0,1], got %f", params.CriterionConcealmentProbability))
	}
}
