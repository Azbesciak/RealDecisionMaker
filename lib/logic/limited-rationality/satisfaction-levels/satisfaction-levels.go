package satisfaction_levels

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type SatisfactionLevels interface {
	Initialize(dmp *model.DecisionMakingParams)
	HasNext() bool
	Next() model.Weights
}

type SatisfactionLevelsSource interface {
	Name() string
	BlankParams() SatisfactionLevels
}

func Find(function string, params interface{}, functions []SatisfactionLevelsSource) SatisfactionLevels {
	if len(function) == 0 {
		panic(fmt.Errorf("satisfaction thresholds function not provided"))
	}
	for _, f := range functions {
		if f.Name() == function {
			functionParams := f.BlankParams()
			utils.DecodeToStruct(params, functionParams)
			return functionParams
		}
	}
	names := make([]string, len(functions))
	for i, f := range functions {
		names[i] = f.Name()
	}
	panic(fmt.Errorf("satisfaction thresholds function '%s' not found in functions %v", function, names))
}
