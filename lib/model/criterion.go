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
	Id   string        `json:"id"`
	Type CriterionType `json:"type"`
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

func (c *Criteria) SortByWeights(weights Weights) *Criteria {
	criteriaCopy := c.ShallowCopy()
	sort.SliceStable(*criteriaCopy, func(i, j int) bool {
		return criteriaCopy.Weight(weights, i) < criteriaCopy.Weight(weights, j)
	})
	return criteriaCopy
}

func (c *Criteria) Weight(weights Weights, criterionIndex int) Weight {
	return weights[c.Get(criterionIndex).Identifier()]
}

func (c *Criteria) First() Criterion {
	return (*c)[0]
}

func (c Criterion) Multiplier() int8 {
	if c.Type == Cost {
		return -1
	} else {
		return 1
	}
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

type Weight = float64

type WeightedCriterion struct {
	Criterion
	Weight Weight `json:"weight"`
}

type Weights map[string]Weight
