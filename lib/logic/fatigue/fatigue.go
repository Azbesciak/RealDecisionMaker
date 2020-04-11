package fatigue

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
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

func (f *Fatigue) Identifier() string {
	return "fatigue"
}

func (f *Fatigue) Apply(
	original, current *model.DecisionMakingParams,
	props *model.HeuristicProps,
	listener *model.HeuristicListener,
) *model.HeuristicResult {
	parsedProps := parseProps(props)
	fun := f.getFatigueFunction(parsedProps)
	funParams := parseFatigueFuncParams(fun, parsedProps)
	fatigueRatio := fun.Evaluate(funParams)
	valueGenerator := f.valueGeneratorSource(parsedProps.RandomSeed)
	signGenerator := f.signGeneratorSource(parsedProps.RandomSeed)
	consideredAlts := blurCriteriaValues(
		current.ConsideredAlternatives, current.Criteria,
		valueGenerator, signGenerator, fatigueRatio,
	)
	notConsideredAlts := blurCriteriaValues(
		current.NotConsideredAlternatives, current.Criteria,
		valueGenerator, signGenerator, fatigueRatio,
	)
	return prepareResult(notConsideredAlts, consideredAlts, current, fatigueRatio)

}

func prepareResult(
	notConsideredAlts []model.AlternativeWithCriteria,
	consideredAlts []model.AlternativeWithCriteria,
	current *model.DecisionMakingParams,
	fatigueRatio float64,
) *model.HeuristicResult {
	return &model.HeuristicResult{
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
	criteria model.Criteria,
	valueGenerator, signGenerator utils.ValueGenerator,
	fatigueRatio float64,
) []model.AlternativeWithCriteria {
	newAlternatives := make([]model.AlternativeWithCriteria, len(alternatives))
	for i, a := range alternatives {
		newWeights := make(model.Weights, len(criteria))
		for _, c := range criteria {
			currentValue := a.CriterionRawValue(&c)
			eps := currentValue * valueGenerator() * fatigueRatio
			sign := 1.0
			if signGenerator() >= 0.5 {
				sign = -1
			}
			newWeights[c.Id] = currentValue + (eps * sign)
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

func parseProps(props *model.HeuristicProps) *FatigueParams {
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
