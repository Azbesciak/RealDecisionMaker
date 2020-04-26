package satisfaction

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/limited-rationality"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

const methodName = "satisfactionHeuristic"

//go:generate easytags $GOFILE json:camel
type Satisfaction struct {
	generator utils.SeededValueGenerator
	functions []SatisfactionLevelsSource
}

func (s *Satisfaction) Identifier() string {
	return methodName
}

type SatisfactionParameters struct {
	Function string            `json:"function"`
	Params   interface{}       `json:"params"`
	Seed     int64             `json:"randomSeed"`
	Current  model.Alternative `json:"currentChoice"`
}

func (s *SatisfactionParameters) CurrentChoice() string {
	return s.Current
}

func (s *SatisfactionParameters) RandomSeed() int64 {
	return s.Seed
}

func (s *Satisfaction) MethodParameters() interface{} {
	return SatisfactionParameters{}
}

func (s *Satisfaction) ParseParams(dm *model.DecisionMaker) interface{} {
	var params SatisfactionParameters
	utils.DecodeToStruct(dm.MethodParameters, &params)
	return params
}

func (s *Satisfaction) Evaluate(dmp *model.DecisionMakingParams) *model.AlternativesRanking {
	params := dmp.MethodParameters.(SatisfactionParameters)
	satisfactionLevels := s.getSatisfactionLevels(&params)
	satisfactionLevels.Initialize(dmp)
	generator := s.generator(params.RandomSeed())
	current, considered := limitedRationality.GetAlternativesSearchOrder(dmp, &params, generator)
	leftToChoice, result, resultIds, resultInsertIndex, thresholdIndex := checkWithinSatisfactionLevels(dmp, current, considered, satisfactionLevels)
	fillRemainingAlternatives(leftToChoice, thresholdIndex, resultInsertIndex, result, resultIds)
	ranking := prepareRanking(result, resultIds)
	return &ranking
}

func fillRemainingAlternatives(
	leftToChoice []model.AlternativeWithCriteria,
	thresholdIndex, resultInsertIndex int,
	result model.AlternativeResults,
	resultIds []model.Alternative,
) {
	if len(leftToChoice) > 0 {
		thresholdIndex += 1
		for _, a := range leftToChoice {
			resultInsertIndex = updateResult(result, resultInsertIndex, a, thresholdIndex, resultIds)
		}
	}
}

func checkWithinSatisfactionLevels(
	dmp *model.DecisionMakingParams,
	current model.AlternativeWithCriteria,
	considered []model.AlternativeWithCriteria,
	satisfactionLevels SatisfactionLevels,
) ([]model.AlternativeWithCriteria, model.AlternativeResults, []model.Alternative, int, int) {
	leftToChoice := append([]model.AlternativeWithCriteria{current}, considered...)
	result := make(model.AlternativeResults, len(leftToChoice))
	resultIds := make([]model.Alternative, len(leftToChoice))
	resultInsertIndex := 0
	thresholdIndex := -1
	for satisfactionLevels.HasNext() {
		thresholdIndex++
		t := satisfactionLevels.Next()
		thresholds := dmp.Criteria.ZipWithWeights(&t)
		for _, a := range leftToChoice {
			if isGoodEnough(a, thresholds) {
				leftToChoice = model.RemoveAlternative(leftToChoice, a)
				resultInsertIndex = updateResult(result, resultInsertIndex, a, thresholdIndex, resultIds)
			}
		}
		if len(leftToChoice) == 0 {
			break
		}
	}
	return leftToChoice, result, resultIds, resultInsertIndex, thresholdIndex
}

func (s *Satisfaction) getSatisfactionLevels(satisfactionParams *SatisfactionParameters) SatisfactionLevels {
	if len(satisfactionParams.Function) == 0 {
		panic(fmt.Errorf("satisfaction thresholds function not provided"))
	}
	for _, f := range s.functions {
		if f.Name() == satisfactionParams.Function {
			functionParams := f.BlankParams()
			utils.DecodeToStruct(satisfactionParams.Params, functionParams)
			return functionParams
		}
	}
	names := make([]string, len(s.functions))
	for i, f := range s.functions {
		names[i] = f.Name()
	}
	panic(fmt.Errorf("satisfaction thresholds function '%s' not found in functions %v", satisfactionParams.Function, names))
}

func updateResult(
	result model.AlternativeResults,
	resultInsertIndex int,
	alternative model.AlternativeWithCriteria,
	alternativeValue int,
	resultIds []model.Alternative,
) int {
	result[resultInsertIndex] = model.AlternativeResult{
		Alternative: alternative,
		Value:       float64(alternativeValue),
	}
	resultIds[resultInsertIndex] = alternative.Id
	return resultInsertIndex + 1
}

func prepareRanking(result model.AlternativeResults, resultIds []model.Alternative) model.AlternativesRanking {
	resultsCount := len(result)
	ranking := make(model.AlternativesRanking, resultsCount)
	for i, r := range result {
		ranking[i] = model.AlternativesRankEntry{
			AlternativeResult:  r,
			BetterThanOrSameAs: resultIds[i+1:],
		}
	}
	return ranking
}

func isGoodEnough(alternative model.AlternativeWithCriteria, thresholds *[]model.WeightedCriterion) bool {
	for _, v := range *thresholds {
		criterionValue := alternative.CriterionValue(&v.Criterion)
		threshold := float64(v.Multiplier()) * v.Weight
		if criterionValue < threshold {
			return false
		}
	}
	return true
}
