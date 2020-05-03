package satisfaction

import (
	satisfaction_levels "github.com/Azbesciak/RealDecisionMaker/lib/logic/limited-rationality/satisfaction-levels"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"reflect"
	"testing"
)

var thresholdListener = satisfaction_levels.SatisfactionLevelsUpdateListeners{
	Listeners: satisfaction_levels.ListenersMap{
		satisfaction_levels.Thresholds: &satisfaction_levels.DecreasingThresholds,
	},
}

var subtractiveListener = satisfaction_levels.SatisfactionLevelsUpdateListeners{
	Listeners: satisfaction_levels.ListenersMap{
		satisfaction_levels.IdealSubtractive: &satisfaction_levels.IdealSubtrCoefficientSatisfaction,
	},
}

var idealCoeffParams = satisfaction_levels.IdealCoefficientSatisfactionLevels{
	Coefficient: 0.4,
	MaxValue:    1,
	MinValue:    0.1,
}

var thresholds3Criteria = satisfaction_levels.ThresholdSatisfactionLevels{
	Thresholds: []model.Weights{{
		"1": 2, "2": 3, "3": 0.5,
	}, {
		"1": 1, "2": 2, "3": 0.25,
	}},
}
var thresholds2Criteria = satisfaction_levels.ThresholdSatisfactionLevels{
	Thresholds: []model.Weights{{
		"1": 2, "2": 3,
	}, {
		"1": 1, "2": 2,
	}},
}
var criteria2 = testUtils.GenerateCriteria(2)
var addedCriterion = model.Criterion{Id: "3", Type: model.Gain}
var thresholdsDif = satisfaction_levels.ThresholdsUpdate{
	Thresholds: []model.Weights{{"3": 0.5}, {"3": 0.25}},
}

var threshold3ParamsWithPointer = SatisfactionParameters{
	Function: satisfaction_levels.Thresholds,
	Params:   &thresholds3Criteria,
	Seed:     123,
	Current:  "xyz",
}

var thresholds2ParamsWithPointer = SatisfactionParameters{
	Function: satisfaction_levels.Thresholds,
	Params:   &thresholds2Criteria,
	Seed:     123,
	Current:  "xyz",
}

var idealSubtractiveParamsWithoutPointer = SatisfactionParameters{
	Function: satisfaction_levels.IdealSubtractive,
	Params:   idealCoeffParams,
	Seed:     432,
	Current:  "",
}

var idealSubtractiveParamsWithPointer = SatisfactionParameters{
	Function: satisfaction_levels.IdealSubtractive,
	Params:   &idealCoeffParams,
	Seed:     432,
	Current:  "",
}

func TestSatisfactionBiasListener_Merge(t *testing.T) {
	type fields struct {
		satisfactionLevelsUpdateListeners satisfaction_levels.SatisfactionLevelsUpdateListeners
	}
	type args struct {
		params   model.MethodParameters
		addition model.MethodParameters
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.MethodParameters
	}{{
		name: "thresholds",
		fields: fields{
			satisfactionLevelsUpdateListeners: thresholdListener,
		},
		args: args{
			params: thresholds2ParamsWithPointer,
			addition: satisfactionAddedCriterion{
				Params: thresholdsDif,
			},
		},
		want: threshold3ParamsWithPointer,
	}, {
		name: "subtractive",
		fields: fields{
			satisfactionLevelsUpdateListeners: subtractiveListener,
		},
		args: args{
			params: idealSubtractiveParamsWithoutPointer,
			addition: satisfactionAddedCriterion{
				Params: nil,
			},
		},
		want: idealSubtractiveParamsWithPointer,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &SatisfactionBiasListener{
				satisfactionLevelsUpdateListeners: tt.fields.satisfactionLevelsUpdateListeners,
			}
			if got := a.Merge(tt.args.params, tt.args.addition); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSatisfactionBiasListener_OnCriteriaRemoved(t *testing.T) {
	type fields struct {
		satisfactionLevelsUpdateListeners satisfaction_levels.SatisfactionLevelsUpdateListeners
	}
	type args struct {
		leftCriteria *model.Criteria
		params       model.MethodParameters
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.MethodParameters
	}{{
		name: "thresholds",
		fields: fields{
			satisfactionLevelsUpdateListeners: thresholdListener,
		},
		args: args{
			leftCriteria: &criteria2,
			params:       threshold3ParamsWithPointer,
		},
		want: thresholds2ParamsWithPointer,
	}, {
		name: "subtractive",
		fields: fields{
			satisfactionLevelsUpdateListeners: subtractiveListener,
		},
		args: args{
			leftCriteria: &criteria2,
			params:       idealSubtractiveParamsWithPointer,
		},
		want: idealSubtractiveParamsWithPointer,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &SatisfactionBiasListener{
				satisfactionLevelsUpdateListeners: tt.fields.satisfactionLevelsUpdateListeners,
			}
			if got := a.OnCriteriaRemoved(tt.args.leftCriteria, tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OnCriteriaRemoved() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSatisfactionBiasListener_OnCriterionAdded(t *testing.T) {
	type fields struct {
		satisfactionLevelsUpdateListeners satisfaction_levels.SatisfactionLevelsUpdateListeners
	}
	type args struct {
		criterion              *model.Criterion
		previousRankedCriteria *model.Criteria
		params                 model.MethodParameters
		generator              utils.ValueGenerator
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.AddedCriterionParams
	}{{
		name: "threshold",
		fields: fields{
			satisfactionLevelsUpdateListeners: thresholdListener,
		},
		args: args{
			criterion:              &addedCriterion,
			previousRankedCriteria: &criteria2,
			params:                 thresholds2ParamsWithPointer,
			generator: func() float64 {
				return 0.25
			},
		},
		want: satisfactionAddedCriterion{
			Params: thresholdsDif,
		},
	}, {
		name: "subtractive",
		fields: fields{
			satisfactionLevelsUpdateListeners: subtractiveListener,
		},
		args: args{
			criterion:              &addedCriterion,
			previousRankedCriteria: &criteria2,
			params:                 idealSubtractiveParamsWithPointer,
			generator: func() float64 {
				return 0.75
			},
		},
		want: satisfactionAddedCriterion{
			Params: nil,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &SatisfactionBiasListener{
				satisfactionLevelsUpdateListeners: tt.fields.satisfactionLevelsUpdateListeners,
			}
			if got := a.OnCriterionAdded(tt.args.criterion, tt.args.previousRankedCriteria, tt.args.params, tt.args.generator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OnCriterionAdded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSatisfactionBiasListener_RankCriteriaAscending(t *testing.T) {
	type fields struct {
		satisfactionLevelsUpdateListeners satisfaction_levels.SatisfactionLevelsUpdateListeners
	}
	type args struct {
		params *model.DecisionMakingParams
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *model.Criteria
	}{{
		name:   "ranking",
		fields: fields{},
		args: args{params: &model.DecisionMakingParams{
			ConsideredAlternatives: []model.AlternativeWithCriteria{{
				Id:       "a",
				Criteria: model.Weights{"1": 1, "2": 0.5},
			}, {
				Id:       "b",
				Criteria: model.Weights{"1": 4, "2": 1.5},
			}},
			Criteria: criteria2,
		}},
		want: &model.Criteria{criteria2[1], criteria2[0]},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &SatisfactionBiasListener{
				satisfactionLevelsUpdateListeners: tt.fields.satisfactionLevelsUpdateListeners,
			}
			if got := a.RankCriteriaAscending(tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RankCriteriaAscending() = %v, want %v", got, tt.want)
			}
		})
	}
}
