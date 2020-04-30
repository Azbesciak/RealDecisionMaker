package satisfaction_levels

import "github.com/Azbesciak/RealDecisionMaker/lib/model"

type SatisfactionLevels interface {
	Initialize(dmp *model.DecisionMakingParams)
	HasNext() bool
	Next() model.Weights
}

type SatisfactionLevelsSource interface {
	Name() string
	BlankParams() SatisfactionLevels
}
