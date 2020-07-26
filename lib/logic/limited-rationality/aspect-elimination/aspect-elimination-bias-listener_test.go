package aspect_elimination

import (
	satisfaction_levels "github.com/Azbesciak/RealDecisionMaker/lib/logic/limited-rationality/satisfaction-levels"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"reflect"
	"testing"
)

var thresholdsListeners = satisfaction_levels.SatisfactionLevelsUpdateListeners{
	Listeners: satisfaction_levels.ListenersMap{
		satisfaction_levels.Thresholds: &satisfaction_levels.IncreasingThresholds,
	},
}

var increasingAdditiveListeners = satisfaction_levels.SatisfactionLevelsUpdateListeners{
	Listeners: satisfaction_levels.ListenersMap{
		satisfaction_levels.IdealAdditive: &satisfaction_levels.IdealAdditiveCoefficientSatisfaction,
	},
}

var thresholdsWeights = model.Weights{"1": 10, "2": 5, "3": 7}
var parsedThresholds = AspectEliminationHeuristicParams{
	Function: satisfaction_levels.Thresholds,
	Params: &satisfaction_levels.ThresholdSatisfactionLevels{
		Thresholds: []model.Weights{
			{"1": 1, "2": 2, "3": 3},
			{"1": 2, "2": 4, "3": 10},
		},
	},
	Weights: thresholdsWeights,
}
var allCriteria = testUtils.GenerateCriteria(3)
var subtractedCriteria = model.Criteria{{Id: "1", Type: model.Gain}, {Id: "3", Type: model.Cost}}
var afterAddCriteria = testUtils.GenerateCriteria(4)
var addedCriterion = model.Criterion{Id: "4", Type: model.Gain}
var parsedSubtractedThresholds = AspectEliminationHeuristicParams{
	Function: satisfaction_levels.Thresholds,
	Params: &satisfaction_levels.ThresholdSatisfactionLevels{
		Thresholds: []model.Weights{
			{"1": 1, "3": 3},
			{"1": 2, "3": 10},
		},
	},
	Weights: model.Weights{"1": 10, "3": 7},
}
var idealAdditiveOriginalParams = AspectEliminationHeuristicParams{
	Function: satisfaction_levels.IdealAdditive,
	Params:   additiveParams,
	Weights:  model.Weights{"1": 10, "2": 5, "3": 7},
}

var rawThresholds = AspectEliminationHeuristicParams{
	Function: satisfaction_levels.Thresholds,
	Params: utils.Map{
		"thresholds": utils.Array{
			utils.Map{"1": 1, "2": 2, "3": 3},
			utils.Map{"1": 2, "2": 4, "3": 10},
		},
	},
	Weights: thresholdsWeights,
}

var thresholdAddition = aspectEliminationAddedCriterion{
	Weights: model.Weights{"4": 15},
	Params: satisfaction_levels.ThresholdsUpdate{
		Thresholds: []model.Weights{{"4": 5}, {"4": 10}},
	},
}
var additiveParams = &satisfaction_levels.IdealCoefficientSatisfactionLevels{
	MaxValue:    0.9,
	MinValue:    0,
	Coefficient: 0.3,
}

