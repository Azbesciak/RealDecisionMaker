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
	Criteria map[string]Weight
}

type AlternativeResult struct {
	Alternative AlternativeWithCriteria
	Value       Weight
}
