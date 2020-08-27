package criteria_concealment

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type AddedCriterion struct {
	Id                 string                 `json:"id"`
	Type               model.CriterionType    `json:"type"`
	ValuesRange        utils.ValueRange       `json:"valuesRange"`
	AlternativesValues model.Weights          `json:"alternativesValues"`
	MethodParameters   model.MethodParameters `json:"methodParameters"`
}

func (c *CriteriaConcealment) addCriterion(
	props *model.BiasProps,
	parsedProps CriteriaConcealmentParams,
	originalParams, resParams *model.DecisionMakingParams,
	listener *model.BiasListener,
) (*model.DecisionMakingParams, []AddedCriterion) {
	generator := c.generatorSource(parsedProps.RandomSeed)
	criterionBase := c.generateNewCriterionBase(listener, parsedProps.NewCriterionScaling, props, originalParams, resParams)
	addResult := generateCriterionValuesForAlternatives(criterionBase.newCriterion, resParams, generator)
	addedCriterionParams := (*listener).OnCriterionAdded(criterionBase.newCriterion, criterionBase.referenceCriterion, originalParams.MethodParameters, generator)
	finalParams := (*listener).Merge(resParams.MethodParameters, addedCriterionParams)
	newCriteria := resParams.Criteria.Add(criterionBase.newCriterion)
	return &model.DecisionMakingParams{
			NotConsideredAlternatives: *addResult.notConsideredAlternatives,
			ConsideredAlternatives:    *addResult.consideredAlternatives,
			Criteria:                  newCriteria,
			MethodParameters:          finalParams,
		}, []AddedCriterion{{
			Id:                 criterionBase.newCriterion.Id,
			Type:               criterionBase.newCriterion.Type,
			AlternativesValues: addResult.alternativesValues,
			MethodParameters:   addedCriterionParams,
			ValuesRange:        *criterionBase.newCriterion.ValuesRange,
		}}
}

const baseConcealedCriterionName = "__concealedCriterion__"

func (c *CriteriaConcealment) generateNewCriterionBase(
	listener *model.BiasListener,
	scaling float64,
	props *model.BiasProps,
	originalParams, currentParams *model.DecisionMakingParams,
) newCriterionBase {
	refCriterionProvider := c.referenceCriterionManager.ForParams(props)
	rankedCriteria := (*listener).RankCriteriaAscending(originalParams)
	referenceCriterion := refCriterionProvider.Provide(rankedCriteria)
	valRange := getCriterionValueRange(originalParams, referenceCriterion, scaling)
	newCriterion := model.Criterion{
		Id:          newConcealedCriterionName(&currentParams.Criteria),
		Type:        model.Gain,
		ValuesRange: valRange,
	}
	return newCriterionBase{
		referenceCriterion: referenceCriterion,
		newCriterion:       &newCriterion,
	}
}

func newConcealedCriterionName(criteria *model.Criteria) string {
	return criteria.NotUsedName(baseConcealedCriterionName)
}

type newCriterionBase struct {
	referenceCriterion *model.Criterion
	newCriterion       *model.Criterion
}

func generateCriterionValuesForAlternatives(
	newCriterion *model.Criterion,
	resParams *model.DecisionMakingParams,
	valueGenerator utils.ValueGenerator,
) *addCriterionResult {
	generator := utils.NewValueInRangeGenerator(valueGenerator, newCriterion.ValuesRange)
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

func getCriterionValueRange(originalParams *model.DecisionMakingParams, referenceCriterion *model.Criterion, scaling float64) *utils.ValueRange {
	allAlternatives := originalParams.AllAlternatives()
	valRange := model.CriteriaValuesRange(&allAlternatives, referenceCriterion).ScaleEqually(scaling)
	return valRange
}

type addCriterionResult struct {
	notConsideredAlternatives *[]model.AlternativeWithCriteria
	consideredAlternatives    *[]model.AlternativeWithCriteria
	alternativesValues        model.Weights
}
