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

type AspectEliminationEvaluation struct {
	NotSatisfiedThreshold model.Weights `json:"notSatisfiedThreshold"`
	ThresholdsIndex       int           `json:"thresholdsIndex"`
}

type AspectEliminationHeuristicParams struct {
	Function string        `json:"function"`
	Params   interface{}   `json:"params"`
	Seed     int64         `json:"randomSeed"`
	Weights  model.Weights `json:"weights"`
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
	generator := a.generator(params.Seed)
	alternatives := model.ShuffleAlternatives(&dmp.ConsideredAlternatives, generator)
	weights := sortCriteria(dmp, params, generator)
	leftToChoice, result, resultIds, thresholdIndex := checkWithinSatisfactionLevels(weights, alternatives, satisfactionLevels)
	fillRemainingAlternatives(leftToChoice, thresholdIndex, result, resultIds)
	ranking := limited_rationality.PrepareSequentialRanking(result, resultIds)
	return &ranking
}

func sortCriteria(dmp *model.DecisionMakingParams, params AspectEliminationHeuristicParams, generator utils.ValueGenerator) []model.WeightedCriterion {
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
	criteria []model.WeightedCriterion,
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
			for _, a := range leftToChoice {
				if isBellowThreshold(&a, &t, &c.Criterion) {
					leftToChoice = model.RemoveAlternative(leftToChoice, a)
					threshold := makeWeightPair(&t, &c.Criterion)
					resultInsertIndex = updateResult(result, resultInsertIndex, a, thresholdIndex, resultIds, &threshold)
				}
				if len(leftToChoice) <= 1 {
					break thresholds
				}
			}
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