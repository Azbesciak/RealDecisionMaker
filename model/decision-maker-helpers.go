package model

type AlternativeWeightFunction func(alternative *AlternativeWithCriteria) *AlternativeResult

func Rank(dm *DecisionMaker, pref AlternativeWeightFunction) *AlternativesRanking {
	results := make(AlternativeResults, len(dm.ChoseToMake))
	for i, r := range dm.ChoseToMake {
		var alternative = dm.Alternative(r)
		results[i] = *pref(&alternative)
	}
	return results.Ranking()
}

func ExtractWeights(dm *DecisionMaker) Weights {
	weights, ok := dm.MethodParameters["weights"]
	if !ok {
		panic("weights not found")
	}
	return weights.(Weights)
}
