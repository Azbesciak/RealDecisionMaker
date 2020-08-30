package anchoring

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
)

//go:generate easytags $GOFILE json:camel
type AnchoringApplier interface {
	FunctionBase
	ApplyAnchoring(
		dmp *model.DecisionMakingParams,
		perReferencePointDiffs *[]ReferencePointsDifference,
		boundingsWithScales BoundingsWithScales,
		params FunctionParams,
		listener *model.BiasListener,
	) (*model.DecisionMakingParams, AnchoringApplierResult)
}

type AnchoringApplierResult = interface{}

type ApplierWithParams struct {
	fun    AnchoringApplier
	params FunctionParams
}

func (a *Anchoring) getAnchoringApplier(params *FunctionDefinition) ApplierWithParams {
	for _, fun := range a.anchoringAppliers {
		if fun.Identifier() == params.Function {
			return ApplierWithParams{
				fun:    fun,
				params: parseFuncParams(fun, params),
			}
		}
	}
	existing := a.knownAnchoringAppliersNames()
	panic(fmt.Errorf("anchoring applier function '%s' not found in %v", params.Function, existing))
}

func (a *Anchoring) knownAnchoringAppliersNames() []string {
	existing := make([]string, len(a.anchoringAppliers))
	for i, id := range a.anchoringAppliers {
		existing[i] = id.Identifier()
	}
	return existing
}
