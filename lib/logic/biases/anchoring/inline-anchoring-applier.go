package anchoring

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
)

//go:generate easytags $GOFILE json:camel
type InlineAnchoringApplier struct {
}

const InlineAnchoringApplierName = "inline"

func (i *InlineAnchoringApplier) Identifier() string {
	return InlineAnchoringApplierName
}

type InlineAnchoringApplierParams struct {
	Unbounded bool `json:"unbounded"`
}

func (i *InlineAnchoringApplier) BlankParams() FunctionParams {
	return &InlineAnchoringApplierParams{}
}

func (i *InlineAnchoringApplier) ApplyAnchoring(
	dmp *model.DecisionMakingParams,
	perReferencePointDiffs *[]ReferencePointsDifference,
	criteriaScaling CriteriaScaling,
	params FunctionParams,
	listener *model.BiasListener,
) *model.DecisionMakingParams {
	parsedParams := params.(*InlineAnchoringApplierParams)
	newAlternatives := make([]model.AlternativeWithCriteria, len(*perReferencePointDiffs))
	for i, p := range *perReferencePointDiffs {
		newWeights := arithmeticAverage(p.ReferencePointsDifference)
		for c, scaling := range criteriaScaling {
			difference := newWeights.Fetch(c)
			value := p.Alternative.Criteria.Fetch(c)
			newValue := value + scaling.ValuesRange.Diff()*difference
			if !parsedParams.Unbounded {
				if newValue > scaling.ValuesRange.Max {
					newValue = scaling.ValuesRange.Max
				} else if newValue < scaling.ValuesRange.Min {
					newValue = scaling.ValuesRange.Min
				}
			}
			(*newWeights)[c] = newValue
		}
		newAlternatives[i] = *p.Alternative.WithCriteriaValues(newWeights)
	}
	return &model.DecisionMakingParams{
		ConsideredAlternatives:    newAlternatives,
		NotConsideredAlternatives: dmp.NotConsideredAlternatives,
		Criteria:                  dmp.Criteria,
		MethodParameters:          dmp.MethodParameters,
	}
}

func arithmeticAverage(points []ReferencePointDifference) *model.Weights {
	newWeights := make(model.Weights, len(points[0].Coefficients))
	for i, a := range points {
		for c, v := range a.Coefficients {
			if i == 0 {
				newWeights[c] = v
			} else {
				oldWeight := newWeights.Fetch(c)
				newWeights[c] = oldWeight + v
			}
		}
	}
	referencePointsCount := float64(len(points))
	if referencePointsCount > 1 {
		for c, v := range newWeights {
			newWeights[c] = v / referencePointsCount
		}
	}
	return &newWeights
}
