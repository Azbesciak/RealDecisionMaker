package satisfaction_levels

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
)

//go:generate easytags $GOFILE json:camel
type thresholdSatisfactionLevels struct {
	Thresholds   []model.Weights `json:"thresholds"`
	currentIndex int
}

func (t *thresholdSatisfactionLevels) Initialize(dmp *model.DecisionMakingParams) {
	t.currentIndex = -1
}

func (t *thresholdSatisfactionLevels) HasNext() bool {
	return t.currentIndex+1 < len(t.Thresholds)
}

func (t *thresholdSatisfactionLevels) Next() model.Weights {
	t.currentIndex += 1
	return t.Thresholds[t.currentIndex]
}

type ThresholdSatisfactionLevelsSource struct {
}

func (t ThresholdSatisfactionLevelsSource) Name() string {
	return "thresholds"
}

func (t ThresholdSatisfactionLevelsSource) BlankParams() SatisfactionLevels {
	return &thresholdSatisfactionLevels{}
}
