package model

import (
	"fmt"
	"github.com/Azbesciak/RealDecisionMaker/lib/utils"
)

//go:generate easytags $GOFILE json:camel

type BiasesParams = []BiasParams
type BiasProps = interface{}
type BiasMap = map[string]Bias
type Biases = []Bias
type BiasesWithProps = []BiasWithProps

type BiasParams struct {
	Name     string    `json:"name"`
	Disabled bool      `json:"disabled"`
	Props    BiasProps `json:"props"`
}

type BiasWithProps struct {
	Bias  *Bias       `json:"bias"`
	Props *BiasParams `json:"props"`
}

type Bias interface {
	utils.Identifiable
	Apply(
		original, current *DecisionMakingParams,
		props *BiasProps,
		listener *BiasListener,
	) *BiasedResult
}
type BiasedResult struct {
	DMP   *DecisionMakingParams `json:"dm"`
	Props BiasProps             `json:"props"`
}

func AsBiasesMap(h *Biases) *BiasMap {
	result := make(BiasMap, len(*h))
	for _, bias := range *h {
		result[bias.Identifier()] = bias
	}
	return &result
}

func ChooseBiases(available *BiasMap, choose *BiasesParams) *BiasesWithProps {
	var result BiasesWithProps
	for _, props := range *choose {
		if props.Disabled {
			continue
		}
		bias, ok := (*available)[props.Name]
		if !ok {
			var keys []string
			for k := range *available {
				keys = append(keys, k)
			}
			panic(fmt.Errorf("bias '%s' not found, available are '%s'", props.Name, keys))
		}
		// required, props reference is mutable in range, during assignment it is copied.
		propsCopy := props
		result = append(result, BiasWithProps{Bias: &bias, Props: &propsCopy})
	}
	return &result
}

func UpdateBiasesProps(oldProps *BiasParams, update BiasProps) *BiasParams {
	return &BiasParams{
		Name:     oldProps.Name,
		Disabled: oldProps.Disabled,
		Props:    update,
	}
}
