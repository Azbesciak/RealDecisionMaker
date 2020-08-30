package criteria_concealment

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/model/criteria-bounding"
	"github.com/Azbesciak/RealDecisionMaker/lib/model/reference-criterion"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

//go:generate easytags $GOFILE json:camel

type CriteriaConcealment struct {
	generatorSource           utils.SeededValueGenerator
	referenceCriterionManager reference_criterion.ReferenceCriteriaManager
}

func NewCriteriaConcealment(
	generatorSource utils.SeededValueGenerator,
	referenceCriterionManager reference_criterion.ReferenceCriteriaManager,
) *CriteriaConcealment {
	return &CriteriaConcealment{
		generatorSource:           generatorSource,
		referenceCriterionManager: referenceCriterionManager,
	}
}

type CriteriaConcealmentParams struct {
	RandomSeed          int64   `json:"randomSeed"`
	NewCriterionScaling float64 `json:"newCriterionScaling"`
}

type CriteriaConcealmentResult struct {
	AddedCriteria []AddedCriterion `json:"addedCriteria"`
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
	bounding := criteria_bounding.FromParams(props)
	resParams, addedCriterion := c.addCriterion(props, parsedProps, original, current, listener, bounding)
	return &model.BiasedResult{
		DMP: resParams,
		Props: CriteriaConcealmentResult{
			AddedCriteria: addedCriterion,
		},
	}
}

func parseProps(props *model.BiasProps) *CriteriaConcealmentParams {
	parsedProps := CriteriaConcealmentParams{NewCriterionScaling: 1}
	utils.DecodeToStruct(*props, &parsedProps)
	if parsedProps.NewCriterionScaling == 0 {
		panic("`concealedCriterionScaling` cannot be 0")
	}
	return &parsedProps
}
