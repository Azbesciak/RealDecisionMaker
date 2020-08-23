package criteria_ordering

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

const RandomCriteria = "random"

type RandomCriteriaOmissionResolver struct {
	Generator utils.SeededValueGenerator
}

func (w *RandomCriteriaOmissionResolver) Identifier() string {
	return RandomCriteria
}

func (w *RandomCriteriaOmissionResolver) OrderCriteria(
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
