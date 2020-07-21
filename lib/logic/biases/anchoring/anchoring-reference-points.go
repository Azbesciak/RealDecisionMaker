package anchoring

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
)

//go:generate easytags $GOFILE json:camel
type ReferencePointsEvaluator interface {
	FunctionBase
	Evaluate(
		params FunctionParams,
		alternatives *[]AnchoringAlternativeWithCriteria,
		criteria *model.Criteria,
	) []model.AlternativeWithCriteria
}

type ReferencePointDifference struct {
	ReferencePoint model.Alternative `json:"referencePoint"`
	Coefficients   model.Weights     `json:"coefficients"`
}

type ReferencePointsDifference struct {
	Alternative               model.AlternativeWithCriteria `json:"alternative"`
	ReferencePointsDifference []ReferencePointDifference    `json:"referencePointsDifference"`
}

func calculateDiffsPerReferencePoint(
	consideredAlternatives []model.AlternativeWithCriteria,
	referencePoints []model.AlternativeWithCriteria,
	criteria *model.Criteria,
	scaleRatios CriteriaScaling,
	loss, gain AnchoringWithParams,
) []ReferencePointsDifference {
	referencePointsDiffs := make([]ReferencePointsDifference, len(consideredAlternatives))
	for ia, a := range consideredAlternatives {
		refPointDiffs := make([]ReferencePointDifference, len(referencePoints))
		for ir, r := range referencePoints {
			refPointDiffs[ir] = calculateReferencePointDiffs(criteria, a, r, scaleRatios, loss, gain)
		}
		referencePointsDiffs[ia] = ReferencePointsDifference{
			Alternative:               a,
			ReferencePointsDifference: refPointDiffs,
		}
	}
	return referencePointsDiffs
}

func calculateReferencePointDiffs(
	criteria *model.Criteria,
	a, r model.AlternativeWithCriteria,
	scaleRatios CriteriaScaling,
	loss, gain AnchoringWithParams,
) ReferencePointDifference {
	referencePointsDiffs := make(model.Weights, len(*criteria))
	for _, c := range *criteria {
		difference := a.CriterionValue(&c) - r.CriterionValue(&c)
		if scaleRatio, ok := scaleRatios[c.Id]; ok {
			scaledDif := difference * scaleRatio.Scale
			var value float64
			if scaledDif > 0 {
				value = gain.fun.Evaluate(gain.params, scaledDif)
			} else {
				value = -loss.fun.Evaluate(loss.params, -scaledDif)
			}
			referencePointsDiffs[c.Id] = value
		} else {
			panic(fmt.Errorf("unknown criterion '%s' in alternative '%s': %v", c.Id, a.Id, a))
		}
	}
	return ReferencePointDifference{
		ReferencePoint: r.Id,
		Coefficients:   referencePointsDiffs,
	}
}
