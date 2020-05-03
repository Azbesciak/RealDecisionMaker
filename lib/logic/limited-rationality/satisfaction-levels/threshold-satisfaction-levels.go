package satisfaction_levels

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
)

//go:generate easytags $GOFILE json:camel
type thresholdSatisfactionLevels struct {
	Thresholds   []model.Weights `json:"thresholds"`
	currentIndex int
}

func (t *thresholdSatisfactionLevels) Initialize(dmp *model.DecisionMakingParams) {
	t.currentIndex = -1
	for i, threshold := range t.Thresholds {
		for _, c := range dmp.Criteria {
			if _, ok := threshold[c.Id]; !ok {
				panic(fmt.Errorf("value of criterion '%s' for threshold %d not found in %v", c.Id, i, threshold.AsKeyValue()))
			}
		}
	}
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

const Thresholds = "thresholds"

func (t *ThresholdSatisfactionLevelsSource) Identifier() string {
	return Thresholds
}

func (t *ThresholdSatisfactionLevelsSource) BlankParams() SatisfactionLevels {
	return &thresholdSatisfactionLevels{}
}
