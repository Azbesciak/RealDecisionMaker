package model

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

//go:generate easytags $GOFILE json:camel

type HeuristicsParams = []HeuristicParams
type HeuristicParams struct {
	Name     string                 `json:"name"`
	Disabled bool                   `json:"disabled"`
	Props    map[string]interface{} `json:"props"`
}
type HeuristicsMap = map[string]Heuristic
type Heuristics = []Heuristic

type Heuristic interface {
	utils.Identifiable
	Process(dm *DecisionMaker) *DecisionMaker
}

func AsHeuristicsMap(h *Heuristics) *HeuristicsMap {
	result := make(HeuristicsMap, len(*h))
	for _, heu := range *h {
		result[heu.Identifier()] = heu
	}
	return &result
}

func ChooseHeuristics(available *HeuristicsMap, choose *HeuristicsParams) *Heuristics {
	var result Heuristics
	for _, k := range *choose {
		if k.Disabled {
			continue
		}
		heu, ok := (*available)[k.Name]
		if !ok {
			var keys []string
			for k := range *available {
				keys = append(keys, k)
			}
			panic(fmt.Errorf("heuristic '%s' not found, available are '%s'", k.Name, keys))
		}
		result = append(result, heu)
	}
	return &result
}
