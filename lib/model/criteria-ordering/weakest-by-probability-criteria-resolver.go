package criteria_ordering

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

const WeakestByProbabilityCriteriaFirst = "weakestByProbability"

type WeakestByProbabilityCriteriaOrderingResolver struct {
	Generator utils.SeededValueGenerator
}

func (w *WeakestByProbabilityCriteriaOrderingResolver) Identifier() string {
	return WeakestByProbabilityCriteriaFirst
}

func (w *WeakestByProbabilityCriteriaOrderingResolver) OrderCriteria(
	params *model.DecisionMakingParams,
	props *model.BiasProps,
	listener *model.BiasListener,
) *model.Criteria {
	parsedProps := parseRandomOrderingProps(props)
	generator := w.Generator(parsedProps.RandomSeed)
	sorted := *(*listener).RankCriteriaAscending(params)
	totalLen := len(sorted)
	result := make(model.Criteria, totalLen)
	if totalLen == 0 {
		return &result
	}
	minWeight := sorted[0].Weight
	maxWeight := sorted[totalLen-1].Weight
	dif := 0.0
	if minWeight <= 1 {
		dif = 1 - minWeight
		minWeight = 1
		maxWeight += dif
	}
	total := 0.0
	for i, s := range sorted {
		tempWeight := s.Weight + dif
		weight := minWeight / tempWeight
		total += weight
		sorted[i].Weight = weight
	}
criterionFind:
	for resultPosition := range result {
		if resultPosition == totalLen-1 {
			result[resultPosition] = sorted[0].Criterion
		}
		randomWeight := generator() * total
		current := 0.0
		for i, c := range sorted {
			current += c.Weight
			if current >= randomWeight {
				result[resultPosition] = c.Criterion
				sorted = append(sorted[:i], sorted[i+1:]...)
				total -= c.Weight
				continue criterionFind
			}
		}
		lastCriterionIndex := totalLen - resultPosition - 1
		lastCriterion := sorted[lastCriterionIndex]
		result[resultPosition] = lastCriterion.Criterion
		sorted = sorted[:len(sorted)-1]
		total -= lastCriterion.Weight
	}
	return &result
}
