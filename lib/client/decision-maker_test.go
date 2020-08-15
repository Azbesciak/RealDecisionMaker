package client

import (
	. "github.com/Azbesciak/RealDecisionMaker/lib/logic/preference-func/electreIII"
	. "github.com/Azbesciak/RealDecisionMaker/lib/logic/preference-func/owa"
	. "github.com/Azbesciak/RealDecisionMaker/lib/logic/preference-func/weighted-sum"
	. "github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"reflect"
	"testing"
)

func TestDecisionMaker_MakeDecision_WeightedSum(t *testing.T) {
	dm := DecisionMaker{
		PreferenceFunction: "weightedSum",
		KnownAlternatives: []AlternativeWithCriteria{
			{"1", Weights{"1": 10, "2": 20}},
			{"2", Weights{"1": 20, "2": 10}},
		},
		ChoseToMake:      []Alternative{"1", "2"},
		Criteria:         Criteria{{Id: "1", Type: "gain"}, {Id: "2", Type: "cost"}},
		MethodParameters: utils.Map{"weights": Weights{"1": 1, "2": 2}},
	}
	weightedSum := &WeightedSumPreferenceFunc{}
	owa := &OWAPreferenceFunc{}
	funcs := PreferenceFunctions{Functions: []PreferenceFunction{weightedSum, owa}}
	ranking := dm.MakeDecision(funcs, BiasListeners{}, &BiasMap{})
	expectedRanking := weightedSum.Evaluate(&DecisionMakingParams{
		NotConsideredAlternatives: nil,
		ConsideredAlternatives:    dm.KnownAlternatives,
		Criteria:                  dm.Criteria,
		MethodParameters:          weightedSum.ParseParams(&dm),
	})

	if !reflect.DeepEqual(ranking.Result, *expectedRanking) {
		t.Errorf("Invalid result, expected '%v', got '%v'", ranking.Result, *expectedRanking)
	}
}

func TestDecisionMaker_MakeDecision_ElectreIII(t *testing.T) {
	dm := DecisionMaker{
		PreferenceFunction: "electreIII",
		KnownAlternatives: []AlternativeWithCriteria{
			{"1", Weights{"1": 10, "2": 20}},
			{"2", Weights{"1": 20, "2": 10}},
		},
		ChoseToMake: []Alternative{"1", "2"},
		Criteria:    Criteria{{Id: "1", Type: "gain"}, {Id: "2", Type: "cost"}},
		MethodParameters: utils.Map{
			"electreCriteria": ElectreCriteria{
				"1": {
					K: 1,
					Q: utils.LinearFunctionParameters{A: 0, B: 2},
					P: utils.LinearFunctionParameters{A: 1, B: 4},
					V: utils.LinearFunctionParameters{A: 0, B: 10},
				},
				"2": {
					K: 2,
					Q: utils.LinearFunctionParameters{A: 0, B: 10},
					P: utils.LinearFunctionParameters{A: 0, B: 12},
					V: utils.LinearFunctionParameters{A: 0, B: 14},
				},
			},
			"electreDistillation": utils.LinearFunctionParameters{
				A: -0.1,
				B: 0.5,
			}},
	}
	weightedSum := &WeightedSumPreferenceFunc{}
	ele := &ElectreIIIPreferenceFunc{}
	funcs := PreferenceFunctions{Functions: []PreferenceFunction{weightedSum, ele}}
	ranking := dm.MakeDecision(funcs, BiasListeners{}, &BiasMap{})
	expectedRanking := ele.Evaluate(&DecisionMakingParams{
		NotConsideredAlternatives: nil,
		ConsideredAlternatives:    dm.KnownAlternatives,
		Criteria:                  dm.Criteria,
		MethodParameters:          ele.ParseParams(&dm),
	})

	if !reflect.DeepEqual(ranking.Result, *expectedRanking) {
		t.Errorf("Invalid result, expected '%v', got '%v'", ranking.Result, *expectedRanking)
	}
}

func TestFetchPreferenceFunc(t *testing.T) {
	functions := PreferenceFunctions{Functions: []PreferenceFunction{
		&ElectreIIIPreferenceFunc{},
		&OWAPreferenceFunc{},
	}}
	defer utils.ExpectError(t, "preference function 'xyz' not found, available are '[electreIII owa]'")()
	functions.Fetch("xyz")
}

