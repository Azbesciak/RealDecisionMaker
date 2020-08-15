package model

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"sort"
)

//go:generate easytags $GOFILE json:camel
type CriterionType string

const (
	Gain CriterionType = "gain"
	Cost CriterionType = "cost"
)

type Criterion struct {
	Id          string            `json:"id"`
	Type        CriterionType     `json:"type"`
	ValuesRange *utils.ValueRange `json:"valuesRange,omitempty"`
}

func (c Criterion) Identifier() string {
	return c.Id
}

type Criteria []Criterion

func (c *Criteria) Len() int {
	return len(*c)
}
func (c *Criteria) Get(index int) utils.Identifiable {
	return (*c)[index]
}

func (c *Criteria) ShallowCopy() *Criteria {
	criteriaCopy := make(Criteria, len(*c))
	copy(criteriaCopy, *c)
	return &criteriaCopy
}

func (c *Criteria) Validate() {
	criteriaSet := make(map[string]bool)
	for i, criterion := range *c {
		if _, ok := criteriaSet[criterion.Id]; ok {
			panic(fmt.Errorf("criterion '%s' [index %d] is not unique", criterion.Id, i))
		}
		if criterion.ValuesRange != nil && criterion.ValuesRange.Max <= criterion.ValuesRange.Min {
			panic(fmt.Errorf("criterion '%s' [index %d] has invalid value range %v: min must be lower than max",
				criterion.Id, i, *criterion.ValuesRange,
			))
		}
		criteriaSet[criterion.Id] = true
	}
}

func (c *Criteria) SortByWeights(weights Weights) *WeightedCriteria {
	result := make(WeightedCriteria, len(weights))
	for i, criterion := range *c {
		result[i] = WeightedCriterion{
			Criterion: criterion,
			Weight:    c.Weight(weights, i),
		}
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Weight < result[j].Weight
	})
	return &result
}

func (c *Criteria) Weight(weights Weights, criterionIndex int) Weight {
	return c.FindWeight(&weights, &(*c)[criterionIndex])
}

func (c *Criteria) FindWeight(weights *Weights, criterion *Criterion) Weight {
	if v, ok := (*weights)[criterion.Id]; !ok {
		criteria := weights.AsKeyValue()
		panic(fmt.Errorf("weight for criterion '%s' not found in criteria %v", criterion.Id, criteria))
	} else {
		return v
	}
}

func (c *Criteria) First() Criterion {
	return (*c)[0]
}

func (c *Criterion) Multiplier() int8 {
	if c.Type == Cost {
		return -1
	} else {
		return 1
	}
}

func (c *Criterion) IsGain() bool {
	return c.Multiplier() == 1
}

func (c *Criteria) Names() *[]string {
	result := make([]string, len(*c))
	for i, crit := range *c {
		result[i] = crit.Id
	}
	return &result
}

func (c *Criteria) Add(criterion *Criterion) Criteria {
	for _, crit := range *c {
		if crit.Id == criterion.Id {
			panic(fmt.Errorf("cannot add criterion '%v' - already exists in criteria: %v", *criterion, *c))
		}
	}
	return append(*c, *criterion)
}

type WeightedCriterion struct {
	Criterion
	Weight Weight `json:"weight"`
}

type WeightedCriteria []WeightedCriterion

func (w *WeightedCriteria) Criteria() *Criteria {
	result := make(Criteria, len(*w))
	for i, c := range *w {
		result[i] = c.Criterion
	}
	return &result
}

func (c *Criteria) ZipWithWeights(weights *Weights) *WeightedCriteria {
	weightedCriteria := make(WeightedCriteria, len(*c))
	for i, crit := range *c {
		value := c.FindWeight(weights, &crit)
		weightedCriteria[i] = WeightedCriterion{
			Criterion: crit,
			Weight:    value,
		}
	}
	return &weightedCriteria
}

func (c *WeightedCriterion) AsWeights() *Weights {
	return &Weights{c.Id: c.Weight}
}
