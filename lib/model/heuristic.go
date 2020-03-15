package model

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

//go:generate easytags $GOFILE json:camel

type HeuristicsParams = []HeuristicParams
type HeuristicProps = interface{}
type HeuristicsMap = map[string]Heuristic
type Heuristics = []Heuristic
type HeuristicsWithProps = []HeuristicWithProps

type HeuristicParams struct {
	Name     string         `json:"name"`
	Disabled bool           `json:"disabled"`
	Props    HeuristicProps `json:"props"`
}

type HeuristicWithProps struct {
	Heuristic *Heuristic      `json:"heuristic"`
	Props     *HeuristicProps `json:"props"`
}

type Heuristic interface {
	utils.Identifiable
	Process(dm *DecisionMaker) *DecisionMaker
}

type HeuristicResult struct {
	Dm    *DecisionMaker  `json:"dm"`
	Props *HeuristicProps `json:"props"`
}

func AsHeuristicsMap(h *Heuristics) *HeuristicsMap {
	result := make(HeuristicsMap, len(*h))
	for _, heu := range *h {
		result[heu.Identifier()] = heu
	}
	return &result
}

func ChooseHeuristics(available *HeuristicsMap, choose *HeuristicsParams) *HeuristicsWithProps {
	var result HeuristicsWithProps
	for _, props := range *choose {
		if props.Disabled {
			continue
		}
		heu, ok := (*available)[props.Name]
		if !ok {
			var keys []string
			for k := range *available {
				keys = append(keys, k)
			}
			panic(fmt.Errorf("heuristic '%s' not found, available are '%s'", props.Name, keys))
		}
		result = append(result, HeuristicWithProps{Heuristic: &heu, Props: &props.Props})
	}
	return &result
}