func TestParsingToModel(t *testing.T) {
	var dm DecisionMaker
	preferenceFun := "abcd"
	alternatives := []utils.Map{{
		"id":       "a",
		"criteria": utils.Map{"1": 1, "2": 2},
	}, {
		"id":       "b",
		"criteria": utils.Map{"1": 3, "2": 4},
	}}
	criteria := []utils.Map{{
		"id":          "1",
		"type":        "gain",
		"valuesRange": utils.Map{"min": 0, "max": 6},
	}, {
		"id":          "2",
		"type":        "cost",
		"valuesRange": utils.Map{"min": -10, "max": 10},
	}}
	chooseToMake := utils.Array{"a", "b"}
	methodParameters := utils.Map{
		"some":   "unique",
		"params": 1,
	}
	model := utils.Map{
		"preferenceFunction": preferenceFun,
		"knownAlternatives":  alternatives,
		"criteria":           criteria,
		"choseToMake":        chooseToMake,
		"methodParameters":   methodParameters,
	}
	utils.DecodeToStruct(model, &dm)
	if preferenceFun != dm.PreferenceFunction {
		t.Errorf("preference function differs, expected %s, got %s", preferenceFun, dm.PreferenceFunction)
	}
	if len(dm.Criteria) != len(criteria) {
		t.Errorf("different criteria len, expected %d, got %d", len(criteria), len(dm.Criteria))
	} else {
		for i, c := range criteria {
			actual := dm.Criteria[i]
			cId := c["id"].(string)
			if actual.Id != cId {
				t.Errorf("different criterion at %d, expected %s, got %s", i, cId, actual.Id)
			}
			crtType := c["type"]
			if string(actual.Type) != crtType {
				t.Errorf("different criterion type at %d, expected %s, got %s", i, crtType, actual.Type)
			}
			if actual.ValuesRange == nil {
				t.Errorf("expected value range for criterion %s (%d)", actual.Id, i)
			} else {
				valRange := c["valuesRange"].(utils.Map)
				minVal := valRange["min"].(int)
				if actual.ValuesRange.Min != float64(minVal) {
					t.Errorf("different min for criterion %s (%d), expected %d, got %f", actual.Id, i, minVal, actual.ValuesRange.Min)
				}
				maxVal := valRange["max"].(int)
				if actual.ValuesRange.Max != float64(maxVal) {
					t.Errorf("different max for criterion %s (%d), expected %d, got %f", actual.Id, i, maxVal, actual.ValuesRange.Max)
				}
			}
		}
	}
	if len(dm.ChoseToMake) != len(chooseToMake) {
		t.Errorf("different len of chooseToMake, expected %d, got %d", len(dm.ChoseToMake), len(chooseToMake))
	} else {
		for i, a := range dm.ChoseToMake {
			if chooseToMake[i] != a {
				t.Errorf("different choose to make at %d, expected %s, got %s", i, chooseToMake[i], a)
			}
		}
	}
	if utils.Differs(methodParameters, dm.MethodParameters) {
		t.Errorf("different method parameters after parse, expected %v, got %v", methodParameters, dm.MethodParameters)
	}
	if len(dm.KnownAlternatives) != len(alternatives) {
		t.Errorf("different len of known alternatives, expected %d, got %d", len(dm.ChoseToMake), len(alternatives))
	} else {
		for i, a := range dm.KnownAlternatives {
			expA := alternatives[i]
			if expA["id"] != a.Id {
				t.Errorf("different alternative at %d, expected %s, got %s", i, expA["id"], a.Id)
			}

			expCriteria := expA["criteria"].(utils.Map)
			if len(expCriteria) != len(a.Criteria) {
				t.Errorf("different criteria len for alternative %s (%d), expected %v, got %v", expA["id"], i, len(expCriteria), len(a.Criteria))
			} else {
				for c, v := range a.Criteria {
					expVal := expCriteria[c].(int)
					if v != float64(expVal) {
						t.Errorf("different value for criterion %s of alternative %s, exp %v, got %v", c, a.Id, expVal, v)
					}
				}
			}
		}
	}
}
