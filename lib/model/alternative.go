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
	Value       Weight                  `json:"value"`
}

func (a *AlternativeResult) Identifier() string {
	return a.Alternative.Id
}

const roundPrecision = 1e8

func (a *AlternativeResult) rounded() *AlternativeResult {
	return &AlternativeResult{
		Alternative: a.Alternative,
		Value:       math.Round(a.Value*roundPrecision) / roundPrecision,
	}
}

type AlternativeResults []AlternativeResult

func (a *AlternativeResults) Len() int {
	return len(*a)
}

func (a *AlternativeResults) Less(i, j int) bool {
	if (*a)[i].Value == (*a)[j].Value {
		return (*a)[i].Alternative.Id < (*a)[j].Alternative.Id
	}
	return (*a)[i].Value > (*a)[j].Value
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
	nextLowerThanAltValue := a.Value
	for _, r := range *allAlternatives {
		if r.Value == a.Value && r.Identifier() != a.Identifier() {
			betterThanOrSameAs = append(betterThanOrSameAs, r.Identifier())
		} else if r.Value < a.Value {
			if !wasLowerValueFound {
				wasLowerValueFound = true
				nextLowerThanAltValue = r.Value
			}
			if r.Value < nextLowerThanAltValue {
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
