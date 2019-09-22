package logic

type DecisionMakerProperties struct {
	Focus   int
	Fatigue int
}
type CriterionType string

const (
	Gain CriterionType = "gain"
	Cost CriterionType = "cost"
)

type Criterion struct {
	Id   string
	Type CriterionType
}

type Identifiable interface {
	Identifier() string
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

func (c Criterion) multiplier() int8 {
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

type Alternative struct {
	Id string
}

type AlternativeWithCriteria struct {
	Alternative
	Criteria Weights
}
type Weights map[string]Weight

type AlternativeResult struct {
	Alternative AlternativeWithCriteria
	Value       Weight
}
