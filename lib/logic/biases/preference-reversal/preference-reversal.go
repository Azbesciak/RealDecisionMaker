package preference_reversal

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/model/criteria-ordering"
	"github.com/Azbesciak/RealDecisionMaker/lib/model/criteria-splitting"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type PreferenceReversal struct {
	orderingResolvers []criteria_ordering.CriteriaOrderingResolver
}

func NewPreferenceReversal(
	orderingResolvers []criteria_ordering.CriteriaOrderingResolver,
) *PreferenceReversal {
	return &PreferenceReversal{
		orderingResolvers: orderingResolvers,
	}
}

// + CriteriaSplitCondition
// + CriteriaOrdering
type PreferenceReversalParams struct {
}

const BiasName = "preferenceReversal"

func (p *PreferenceReversal) Identifier() string {
	return BiasName
}

func (p *PreferenceReversal) Apply(
	original, current *model.DecisionMakingParams,
	props *model.BiasProps,
	listener *model.BiasListener,
) *model.BiasedResult {
	ordering, condition := parseProps(props)
	resolver := criteria_ordering.FetchOrderingResolver(&p.orderingResolvers, ordering)
	sortedCriteria := resolver.OrderCriteria(current, props, listener)
	splitted := condition.SplitCriteriaByOrdering(sortedCriteria)
	criteriaToReverse := getCriteriaToReverse(splitted.Left, original, current)
	updatedAlternatives := updateAlternativesWithReversedCriteriaValues(criteriaToReverse, current)
	result := prepareReverseResult(criteriaToReverse, updatedAlternatives)
	return &model.BiasedResult{
		DMP: &model.DecisionMakingParams{
			NotConsideredAlternatives: *updatedAlternatives.notConsideredAlternatives,
			ConsideredAlternatives:    *updatedAlternatives.consideredAlternatives,
			Criteria:                  current.Criteria,
			MethodParameters:          current.MethodParameters,
		},
		Props: PreferenceReversalResult{
			ReversedPreferenceCriteria: result,
		},
	}
}

func prepareReverseResult(
	criteriaToReverse *[]criterionToReverse,
	reverseResult *criterionReversalResult,
) []ReversedPreferenceCriterion {
	result := make([]ReversedPreferenceCriterion, len(*criteriaToReverse))
	for i, c := range *criteriaToReverse {
		reversal := (*reverseResult.alternativesValues)[i]
		result[i] = ReversedPreferenceCriterion{
			Id:                 c.criterion.Id,
			Type:               c.criterion.Type,
			AlternativesValues: reversal,
			ValuesRange:        *c.valRange,
		}
	}
	return result
}

func parseProps(props *model.BiasProps) (*criteria_ordering.CriteriaOrdering, *criteria_splitting.CriteriaSplitCondition) {
	ordering := criteria_ordering.Parse(props)
	splittingProps := criteria_splitting.Parse(props)
	return ordering, splittingProps
}

func getCriteriaToReverse(
	criteriaToReverse *model.Criteria,
	_, currentParams *model.DecisionMakingParams,
) *[]criterionToReverse {
	allAlternatives := currentParams.AllAlternatives()
	result := make([]criterionToReverse, len(*criteriaToReverse))
	for i, c := range *criteriaToReverse {
		valRange := model.CriteriaValuesRange(&allAlternatives, &c)
		result[i] = criterionToReverse{
			criterion: c,
			valRange:  valRange,
		}
	}
	return &result
}

type criterionToReverse struct {
	criterion model.Criterion
	valRange  *utils.ValueRange
}

func updateAlternativesWithReversedCriteriaValues(
	criteriaToReverse *[]criterionToReverse,
	resParams *model.DecisionMakingParams,
) *criterionReversalResult {
	sortedAlternatives, alternativesValues := reverseCriteriaForEachAlternative(criteriaToReverse, resParams)
	return &criterionReversalResult{
		notConsideredAlternatives: model.UpdateAlternatives(&resParams.NotConsideredAlternatives, sortedAlternatives),
		consideredAlternatives:    model.UpdateAlternatives(&resParams.ConsideredAlternatives, sortedAlternatives),
		alternativesValues:        alternativesValues,
	}
}

func reverseCriteriaForEachAlternative(
	criteriaToReverse *[]criterionToReverse,
	resParams *model.DecisionMakingParams,
) (*[]model.AlternativeWithCriteria, *[]model.Weights) {
	allAlternatives := resParams.AllAlternatives()
	alternativesValues := make([]model.Weights, len(*criteriaToReverse))
	for i := range alternativesValues {
		alternativesValues[i] = make(model.Weights, len(allAlternatives))
	}
	for i, a := range allAlternatives {
		newCriteria := a.Criteria.Copy()
		for ic, c := range *criteriaToReverse {
			currentValue := newCriteria.Fetch(c.criterion.Id)
			newValue := c.valRange.Max - currentValue + c.valRange.Min
			alternativesValues[ic][a.Id] = newValue
			(*newCriteria)[c.criterion.Id] = newValue
		}
		allAlternatives[i] = *a.WithCriteriaValues(newCriteria)
	}
	return &allAlternatives, &alternativesValues
}

func (p *PreferenceReversal) getCriterionValueRange(originalParams *model.DecisionMakingParams, referenceCriterion *model.Criterion) *utils.ValueRange {
	allAlternatives := originalParams.AllAlternatives()
	valRange := model.CriteriaValuesRange(&allAlternatives, referenceCriterion)
	return valRange
}

type PreferenceReversalResult struct {
	ReversedPreferenceCriteria []ReversedPreferenceCriterion `json:"reversedPreferenceCriteria"`
}

type ReversedPreferenceCriterion struct {
	Id                 string              `json:"id"`
	Type               model.CriterionType `json:"type"`
	ValuesRange        utils.ValueRange    `json:"valuesRange"`
	AlternativesValues model.Weights       `json:"alternativesValues"`
}

type criterionReversalResult struct {
	notConsideredAlternatives *[]model.AlternativeWithCriteria
	consideredAlternatives    *[]model.AlternativeWithCriteria
	alternativesValues        *[]model.Weights
}
