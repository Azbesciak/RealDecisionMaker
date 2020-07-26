package satisfaction_levels

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
	"sort"
)

//go:generate easytags $GOFILE json:camel
type ThresholdSatisfactionLevels struct {
	Thresholds   []model.Weights `json:"thresholds"`
	currentIndex int
}

func (t *ThresholdSatisfactionLevels) Initialize(dmp *model.DecisionMakingParams) {
	t.currentIndex = -1
	for i, threshold := range t.Thresholds {
		for _, c := range dmp.Criteria {
			if _, ok := threshold[c.Id]; !ok {
				panic(fmt.Errorf("value of criterion '%s' for threshold %d not found in %v", c.Id, i, threshold.AsKeyValue()))
			}
		}
	}
}

func (t *ThresholdSatisfactionLevels) HasNext() bool {
	return t.currentIndex+1 < len(t.Thresholds)
}

func (t *ThresholdSatisfactionLevels) Next() model.Weights {
	t.currentIndex += 1
	return t.Thresholds[t.currentIndex]
}

type ThresholdSatisfactionLevelsSource struct {
	ascending bool
}

type ThresholdsUpdate struct {
	Thresholds []model.Weights `json:"thresholds"`
}

func (t *ThresholdSatisfactionLevelsSource) OnCriterionAdded(
	criterion *model.Criterion,
	referenceCriterion *model.Criterion,
	params SatisfactionLevels,
	generator utils.ValueGenerator,
) ParamsAddition {
	pParams := fetchParams(params)
	thresholdsValues := assignNewThresholds(pParams, referenceCriterion, generator)
	sortThresholds(thresholdsValues, t.ascending)
	thresholds := mapThresholdsToEntries(criterion, thresholdsValues)
	return ThresholdsUpdate{Thresholds: thresholds}
}

func mapThresholdsToEntries(criterion *model.Criterion, thresholdsValues []model.Weight) []model.Weights {
	thresholds := make([]model.Weights, len(thresholdsValues))
	for i, threshold := range thresholdsValues {
		thresholds[i] = model.Weights{criterion.Id: threshold}
	}
	return thresholds
}

func assignNewThresholds(params *ThresholdSatisfactionLevels, referenceCriterion *model.Criterion, generator utils.ValueGenerator) []model.Weight {
	thresholds := make([]model.Weight, len(params.Thresholds))
	for i, threshold := range params.Thresholds {
		thresholds[i] = threshold.Fetch(referenceCriterion.Id) * generator()
	}
	return thresholds
}

func sortThresholds(thresholds []model.Weight, ascending bool) {
	sort.Slice(thresholds, func(i, j int) bool {
		less := thresholds[i] < thresholds[j]
		if ascending {
			return less
		} else {
			return !less
		}
	})
}

func (t *ThresholdSatisfactionLevelsSource) OnCriteriaRemoved(leftCriteria *model.Criteria, params SatisfactionLevels) SatisfactionLevels {
	pParams := fetchParams(params)
	thresholds := pParams.preserveLeftThresholds(leftCriteria)
	return &ThresholdSatisfactionLevels{
		Thresholds:   thresholds,
		currentIndex: pParams.currentIndex,
	}
}

func (t *ThresholdSatisfactionLevels) preserveLeftThresholds(leftCriteria *model.Criteria) []model.Weights {
	thresholds := make([]model.Weights, len(t.Thresholds))
	for i, threshold := range t.Thresholds {
		thresholds[i] = *threshold.PreserveOnly(leftCriteria)
	}
	return thresholds
}

func (t *ThresholdSatisfactionLevelsSource) Merge(params SatisfactionLevels, addition ParamsAddition) SatisfactionLevels {
	pParams := fetchParams(params)
	add := addition.(ThresholdsUpdate)
	newThresholds := pParams.merge(add)
	return &ThresholdSatisfactionLevels{
		Thresholds:   newThresholds,
		currentIndex: pParams.currentIndex,
	}
}

func (t *ThresholdSatisfactionLevels) merge(add ThresholdsUpdate) []model.Weights {
	newThresholds := make([]model.Weights, len(t.Thresholds))
	for i, thresholds := range t.Thresholds {
		newThresholds[i] = *thresholds.Merge(&add.Thresholds[i])
	}
	return newThresholds
}

func fetchParams(params SatisfactionLevels) *ThresholdSatisfactionLevels {
	if p, ok := params.(*ThresholdSatisfactionLevels); !ok {
		panic(fmt.Errorf("threshold params shold be instance of ThresholdSatisfactionLevels"))
	} else {
		return p
	}
}

const Thresholds = "thresholds"

func (t *ThresholdSatisfactionLevelsSource) Identifier() string {
	return Thresholds
}

func (t *ThresholdSatisfactionLevelsSource) BlankParams() SatisfactionLevels {
	return &ThresholdSatisfactionLevels{}
}

var IncreasingThresholds = ThresholdSatisfactionLevelsSource{
	ascending: true,
}

var DecreasingThresholds = ThresholdSatisfactionLevelsSource{
	ascending: false,
}