func TestAspectEliminationBiasListener_Merge(t *testing.T) {
	type fields struct {
		satisfactionLevelsUpdateListeners satisfaction_levels.SatisfactionLevelsUpdateListeners
	}
	type args struct {
		params   model.MethodParameters
		addition model.MethodParameters
	}

	mergedThresholdParams := AspectEliminationHeuristicParams{
		Function: satisfaction_levels.Thresholds,
		Params: &satisfaction_levels.ThresholdSatisfactionLevels{
			Thresholds: []model.Weights{
				{"1": 1, "2": 2, "3": 3, "4": 5},
				{"1": 2, "2": 4, "3": 10, "4": 10},
			},
		},
		Weights: model.Weights{"1": 10, "2": 5, "3": 7, "4": 15},
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.MethodParameters
	}{{
		name:   "thresholds",
		fields: fields{satisfactionLevelsUpdateListeners: thresholdsListeners},
		args: args{
			params:   rawThresholds,
			addition: thresholdAddition,
		},
		want: mergedThresholdParams,
	}, {
		name:   "thresholds_parsed",
		fields: fields{satisfactionLevelsUpdateListeners: thresholdsListeners},
		args: args{
			params:   parsedThresholds,
			addition: thresholdAddition,
		},
		want: mergedThresholdParams,
	}, {
		name:   "increasing additive auto thresholds",
		fields: fields{satisfactionLevelsUpdateListeners: increasingAdditiveListeners},
		args: args{
			params: AspectEliminationHeuristicParams{
				Function: satisfaction_levels.IdealAdditive,
				Params:   additiveParams,
				Weights:  model.Weights{"1": 10, "2": 5, "3": 7},
			},
			addition: aspectEliminationAddedCriterion{
				Weights: model.Weights{"4": 15},
				Params:  nil,
			},
		},
		want: AspectEliminationHeuristicParams{
			Function: satisfaction_levels.IdealAdditive,
			Params:   additiveParams,
			Weights:  model.Weights{"1": 10, "2": 5, "3": 7, "4": 15},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AspectEliminationBiasListener{
				satisfactionLevelsUpdateListeners: tt.fields.satisfactionLevelsUpdateListeners,
			}
			if got := a.Merge(tt.args.params, tt.args.addition); utils.Differs(got, tt.want) {
				t.Errorf("Merge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAspectEliminationBiasListener_OnCriteriaRemoved(t *testing.T) {
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
		name: "thresholds raw",
		fields: fields{
			satisfactionLevelsUpdateListeners: thresholdsListeners,
		},
		args: args{
			leftCriteria: &subtractedCriteria,
			params:       rawThresholds,
		},
		want: parsedSubtractedThresholds,
	}, {
		name: "thresholds parsed",
		fields: fields{
			satisfactionLevelsUpdateListeners: thresholdsListeners,
		},
		args: args{
			leftCriteria: &subtractedCriteria,
			params:       parsedThresholds,
		},
		want: parsedSubtractedThresholds,
	}, {
		name: "additive",
		fields: fields{
			satisfactionLevelsUpdateListeners: increasingAdditiveListeners,
		},
		args: args{
			leftCriteria: &subtractedCriteria,
			params:       idealAdditiveOriginalParams,
		},
		want: AspectEliminationHeuristicParams{
			Function: satisfaction_levels.IdealAdditive,
			Params:   additiveParams,
			Weights:  model.Weights{"1": 10, "3": 7},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AspectEliminationBiasListener{
				satisfactionLevelsUpdateListeners: tt.fields.satisfactionLevelsUpdateListeners,
			}
			if got := a.OnCriteriaRemoved(tt.args.leftCriteria, tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OnCriteriaRemoved() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAspectEliminationBiasListener_OnCriterionAdded(t *testing.T) {
	type fields struct {
		satisfactionLevelsUpdateListeners satisfaction_levels.SatisfactionLevelsUpdateListeners
	}
	type args struct {
		criterion          *model.Criterion
		referenceCriterion *model.Criterion
		params             model.MethodParameters
		generator          utils.ValueGenerator
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.AddedCriterionParams
	}{{
		name: "thresholds",
		fields: fields{
			satisfactionLevelsUpdateListeners: thresholdsListeners,
		},
		args: args{
			criterion:          &addedCriterion,
			referenceCriterion: &allCriteria[0],
			params:             rawThresholds,
			generator: func() float64 {
				return 0.5
			},
		},
		want: aspectEliminationAddedCriterion{
			Weights: model.Weights{"4": 5},
			Params: satisfaction_levels.ThresholdsUpdate{
				Thresholds: []model.Weights{{"4": .5}, {"4": 1}},
			},
		},
	}, {
		name: "thresholds processed",
		fields: fields{
			satisfactionLevelsUpdateListeners: thresholdsListeners,
		},
		args: args{
			criterion:          &addedCriterion,
			referenceCriterion: &allCriteria[0],
			params:             parsedThresholds,
			generator: func() float64 {
				return 0.5
			},
		},
		want: aspectEliminationAddedCriterion{
			Weights: model.Weights{"4": 5},
			Params: satisfaction_levels.ThresholdsUpdate{
				Thresholds: []model.Weights{{"4": .5}, {"4": 1}},
			},
		},
	}, {
		name: "additive",
		fields: fields{
			satisfactionLevelsUpdateListeners: increasingAdditiveListeners,
		},
		args: args{
			criterion:          &addedCriterion,
			referenceCriterion: &allCriteria[0],
			params: AspectEliminationHeuristicParams{
				Function: satisfaction_levels.IdealAdditive,
				Params:   additiveParams,
				Weights:  thresholdsWeights,
			},
			generator: func() float64 {
				return 0.5
			},
		},
		want: aspectEliminationAddedCriterion{
			Weights: model.Weights{"4": 5},
			Params:  nil,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AspectEliminationBiasListener{
				satisfactionLevelsUpdateListeners: tt.fields.satisfactionLevelsUpdateListeners,
			}
			if got := a.OnCriterionAdded(tt.args.criterion, tt.args.referenceCriterion, tt.args.params, tt.args.generator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OnCriterionAdded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAspectEliminationBiasListener_RankCriteriaAscending(t *testing.T) {
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
		want   *model.WeightedCriteria
	}{{
		name:   "different weights",
		fields: fields{},
		args: args{
			params: &model.DecisionMakingParams{
				Criteria: allCriteria,
				MethodParameters: AspectEliminationHeuristicParams{
					Weights: model.Weights{"1": 10, "2": 5, "3": 7},
				},
			},
		},
		want: &model.WeightedCriteria{
			{Criterion: allCriteria[1], Weight: 5},
			{Criterion: allCriteria[2], Weight: 7},
			{Criterion: allCriteria[0], Weight: 10},
		},
	}, {
		name:   "same weights",
		fields: fields{},
		args: args{
			params: &model.DecisionMakingParams{
				Criteria: allCriteria,
				MethodParameters: AspectEliminationHeuristicParams{
					Weights: model.Weights{"1": 10, "2": 7, "3": 7},
				},
			},
		},
		want: &model.WeightedCriteria{
			{Criterion: allCriteria[1], Weight: 7},
			{Criterion: allCriteria[2], Weight: 7},
			{Criterion: allCriteria[0], Weight: 10},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AspectEliminationBiasListener{
				satisfactionLevelsUpdateListeners: tt.fields.satisfactionLevelsUpdateListeners,
			}
			if got := a.RankCriteriaAscending(tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RankCriteriaAscending() = %v, want %v", got, tt.want)
			}
		})
	}
}
