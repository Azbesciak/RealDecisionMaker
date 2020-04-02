package criteria_omission

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"math"
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

type AddedCriterion struct {
	Type                model.CriterionType    `json:"type"`
	AlternativesValues  model.Weights          `json:"alternativesValues"`
	MethodParameters    model.MethodParameters `json:"methodParameters"`
	CriterionValueRange utils.ValueRange       `json:"criterionValueRange"`
}

func (c *CriteriaOmission) Identifier() string {
	return "criteriaOmission"
}

func (c *CriteriaOmission) Apply(
	params *model.DecisionMakingParams,
	props *model.HeuristicProps,
	listener *model.HeuristicListener,
) *model.HeuristicResult {
	parsedProps := *parseProps(props)
	if parsedProps.OmittedCriteriaRatio == 0 && parsedProps.AddCriterionProbability == 0 {
		return &model.HeuristicResult{DMP: params, Props: CriteriaOmissionResult{}}
	}
	paramsWithSortedCriteria := paramsWithSortedCriteria(params, listener)
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

func (c *CriteriaOmission) addCriterion(
	parsedProps CriteriaOmissionParams,
	originalParams *model.DecisionMakingParams,
	resParams *model.DecisionMakingParams,
	listener *model.HeuristicListener,
) (*model.DecisionMakingParams, []AddedCriterion) {
	if parsedProps.AddCriterionProbability > 0 {
		generator := c.generatorSource(parsedProps.RandomSeed)
		if parsedProps.AddCriterionProbability > generator() {
			return c.generateNewCriterion(listener, originalParams, resParams, generator)
		}
	}
	return resParams, []AddedCriterion{}
}

func addedCriterionName() string {
	return "__addedCriterion__"
}

func (c *CriteriaOmission) generateNewCriterion(
	listener *model.HeuristicListener,
	originalParams *model.DecisionMakingParams,
	resParams *model.DecisionMakingParams,
	valueGenerator utils.ValueGenerator,
) (*model.DecisionMakingParams, []AddedCriterion) {
	newCriterion := model.Criterion{Id: addedCriterionName(), Type: model.Gain}
	valRange := c.getCriterionValueRange(originalParams)
	addResult := generateCriterionValuesForAlternatives(&newCriterion, valRange, resParams, valueGenerator)
	addedCriterionParams := (*listener).OnCriterionAdded(&newCriterion, &resParams.Criteria, resParams.MethodParameters, valueGenerator)
	finalParams := (*listener).Merge(resParams.MethodParameters, addedCriterionParams)
	newCriteria := resParams.Criteria.Add(&newCriterion)
	return &model.DecisionMakingParams{
			NotConsideredAlternatives: *addResult.notConsideredAlternatives,
			ConsideredAlternatives:    *addResult.consideredAlternatives,
			Criteria:                  newCriteria,
			MethodParameters:          finalParams,
		}, []AddedCriterion{{
			Type:                newCriterion.Type,
			AlternativesValues:  addResult.alternativesValues,
			MethodParameters:    addedCriterionParams,
			CriterionValueRange: *valRange,
		}}
}

func generateCriterionValuesForAlternatives(
	newCriterion *model.Criterion,
	criterionValueRange *utils.ValueRange,
	resParams *model.DecisionMakingParams,
	valueGenerator utils.ValueGenerator,
) *addCriterionResult {
	generator := utils.NewValueInRangeGenerator(valueGenerator, criterionValueRange)
	sortedAlternatives, alternativesValues := assignNewCriterionToAlternatives(resParams, generator, newCriterion)
	return &addCriterionResult{
		notConsideredAlternatives: model.UpdateAlternatives(&resParams.NotConsideredAlternatives, sortedAlternatives),
		consideredAlternatives:    model.UpdateAlternatives(&resParams.ConsideredAlternatives, sortedAlternatives),
		alternativesValues:        alternativesValues,
	}
}

func assignNewCriterionToAlternatives(
	resParams *model.DecisionMakingParams,
	generator utils.ValueGenerator,
	newCriterion *model.Criterion,
) (*[]model.AlternativeWithCriteria, model.Weights) {
	allAlternatives := resParams.AllAlternatives()
	sortedAlternatives := model.SortAlternativesByName(&allAlternatives)
	alternativesValues := make(model.Weights, len(*sortedAlternatives))
	sortedAlternatives = model.AddCriterionToAlternatives(sortedAlternatives, newCriterion,
		func(a *model.AlternativeWithCriteria) model.Weight {
			newValue := generator()
			alternativesValues[a.Id] = newValue
			return newValue
		})
	return sortedAlternatives, alternativesValues
}

func (c *CriteriaOmission) getCriterionValueRange(originalParams *model.DecisionMakingParams) *utils.ValueRange {
	weakestCriterion := originalParams.Criteria[0]
	allAlternatives := originalParams.AllAlternatives()
	valRange := model.CriteriaValuesRange(&allAlternatives, &weakestCriterion).ScaleEqually(c.newCriterionValueScalar)
	return valRange
}

type addCriterionResult struct {
	notConsideredAlternatives *[]model.AlternativeWithCriteria
	consideredAlternatives    *[]model.AlternativeWithCriteria
	alternativesValues        model.Weights
}

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
