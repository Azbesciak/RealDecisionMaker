package testUtils

import (
	. "github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"strconv"
	"testing"
)

func compareSize(expectedRanking, ranking *AlternativesRanking, t *testing.T) bool {
	expectedRankingSize := len(*expectedRanking)
	receivedRankingSize := len(*ranking)
	if receivedRankingSize != expectedRankingSize {
		t.Errorf("Expected ranking of length %d , got %d", expectedRankingSize, receivedRankingSize)
		return false
	}
	return true
}

func CompareRankings(expected, received *AlternativesRanking, t *testing.T) {
	if !compareSize(expected, received, t) {
		return
	}
	for i, e := range *expected {
		rec := (*received)[i]
		if e.Alternative.Id != rec.Alternative.Id {
			t.Errorf("Expected id of '%s' at position %d, got '%s'", e.Alternative.Id, i, rec.Alternative.Id)
		}

		if !cmp.Equal(e.BetterThanOrSameAs, rec.BetterThanOrSameAs) {
			t.Errorf(
				"Invalid Preferrence of id '%s' at position %d, expected '%s', got '%s'",
				e.Alternative.Id, i, e.BetterThanOrSameAs, rec.BetterThanOrSameAs,
			)
		}

		if !cmp.Equal(e.Evaluation, rec.Evaluation, cmpopts.EquateApprox(0, 1e-8)) {
			t.Errorf(
				"Invalid evaluation for id '%s' at position %d, expected '%v', got '%v'",
				e.Alternative.Id, i, e.Evaluation, rec.Evaluation,
			)
		}
	}
}

type AltsMap = *map[string]AlternativeResult

func AlternativesResultToMap(a *AlternativeResults) AltsMap {
	altsById := map[string]AlternativeResult{}
	for _, alt := range *a {
		altsById[alt.Alternative.Id] = alt
	}
	return &altsById
}

func ExtractAlternativesFromResults(a *AlternativeResults) *[]AlternativeWithCriteria {
	alts := make([]AlternativeWithCriteria, len(*a))
	for i, alt := range *a {
		alts[i] = alt.Alternative
	}
	return &alts
}

func DummyRankingEntry(alts AltsMap, thisAlt string, betterThanOrSameAs ...Alternative) AlternativesRankEntry {
	var betThOrSamAs = make(Alternatives, len(betterThanOrSameAs))
	for i, a := range betterThanOrSameAs {
		betThOrSamAs[i] = a
	}
	alt := (*alts)[thisAlt]
	return AlternativesRankEntry{
		AlternativeResult:  alt,
		BetterThanOrSameAs: betThOrSamAs,
	}
}

func DummyAlternative(id string, value Weight) AlternativeResult {
	return *ValueAlternativeResult(&AlternativeWithCriteria{Id: id}, value)
}

func ValidateWeights(t *testing.T, name string, expected, actual Weights) {
	if len(expected) != len(actual) {
		t.Errorf("%s: expected %d elements, got %d", name, len(expected), len(actual))
		return
	}
	for k, expValue := range expected {
		actValue, ok := actual[k]
		if !ok {
			t.Errorf("%s: no value for key '%s'", name, k)
		} else if !utils.FloatsAreEqual(expValue, actValue, 1e-6) {
			t.Errorf("%s: weights differ for %s: exp %f vs act %f", name, k, expValue, actValue)
		}
	}
}

func CyclicRandomGenerator(start, cycle int64) func(seed int64) utils.ValueGenerator {
	return func(seed int64) utils.ValueGenerator {
		counter := start
		return func() float64 {
			counter = (counter + 1) % cycle
			return float64(counter) / float64(cycle)
		}
	}
}

func GenerateCriteria(number int) Criteria {
	criteria := make(Criteria, number)
	for i := 0; i < number; i++ {
		criteria[i] = Criterion{Id: strconv.Itoa(i + 1), Type: Gain}
	}
	return criteria
}
