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
	Unbounded            bool `json:"unbounded"`
	ApplyOnNotConsidered bool `json:"applyOnNotConsidered"`
}

func (i *InlineAnchoringApplier) BlankParams() FunctionParams {
	return &InlineAnchoringApplierParams{}
}

type InlineAnchoringApplierResult struct {
	AppliedDifferences []model.AlternativeWithCriteria `json:"appliedDifferences"`
}

func (i *InlineAnchoringApplier) ApplyAnchoring(
	dmp *model.DecisionMakingParams,
	perReferencePointDiffs *[]ReferencePointsDifference,
	criteriaScaling CriteriaScaling,
	params FunctionParams,
	listener *model.BiasListener,
) (*model.DecisionMakingParams, AnchoringApplierResult) {
	parsedParams := params.(*InlineAnchoringApplierParams)
	newAlternatives := make([]model.AlternativeWithCriteria, len(*perReferencePointDiffs))
	appliedDifferences := make([]model.AlternativeWithCriteria, len(*perReferencePointDiffs))
	for i, p := range *perReferencePointDiffs {
		newWeights := *arithmeticAverage(p.ReferencePointsDifference)
		differences := make(model.Weights, len(criteriaScaling))
		for c, scaling := range criteriaScaling {
			difference := newWeights.Fetch(c)
			value := p.Alternative.Criteria.Fetch(c)
			newValue := value + scaling.ValuesRange.Diff()*difference
			newValue = boundIfRequested(parsedParams.Unbounded, newValue, scaling)
			differences[c] = newValue - value
			newWeights[c] = newValue
		}
		newAlternatives[i] = *p.Alternative.WithCriteriaValues(&newWeights)
		appliedDifferences[i] = *p.Alternative.WithCriteriaValues(&differences)
	}
	notConsidered := dmp.NotConsideredAlternatives
	result := InlineAnchoringApplierResult{
		AppliedDifferences: appliedDifferences,
	}
	if parsedParams.ApplyOnNotConsidered {
		notConsidered = *model.UpdateAlternatives(&notConsidered, &newAlternatives)
	} else {
		result.AppliedDifferences = *model.UpdateAlternatives(&dmp.ConsideredAlternatives, &appliedDifferences)
	}
	return &model.DecisionMakingParams{
		ConsideredAlternatives:    *model.UpdateAlternatives(&dmp.ConsideredAlternatives, &newAlternatives),
		NotConsideredAlternatives: notConsidered,
		Criteria:                  dmp.Criteria,
		MethodParameters:          dmp.MethodParameters,
	}, result
}

func boundIfRequested(unbounded bool, newValue model.Weight, scaling ScaleWithValueRange) model.Weight {
	if !unbounded {
		if newValue > scaling.ValuesRange.Max {
			newValue = scaling.ValuesRange.Max
		} else if newValue < scaling.ValuesRange.Min {
			newValue = scaling.ValuesRange.Min
		}
	}
	return newValue
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
