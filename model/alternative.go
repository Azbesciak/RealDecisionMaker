package model

import (
	"math"
	"sort"
)

type AlternativeResult struct {
	Alternative AlternativeWithCriteria
	Value       Weight
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
	Id       Alternative
	Criteria Weights
}

type Alternatives []Alternative

type AlternativesRankEntry struct {
	AlternativeResult
	BetterThanOrSameAs Alternatives // preference >=
}
type AlternativesRanking []AlternativesRankEntry
