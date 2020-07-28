package anchoring

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

//go:generate easytags $GOFILE json:camel

type AnchoringParams struct {
	AnchoringAlternatives []AnchoringAlternative `json:"anchoringAlternatives"`
	Loss                  FunctionDefinition     `json:"loss"`
	Gain                  FunctionDefinition     `json:"gain"`
	ReferencePoints       FunctionDefinition     `json:"referencePoints"`
	Applier               FunctionDefinition     `json:"applier"`
}

type AnchoringAlternative struct {
	Alternative model.Alternative `json:"alternative"`
	Coefficient float64           `json:"coefficient"`
}

type AnchoringAlternativeWithCriteria struct {
	Alternative model.AlternativeWithCriteria `json:"alternative"`
	Coefficient float64                       `json:"coefficient"`
}

type Anchoring struct {
	anchoringEvaluators       []AnchoringEvaluator
	referencePointsEvaluators []ReferencePointsEvaluator
	anchoringAppliers         []AnchoringApplier
}

func NewAnchoring(
	anchoringEvaluators []AnchoringEvaluator,
	referencePointsEvaluators []ReferencePointsEvaluator,
	anchoringAppliers []AnchoringApplier,
) *Anchoring {
	return &Anchoring{
		anchoringEvaluators:       anchoringEvaluators,
		referencePointsEvaluators: referencePointsEvaluators,
		anchoringAppliers:         anchoringAppliers,
	}
}

type AnchoringWithParams struct {
	fun    AnchoringEvaluator
	params FunctionParams
}

const BiasName = "anchoring"

func (a *Anchoring) Identifier() string {
	return BiasName
}

func (a *Anchoring) Apply(
	_, current *model.DecisionMakingParams,
	props *model.BiasProps,
	listener *model.BiasListener,
) *model.BiasedResult {
	parsedProps := parseProps(props)
	loss := a.getAnchoringEvaluatorFunction(&parsedProps.Loss, "loss")
	gain := a.getAnchoringEvaluatorFunction(&parsedProps.Gain, "gain")
	applier := a.getAnchoringApplier(&parsedProps.Applier)
	allAlternatives := current.AllAlternatives()
	referencePoints := a.evaluateAnchoringAlternatives(allAlternatives, parsedProps, &current.Criteria)
	criteriaScaling := evaluatePerCriterionNormalizationScaleRatio(&current.Criteria, allAlternatives)
	perReferencePointsDiffs := calculateDiffsPerReferencePoint(
		current.ConsideredAlternatives, referencePoints, &current.Criteria, criteriaScaling, loss, gain,
	)
	newDmp := applier.fun.ApplyAnchoring(current, &perReferencePointsDiffs, criteriaScaling, applier.params, listener)
	return &model.BiasedResult{
		DMP: newDmp,
		Props: AnchoringResult{
			ReferencePoints:               referencePoints,
			CriteriaScaling:               criteriaScaling,
			PerReferencePointsDifferences: perReferencePointsDiffs,
		},
	}
}

type AnchoringResult struct {
	ReferencePoints               []model.AlternativeWithCriteria
	CriteriaScaling               CriteriaScaling
	PerReferencePointsDifferences []ReferencePointsDifference
}

func parseProps(props *model.BiasProps) *AnchoringParams {
	parsedProps := AnchoringParams{}
	utils.DecodeToStruct(*props, &parsedProps)
	checkAnchoringAlternatives(props, &parsedProps)
	return &parsedProps
}

func checkAnchoringAlternatives(props *model.BiasProps, params *AnchoringParams) {
	if len(params.AnchoringAlternatives) == 0 {
		panic(fmt.Errorf("no anchoring alternatives passed"))
	}
	asMap := (*props).(utils.Map)
	anchoringAltsObj := asMap["anchoringAlternatives"]
	anchoring, ok := anchoringAltsObj.([]map[string]interface{})
	if !ok {
		for _, p := range params.AnchoringAlternatives {
			p.Coefficient = 1
		}
	} else {
		for i, a := range anchoring {
			if _, ok := a["coefficient"]; !ok {
				(*params).AnchoringAlternatives[i].Coefficient = 1
			}
		}
	}
}

func fetchAnchoringAlternativesWithCriteria(alternatives *[]model.AlternativeWithCriteria, anchoringAlternatives *[]AnchoringAlternative) *[]AnchoringAlternativeWithCriteria {
	alternativesCount := len(*anchoringAlternatives)
	result := make([]AnchoringAlternativeWithCriteria, alternativesCount)
	for i, a := range *anchoringAlternatives {
		alternativeWithCriteria := model.FetchAlternative(alternatives, a.Alternative)
		result[i] = AnchoringAlternativeWithCriteria{
			Alternative: alternativeWithCriteria,
			Coefficient: a.Coefficient,
		}
	}
	return &result
}
