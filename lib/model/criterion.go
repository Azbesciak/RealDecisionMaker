package model

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
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

func (c Criterion) Multiplier() int8 {
	if c.Type == Cost {
		return -1
	} else {
		return 1
	}
}

type Weight = float64

type WeightedCriterion struct {
	Criterion
	Weight Weight `json:"weight"`
}

type Weights map[string]Weight
