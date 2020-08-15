package model

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"math"
	"sort"
)

//go:generate easytags $GOFILE json:camel
type AlternativeResult struct {
	Alternative AlternativeWithCriteria `json:"alternative"`
	Evaluation  interface{}             `json:"evaluation"`
}

func ValueAlternativeResult(alternative *AlternativeWithCriteria, value float64) *AlternativeResult {
	return &AlternativeResult{
		Alternative: *alternative,
		Evaluation:  EvaluationSingleValue{Value: value},
	}
}

type EvaluationSingleValue struct {
	Value float64 `json:"value"`
}

func (a *AlternativeResult) String() string {
	return fmt.Sprintf("alternative: %s, evaluation: %v", a.Alternative.Id, a.Evaluation)
}

func (a *AlternativeResult) Identifier() string {
	return a.Alternative.Id
}

const roundPrecision = 1e8

func (a *AlternativeResult) rounded() *AlternativeResult {
	return &AlternativeResult{
		Alternative: a.Alternative,
		Evaluation:  EvaluationSingleValue{math.Round(a.Value()*roundPrecision) / roundPrecision},
	}
}

func (a *AlternativeResult) Value() float64 {
	value, ok := a.Evaluation.(EvaluationSingleValue)
	if !ok {
		panic(fmt.Errorf("evaluation must be instance of "))
	}
	return value.Value
}

type AlternativeResults []AlternativeResult

func (a *AlternativeResults) Len() int {
	return len(*a)
}

func (a *AlternativeResults) Less(i, j int) bool {
	a1, a2 := (*a)[i], (*a)[j]
	v1, v2 := a1.Value(), a2.Value()
	if v1 == v2 {
		return a1.Alternative.Id < a2.Alternative.Id
	}
	return v1 > v2
}

func (a *AlternativeResults) Swap(i, j int) {
	(*a)[i], (*a)[j] = (*a)[j], (*a)[i]
}

func (a *AlternativeResults) Ranking() *AlternativesRanking {
	alternativesNum := len(*a)
	alternativeResults := make(AlternativeResults, alternativesNum)
	for i, alt := range *a {
		alternativeResults[i] = *alt.rounded()
	}
	sort.Sort(&alternativeResults)
	ranking := make(AlternativesRanking, alternativesNum)
	for i, r := range alternativeResults {
		ranking[i] = *r.positionInRanking(&alternativeResults)
	}
	return &ranking
}

func (a *AlternativeResult) positionInRanking(allAlternatives *AlternativeResults) *AlternativesRankEntry {
	var betterThanOrSameAs = Alternatives{}
	wasLowerValueFound := false
	nextLowerThanAltValue := a.Value()
	for _, r := range *allAlternatives {
		if r.Value() == a.Value() && r.Identifier() != a.Identifier() {
			betterThanOrSameAs = append(betterThanOrSameAs, r.Identifier())
		} else if r.Value() < a.Value() {
			if !wasLowerValueFound {
				wasLowerValueFound = true
				nextLowerThanAltValue = r.Value()
			}
			if r.Value() < nextLowerThanAltValue {
				break
			}
			betterThanOrSameAs = append(betterThanOrSameAs, r.Identifier())
		}
	}
	return &AlternativesRankEntry{
		AlternativeResult:  *a,
		BetterThanOrSameAs: betterThanOrSameAs,
	}
}

type Alternative = string

type AlternativeWithCriteria struct {
	Id       Alternative `json:"id"`
	Criteria Weights     `json:"criteria"`
}

func (a *AlternativeWithCriteria) String() string {
	return fmt.Sprintf("alternative: %s, criteria: %v", a.Id, a.Criteria)
}

func (a *AlternativeWithCriteria) CriterionValue(criterion *Criterion) Weight {
	return a.CriterionRawValue(criterion) * Weight(criterion.Multiplier())
}

func (a *AlternativeWithCriteria) WithCriteriaOnly(criteria *Criteria) *AlternativeWithCriteria {
	newCriteria := make(Weights, len(*criteria))
	for _, c := range *criteria {
		newCriteria[c.Id] = a.CriterionRawValue(&c)
	}
	return a.WithCriteriaValues(&newCriteria)
}

func (a *AlternativeWithCriteria) WithCriteriaValues(criteriaValues *Weights) *AlternativeWithCriteria {
	return &AlternativeWithCriteria{
		Id:       a.Id,
		Criteria: *criteriaValues,
	}
}

