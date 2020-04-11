package fatigue

//go:generate easytags $GOFILE json:camel

type ConstFatigueParams struct {
	Value float64 `json:"value"`
}

type ConstFatigueFunction struct {
}

const FatConstFunc = "const"

func (c *ConstFatigueFunction) Name() string {
	return FatConstFunc
}

func (c *ConstFatigueFunction) BlankParams() FatigueFunctionParams {
	return &ConstFatigueParams{}
}

func (c *ConstFatigueFunction) Evaluate(params FatigueFunctionParams) float64 {
	return params.(*ConstFatigueParams).Value
}
