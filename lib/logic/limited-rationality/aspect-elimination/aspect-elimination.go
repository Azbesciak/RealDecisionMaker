package aspect_elimination

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/limited-rationality"
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/limited-rationality/satisfaction-levels"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"sort"
)

//go:generate easytags $GOFILE json:camel
type AspectEliminationHeuristic struct {
	functions []satisfaction_levels.SatisfactionLevelsSource
	generator utils.SeededValueGenerator
}

func NewAspectEliminationHeuristic(
	functions []satisfaction_levels.SatisfactionLevelsSource,
	generator utils.SeededValueGenerator,
) *AspectEliminationHeuristic {
	return &AspectEliminationHeuristic{
		functions: functions,
		generator: generator,
	}
}

type AspectEliminationEvaluation struct {
	NotSatisfiedThreshold model.Weights `json:"notSatisfiedThreshold"`
	ThresholdsIndex       int           `json:"thresholdsIndex"`
}

type AspectEliminationHeuristicParams struct {
	Function                   string        `json:"function"`
	Params                     interface{}   `json:"params"`
	RandomSeed                 int64         `json:"randomSeed"`
	Weights                    model.Weights `json:"weights"`
	RandomAlternativesOrdering bool          `json:"randomAlternativesOrdering"`
}

func (a *AspectEliminationHeuristicParams) with(params interface{}, weights *model.Weights) AspectEliminationHeuristicParams {
	return AspectEliminationHeuristicParams{
		Function:                   a.Function,
		Params:                     params,
		RandomSeed:                 a.RandomSeed,
		Weights:                    *weights,
		RandomAlternativesOrdering: a.RandomAlternativesOrdering,
	}
}

const methodName = "aspectEliminationHeuristic"

func (a *AspectEliminationHeuristic) Identifier() string {
	return methodName
}

func (a *AspectEliminationHeuristic) MethodParameters() interface{} {
	return AspectEliminationHeuristic{}
}

func (a *AspectEliminationHeuristic) ParseParams(dm *model.DecisionMaker) interface{} {
	var params AspectEliminationHeuristicParams
	utils.DecodeToStruct(dm.MethodParameters, &params)
	return params
}

func (a *AspectEliminationHeuristic) Evaluate(dmp *model.DecisionMakingParams) *model.AlternativesRanking {
	params := dmp.MethodParameters.(AspectEliminationHeuristicParams)
	satisfactionLevels := satisfaction_levels.Find(params.Function, params.Params, a.functions)
	satisfactionLevels.Initialize(dmp)
	generator := a.generator(params.RandomSeed)
	alternatives := limited_rationality.OrderAlternatives(params.RandomAlternativesOrdering, &dmp.ConsideredAlternatives, generator)
	weights := sortCriteria(dmp, params, generator)
	leftToChoice, result, resultIds, thresholdIndex := checkWithinSatisfactionLevels(weights, alternatives, satisfactionLevels)
	fillRemainingAlternatives(leftToChoice, thresholdIndex, result, resultIds)
	ranking := limited_rationality.PrepareSequentialRanking(result, resultIds)
	return &ranking
}

func sortCriteria(dmp *model.DecisionMakingParams, params AspectEliminationHeuristicParams, generator utils.ValueGenerator) model.WeightedCriteria {
	weights := *dmp.Criteria.ZipWithWeights(&params.Weights)
	sort.Slice(weights, func(i, j int) bool {
		w1, w2 := weights[i], weights[j]
		if w1.Weight != w2.Weight {
			return w1.Weight > w2.Weight
		}
		return generator() < 0.5
	})
	return weights
}

func checkWithinSatisfactionLevels(
	criteria model.WeightedCriteria,
	considered *[]model.AlternativeWithCriteria,
	satisfactionLevels satisfaction_levels.SatisfactionLevels,
) ([]model.AlternativeWithCriteria, model.AlternativeResults, []model.Alternative, int) {
	leftToChoice := *considered
	result := make(model.AlternativeResults, len(leftToChoice))
	resultIds := make([]model.Alternative, len(leftToChoice))
	resultInsertIndex := len(leftToChoice) - 1
	thresholdIndex := -1
	if len(leftToChoice) <= 1 {
		return leftToChoice, result, resultIds, thresholdIndex
	}
thresholds:
	for satisfactionLevels.HasNext() {
		thresholdIndex++
		t := satisfactionLevels.Next()
		for _, c := range criteria {
			tempAlternatives := *model.CopyAlternatives(&leftToChoice)
			for _, a := range leftToChoice {
				if isBellowThreshold(&a, &t, &c.Criterion) {
					tempAlternatives = model.RemoveAlternative(tempAlternatives, a)
					threshold := makeWeightPair(&t, &c.Criterion)
					resultInsertIndex = updateResult(result, resultInsertIndex, a, thresholdIndex, resultIds, &threshold)
				}
				if len(tempAlternatives) <= 1 {
					leftToChoice = tempAlternatives
					break thresholds
				}
			}
			leftToChoice = tempAlternatives
		}
	}
	return leftToChoice, result, resultIds, thresholdIndex
}

func makeWeightPair(weights *model.Weights, criterion *model.Criterion) model.Weights {
	res := make(model.Weights, 1)
	res[criterion.Id] = (*weights)[criterion.Id]
	return res
}

func isBellowThreshold(a *model.AlternativeWithCriteria, thresholds *model.Weights, criterion *model.Criterion) bool {
	return a.CriterionValue(criterion) < (*thresholds)[criterion.Id]*float64(criterion.Multiplier())
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
		Evaluation: AspectEliminationEvaluation{
			ThresholdsIndex:       alternativeValue,
			NotSatisfiedThreshold: *thresholds,
		},
	}
	resultIds[resultInsertIndex] = alternative.Id
	return resultInsertIndex - 1
}

func fillRemainingAlternatives(
	leftToChoice []model.AlternativeWithCriteria,
	thresholdIndex int,
	result model.AlternativeResults,
	resultIds []model.Alternative,
) {
	if len(leftToChoice) > 0 {
		thresholdIndex += 1
		lowestThresholds := make(model.Weights, 0)
		for i, a := range leftToChoice {
			updateResult(result, i, a, thresholdIndex, resultIds, &lowestThresholds)
		}
	}
}