func (a *AlternativeWithCriteria) CriterionRawValue(criterion *Criterion) Weight {
	weight, ok := a.Criteria[criterion.Id]
	if !ok {
		panic(fmt.Errorf("alternative '%s' does not have value for criterion '%s'", a.Id, criterion.Id))
	}
	return weight
}

func SortAlternativesByName(alternatives *[]AlternativeWithCriteria) *[]AlternativeWithCriteria {
	res := make([]AlternativeWithCriteria, len(*alternatives))
	copy(res, *alternatives)
	sort.Slice(res, func(i, j int) bool {
		return res[i].Id < res[j].Id
	})
	return &res
}

func ShuffleAlternatives(alternatives *[]AlternativeWithCriteria, generator utils.ValueGenerator) *[]AlternativeWithCriteria {
	alternativesCount := len(*alternatives)
	copied := make([]AlternativeWithCriteria, alternativesCount)
	copy(copied, *alternatives)
	for i := alternativesCount - 1; i > 0; i-- { // Fisher–Yates shuffle
		j := int(generator() * float64(i))
		copied[i], copied[j] = copied[j], copied[i]
	}
	return &copied
}

func CopyAlternatives(alternatives *[]AlternativeWithCriteria) *[]AlternativeWithCriteria {
	result := make([]AlternativeWithCriteria, len(*alternatives))
	copy(result, *alternatives)
	return &result
}

func AddCriterionToAlternatives(
	alternatives *[]AlternativeWithCriteria,
	newCriterion *Criterion,
	valueProvider func(alt *AlternativeWithCriteria) Weight,
) *[]AlternativeWithCriteria {
	newAlts := make([]AlternativeWithCriteria, len(*alternatives))
	for i, a := range *alternatives {
		newValue := valueProvider(&a)
		newAlts[i] = *a.WithCriterion(newCriterion.Id, newValue)
	}
	return &newAlts
}

func (a *AlternativeWithCriteria) WithCriterion(name string, value Weight) *AlternativeWithCriteria {
	if _, ok := a.Criteria[name]; ok {
		panic(fmt.Errorf("cannot add new criterion '%s' because it already exist in alternative %v", name, *a))
	}
	criteria := make(Weights, len(a.Criteria)+1)
	for k, v := range a.Criteria {
		criteria[k] = v
	}
	criteria[name] = value
	return a.WithCriteriaValues(&criteria)
}

func CriteriaValuesRange(alternatives *[]AlternativeWithCriteria, criterion *Criterion) *utils.ValueRange {
	if criterion.ValuesRange != nil {
		return criterion.ValuesRange
	}
	valRange := utils.NewValueRange()
	for i, a := range *alternatives {
		value := a.CriterionRawValue(criterion)
		if i == 0 {
			valRange.Max = value
			valRange.Min = value
		} else {
			if valRange.Min > value {
				valRange.Min = value
			}
			if valRange.Max < value {
				valRange.Max = value
			}
		}
	}
	return valRange
}

func PreserveCriteriaForAlternatives(alternatives *[]AlternativeWithCriteria, criteria *Criteria) *[]AlternativeWithCriteria {
	result := make([]AlternativeWithCriteria, len(*alternatives))
	for i, a := range *alternatives {
		result[i] = *a.WithCriteriaOnly(criteria)
	}
	return &result
}

type Alternatives []Alternative

type AlternativesRankEntry struct {
	AlternativeResult
	BetterThanOrSameAs Alternatives `json:"betterThanOrSameAs"` // preference >=
}
type AlternativesRanking []AlternativesRankEntry

func (r *AlternativesRanking) ReverseOrder() {
	for i, j := 0, len(*r)-1; i < j; i, j = i+1, j-1 {
		(*r)[i], (*r)[j] = (*r)[j], (*r)[i]
	}
}

func RemoveAlternative(alternatives []AlternativeWithCriteria, alternative AlternativeWithCriteria) []AlternativeWithCriteria {
	for i, v := range alternatives {
		if v.Id == alternative.Id {
			return append(alternatives[:i], alternatives[i+1:]...)
		}
	}
	return alternatives
}

func RemoveAlternativeAt(alternatives []AlternativeWithCriteria, index int) []AlternativeWithCriteria {
	return append(alternatives[:index], alternatives[index+1:]...)
}
