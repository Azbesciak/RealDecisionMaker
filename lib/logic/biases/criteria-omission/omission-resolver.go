package criteria_omission

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type OmissionResolver interface {
	utils.Identifiable
	CriteriaOmissionOrder(
		params *model.DecisionMakingParams,
		props *model.BiasProps,
		listener *model.BiasListener,
	) *model.Criteria
}

const WeakestCriteriaFirst = "weakest"

type WeakestCriteriaOmissionResolver struct {
}

func (w *WeakestCriteriaOmissionResolver) Identifier() string {
	return WeakestCriteriaFirst
}

func (w *WeakestCriteriaOmissionResolver) CriteriaOmissionOrder(
	params *model.DecisionMakingParams,
	_ *model.BiasProps,
	listener *model.BiasListener,
) *model.Criteria {
	return (*listener).RankCriteriaAscending(params).Criteria()
}

const StrongestCriteriaFirst = "strongest"

type StrongestCriteriaOmissionResolver struct {
}

func (w *StrongestCriteriaOmissionResolver) Identifier() string {
	return StrongestCriteriaFirst
}

func (w *StrongestCriteriaOmissionResolver) CriteriaOmissionOrder(
	params *model.DecisionMakingParams,
	_ *model.BiasProps,
	listener *model.BiasListener,
) *model.Criteria {
	ascending := (*listener).RankCriteriaAscending(params).Criteria()
	totalCount := len(*ascending)
	descending := make(model.Criteria, totalCount)
	for i, a := range *ascending {
		descending[totalCount-i-1] = a
	}
	return &descending
}

const RandomCriteria = "random"

type RandomCriteriaOmissionResolver struct {
	Generator utils.SeededValueGenerator
}

func (w *RandomCriteriaOmissionResolver) Identifier() string {
	return RandomCriteria
}

func (w *RandomCriteriaOmissionResolver) CriteriaOmissionOrder(
	params *model.DecisionMakingParams,
	props *model.BiasProps,
	_ *model.BiasListener,
) *model.Criteria {
	parsedProps := parseRandomOmissionProps(props)
	generator := w.Generator(parsedProps.RandomSeed)
	return shuffleCriteria(&params.Criteria, generator)
}

type randomProps struct {
	RandomSeed int64 `json:"randomSeed"`
}

func parseRandomOmissionProps(props *model.BiasProps) *randomProps {
	parsedProps := randomProps{}
	utils.DecodeToStruct(*props, &parsedProps)
	return &parsedProps
}

func shuffleCriteria(criteria *model.Criteria, generator utils.ValueGenerator) *model.Criteria {
	criteriaCount := len(*criteria)
	copied := make(model.Criteria, criteriaCount)
	copy(copied, *criteria)
	for i := criteriaCount - 1; i > 0; i-- { // Fisher–Yates shuffle
		j := int(generator() * float64(i))
		copied[i], copied[j] = copied[j], copied[i]
	}
	return &copied
}

const WeakestByProbabilityCriteriaFirst = "weakestByProbability"

type WeakestByProbabilityCriteriaOmissionResolver struct {
	Generator utils.SeededValueGenerator
}

func (w *WeakestByProbabilityCriteriaOmissionResolver) Identifier() string {
	return WeakestByProbabilityCriteriaFirst
}

func (w *WeakestByProbabilityCriteriaOmissionResolver) CriteriaOmissionOrder(
	params *model.DecisionMakingParams,
	props *model.BiasProps,
	listener *model.BiasListener,
) *model.Criteria {
	parsedProps := parseRandomOmissionProps(props)
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
