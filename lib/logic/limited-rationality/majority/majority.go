package majority

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/limited-rationality"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

//go:generate easytags $GOFILE json:camel

type Majority struct {
	generator utils.SeededValueGenerator
}

const eps = 1e-6
const methodName = "majorityHeuristic"

type MajorityHeuristicParams struct {
	weights       model.Weights
	currentChoice model.Alternative
	randomSeed    int64
}

func (m *MajorityHeuristicParams) CurrentChoice() string {
	return m.currentChoice
}

func (m *MajorityHeuristicParams) RandomSeed() int64 {
	return m.randomSeed
}

func (m *Majority) Identifier() string {
	return methodName
}

func (m *Majority) MethodParameters() interface{} {
	return model.WeightsParamOnly()
}

func (m *Majority) Evaluate(dm *model.DecisionMakingParams) *model.AlternativesRanking {
	params := dm.MethodParameters.(MajorityHeuristicParams)
	criteriaWithWeights := dm.Criteria.ZipWithWeights(&params.weights)
	generator := m.generator(params.randomSeed)
	current, considered := limitedRationality.GetAlternativesSearchOrder(dm, &params, generator)
	var sameBuffer []model.AlternativeResult
	var worseThanCurrent [][]model.AlternativeResult
	var currentEvaluation model.Weight = 0
	for _, another := range considered {
		s1, s2 := compare(criteriaWithWeights, &current, &another)
		worseThanCurrent, sameBuffer, current, currentEvaluation =
			takeBetter(s1, s2, sameBuffer, another, current, worseThanCurrent)
	}
	sameBuffer = append(sameBuffer, model.AlternativeResult{
		Alternative: current,
		Value:       currentEvaluation,
	})
	worseThanCurrent = append(worseThanCurrent, sameBuffer)
	return prepareRanking(worseThanCurrent)
}

func prepareRanking(ranking [][]model.AlternativeResult) *model.AlternativesRanking {
	worseThanCurrent := make([]string, 0)
	var result model.AlternativesRanking
	for _, equivalentEntries := range ranking {
		var sameAlternativesId []string
		for i, r := range equivalentEntries {
			var thisAlternativeWorse = worseThanCurrent
			sameAlternativesId = append(sameAlternativesId, r.Alternative.Id)
			for j, a := range equivalentEntries {
				if i != j {
					thisAlternativeWorse = append(thisAlternativeWorse, a.Alternative.Id)
				}
			}
			result = append(result, model.AlternativesRankEntry{
				AlternativeResult:  r,
				BetterThanOrSameAs: thisAlternativeWorse,
			})
		}
		worseThanCurrent = append(worseThanCurrent, sameAlternativesId...)
	}
	return &result
}

func takeBetter(s1, s2 model.Weight, sameBuffer []model.AlternativeResult,
	another, current model.AlternativeWithCriteria,
	worseThanCurrent [][]model.AlternativeResult,
) ([][]model.AlternativeResult, []model.AlternativeResult, model.AlternativeWithCriteria, model.Weight) {
	currentEvaluation := s1
	if utils.FloatsAreEqual(s1, s2, eps) {
		sameBuffer = append(sameBuffer, model.AlternativeResult{
			Alternative: another,
			Value:       s2,
		})
	} else if s2 < s1 {
		worseThanCurrent = append(worseThanCurrent, []model.AlternativeResult{{
			Alternative: another,
			Value:       s2,
		}})
	} else {
		currentEvaluation = s2
		sameBuffer = append(sameBuffer, model.AlternativeResult{
			Alternative: current,
			Value:       s1,
		})
		current = another
		worseThanCurrent = append(worseThanCurrent, sameBuffer)
		sameBuffer = make([]model.AlternativeResult, 0)
	}
	return worseThanCurrent, sameBuffer, current, currentEvaluation
}

func compare(criteriaWithWeights *[]model.WeightedCriterion, a1, a2 *model.AlternativeWithCriteria) (model.Weight, model.Weight) {
	a1Score := 0.0
	a2Score := 0.0
	for _, criterion := range *criteriaWithWeights {
		v1 := a1.CriterionValue(&criterion.Criterion)
		v2 := a2.CriterionValue(&criterion.Criterion)
		if utils.FloatsAreEqual(v1, v2, eps) {
			continue
		} else if v1 > v2 {
			a1Score += criterion.Weight
		} else {
			a2Score += criterion.Weight
		}
	}
	return a1Score, a2Score
}

func (m *Majority) ParseParams(dm *model.DecisionMaker) interface{} {
	var params MajorityHeuristicParams
	utils.DecodeToStruct(dm.MethodParameters, &params)
	return params
}
