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
	Weights model.Weights     `json:"weights"`
	Current model.Alternative `json:"currentChoice"`
	Seed    int64             `json:"randomSeed"`
}

func (m *MajorityHeuristicParams) CurrentChoice() string {
	return m.Current
}

func (m *MajorityHeuristicParams) RandomSeed() int64 {
	return m.Seed
}

func (m *Majority) Identifier() string {
	return methodName
}

func (m *Majority) MethodParameters() interface{} {
	return model.WeightsParamOnly()
}

func (m *Majority) Evaluate(dm *model.DecisionMakingParams) *model.AlternativesRanking {
	params := dm.MethodParameters.(MajorityHeuristicParams)
	criteriaWithWeights := dm.Criteria.ZipWithWeights(&params.Weights)
	generator := m.generator(params.Seed)
	current, considered := limited_rationality.GetAlternativesSearchOrder(dm, &params, generator)
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
		Evaluation: MajorityEvaluation{
			Value: currentEvaluation,
		},
	})
	worseThanCurrent = append(worseThanCurrent, sameBuffer)
	return prepareRanking(worseThanCurrent)
}

type MajorityEvaluation struct {
	Value                    float64           `json:"value"`
	ComparedWith             model.Alternative `json:"comparedWith"`
	ComparedAlternativeValue float64           `json:"comparedAlternativeValue"`
}

func prepareRanking(ranking [][]model.AlternativeResult) *model.AlternativesRanking {
	worseThanCurrent := make([]string, 0)
	var result = make(model.AlternativesRanking, 0)
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
	result.ReverseOrder()
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
			Evaluation: MajorityEvaluation{
				Value:                    s2,
				ComparedWith:             current.Id,
				ComparedAlternativeValue: s1,
			},
		})
	} else if s2 < s1 {
		worseThanCurrent = append(worseThanCurrent, []model.AlternativeResult{{
			Alternative: another,
			Evaluation: MajorityEvaluation{
				Value:                    s2,
				ComparedWith:             current.Id,
				ComparedAlternativeValue: s1,
			},
		}})
	} else {
		currentEvaluation = s2
		sameBuffer = append(sameBuffer, model.AlternativeResult{
			Alternative: current,
			Evaluation: MajorityEvaluation{
				Value:                    s1,
				ComparedWith:             another.Id,
				ComparedAlternativeValue: s2,
			},
		})
		current = another
		worseThanCurrent = append(worseThanCurrent, sameBuffer)
		sameBuffer = make([]model.AlternativeResult, 0)
	}
	return worseThanCurrent, sameBuffer, current, currentEvaluation
}

func compare(criteriaWithWeights *model.WeightedCriteria, a1, a2 *model.AlternativeWithCriteria) (model.Weight, model.Weight) {
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
