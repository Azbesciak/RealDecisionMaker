package majority

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/logic/limited-rationality"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

//go:generate easytags $GOFILE json:camel

type Majority struct {
	generator             utils.SeededValueGenerator
	drawResolvers         []DrawResolver
	currentWinnerResolver CurrentIsWinnerDrawResolver
	newerIsWinnerResolver NewerIsWinnerResolver
}

func NewMajority(generator utils.SeededValueGenerator, drawResolvers []DrawResolver) *Majority {
	if len(drawResolvers) == 0 {
		panic("no draw resolvers for majority heuristic!")
	}
	return &Majority{
		generator:             generator,
		drawResolvers:         drawResolvers,
		newerIsWinnerResolver: NewerIsWinnerResolver{},
		currentWinnerResolver: CurrentIsWinnerDrawResolver{},
	}
}

const eps = 1e-6
const methodName = "majorityHeuristic"

type MajorityHeuristicParams struct {
	Weights                    model.Weights     `json:"weights"`
	CurrentChoice              model.Alternative `json:"currentChoice"`
	RandomSeed                 int64             `json:"randomSeed"`
	RandomAlternativesOrdering bool              `json:"randomAlternativesOrdering"`
	DrawResolution             string            `json:"drawResolution"`
}

func (m *MajorityHeuristicParams) GetCurrentChoice() string {
	return m.CurrentChoice
}

func (m *MajorityHeuristicParams) GetRandomSeed() int64 {
	return m.RandomSeed
}

func (m *MajorityHeuristicParams) IsRandomAlternativesOrdering() bool {
	return m.RandomAlternativesOrdering
}

func (m *Majority) Identifier() string {
	return methodName
}

func (m *Majority) MethodParameters() interface{} {
	return MajorityHeuristicParams{}
}

func (m *Majority) drawResolver(params *MajorityHeuristicParams) DrawResolver {
	if len(params.DrawResolution) == 0 {
		return m.drawResolvers[0]
	}
	for _, r := range m.drawResolvers {
		if r.Identifier() == params.DrawResolution {
			return r
		}
	}
	names := make([]string, len(m.drawResolvers))
	for i, r := range m.drawResolvers {
		names[i] = r.Identifier()
	}
	panic(fmt.Errorf("draw resolution '%s' not found in %v", params.DrawResolution, names))
}

func (m *Majority) Evaluate(dm *model.DecisionMakingParams) *model.AlternativesRanking {
	params := dm.MethodParameters.(MajorityHeuristicParams)
	criteriaWithWeights := dm.Criteria.ZipWithWeights(&params.Weights)
	generator := m.generator(params.RandomSeed)
	current, considered := limited_rationality.GetAlternativesSearchOrder(dm, &params, generator)
	var sameBuffer []model.AlternativeResult
	var worseThanCurrent [][]model.AlternativeResult
	var currentEvaluation model.Weight = 0
	drawResolver := m.drawResolver(&params)
	for _, another := range considered {
		s1, s2 := compare(criteriaWithWeights, &current, &another)
		worseThanCurrent, sameBuffer, current, currentEvaluation =
			m.takeBetter(s1, s2, sameBuffer, another, current, worseThanCurrent, drawResolver, generator)
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
	worseOneLevelThanCurrent := make([]string, 0)
	var result = make(model.AlternativesRanking, 0)
	for _, equivalentEntries := range ranking {
		var sameAlternativesId []string
		for i, r := range equivalentEntries {
			var thisAlternativeWorse = worseOneLevelThanCurrent
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
		worseOneLevelThanCurrent = sameAlternativesId
	}
	result.ReverseOrder()
	return &result
}

func (m *Majority) takeBetter(s1, s2 model.Weight, sameBuffer []model.AlternativeResult,
	another, current model.AlternativeWithCriteria,
	worseThanCurrent [][]model.AlternativeResult,
	resolver DrawResolver,
	generator utils.ValueGenerator,
) ([][]model.AlternativeResult, []model.AlternativeResult, model.AlternativeWithCriteria, model.Weight) {
	currentEvaluation := s1
	if utils.FloatsAreEqual(s1, s2, eps) {
		resolution := resolver.Resolve(s1, s2, sameBuffer, worseThanCurrent, current, another, generator)
		current = resolution.current
		worseThanCurrent = resolution.worseThanCurrent
		sameBuffer = resolution.sameBuffer
	} else if s2 < s1 {
		worseThanCurrent = m.currentWinnerResolver.Resolve(
			s1, s2, sameBuffer, worseThanCurrent, current, another, generator,
		).worseThanCurrent
	} else {
		currentEvaluation = s2
		resolution := m.newerIsWinnerResolver.Resolve(
			s1, s2, sameBuffer, worseThanCurrent, current, another, generator,
		)
		current = resolution.current
		worseThanCurrent = resolution.worseThanCurrent
		sameBuffer = resolution.sameBuffer
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
