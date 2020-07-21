package anchoring

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type ScaleWithValueRange struct {
	Scale       float64          `json:"scale"`
	ValuesRange utils.ValueRange `json:"valuesRange"`
}

type CriteriaScaling = map[string]ScaleWithValueRange

func evaluatePerCriterionNormalizationScaleRatio(criteria *model.Criteria, allAlternatives []model.AlternativeWithCriteria) CriteriaScaling {
	scaleRatios := make(CriteriaScaling, len(*criteria))
	for _, c := range *criteria {
		criterionRange := model.CriteriaValuesRange(&allAlternatives, &c)
		scale := model.GetNormalScaleRatio(criterionRange)
		scaleRatios[c.Id] = ScaleWithValueRange{
			Scale:       scale,
			ValuesRange: *criterionRange,
		}
	}
	return scaleRatios
}
