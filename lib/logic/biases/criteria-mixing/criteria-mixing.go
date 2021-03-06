package criteria_mixing

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/model/reference-criterion"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

//go:generate easytags $GOFILE json:camel

type CriteriaMixingParams struct {
	// TODO interaction factor!! how to implement
	RandomSeed  int64   `json:"randomSeed"`
	MixingRatio float64 `json:"mixingRatio"`
}

// + referenceCriterionParams
type MixedCriterion struct {
	Component1   CriterionComponent     `json:"component1"`
	Component2   CriterionComponent     `json:"component2"`
	NewCriterion CriterionComponent     `json:"newCriterion"`
	Params       model.MethodParameters `json:"params"`
}

type CriterionComponent struct {
	Id           string              `json:"id"`
	Type         model.CriterionType `json:"type"`
	ScaledValues model.Weights       `json:"scaledValues"`
}

type MixedCriterionValue struct {
	Value model.Weight `json:"value"`
}

type CriteriaMixingResult struct {
}

func (p *CriteriaMixingParams) validate() {
	if !utils.IsProbability(p.MixingRatio) {
		panic(fmt.Errorf("mixingRatio should be in range [0,1]"))
	}
}

type CriteriaMixing struct {
	generatorSource          utils.SeededValueGenerator
	referenceCriteriaManager reference_criterion.ReferenceCriteriaManager
}

func NewCriteriaMixing(
	generatorSource utils.SeededValueGenerator,
	referenceCriteriaManager reference_criterion.ReferenceCriteriaManager,
) *CriteriaMixing {
	return &CriteriaMixing{
		generatorSource:          generatorSource,
		referenceCriteriaManager: referenceCriteriaManager,
	}
}

const BiasName = "criteriaMixing"

func (c *CriteriaMixing) Identifier() string {
	return BiasName
}

func (c *CriteriaMixing) Apply(
	original, current *model.DecisionMakingParams,
	props *model.BiasProps,
	listener *model.BiasListener,
) *model.BiasedResult {
	if current.Criteria.Len() < 2 {
		return &model.BiasedResult{DMP: current}
	}
	parsedProps := parseProps(props)
	generator := c.generatorSource(parsedProps.RandomSeed)
	c2m := selectCriteriaToMix(original, generator)
	allAlternatives := original.AllAlternatives()
	referenceCriterionProvider := c.referenceCriteriaManager.ForParams(props)
	referenceCriterion := referenceCriterion(original, listener, referenceCriterionProvider)
	targetValRange := model.ValuesRangeWithGroundZero(&allAlternatives, referenceCriterion)
	mixResult := c2m.mix(&allAlternatives, targetValRange, parsedProps)
	newCriterion := c2m.criterion(targetValRange)
	criterionParams := (*listener).OnCriterionAdded(&newCriterion, referenceCriterion, current.MethodParameters, generator)
	newMethodParams := (*listener).Merge(current.MethodParameters, criterionParams)
	newAlternatives := updateAlternatives(allAlternatives, newCriterion, mixResult)
	newParams := updateDMParams(current, newAlternatives, newCriterion, newMethodParams)
	return &model.BiasedResult{
		DMP:   &newParams,
		Props: prepareMixedCriterion(c2m, mixResult, newCriterion, criterionParams),
	}
}

func updateDMParams(
	params *model.DecisionMakingParams,
	newAlternatives *[]model.AlternativeWithCriteria,
	newCriterion model.Criterion,
	newMethodParams model.MethodParameters,
) model.DecisionMakingParams {
	return model.DecisionMakingParams{
		NotConsideredAlternatives: *model.UpdateAlternatives(&params.NotConsideredAlternatives, newAlternatives),
		ConsideredAlternatives:    *model.UpdateAlternatives(&params.ConsideredAlternatives, newAlternatives),
		Criteria:                  params.Criteria.Add(&newCriterion),
		MethodParameters:          newMethodParams,
	}
}

func updateAlternatives(allAlternatives []model.AlternativeWithCriteria, newCriterion model.Criterion, mixResult *mixResult) *[]model.AlternativeWithCriteria {
	return model.AddCriterionToAlternatives(&allAlternatives, &newCriterion, func(alt *model.AlternativeWithCriteria) model.Weight {
		return mixResult.result[alt.Id]
	})
}

func prepareMixedCriterion(c2m criteriaToMix, mixResult *mixResult, newCriterion model.Criterion, criterionParams model.AddedCriterionParams) MixedCriterion {
	return MixedCriterion{
		Component1: CriterionComponent{
			Id:           c2m.c1.Id,
			Type:         c2m.c1.Type,
			ScaledValues: mixResult.c1,
		},
		Component2: CriterionComponent{
			Id:           c2m.c2.Id,
			Type:         c2m.c2.Type,
			ScaledValues: mixResult.c2,
		},
		NewCriterion: CriterionComponent{
			Id:           newCriterion.Id,
			Type:         newCriterion.Type,
			ScaledValues: mixResult.result,
		},
		Params: criterionParams,
	}
}

func selectCriteriaToMix(
	params *model.DecisionMakingParams,
	generator utils.ValueGenerator,
) criteriaToMix {
	criteriaNum := len(params.Criteria)
	i1 := int(generator() * float64(criteriaNum))
	offset := int(generator()*float64(criteriaNum-2)) + 1
	return criteriaToMix{params.Criteria[i1], params.Criteria[(i1+offset)%criteriaNum]}
}

type criteriaToMix struct {
	c1, c2 model.Criterion
}

func referenceCriterion(
	params *model.DecisionMakingParams,
	listener *model.BiasListener,
	refCriterionProvider reference_criterion.ReferenceCriterionProvider,
) *model.Criterion {
	ranked := (*listener).RankCriteriaAscending(params)
	return refCriterionProvider.Provide(ranked)
}

func (c *criteriaToMix) mix(
	allAlternatives *[]model.AlternativeWithCriteria,
	targetValuesRange *utils.ValueRange,
	props *CriteriaMixingParams,
) *mixResult {
	c1Values := model.RescaleCriterion(&c.c1, allAlternatives, targetValuesRange)
	c2Values := model.RescaleCriterion(&c.c2, allAlternatives, targetValuesRange)
	resultValues := make(model.Weights, len(c2Values))
	for a, c1Value := range c1Values {
		c2Value, ok := c2Values[a]
		if !ok {
			panic(fmt.Errorf("criterion value for '%s' not found for alternative '%s'", c.c2, a))
		}
		value := c1Value*props.MixingRatio + c2Value*(1-props.MixingRatio)
		resultValues[a] = value
	}
	return &mixResult{
		c1:     c1Values,
		c2:     c2Values,
		result: resultValues,
	}
}

func (c *criteriaToMix) criterion(valRange *utils.ValueRange) model.Criterion {
	return model.Criterion{
		Id:          "__" + c.c1.Id + "+" + c.c2.Id + "__",
		Type:        model.Gain,
		ValuesRange: valRange,
	}
}

type mixResult struct {
	c1, c2, result model.Weights
}

func parseProps(props *model.BiasProps) *CriteriaMixingParams {
	parsedProps := CriteriaMixingParams{MixingRatio: 0.5}
	utils.DecodeToStruct(*props, &parsedProps)
	parsedProps.validate()
	return &parsedProps
}
