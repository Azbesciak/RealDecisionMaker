package satisfaction

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/limited-rationality"
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/limited-rationality/satisfaction-levels"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

const methodName = "satisfactionHeuristic"

//go:generate easytags $GOFILE json:camel
type Satisfaction struct {
	generator utils.SeededValueGenerator
	functions []satisfaction_levels.SatisfactionLevelsSource
}

func NewSatisfaction(
	generator utils.SeededValueGenerator,
	functions []satisfaction_levels.SatisfactionLevelsSource,
) *Satisfaction {
	return &Satisfaction{generator: generator, functions: functions}
}

func (s *Satisfaction) Identifier() string {
	return methodName
}

type SatisfactionParameters struct {
	Function                   string            `json:"function"`
	Params                     interface{}       `json:"params"`
	RandomSeed                 int64             `json:"randomSeed"`
	CurrentChoice              model.Alternative `json:"currentChoice"`
	RandomAlternativesOrdering bool              `json:"randomAlternativesOrdering"`
}

func (s *SatisfactionParameters) with(params interface{}) SatisfactionParameters {
	return SatisfactionParameters{
		Function:      s.Function,
		Params:        params,
		RandomSeed:    s.RandomSeed,
		CurrentChoice: s.CurrentChoice,
	}
}

type SatisfactionEvaluation struct {
	SatisfiedThresholds model.Weights `json:"satisfiedThresholds"`
	ThresholdsIndex     int           `json:"thresholdsIndex"`
}

func (s *SatisfactionParameters) GetCurrentChoice() string {
	return s.CurrentChoice
}

func (s *SatisfactionParameters) GetRandomSeed() int64 {
	return s.RandomSeed
}

func (s *SatisfactionParameters) IsRandomAlternativesOrdering() bool {
	return s.RandomAlternativesOrdering
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
	satisfactionLevels := satisfaction_levels.Find(params.Function, params.Params, s.functions)
	satisfactionLevels.Initialize(dmp)
	generator := s.generator(params.GetRandomSeed())
	current, considered := limited_rationality.GetAlternativesSearchOrder(dmp, &params, generator)
	leftToChoice, result, resultIds, resultInsertIndex, thresholdIndex := checkWithinSatisfactionLevels(dmp, current, considered, satisfactionLevels)
	fillRemainingAlternatives(leftToChoice, thresholdIndex, resultInsertIndex, result, resultIds, weightsSupplier(dmp))
	ranking := limited_rationality.PrepareSequentialRanking(result, resultIds)
	return &ranking
}

func weightsSupplier(dmp *model.DecisionMakingParams) func() model.Weights {
	return func() model.Weights {
		alternatives := dmp.AllAlternatives()
		weights := make(model.Weights, len(dmp.Criteria))
		for _, c := range dmp.Criteria {
			valRange := model.CriteriaValuesRange(&alternatives, &c)
			if c.IsGain() {
				weights[c.Id] = valRange.Min
			} else {
				weights[c.Id] = valRange.Max
			}
		}
		return weights
	}
}

func fillRemainingAlternatives(
	leftToChoice []model.AlternativeWithCriteria,
	thresholdIndex, resultInsertIndex int,
	result model.AlternativeResults,
	resultIds []model.Alternative,
	lowestThresholdSup func() model.Weights,
) {
	if len(leftToChoice) > 0 {
		thresholdIndex += 1
		lowestThresholds := lowestThresholdSup()
		for _, a := range leftToChoice {
			resultInsertIndex = updateResult(result, resultInsertIndex, a, thresholdIndex, resultIds, &lowestThresholds)
		}
	}
}

func checkWithinSatisfactionLevels(
	dmp *model.DecisionMakingParams,
	current model.AlternativeWithCriteria,
	considered []model.AlternativeWithCriteria,
	satisfactionLevels satisfaction_levels.SatisfactionLevels,
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
				resultInsertIndex = updateResult(result, resultInsertIndex, a, thresholdIndex, resultIds, &t)
			}
		}
		if len(leftToChoice) == 0 {
			break
		}
	}
	return leftToChoice, result, resultIds, resultInsertIndex, thresholdIndex
}

func updateResult(
	result model.AlternativeResults,
	resultInsertIndex int,
	alternative model.AlternativeWithCriteria,
	alternativeValue int,
	resultIds []model.Alternative,
	thresholds *model.Weights,
) int {
	result[resultInsertIndex] = model.AlternativeResult{
		Alternative: alternative,
		Evaluation: SatisfactionEvaluation{
			ThresholdsIndex:     alternativeValue,
			SatisfiedThresholds: *thresholds,
		},
	}
	resultIds[resultInsertIndex] = alternative.Id
	return resultInsertIndex + 1
}

func isGoodEnough(alternative model.AlternativeWithCriteria, thresholds *model.WeightedCriteria) bool {
	for _, v := range *thresholds {
		criterionValue := alternative.CriterionValue(&v.Criterion)
		threshold := float64(v.Multiplier()) * v.Weight
		if criterionValue < threshold {
			return false
		}
	}
	return true
}
