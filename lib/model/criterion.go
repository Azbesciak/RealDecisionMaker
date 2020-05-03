package model

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"sort"
	"strings"
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
	return c.FindWeight(&weights, &(*c)[criterionIndex])
}

type namedWeight struct {
	name   string
	weight Weight
}

func (c *Criteria) FindWeight(weights *Weights, criterion *Criterion) Weight {
	if v, ok := (*weights)[criterion.Id]; !ok {
		criteria := weights.AsKeyValue()
		panic(fmt.Errorf("weight for criterion '%s' not found in criteria %v", criterion.Id, criteria))
	} else {
		return v
	}
}

func (w *Weights) AsKeyValue() []namedWeight {
	criteria := make([]namedWeight, len(*w))
	i := 0
	for crit, value := range *w {
		criteria[i] = namedWeight{crit, value}
		i++
	}
	sort.SliceStable(criteria, func(i, j int) bool {
		return strings.Compare(criteria[i].name, criteria[j].name) < 0
	})
	return criteria
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

type Weight = float64

type WeightedCriterion struct {
	Criterion
	Weight Weight `json:"weight"`
}

func (c *Criteria) ZipWithWeights(weights *Weights) *[]WeightedCriterion {
	weightedCriteria := make([]WeightedCriterion, len(*c))
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

type Weights map[string]Weight

func (w *Weights) Fetch(key string) Weight {
	if v, ok := (*w)[key]; !ok {
		values := w.AsKeyValue()
		panic(fmt.Errorf("criterion %s not found in %v", key, values))
	} else {
		return v
	}
}

func (w *Weights) PreserveOnly(criteria *Criteria) *Weights {
	cpy := make(Weights, len(*criteria))
	for _, c := range *criteria {
		cpy[c.Id] = w.Fetch(c.Id)
	}
	return &cpy
}

func (w *Weights) Merge(other *Weights) *Weights {
	result := make(Weights, len(*other)+len(*w))
	for cryt, weight := range *w {
		result[cryt] = weight
	}
	for cryt, weight := range *other {
		if _, ok := result[cryt]; ok {
			oldWeights := w.AsKeyValue()
			newWeights := other.AsKeyValue()
			panic(fmt.Errorf("criterion '%s' from %v already exists in %v", cryt, oldWeights, newWeights))
		}
		result[cryt] = weight
	}
	return &result
}
