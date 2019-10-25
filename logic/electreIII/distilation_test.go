package electreIII

import (
	"reflect"
	"testing"
)
import "../../model"

var ele1DistillationExample = AlternativesMatrix{
	Alternatives: &model.Alternatives{"1", "2", "3", "4", "5"},
	Values: &Matrix{
		Size: 5,
		Data: []float64{
			1, 0, 1, 0.8, 1,
			0, 1, 0, 0.9, 0.67,
			0.6, 0, 1, 0.6, 0.8,
			0.25, 0.8, 0.67, 1, 0.85,
			0.67, 0, 0.8, 0.8, 1,
		},
	},
}

var powerStationExample = AlternativesMatrix{
	Alternatives: &model.Alternatives{"ITA", "BEL", "GER", "AUT", "FRA"},
	Values: &Matrix{
		Size: 5,
		Data: []float64{
			1, 0.6, 0.3, 0.7, 0.55,
			0, 1, 0.2, 0, 0,
			0.58, 0.6, 1, 0.85, 0,
			0.3, 0, 0.6, 1, 0.58,
			0.6, 0, 0.6, 0.7, 1,
		},
	},
}

func TestRankAscending_DistillationExample(t *testing.T) {
	check(t, &ele1DistillationExample, RankAscending, []int{1, 2, 3, 3, 3}, "ascending")
}

func TestRankDescending_DistillationExample(t *testing.T) {
	check(t, &ele1DistillationExample, RankDescending, []int{1, 1, 3, 2, 3}, "descending")
}

func TestRankAscending_PowerStationExample(t *testing.T) {
	check(t, &powerStationExample, RankAscending, []int{2, 3, 1, 3, 3}, "ascending")
}

func TestRankDescending_PowerStationExample(t *testing.T) {
	check(t, &powerStationExample, RankDescending, []int{3, 4, 2, 5, 1}, "descending")
}

func check(
	t *testing.T,
	data *AlternativesMatrix,
	rankingProvider func(matrix *AlternativesMatrix,
		distillationFun *LinearFunctionParameters) *[]int,
	expected []int, rankingType string,
) {
	actual := rankingProvider(data, &DefaultDistillationFunc)
	checkRanking(t, &expected, actual, rankingType)
}

func checkRanking(t *testing.T, expected, actual *[]int, rankingType string) {
	if !reflect.DeepEqual(*expected, *actual) {
		t.Errorf("expected %s ranking to be %v, got %v", rankingType, *expected, *actual)
	}
}

func Test_quality(t *testing.T) {
	mat := Matrix{
		Size: 5,
		Data: []float64{
			0, 0, 1, 0, 1,
			0, 0, 0, 0, 0,
			0, 0, 0, 0, 0,
			0, 0, 0, 0, 0,
			0, 0, 0, 0, 0,
		},
	}
	expected := []int{2, 0, -1, 0, -1}
	quality := computeQuality(&mat)
	if !reflect.DeepEqual(expected, *quality) {
		t.Errorf("expected quality %v, got %v", expected, *quality)
	}
}
