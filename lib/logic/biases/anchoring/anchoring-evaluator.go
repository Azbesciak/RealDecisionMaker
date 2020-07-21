package anchoring

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
)

//go:generate easytags $GOFILE json:camel
type AnchoringEvaluator interface {
	FunctionBase
	Evaluate(params FunctionParams, difference float64) float64
}

func (a *Anchoring) getAnchoringEvaluatorFunction(params *FunctionDefinition, anchoringType string) AnchoringWithParams {
	for _, fun := range a.anchoringEvaluators {
		if fun.Identifier() == params.Function {
			return AnchoringWithParams{
				fun:    fun,
				params: parseFuncParams(fun, params),
			}
		}
	}
	existing := a.knownAnchoringEvaluatorsNames()
	panic(fmt.Errorf("%s anchoring function '%s' not found in %v", anchoringType, params.Function, existing))
}

func (a *Anchoring) knownAnchoringEvaluatorsNames() []string {
	existing := make([]string, len(a.anchoringEvaluators))
	for i, id := range a.anchoringEvaluators {
		existing[i] = id.Identifier()
	}
	return existing
}

func (a *Anchoring) evaluateAnchoringAlternatives(
	allAlternatives []model.AlternativeWithCriteria,
	parsedProps *AnchoringParams,
	criteria *model.Criteria,
) []model.AlternativeWithCriteria {
	anchoringAlternatives := fetchAnchoringAlternativesWithCriteria(&allAlternatives, &parsedProps.AnchoringAlternatives)
	referencePointsEvaluator := a.getReferencePointsFunction(&parsedProps.ReferencePoints)
	referencePoints := referencePointsEvaluator.Evaluate(parsedProps.ReferencePoints, anchoringAlternatives, criteria)
	return referencePoints
}

func (a *Anchoring) getReferencePointsFunction(params *FunctionDefinition) ReferencePointsEvaluator {
	for _, fun := range a.referencePointsEvaluators {
		if fun.Identifier() == params.Function {
			return fun
		}
	}
	knownReferencePointsEvaluators := a.knownReferencePointsEvaluatorsNames()
	panic(fmt.Errorf(
		"reference points function type '%s' not found in %v",
		params.Function, knownReferencePointsEvaluators,
	))
}

func (a *Anchoring) knownReferencePointsEvaluatorsNames() []string {
	existing := make([]string, len(a.referencePointsEvaluators))
	for i, id := range a.referencePointsEvaluators {
		existing[i] = id.Identifier()
	}
	return existing
}
