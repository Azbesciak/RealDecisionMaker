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

func TestIncreasingMulCoefficientScalingHigh(t *testing.T) {
	params := IdealIncreasingMulCoefficientSatisfaction.BlankParams().(*IdealCoefficientSatisfactionLevels)
	params.Coefficient = 0.3
	params.MinValue = 0.2
	params.MaxValue = 0.6
	expected := []model.Weights{{
		"a": 1.2, "b": 4.4,
	}, {
		"a": 1.56, "b": 5.12,
	}}
	validateThresholds(t, params, expected)
}

func TestIncreasingMulCoefficientScalingLow(t *testing.T) {
	params := IdealIncreasingMulCoefficientSatisfaction.BlankParams().(*IdealCoefficientSatisfactionLevels)
	params.Coefficient = 0.2
	params.MinValue = 0.1
	params.MaxValue = 0.95
	expected := []model.Weights{{
		"a": 1.1, "b": 4.2,
	}, {
		"a": 1.32, "b": 4.64,
	}, {
		"a": 1.584, "b": 5.168,
	}, {
		"a": 1.9008, "b": 5.8016,
	}}
	validateThresholds(t, params, expected)
}

func TestIncreasingAdditiveScaling(t *testing.T) {
	params := IdealAdditiveCoefficientSatisfaction.BlankParams().(*IdealCoefficientSatisfactionLevels)
	params.Coefficient = 0.2
	params.MinValue = 0.1
	params.MaxValue = 0.95
	expected := []model.Weights{{
		"a": 1.1, "b": 4.2,
	}, {
		"a": 1.3, "b": 4.6,
	}, {
		"a": 1.5, "b": 5,
	}, {
		"a": 1.7, "b": 5.4,
	}, {
		"a": 1.9, "b": 5.8,
	}}
	validateThresholds(t, params, expected)
}

func TestIncreasingAdditiveScaling_highStep(t *testing.T) {
	params := IdealAdditiveCoefficientSatisfaction.BlankParams().(*IdealCoefficientSatisfactionLevels)
	params.Coefficient = 0.4
	params.MinValue = 0.0
	params.MaxValue = 0.99
	expected := []model.Weights{{
		"a": 1, "b": 4,
	}, {
		"a": 1.4, "b": 4.8,
	}, {
		"a": 1.8, "b": 5.6,
	}}
	validateThresholds(t, params, expected)
}

func TestSubstrScaling(t *testing.T) {
	params := IdealSubtrCoefficientSatisfaction.BlankParams().(*IdealCoefficientSatisfactionLevels)
	params.Coefficient = 0.2
	params.MinValue = 0.1
	params.MaxValue = 0.9
	expected := []model.Weights{{
		"a": 1.9, "b": 5.8,
	}, {
		"a": 1.7, "b": 5.4,
	}, {
		"a": 1.5, "b": 5,
	}, {
		"a": 1.3, "b": 4.6,
	}}
	validateThresholds(t, params, expected)
}

func TestDecreasingMulScaling(t *testing.T) {
	params := IdealDecreasingMulCoefficientSatisfaction.BlankParams().(*IdealCoefficientSatisfactionLevels)
	params.Coefficient = 0.8
	params.MinValue = 0.4
	params.MaxValue = 0.9
	expected := []model.Weights{{
		"a": 1.9, "b": 5.8,
	}, {
		"a": 1.72, "b": 5.44,
	}, {
		"a": 1.576, "b": 5.152,
	}, {
		"a": 1.4608, "b": 4.9216,
	}}
	validateThresholds(t, params, expected)
}

func validateThresholds(t *testing.T, params *IdealCoefficientSatisfactionLevels, expected []model.Weights) {
	dmp := model.DecisionMakingParams{
		ConsideredAlternatives: []model.AlternativeWithCriteria{
			{Id: "1", Criteria: model.Weights{"a": 1, "b": 4}},
			{Id: "2", Criteria: model.Weights{"a": 2, "b": 6}},
		},
		Criteria: model.Criteria{{Id: "a"}, {Id: "b"}},
	}
	params.Initialize(&dmp)
	var actual []model.Weights
	for params.HasNext() {
		actual = append(actual, params.Next())
	}
	if utils.Differs(expected, actual) {
		t.Errorf("different thresholds, "+
			"\nexpected %v, "+
			"\n     got %v", expected, actual)
	}
}

