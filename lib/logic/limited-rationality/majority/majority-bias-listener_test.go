package majority

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

func TestMajorityBiasListener_RankCriteriaAscending(t *testing.T) {
	listener := MajorityBiasListener{}
	actual := *listener.RankCriteriaAscending(&model.DecisionMakingParams{
		ConsideredAlternatives: nil,
		Criteria:               testUtils.GenerateCriteria(4),
		MethodParameters: MajorityHeuristicParams{
			Weights: model.Weights{"1": 10, "2": 20, "3": 5, "4": 100},
		},
	})
	expected := model.WeightedCriteria{
		{Criterion: model.Criterion{Id: "3", Type: model.Gain}, Weight: 5},
		{Criterion: model.Criterion{Id: "1", Type: model.Gain}, Weight: 10},
		{Criterion: model.Criterion{Id: "2", Type: model.Gain}, Weight: 20},
		{Criterion: model.Criterion{Id: "4", Type: model.Gain}, Weight: 100},
	}
	if utils.Differs(expected, actual) {
		t.Errorf("Invalid order returned by majority criteria ranking, expected %v, got %v", expected, actual)
	}
}
