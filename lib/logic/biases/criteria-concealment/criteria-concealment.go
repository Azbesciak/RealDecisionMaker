package criteria_concealment

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/model/reference-criterion"
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
	RandomSeed int64 `json:"randomSeed"`
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
	return &parsedProps
}
