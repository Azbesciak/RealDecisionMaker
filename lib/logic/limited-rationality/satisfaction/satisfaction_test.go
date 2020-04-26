package satisfaction

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

var _satisfaction = Satisfaction{
	generator: func(seed int64) utils.ValueGenerator {
		return func() float64 {
			return 1
		}
	},
	functions: []SatisfactionLevelsSource{
		&IdealSubtrCoefficientSatisfaction,
		&IdealMulCoefficientSatisfaction,
	},
}

var criteria = testUtils.GenerateCriteria(3)

var a1 = model.AlternativeWithCriteria{
	Id:       "a",
	Criteria: model.Weights{"1": 3, "2": 4, "3": 5},
}

var a2 = model.AlternativeWithCriteria{
	Id:       "b",
	Criteria: model.Weights{"1": 6, "2": 5, "3": 4},
}

var a3 = model.AlternativeWithCriteria{
	Id:       "c",
	Criteria: model.Weights{"1": 3, "2": 6, "3": 3},
}

var a4 = model.AlternativeWithCriteria{
	Id:       "d",
	Criteria: model.Weights{"1": 10, "2": 10, "3": 10}, //not considered but as a reference for ideal
}

var cur = model.AlternativeWithCriteria{
	Id:       "cur",
	Criteria: model.Weights{"1": 0, "2": 0, "3": 0}, //not considered but as a reference for ideal
}

func TestSatisfaction_Evaluate(t *testing.T) {
	dm := model.DecisionMakingParams{
		ConsideredAlternatives:    []model.AlternativeWithCriteria{a1, a2, a3},
		NotConsideredAlternatives: []model.AlternativeWithCriteria{a4, cur},
		Criteria:                  criteria,
		MethodParameters: SatisfactionParameters{
			Function: IdealSubtrCoefficientSatisfaction.name,
			Params: map[string]interface{}{
				"Coefficient": 0.1,
				"MinValue":    0.1,
			},
			Seed:    0,
			Current: "cur",
		},
	}
	actual := _satisfaction.Evaluate(&dm)
	expected := model.AlternativesRanking{
		{
			AlternativeResult: model.AlternativeResult{
				Alternative: a2,
				Value:       7,
			},
			BetterThanOrSameAs: []string{a1.Id, a3.Id, cur.Id},
		},
		{
			AlternativeResult: model.AlternativeResult{
				Alternative: a1,
				Value:       8,
			},
			BetterThanOrSameAs: []string{a3.Id, cur.Id},
		},
		{
			AlternativeResult: model.AlternativeResult{
				Alternative: a3,
				Value:       8,
			},
			BetterThanOrSameAs: []string{cur.Id},
		},
		{
			AlternativeResult: model.AlternativeResult{
				Alternative: cur,
				Value:       10,
			},
			BetterThanOrSameAs: []string{},
		},
	}
	testUtils.CompareRankings(&expected, actual, t)
}
