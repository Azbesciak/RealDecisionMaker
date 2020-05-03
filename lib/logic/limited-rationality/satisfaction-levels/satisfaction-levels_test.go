package satisfaction_levels

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"testing"
)

func TestThresholdSatisfactionLevels_Initialize(t *testing.T) {
	criteria := testUtils.GenerateCriteria(3)
	valid := ThresholdSatisfactionLevels{
		Thresholds: []model.Weights{{
			"1": 1, "2": 2, "3": 3,
		}, {
			"1": 0, "2": 1, "3": 9,
		}},
		currentIndex: 0,
	}
	dmp := &model.DecisionMakingParams{Criteria: criteria}
	valid.Initialize(dmp)
	invalid := ThresholdSatisfactionLevels{
		Thresholds: []model.Weights{{
			"1": 1, "2": 2, "3": 3,
		}, {
			"1": 0, "2": 1,
		}},
		currentIndex: 0,
	}
	defer utils.ExpectError(t, "value of criterion '3' for threshold 1 not found in [{1 0} {2 1}]")()
	invalid.Initialize(dmp)
}
