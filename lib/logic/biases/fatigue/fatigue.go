package fatigue

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/model/criteria-bounding"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

//go:generate easytags $GOFILE json:camel

type FatigueFunctionParams = interface{}

type FatigueResult struct {
	EffectiveFatigueRatio     float64                         `json:"effectiveFatigueRatio"`
	ConsideredAlternatives    []model.AlternativeWithCriteria `json:"consideredAlternatives"`
	NotConsideredAlternatives []model.AlternativeWithCriteria `json:"notConsideredAlternatives"`
}

type FatigueFunction interface {
	Name() string
	// need to return pointer - it will ve filled later
	BlankParams() FatigueFunctionParams
	Evaluate(params FatigueFunctionParams) float64
}

type FatigueParams struct {
	Function   string      `json:"function"`
	Params     interface{} `json:"params"`
	RandomSeed int64       `json:"randomSeed"`
}

type Fatigue struct {
	valueGeneratorSource utils.SeededValueGenerator
	signGeneratorSource  utils.SeededValueGenerator
	functions            []FatigueFunction
}

func NewFatigue(
	valueGeneratorSource utils.SeededValueGenerator,
	signGeneratorSource utils.SeededValueGenerator,
	functions []FatigueFunction,
) *Fatigue {
	return &Fatigue{valueGeneratorSource: valueGeneratorSource, signGeneratorSource: signGeneratorSource, functions: functions}
}

const BiasName = "fatigue"

func (f *Fatigue) Identifier() string {
	return BiasName
}

func (f *Fatigue) Apply(
	_, current *model.DecisionMakingParams,
	props *model.BiasProps,
	_ *model.BiasListener,
) *model.BiasedResult {
	parsedProps := parseProps(props)
	fun := f.getFatigueFunction(parsedProps)
	funParams := parseFatigueFuncParams(fun, parsedProps)
	fatigueRatio := fun.Evaluate(funParams)
	valueGenerator := f.valueGeneratorSource(parsedProps.RandomSeed)
	signGenerator := f.signGeneratorSource(parsedProps.RandomSeed)
	criteria := matchCriteriaWithBoundings(current, props)
	consideredAlts := blurCriteriaValues(
		current.ConsideredAlternatives, criteria,
		valueGenerator, signGenerator, fatigueRatio,
	)
	notConsideredAlts := blurCriteriaValues(
		current.NotConsideredAlternatives, criteria,
		valueGenerator, signGenerator, fatigueRatio,
	)
	return prepareResult(notConsideredAlts, consideredAlts, current, fatigueRatio)
}

type CriterionWithBounding struct {
	criterion model.Criterion
	bounding  *criteria_bounding.CriteriaInRangeBounding
}

func matchCriteriaWithBoundings(dmp *model.DecisionMakingParams, props *model.BiasProps) []CriterionWithBounding {
	bounding := criteria_bounding.FromParams(props)
	result := make([]CriterionWithBounding, len(dmp.Criteria))
	alternatives := dmp.AllAlternatives()
	for i, c := range dmp.Criteria {
		valuesRange := model.CriteriaValuesRange(&alternatives, &c)
		result[i] = CriterionWithBounding{
			criterion: c,
			bounding:  bounding.WithRange(valuesRange),
		}
	}
	return result
}

func prepareResult(
	notConsideredAlts []model.AlternativeWithCriteria,
	consideredAlts []model.AlternativeWithCriteria,
	current *model.DecisionMakingParams,
	fatigueRatio float64,
) *model.BiasedResult {
	return &model.BiasedResult{
		DMP: &model.DecisionMakingParams{
			NotConsideredAlternatives: notConsideredAlts,
			ConsideredAlternatives:    consideredAlts,
			Criteria:                  current.Criteria,
			MethodParameters:          current.MethodParameters,
		},
		Props: FatigueResult{
			EffectiveFatigueRatio:     fatigueRatio,
			ConsideredAlternatives:    consideredAlts,
			NotConsideredAlternatives: notConsideredAlts,
		},
	}
}

func blurCriteriaValues(
	alternatives []model.AlternativeWithCriteria,
	criteria []CriterionWithBounding,
	valueGenerator, signGenerator utils.ValueGenerator,
	fatigueRatio float64,
) []model.AlternativeWithCriteria {
	newAlternatives := make([]model.AlternativeWithCriteria, len(alternatives))
	for i, a := range alternatives {
		newWeights := make(model.Weights, len(criteria))
		for _, c := range criteria {
			currentValue := a.CriterionRawValue(&c.criterion)
			eps := currentValue * valueGenerator() * fatigueRatio
			sign := 1.0
			if signGenerator() >= 0.5 {
				sign = -1
			}
			blurredValue := currentValue + (eps * sign)
			boundedBlurredValue := c.bounding.BoundValue(blurredValue)
			newWeights[c.criterion.Id] = boundedBlurredValue
		}
		newAlternatives[i] = *a.WithCriteriaValues(&newWeights)
	}
	return newAlternatives
}

func parseFatigueFuncParams(fun FatigueFunction, parsedProps *FatigueParams) FatigueFunctionParams {
	funParams := fun.BlankParams()
	utils.DecodeToStruct(parsedProps.Params, funParams)
	return funParams
}

func parseProps(props *model.BiasProps) *FatigueParams {
	parsedProps := FatigueParams{}
	utils.DecodeToStruct(*props, &parsedProps)
	return &parsedProps
}

func (f *Fatigue) getFatigueFunction(params *FatigueParams) FatigueFunction {
	for _, fun := range f.functions {
		if fun.Name() == params.Function {
			return fun
		}
	}
	panic(fmt.Errorf("function type for fatigue not defined"))
}