func TestDecreasingCoefficientManager_Validate(t *testing.T) {
	manager := DecreasingCoefficientManager{}
	func() {
		defer utils.ExpectError(t, "satisfaction coefficient degradation level must be in range (0, 1), got 0.000000")()
		manager.Validate(&IdealCoefficientSatisfactionLevels{Coefficient: 0})
		t.Errorf("expected error")
	}()
	func() {
		defer utils.ExpectError(t, "satisfaction coefficient degradation level must be in range (0, 1), got 1.000000")()
		manager.Validate(&IdealCoefficientSatisfactionLevels{Coefficient: 1})
		t.Errorf("expected error")
	}()
	func() {
		defer utils.ExpectError(t, "min satisfaction coefficient level must be in range (0, 1], got 0.000000")()
		manager.Validate(&IdealCoefficientSatisfactionLevels{Coefficient: 0.5, MinValue: 0})
		t.Errorf("expected error")
	}()
	func() {
		defer utils.ExpectError(t, "min satisfaction coefficient level must be in range (0, 1], got 1.100000")()
		manager.Validate(&IdealCoefficientSatisfactionLevels{Coefficient: 0.5, MinValue: 1.1})
		t.Errorf("expected error")
	}()
	func() {
		defer utils.ExpectError(t, "max satisfaction coefficient level must be in range (0, 1], got 0.000000")()
		manager.Validate(&IdealCoefficientSatisfactionLevels{Coefficient: 0.5, MinValue: 0.5, MaxValue: 0})
		t.Errorf("expected error")
	}()
	func() {
		defer utils.ExpectError(t, "max satisfaction coefficient level must be in range (0, 1], got 1.100000")()
		manager.Validate(&IdealCoefficientSatisfactionLevels{Coefficient: 0.5, MinValue: 0.5, MaxValue: 1.1})
		t.Errorf("expected error")
	}()
	manager.Validate(&IdealCoefficientSatisfactionLevels{Coefficient: 0.5, MinValue: 0.5, MaxValue: 0.6})
}

func TestIncreasingCoefficientManager_Validate(t *testing.T) {
	manager := IncreasingCoefficientManager{}
	func() {
		defer utils.ExpectError(t, "satisfaction coefficient increasing level must be in range (0, 1), got 0.000000")()
		manager.Validate(&IdealCoefficientSatisfactionLevels{Coefficient: 0})
		t.Errorf("expected error")
	}()
	func() {
		defer utils.ExpectError(t, "satisfaction coefficient increasing level must be in range (0, 1), got 1.000000")()
		manager.Validate(&IdealCoefficientSatisfactionLevels{Coefficient: 1})
		t.Errorf("expected error")
	}()
	manager.Validate(&IdealCoefficientSatisfactionLevels{Coefficient: 0.5, MinValue: 0})
	func() {
		defer utils.ExpectError(t, "min satisfaction coefficient level must be in range [0, 1], got -0.100000")()
		manager.Validate(&IdealCoefficientSatisfactionLevels{Coefficient: 0.5, MinValue: -0.1})
		t.Errorf("expected error")
	}()
	func() {
		defer utils.ExpectError(t, "min satisfaction coefficient level must be in range [0, 1], got 1.100000")()
		manager.Validate(&IdealCoefficientSatisfactionLevels{Coefficient: 0.5, MinValue: 1.1})
		t.Errorf("expected error")
	}()
	func() {
		defer utils.ExpectError(t, "max satisfaction coefficient level must be in range [0, 1], got -0.100000")()
		manager.Validate(&IdealCoefficientSatisfactionLevels{Coefficient: 0.5, MinValue: 0.5, MaxValue: -.1})
		t.Errorf("expected error")
	}()
	func() {
		defer utils.ExpectError(t, "max satisfaction coefficient level must be in range [0, 1], got 1.100000")()
		manager.Validate(&IdealCoefficientSatisfactionLevels{Coefficient: 0.5, MinValue: 0.1, MaxValue: 1.1})
		t.Errorf("expected error")
	}()
}
