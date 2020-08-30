package anchoring

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/model/reference-criterion"
	"github.com/Azbesciak/RealDecisionMaker/lib/testUtils"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"reflect"
	"testing"
)

func generateProps(applierName string, applierParams utils.Map) *model.BiasProps {
	biasProps := model.BiasProps(utils.Map{
		"anchoringAlternatives": utils.Array{
			utils.Map{"alternative": "1", "coefficient": 3},
			utils.Map{"alternative": "4", "coefficient": 1},
		},
		"loss": utils.Map{
			"function": "linear",
			"params":   utils.Map{"a": 2},
		},
		"gain": utils.Map{
			"function": "linear",
			"params":   utils.Map{"a": 1},
		},
		"referencePoints": utils.Map{
			"function": "ideal",
		},
		"applier": utils.Map{
			"function": applierName,
			"params":   applierParams,
		},
	})
	return &biasProps
}

func TestAnchoring_Apply(t *testing.T) {
	type fields struct {
		anchoringEvaluators       []AnchoringEvaluator
		referencePointsEvaluators []ReferencePointsEvaluator
		anchoringAppliers         []AnchoringApplier
	}
	type args struct {
		current  *model.DecisionMakingParams
		props    *model.BiasProps
		listener *model.BiasListener
	}
	listener := model.BiasListener(&testUtils.DummyBiasListener{})
	consideredAlternatives := []model.AlternativeWithCriteria{{
		Id: "1", Criteria: model.Weights{"a": 1, "b": 2, "c": 3},
	}, {
		Id: "2", Criteria: model.Weights{"a": 2, "b": 0, "c": 2},
	}, {
		Id: "3", Criteria: model.Weights{"a": 4, "b": 3, "c": 4},
	}}
	criteria := model.Criteria{
		{Id: "a", Type: model.Gain},
		{Id: "b", Type: model.Gain},
		{Id: "c", Type: model.Cost},
	}
	notConsideredCriteria := []model.AlternativeWithCriteria{{
		Id: "4", Criteria: model.Weights{"a": 5, "b": 4, "c": 9},
	}, {
		Id: "5", Criteria: model.Weights{"a": 3, "b": 1, "c": 1},
	}}
	methodParams := testUtils.DummyMethodParameters{
		Criteria: []string{"a", "b", "c"},
	}
	anchoringFields := fields{
		anchoringEvaluators:       []AnchoringEvaluator{&LinearAnchoringEvaluator{}, &ExpFromZeroAnchoringEvaluator{}},
		referencePointsEvaluators: []ReferencePointsEvaluator{&IdealReferenceAlternativeEvaluator{}, &NadirReferenceAlternativeEvaluator{}},
		anchoringAppliers: []AnchoringApplier{&InlineAnchoringApplier{}, NewNewCriterionAnchoringApplier(
			func(seed int64) utils.ValueGenerator {
				return func() float64 {
					return 1
				}
			},
			*reference_criterion.NewReferenceCriteriaManager([]reference_criterion.ReferenceCriterionFactory{
				&reference_criterion.ImportanceRatioReferenceCriterionManager{},
			}),
		)},
	}
	dmpParams := model.DecisionMakingParams{
		ConsideredAlternatives:    consideredAlternatives,
		NotConsideredAlternatives: notConsideredCriteria,
		Criteria:                  criteria,
		MethodParameters:          methodParams,
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *model.BiasedResult
	}{{
		name:   "inline anchoring - not applied for not considered",
		fields: anchoringFields,
		args: args{
			current:  &dmpParams,
			props:    generateProps("inline", utils.Map{"allowedValuesRangeScaling": 1}),
			listener: &listener,
		},
		want: &model.BiasedResult{
			DMP: &model.DecisionMakingParams{
				NotConsideredAlternatives: notConsideredCriteria,
				ConsideredAlternatives: []model.AlternativeWithCriteria{{
					Id: "1", Criteria: model.Weights{"a": 1, "b": 2, "c": 3},
				}, {
					Id: "2", Criteria: model.Weights{"a": 1, "b": 0, "c": 3},
				}, {
					Id: "3", Criteria: model.Weights{"a": 2, "b": 4, "c": 2},
				}},
				Criteria:         criteria,
				MethodParameters: methodParams,
			},
			Props: AnchoringResult{
				ReferencePoints: []model.AlternativeWithCriteria{{
					Id:       "ideal",
					Criteria: model.Weights{"a": 5, "b": 2, "c": 3},
				}},
				CriteriaScaling: CriteriaScaling{
					"a": {
						Scale:       1 / 4.0,
						ValuesRange: utils.ValueRange{Min: 1, Max: 5},
					},
					"b": {
						Scale:       1 / 4.0,
						ValuesRange: utils.ValueRange{Min: 0, Max: 4},
					},
					"c": {
						Scale:       1 / 8.0,
						ValuesRange: utils.ValueRange{Min: 1, Max: 9},
					},
				},
				//ideal a: 5, b: 2, c: 3
				PerReferencePointsDifferences: []ReferencePointsDifference{{
					Alternative: consideredAlternatives[0],
					ReferencePointsDifference: []ReferencePointDifference{{
						ReferencePoint: "ideal",
						Coefficients:   model.Weights{"a": -2, "b": 0, "c": 0},
					}},
				}, {
					Alternative: consideredAlternatives[1],
					ReferencePointsDifference: []ReferencePointDifference{{
						ReferencePoint: "ideal",
						Coefficients:   model.Weights{"a": -1.5, "b": -1, "c": 0.125},
					}},
				}, {
					Alternative: consideredAlternatives[2],
					ReferencePointsDifference: []ReferencePointDifference{{
						ReferencePoint: "ideal",
						Coefficients:   model.Weights{"a": -0.5, "b": 0.25, "c": -0.25},
					}},
				}, {
					Alternative: notConsideredCriteria[0],
					ReferencePointsDifference: []ReferencePointDifference{{
						ReferencePoint: "ideal",
						Coefficients:   model.Weights{"a": 0, "b": 0.5, "c": -1.5},
					}},
				}, {
					Alternative: notConsideredCriteria[1],
					ReferencePointsDifference: []ReferencePointDifference{{
						ReferencePoint: "ideal",
						Coefficients:   model.Weights{"a": -1, "b": -0.5, "c": 0.25},
					}},
				}},
				ApplierResult: InlineAnchoringApplierResult{
					AppliedDifferences: []model.AlternativeWithCriteria{{
						Id:       consideredAlternatives[0].Id,
						Criteria: model.Weights{"a": -0, "b": 0, "c": 0},
					}, {
						Id:       consideredAlternatives[1].Id,
						Criteria: model.Weights{"a": -1, "b": 0, "c": 1},
					}, {
						Id:       consideredAlternatives[2].Id,
						Criteria: model.Weights{"a": -2, "b": 1, "c": -2},
					}},
				},
			},
		},
	}, {
		name:   "inline anchoring - applied for not considered",
		fields: anchoringFields,
		args: args{
			current: &dmpParams,
			props: generateProps("inline", utils.Map{
				"allowedValuesRangeScaling": 1,
				"applyOnNotConsidered":      true,
			}),
			listener: &listener,
		},
		want: &model.BiasedResult{
			DMP: &model.DecisionMakingParams{
				NotConsideredAlternatives: []model.AlternativeWithCriteria{{
					Id: notConsideredCriteria[0].Id, Criteria: model.Weights{"a": 5, "b": 4, "c": 1},
				}, {
					Id: notConsideredCriteria[1].Id, Criteria: model.Weights{"a": 1, "b": 0, "c": 3},
				}},
				ConsideredAlternatives: []model.AlternativeWithCriteria{{
					Id: consideredAlternatives[0].Id, Criteria: model.Weights{"a": 1, "b": 2, "c": 3},
				}, {
					Id: consideredAlternatives[1].Id, Criteria: model.Weights{"a": 1, "b": 0, "c": 3},
				}, {
					Id: consideredAlternatives[2].Id, Criteria: model.Weights{"a": 2, "b": 4, "c": 2},
				}},
				Criteria:         criteria,
				MethodParameters: methodParams,
			},
			Props: AnchoringResult{
				ReferencePoints: []model.AlternativeWithCriteria{{
					Id:       "ideal",
					Criteria: model.Weights{"a": 5, "b": 2, "c": 3},
				}},
				CriteriaScaling: CriteriaScaling{
					"a": {
						Scale:       1 / 4.0,
						ValuesRange: utils.ValueRange{Min: 1, Max: 5},
					},
					"b": {
						Scale:       1 / 4.0,
						ValuesRange: utils.ValueRange{Min: 0, Max: 4},
					},
					"c": {
						Scale:       1 / 8.0,
						ValuesRange: utils.ValueRange{Min: 1, Max: 9},
					},
				},
				//ideal a: 5, b: 2, c: 3
				PerReferencePointsDifferences: []ReferencePointsDifference{{
					Alternative: consideredAlternatives[0],
					ReferencePointsDifference: []ReferencePointDifference{{
						ReferencePoint: "ideal",
						Coefficients:   model.Weights{"a": -2, "b": 0, "c": 0},
					}},
				}, {
					Alternative: consideredAlternatives[1],
					ReferencePointsDifference: []ReferencePointDifference{{
						ReferencePoint: "ideal",
						Coefficients:   model.Weights{"a": -1.5, "b": -1, "c": 0.125},
					}},
				}, {
					Alternative: consideredAlternatives[2],
					ReferencePointsDifference: []ReferencePointDifference{{
						ReferencePoint: "ideal",
						Coefficients:   model.Weights{"a": -0.5, "b": 0.25, "c": -0.25},
					}},
				}, {
					Alternative: notConsideredCriteria[0],
					ReferencePointsDifference: []ReferencePointDifference{{
						ReferencePoint: "ideal",
						Coefficients:   model.Weights{"a": 0, "b": 0.5, "c": -1.5},
					}},
				}, {
					Alternative: notConsideredCriteria[1],
					ReferencePointsDifference: []ReferencePointDifference{{
						ReferencePoint: "ideal",
						Coefficients:   model.Weights{"a": -1, "b": -0.5, "c": 0.25},
					}},
				}},
				ApplierResult: InlineAnchoringApplierResult{
					AppliedDifferences: []model.AlternativeWithCriteria{{
						Id:       consideredAlternatives[0].Id,
						Criteria: model.Weights{"a": -0, "b": 0, "c": 0},
					}, {
						Id:       consideredAlternatives[1].Id,
						Criteria: model.Weights{"a": -1, "b": 0, "c": 1},
					}, {
						Id:       consideredAlternatives[2].Id,
						Criteria: model.Weights{"a": -2, "b": 1, "c": -2},
					}, {
						Id:       notConsideredCriteria[0].Id,
						Criteria: model.Weights{"a": 0, "b": 0, "c": -8},
					}, {
						Id:       notConsideredCriteria[1].Id,
						Criteria: model.Weights{"a": -2, "b": -1, "c": 2},
					}},
				},
			},
		},
	}, {
		name:   "new criterion anchoring",
		fields: anchoringFields,
		args: args{
			current: &dmpParams,
			props: generateProps("newCriterion", utils.Map{
				"referenceCriterionType":    "importanceRatio",
				"allowedValuesRangeScaling": 1,
				"newCriterionImportance":    0,
			}),
			listener: &listener,
		},
		want: &model.BiasedResult{
			DMP: &model.DecisionMakingParams{
				NotConsideredAlternatives: []model.AlternativeWithCriteria{{
					Id: notConsideredCriteria[0].Id, Criteria: model.Weights{"a": 5, "b": 4, "c": 9, "__anchoring_criterion_ideal": 1},
				}, {
					Id: notConsideredCriteria[1].Id, Criteria: model.Weights{"a": 3, "b": 1, "c": 1, "__anchoring_criterion_ideal": 1.5833333333333335},
				}},
				ConsideredAlternatives: []model.AlternativeWithCriteria{{
					Id: consideredAlternatives[0].Id, Criteria: model.Weights{"a": 1, "b": 2, "c": 3, "__anchoring_criterion_ideal": 1.3333333333333335},
				}, {
					Id: consideredAlternatives[1].Id, Criteria: model.Weights{"a": 2, "b": 0, "c": 2, "__anchoring_criterion_ideal": 1},
				}, {
					Id: consideredAlternatives[2].Id, Criteria: model.Weights{"a": 4, "b": 3, "c": 4, "__anchoring_criterion_ideal": 1.75},
				}},
				Criteria: append(criteria, model.Criterion{Id: "__anchoring_criterion_ideal", Type: model.Gain}),
				MethodParameters: testUtils.DummyMethodParameters{
					Criteria: append(methodParams.Criteria, "__anchoring_criterion_ideal"),
				},
			},
			Props: AnchoringResult{
				ReferencePoints: []model.AlternativeWithCriteria{{
					Id:       "ideal",
					Criteria: model.Weights{"a": 5, "b": 2, "c": 3},
				}},
				CriteriaScaling: CriteriaScaling{
					"a": {
						Scale:       1 / 4.0,
						ValuesRange: utils.ValueRange{Min: 1, Max: 5},
					},
					"b": {
						Scale:       1 / 4.0,
						ValuesRange: utils.ValueRange{Min: 0, Max: 4},
					},
					"c": {
						Scale:       1 / 8.0,
						ValuesRange: utils.ValueRange{Min: 1, Max: 9},
					},
				},
				//ideal a: 5, b: 2, c: 3
				PerReferencePointsDifferences: []ReferencePointsDifference{{
					Alternative: consideredAlternatives[0],
					ReferencePointsDifference: []ReferencePointDifference{{
						ReferencePoint: "ideal",
						Coefficients:   model.Weights{"a": -2, "b": 0, "c": 0},
					}},
				}, {
					Alternative: consideredAlternatives[1],
					ReferencePointsDifference: []ReferencePointDifference{{
						ReferencePoint: "ideal",
						Coefficients:   model.Weights{"a": -1.5, "b": -1, "c": 0.125},
					}},
				}, {
					Alternative: consideredAlternatives[2],
					ReferencePointsDifference: []ReferencePointDifference{{
						ReferencePoint: "ideal",
						Coefficients:   model.Weights{"a": -0.5, "b": 0.25, "c": -0.25},
					}},
				}, {
					Alternative: notConsideredCriteria[0],
					ReferencePointsDifference: []ReferencePointDifference{{
						ReferencePoint: "ideal",
						Coefficients:   model.Weights{"a": 0, "b": 0.5, "c": -1.5},
					}},
				}, {
					Alternative: notConsideredCriteria[1],
					ReferencePointsDifference: []ReferencePointDifference{{
						ReferencePoint: "ideal",
						Coefficients:   model.Weights{"a": -1, "b": -0.5, "c": 0.25},
					}},
				}},
				ApplierResult: NewCriterionAnchoringApplierResult{
					ReferenceCriterion: criteria[0],
					AddedCriteria: []AddedCriterion{{
						Id:   "__anchoring_criterion_ideal",
						Type: model.Gain,
						MethodParameters: testUtils.DummyMethodParameters{
							Criteria: []string{"__anchoring_criterion_ideal"},
						},
						AlternativesValues: model.Weights{
							"1": 1.3333333333333335, "2": 1, "3": 1.75, "4": 1, "5": 1.5833333333333335,
						},
						ValuesRange: utils.ValueRange{
							Min: 1,
							Max: 1.75,
						},
					}},
				},
			},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Anchoring{
				anchoringEvaluators:       tt.fields.anchoringEvaluators,
				referencePointsEvaluators: tt.fields.referencePointsEvaluators,
				anchoringAppliers:         tt.fields.anchoringAppliers,
			}
			if got := a.Apply(nil, tt.args.current, tt.args.props, tt.args.listener); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Apply():\n dmp   %v \n props %v,\nwant:\n dmp   %v\n props %v", got.DMP, got.Props, tt.want.DMP, tt.want.Props)
			}
		})
	}
}

