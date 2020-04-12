package criteria_omission

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

func (c *CriteriaOmission) addCriterion(
	parsedProps CriteriaOmissionParams,
	originalParams *model.DecisionMakingParams,
	resParams *model.DecisionMakingParams,
	listener *model.BiasListener,
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
	listener *model.BiasListener,
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
