package criteria_splitting

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"math"
)

type CriteriaSplitCondition struct {
	Ratio float64 `json:"ratio"`
	Min   int     `json:"min"`
	Max   int     `json:"max"`
}

func Parse(props *interface{}) *CriteriaSplitCondition {
	parsedProps := CriteriaSplitCondition{Max: math.MaxInt64}
	utils.DecodeToStruct(*props, &parsedProps)
	parsedProps.validate()
	return &parsedProps
}

func (c *CriteriaSplitCondition) validate() {
	if !utils.IsProbability(c.Ratio) {
		panic(fmt.Errorf("'ratio' need to be in range [0,1], got %f", c.Ratio))
	}
	if c.Max < c.Min {
		panic(fmt.Errorf("'max' (%v) is lower than 'min' (%v)", c.Max, c.Min))
	}
}

func (c *CriteriaSplitCondition) SplitCriteriaByOrdering(sortedCriteria *model.Criteria) *CriteriaPartition {
	criteriaCount := len(*sortedCriteria)
	pivot := int(math.Floor(float64(criteriaCount) * c.Ratio))
	if pivot < c.Min {
		pivot = c.Min
	} else if pivot > c.Max {
		pivot = c.Max
	}
	left := (*sortedCriteria)[0:pivot]
	right := (*sortedCriteria)[pivot:]
	return &CriteriaPartition{
		Left:  &left,
		Right: &right,
	}
}

type CriteriaPartition struct {
	Left  *model.Criteria
	Right *model.Criteria
}