func TestAnchoring_evaluateAnchoringAlternatives(t *testing.T) {
	type fields struct {
		anchoringEvaluators       []AnchoringEvaluator
		referencePointsEvaluators []ReferencePointsEvaluator
		anchoringAppliers         []AnchoringApplier
	}
	type args struct {
		allAlternatives []model.AlternativeWithCriteria
		parsedProps     *AnchoringParams
		criteria        *model.Criteria
	}

	calculateTestArgs := func(name string) args {
		return args{
			allAlternatives: []model.AlternativeWithCriteria{{
				Id:       "1",
				Criteria: model.Weights{"a": 1, "b": 2, "c": 4},
			}, {
				Id:       "2",
				Criteria: model.Weights{"a": 0, "b": 3, "c": 2},
			}, {
				Id:       "3",
				Criteria: model.Weights{"a": 5, "b": 1, "c": 1},
			}},
			parsedProps: &AnchoringParams{
				AnchoringAlternatives: []AnchoringAlternative{{Alternative: "1"}, {Alternative: "2"}, {Alternative: "3"}},
				ReferencePoints:       FunctionDefinition{Function: name},
			},
			criteria: &model.Criteria{{
				Id: "a", Type: model.Gain,
			}, {
				Id: "b", Type: model.Gain,
			}, {
				Id: "c", Type: model.Cost,
			}},
		}
	}
	referenceEvaluators := fields{
		referencePointsEvaluators: []ReferencePointsEvaluator{
			&IdealReferenceAlternativeEvaluator{},
			&NadirReferenceAlternativeEvaluator{},
		},
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []model.AlternativeWithCriteria
	}{{
		name:   "idealAlternative",
		fields: referenceEvaluators,
		args:   calculateTestArgs(IdealReferenceAltEvaluator),
		want: []model.AlternativeWithCriteria{{
			Id:       IdealReferenceAltEvaluator,
			Criteria: model.Weights{"a": 5, "b": 3, "c": 1},
		}},
	}, {
		name:   "nadirAlternative",
		fields: referenceEvaluators,
		args:   calculateTestArgs(NadirReferenceAltEvaluator),
		want: []model.AlternativeWithCriteria{{
			Id:       NadirReferenceAltEvaluator,
			Criteria: model.Weights{"a": 0, "b": 1, "c": 4},
		}},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Anchoring{
				referencePointsEvaluators: tt.fields.referencePointsEvaluators,
			}
			if got := a.evaluateAnchoringAlternatives(tt.args.allAlternatives, tt.args.parsedProps, tt.args.criteria); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("evaluateAnchoringAlternatives() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateDiffsPerReferencePoint(t *testing.T) {
	type args struct {
		consideredAlternatives []model.AlternativeWithCriteria
		referencePoints        []model.AlternativeWithCriteria
		criteria               *model.Criteria
		scaleRatios            CriteriaScaling
		loss                   AnchoringWithParams
		gain                   AnchoringWithParams
	}
	alternatives := []model.AlternativeWithCriteria{{
		Id: "a", Criteria: model.Weights{"1": 5, "2": 14},
	}, {
		Id: "b", Criteria: model.Weights{"1": 3.5, "2": 10},
	}, {
		Id: "c", Criteria: model.Weights{"1": 3, "2": 11},
	}}

	prepareDifAlternatives := func(refAlt model.AlternativeWithCriteria) args {
		return args{
			consideredAlternatives: alternatives,
			referencePoints:        []model.AlternativeWithCriteria{refAlt},
			criteria:               &model.Criteria{{Id: "1", Type: model.Gain}, {Id: "2", Type: model.Cost}},
			scaleRatios: CriteriaScaling{
				"1": ScaleWithValueRange{
					Scale:       1 / 2.0,
					ValuesRange: utils.ValueRange{Min: 3, Max: 5},
				},
				"2": ScaleWithValueRange{
					Scale:       1 / 4.0,
					ValuesRange: utils.ValueRange{Min: 10, Max: 14},
				},
			},
			loss: AnchoringWithParams{
				fun:    &LinearAnchoringEvaluator{},
				params: &utils.LinearFunctionParameters{A: 2, B: 0},
			},
			gain: AnchoringWithParams{
				fun:    &LinearAnchoringEvaluator{},
				params: &utils.LinearFunctionParameters{A: 1, B: 0},
			},
		}
	}

	tests := []struct {
		name string
		args args
		want []ReferencePointsDifference
	}{{
		name: "ref point to ideal alternative",
		args: prepareDifAlternatives(model.AlternativeWithCriteria{
			Id:       "ideal",
			Criteria: model.Weights{"1": 5, "2": 10},
		}),
		want: []ReferencePointsDifference{{
			Alternative: alternatives[0],
			ReferencePointsDifference: []ReferencePointDifference{
				{
					ReferencePoint: "ideal",
					Coefficients:   model.Weights{"1": 0, "2": -2},
				},
			},
		}, {
			Alternative: alternatives[1],
			ReferencePointsDifference: []ReferencePointDifference{
				{
					ReferencePoint: "ideal",
					Coefficients:   model.Weights{"1": -1.5, "2": 0},
				},
			},
		}, {
			Alternative: alternatives[2],
			ReferencePointsDifference: []ReferencePointDifference{
				{
					ReferencePoint: "ideal",
					Coefficients:   model.Weights{"1": -2, "2": -0.5},
				},
			},
		}},
	}, {
		name: "ref point to nadir alternative",
		args: prepareDifAlternatives(model.AlternativeWithCriteria{
			Id:       "nadir",
			Criteria: model.Weights{"1": 3, "2": 14},
		}),
		want: []ReferencePointsDifference{{
			Alternative: alternatives[0],
			ReferencePointsDifference: []ReferencePointDifference{
				{
					ReferencePoint: "nadir",
					Coefficients:   model.Weights{"1": 1, "2": 0},
				},
			},
		}, {
			Alternative: alternatives[1],
			ReferencePointsDifference: []ReferencePointDifference{
				{
					ReferencePoint: "nadir",
					Coefficients:   model.Weights{"1": 0.25, "2": 1},
				},
			},
		}, {
			Alternative: alternatives[2],
			ReferencePointsDifference: []ReferencePointDifference{
				{
					ReferencePoint: "nadir",
					Coefficients:   model.Weights{"1": 0, "2": 0.75},
				},
			},
		}},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateDiffsPerReferencePoint(
				tt.args.consideredAlternatives, tt.args.referencePoints, tt.args.criteria,
				tt.args.scaleRatios, tt.args.loss, tt.args.gain,
			); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculateDiffsPerReferencePoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_evaluatePerCriterionNormalizationScaleRatio(t *testing.T) {
	type args struct {
		criteria        *model.Criteria
		allAlternatives []model.AlternativeWithCriteria
	}
	tests := []struct {
		name string
		args args
		want CriteriaScaling
	}{
		{
			name: "all",
			args: args{
				criteria: &model.Criteria{{
					Id:   "1",
					Type: model.Gain,
				}, {
					Id:   "2",
					Type: model.Gain,
				}, {
					Id:   "3",
					Type: model.Cost, // <---- !!!! does not make a difference
				}},
				allAlternatives: []model.AlternativeWithCriteria{{
					Id:       "a",
					Criteria: model.Weights{"1": 1, "2": 4, "3": 3},
				}, {
					Id:       "b",
					Criteria: model.Weights{"1": 0.5, "2": 0.75, "3": 4},
				}, {
					Id:       "c",
					Criteria: model.Weights{"1": 2, "2": 0.25, "3": 2},
				}},
			},
			want: CriteriaScaling{
				"1": ScaleWithValueRange{
					Scale:       1 / 1.5,
					ValuesRange: utils.ValueRange{Min: 0.5, Max: 2},
				},
				"2": ScaleWithValueRange{
					Scale:       1 / 3.75,
					ValuesRange: utils.ValueRange{Min: 0.25, Max: 4},
				},
				"3": ScaleWithValueRange{
					Scale:       1 / 2.0,
					ValuesRange: utils.ValueRange{Min: 2, Max: 4},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := evaluatePerCriterionNormalizationScaleRatio(tt.args.criteria, tt.args.allAlternatives); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("evaluatePerCriterionNormalizationScaleRatio() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fetchAnchoringAlternativesWithCriteria(t *testing.T) {
	type args struct {
		alternatives          *[]model.AlternativeWithCriteria
		anchoringAlternatives *[]AnchoringAlternative
	}
	alternatives := []model.AlternativeWithCriteria{
		{
			Id:       "1",
			Criteria: model.Weights{"a": 1, "b": 2, "c": 3},
		}, {
			Id:       "2",
			Criteria: model.Weights{"a": 12, "b": 23, "c": 9},
		}, {
			Id:       "3",
			Criteria: model.Weights{"a": 4, "b": 6, "c": 8},
		},
	}
	tests := []struct {
		name string
		args args
		want *[]AnchoringAlternativeWithCriteria
	}{{
		name: "merging consideredAlternatives",
		args: args{
			alternatives: &alternatives,
			anchoringAlternatives: &[]AnchoringAlternative{{
				Alternative: "1",
				Coefficient: 0.5,
			}, {
				Alternative: "3",
				Coefficient: 0.75,
			}},
		},
		want: &[]AnchoringAlternativeWithCriteria{{
			Alternative: alternatives[0],
			Coefficient: 0.5,
		}, {
			Alternative: alternatives[2],
			Coefficient: 0.75,
		}},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fetchAnchoringAlternativesWithCriteria(tt.args.alternatives, tt.args.anchoringAlternatives); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetchAnchoringAlternativesWithCriteria() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseProps(t *testing.T) {
	type args struct {
		props *model.BiasProps
	}
	var exampleProps model.BiasProps = utils.Map{
		"anchoringAlternatives": []utils.Map{{
			"alternative": "1", "coefficient": 0.5,
		}, {
			"alternative": "xyz", "coefficient": 1,
		},
		},
		"gain": utils.Map{
			"function": "linear",
			"params":   utils.Map{"a": 1, "b": 2},
		},
		"loss": utils.Map{
			"function": "expFromZero",
			"params":   utils.Map{"alpha": 0.25, "multiplier": 0.5},
		},
		"referencePoints": utils.Map{
			"function": "ideal",
		},
		"applier": utils.Map{
			"function": "inline",
		},
	}
	tests := []struct {
		name string
		args args
		want *AnchoringParams
	}{{
		name: "deserialize from camelCase json",
		args: args{
			props: &exampleProps,
		},
		want: &AnchoringParams{
			AnchoringAlternatives: []AnchoringAlternative{{
				Alternative: "1",
				Coefficient: 0.5,
			}, {
				Alternative: "xyz",
				Coefficient: 1,
			}},
			Gain: FunctionDefinition{
				Function: "linear",
				Params:   utils.Map{"a": 1, "b": 2},
			},
			Loss: FunctionDefinition{
				Function: "expFromZero",
				Params:   utils.Map{"alpha": 0.25, "multiplier": 0.5},
			},
			ReferencePoints: FunctionDefinition{
				Function: "ideal",
				Params:   nil,
			},
			Applier: FunctionDefinition{
				Function: "inline",
				Params:   nil,
			},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseProps(tt.args.props); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseProps() = %v, want %v", got, tt.want)
			}
		})
	}
}
