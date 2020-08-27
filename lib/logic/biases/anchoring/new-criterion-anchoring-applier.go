package anchoring

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/model/reference-criterion"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

//go:generate easytags $GOFILE json:camel
type NewCriterionAnchoringApplier struct {
	referenceCriterionManager reference_criterion.ReferenceCriteriaManager
	generator                 utils.SeededValueGenerator
}

func NewNewCriterionAnchoringApplier(
	referenceCriterionManager reference_criterion.ReferenceCriteriaManager,
	generator utils.SeededValueGenerator,
) *NewCriterionAnchoringApplier {
	return &NewCriterionAnchoringApplier{
		referenceCriterionManager: referenceCriterionManager,
		generator:                 generator,
	}
}

const NewCriterionAnchoringApplierName = "newCriterion"

func (n *NewCriterionAnchoringApplier) Identifier() string {
	return NewCriterionAnchoringApplierName
}

type NewCriterionAnchoringApplierParams struct {
	Unbounded  bool  `json:"unbounded"`
	RandomSeed int64 `json:"randomSeed"`
}

func (n *NewCriterionAnchoringApplier) BlankParams() FunctionParams {
	return &utils.Map{
		"unbounded":  false,
		"randomSeed": 0,
	}
}

func addedCriterionName(criteria *model.Criteria, refPointDif model.Alternative) string {
	return criteria.NotUsedName("__anchoring_criterion_" + refPointDif)
}

type additionalCriterionAnchoringState struct {
	listener           *model.BiasListener
	referenceCriterion *model.Criterion
	generator          utils.SeededValueGenerator
	params             NewCriterionAnchoringApplierParams
	currentCriteria    model.Criteria
	methodParams       model.MethodParameters
	addedCriteria      []AddedCriterion
}

func (a *additionalCriterionAnchoringState) newCriterion(ri int, r ReferencePointDifference) *AddedCriterion {
	if len(a.addedCriteria) == ri {
		generator := a.generator(a.params.RandomSeed + int64(ri))
		newCriterionName := addedCriterionName(&a.currentCriteria, r.ReferencePoint)
		criterion := model.Criterion{
			Id:          newCriterionName,
			Type:        a.referenceCriterion.Type,
			ValuesRange: a.referenceCriterion.ValuesRange,
		}
		a.currentCriteria = a.currentCriteria.Add(&criterion)
		newCriterionParams := (*a.listener).OnCriterionAdded(&criterion, a.referenceCriterion, a.methodParams, generator)
		addedCriterion := AddedCriterion{
			Id:                 criterion.Id,
			Type:               criterion.Type,
			ValuesRange:        utils.ValueRange{},
			MethodParameters:   newCriterionParams,
			AlternativesValues: model.Weights{},
		}
		a.addedCriteria = append(a.addedCriteria, addedCriterion)
		a.methodParams = (*a.listener).Merge(a.methodParams, newCriterionParams)
	}
	// cannot return addedCriterion in if because append makes a copy of it.
	return &a.addedCriteria[ri]
}

type AddedCriterion struct {
	Id                 string                 `json:"id"`
	Type               model.CriterionType    `json:"type"`
	ValuesRange        utils.ValueRange       `json:"valuesRange"`
	MethodParameters   model.MethodParameters `json:"methodParameters"`
	AlternativesValues model.Weights          `json:"alternativesValues"`
}

type NewCriterionAnchoringApplierResult struct {
	ReferenceCriterion model.Criterion  `json:"referenceCriterion"`
	AddedCriteria      []AddedCriterion `json:"addedCriteria"`
}

func (n *NewCriterionAnchoringApplier) ApplyAnchoring(
	dmp *model.DecisionMakingParams,
	perReferencePointDiffs *[]ReferencePointsDifference,
	criteriaScaling CriteriaScaling,
	params FunctionParams,
	listener *model.BiasListener,
) (*model.DecisionMakingParams, AnchoringApplierResult) {
	parsedParams := NewCriterionAnchoringApplierParams{}
	utils.DecodeToStruct(params, &parsedParams)
	referenceCriterionProvider := n.referenceCriterionManager.ForParams(&params)
	criteria := *(*listener).RankCriteriaAscending(dmp)
	state := additionalCriterionAnchoringState{
		listener:           listener,
		referenceCriterion: referenceCriterionProvider.Provide(&criteria),
		generator:          n.generator,
		params:             parsedParams,
		addedCriteria:      []AddedCriterion{},
		currentCriteria:    *dmp.Criteria.ShallowCopy(),
		methodParams:       dmp.MethodParameters,
	}
	normalizeCriteriaByTotalValue(criteria)
	if scaling, ok := criteriaScaling[state.referenceCriterion.Id]; !ok {
		panic(fmt.Errorf("scaling for criterion '%s' not found", state.referenceCriterion.Id))
	} else {
		newAlternatives := addAnchoringCriteriaToAlternatives(perReferencePointDiffs, &state, &criteria, &parsedParams, &scaling)
		result := NewCriterionAnchoringApplierResult{
			ReferenceCriterion: *state.referenceCriterion,
			AddedCriteria:      state.addedCriteria,
		}
		return &model.DecisionMakingParams{
			ConsideredAlternatives:    *model.UpdateAlternatives(&dmp.ConsideredAlternatives, &newAlternatives),
			NotConsideredAlternatives: *model.UpdateAlternatives(&dmp.NotConsideredAlternatives, &newAlternatives),
			Criteria:                  state.currentCriteria,
			MethodParameters:          state.methodParams,
		}, result
	}
}

func addAnchoringCriteriaToAlternatives(
	perReferencePointDiffs *[]ReferencePointsDifference,
	state *additionalCriterionAnchoringState,
	criteria *model.WeightedCriteria,
	parsedParams *NewCriterionAnchoringApplierParams,
	scaling *ScaleWithValueRange,
) []model.AlternativeWithCriteria {
	diff := scaling.ValuesRange.Diff() / 2
	newAlternatives := make([]model.AlternativeWithCriteria, len(*perReferencePointDiffs))
	for i, p := range *perReferencePointDiffs {
		alt := p.Alternative
		for ri, r := range p.ReferencePointsDifference {
			anchoringCriterion := state.newCriterion(ri, r)
			criterionValue := 0.0
			for _, c := range *criteria {
				value := r.Coefficients.Fetch(c.Id)
				criterionValue += value * c.Weight
			}
			newValue := diff + diff*criterionValue
			newValue = boundIfRequested(parsedParams.Unbounded, newValue, scaling)
			alt = *alt.WithCriterion(anchoringCriterion.Id, newValue)
			anchoringCriterion.AlternativesValues[alt.Id] = newValue
			valuesRange := &anchoringCriterion.ValuesRange
			if i == 0 {
				valuesRange.Min = newValue
				valuesRange.Max = newValue
			} else {
				if valuesRange.Min >= newValue {
					valuesRange.Min = newValue
				}
				if valuesRange.Max <= newValue {
					valuesRange.Max = newValue
				}
			}
		}
		newAlternatives[i] = alt
	}
	return newAlternatives
}

const _minAllowedWeight = 0.01

func normalizeCriteriaByTotalValue(criteria model.WeightedCriteria) {
	minWeight := criteria[0].Weight
	dif := 0.0
	if minWeight < _minAllowedWeight {
		dif = _minAllowedWeight - minWeight
	}
	total := 0.0
	for i, c := range criteria {
		weight := c.Weight + dif
		total += weight
		criteria[i].Weight = weight
	}
	for i, c := range criteria {
		criteria[i].Weight = c.Weight / total
	}
}
