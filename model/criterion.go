package model

type CriterionType string

const (
	Gain CriterionType = "gain"
	Cost CriterionType = "cost"
)

type Criterion struct {
	Id   string
	Type CriterionType
}

func (c Criterion) Identifier() string {
	return c.Id
}

type Criteria []Criterion

func (c Criteria) Len() int {
	return len(c)
}
func (c Criteria) Get(index int) Identifiable {
	return c[index]
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
	Weight Weight
}

type Weights map[string]Weight
