package criteria_concealment

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type AddedCriterion struct {
	Type                model.CriterionType    `json:"type"`
	AlternativesValues  model.Weights          `json:"alternativesValues"`
	MethodParameters    model.MethodParameters `json:"methodParameters"`
	CriterionValueRange utils.ValueRange       `json:"criterionValueRange"`
}

func (c *CriteriaConcealment) addCriterion(
	props *model.BiasProps,
	parsedProps CriteriaConcealmentParams,
	originalParams, resParams *model.DecisionMakingParams,
	listener *model.BiasListener,
) (*model.DecisionMakingParams, []AddedCriterion) {
	if parsedProps.AddCriterionProbability > 0 {
		generator := c.generatorSource(parsedProps.RandomSeed)
		if parsedProps.AddCriterionProbability > generator() {
			return c.generateNewCriterion(listener, props, originalParams, resParams, generator)
		}
	}
	return resParams, []AddedCriterion{}
}

func addedCriterionName() string {
	return "__addedCriterion__"
}

func (c *CriteriaConcealment) generateNewCriterion(
	listener *model.BiasListener,
	props *model.BiasProps,
	originalParams, resParams *model.DecisionMakingParams,
	valueGenerator utils.ValueGenerator,
) (*model.DecisionMakingParams, []AddedCriterion) {
	criterionBase := c.generateNewCriterionBase(listener, props, originalParams)
	addResult := generateCriterionValuesForAlternatives(criterionBase.newCriterion, criterionBase.valuesRange, resParams, valueGenerator)
	addedCriterionParams := (*listener).OnCriterionAdded(criterionBase.newCriterion, criterionBase.referenceCriterion, originalParams.MethodParameters, valueGenerator)
	finalParams := (*listener).Merge(resParams.MethodParameters, addedCriterionParams)
	newCriteria := resParams.Criteria.Add(criterionBase.newCriterion)
	return &model.DecisionMakingParams{
			NotConsideredAlternatives: *addResult.notConsideredAlternatives,
			ConsideredAlternatives:    *addResult.consideredAlternatives,
			Criteria:                  newCriteria,
			MethodParameters:          finalParams,
		}, []AddedCriterion{{
			Type:                criterionBase.newCriterion.Type,
			AlternativesValues:  addResult.alternativesValues,
			MethodParameters:    addedCriterionParams,
			CriterionValueRange: *criterionBase.valuesRange,
		}}
}

func (c *CriteriaConcealment) generateNewCriterionBase(
	listener *model.BiasListener,
	props *model.BiasProps,
	originalParams *model.DecisionMakingParams,
) newCriterionBase {
	refCriterionProvider := c.referenceCriterionManager.ForParams(props)
	rankedCriteria := (*listener).RankCriteriaAscending(originalParams)
	referenceCriterion := refCriterionProvider.Provide(rankedCriteria)
	newCriterion := model.Criterion{Id: addedCriterionName(), Type: model.Gain}
	valRange := c.getCriterionValueRange(originalParams, referenceCriterion)
	return newCriterionBase{
		referenceCriterion: referenceCriterion,
		newCriterion:       &newCriterion,
		valuesRange:        valRange,
	}
}

type newCriterionBase struct {
	referenceCriterion *model.Criterion
	newCriterion       *model.Criterion
	valuesRange        *utils.ValueRange
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

func (c *CriteriaConcealment) getCriterionValueRange(originalParams *model.DecisionMakingParams, referenceCriterion *model.Criterion) *utils.ValueRange {
	allAlternatives := originalParams.AllAlternatives()
	valRange := model.CriteriaValuesRange(&allAlternatives, referenceCriterion).ScaleEqually(c.newCriterionValueScalar)
	return valRange
}

type addCriterionResult struct {
	notConsideredAlternatives *[]model.AlternativeWithCriteria
	consideredAlternatives    *[]model.AlternativeWithCriteria
	alternativesValues        model.Weights
}
