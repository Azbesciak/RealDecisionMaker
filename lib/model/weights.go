package model

import (
	"fmt"
	"sort"
	"strings"
)

type Weight = float64

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

type namedWeight struct {
	name   string
	weight Weight
}
