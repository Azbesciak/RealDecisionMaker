package majority

import (
	"github.com/Azbesciak/RealDecisionMaker/lib/model"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

type DrawResolver interface {
	utils.Identifiable
	Resolve(
		currentEval, newEval model.Weight,
		sameBuffer []model.AlternativeResult,
		worseThanCurrent [][]model.AlternativeResult,
		current, another model.AlternativeWithCriteria,
		generator utils.ValueGenerator,
	) *DrawResolution
}

type DrawResolution struct {
	sameBuffer       []model.AlternativeResult
	worseThanCurrent [][]model.AlternativeResult
	current          model.AlternativeWithCriteria
}

const DrawAllowedResolverName = "allow"

type DrawAllowedResolver struct {
}

func (d *DrawAllowedResolver) Identifier() string {
	return DrawAllowedResolverName
}

func (d *DrawAllowedResolver) Resolve(
	currentEval, newEval model.Weight,
	sameBuffer []model.AlternativeResult,
	worseThanCurrent [][]model.AlternativeResult,
	current, another model.AlternativeWithCriteria,
	_ utils.ValueGenerator,
) *DrawResolution {
	sameBuffer = append(sameBuffer, model.AlternativeResult{
		Alternative: another,
		Evaluation: MajorityEvaluation{
			Value:                    newEval,
			ComparedWith:             current.Id,
			ComparedAlternativeValue: currentEval,
		},
	})
	return &DrawResolution{
		sameBuffer:       sameBuffer,
		worseThanCurrent: worseThanCurrent,
		current:          current,
	}
}

const CurrentIsWinnerResolverName = "current"

type CurrentIsWinnerDrawResolver struct {
}

func (d *CurrentIsWinnerDrawResolver) Identifier() string {
	return CurrentIsWinnerResolverName
}

func (d *CurrentIsWinnerDrawResolver) Resolve(
	currentEval, newEval model.Weight,
	sameBuffer []model.AlternativeResult,
	worseThanCurrent [][]model.AlternativeResult,
	current, another model.AlternativeWithCriteria,
	_ utils.ValueGenerator,
) *DrawResolution {
	worseThanCurrent = append(worseThanCurrent, []model.AlternativeResult{{
		Alternative: another,
		Evaluation: MajorityEvaluation{
			Value:                    newEval,
			ComparedWith:             current.Id,
			ComparedAlternativeValue: currentEval,
		},
	}})
	return &DrawResolution{
		sameBuffer:       sameBuffer,
		worseThanCurrent: worseThanCurrent,
		current:          current,
	}
}

const NewerIsWinnerResolverName = "newer"

type NewerIsWinnerResolver struct {
}

func (d *NewerIsWinnerResolver) Identifier() string {
	return NewerIsWinnerResolverName
}

func (d *NewerIsWinnerResolver) Resolve(
	currentEval, newEval model.Weight,
	sameBuffer []model.AlternativeResult,
	worseThanCurrent [][]model.AlternativeResult,
	current, another model.AlternativeWithCriteria,
	_ utils.ValueGenerator,
) *DrawResolution {
	sameBuffer = append(sameBuffer, model.AlternativeResult{
		Alternative: current,
		Evaluation: MajorityEvaluation{
			Value:                    currentEval,
			ComparedWith:             another.Id,
			ComparedAlternativeValue: newEval,
		},
	})
	current = another
	worseThanCurrent = append(worseThanCurrent, sameBuffer)
	sameBuffer = make([]model.AlternativeResult, 0)
	return &DrawResolution{
		sameBuffer:       sameBuffer,
		worseThanCurrent: worseThanCurrent,
		current:          current,
	}
}

const RandomIsWinnerResolverName = "random"

type RandomWinnerResolver struct {
	newer   NewerIsWinnerResolver
	current CurrentIsWinnerDrawResolver
}

func (d *RandomWinnerResolver) Identifier() string {
	return RandomIsWinnerResolverName
}

func (d *RandomWinnerResolver) Resolve(
	currentEval, newEval model.Weight,
	sameBuffer []model.AlternativeResult,
	worseThanCurrent [][]model.AlternativeResult,
	current, another model.AlternativeWithCriteria,
	generator utils.ValueGenerator,
) *DrawResolution {
	if generator() < 0.5 {
		return d.current.Resolve(currentEval, newEval, sameBuffer, worseThanCurrent, current, another, generator)
	} else {
		return d.newer.Resolve(currentEval, newEval, sameBuffer, worseThanCurrent, current, another, generator)
	}
}
