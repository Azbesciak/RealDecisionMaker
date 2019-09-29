package model

type DecisionMakerProperties struct {
	Focus   int
	Fatigue int
}

type Identifiable interface {
	Identifier() string
}

type AlternativeResult struct {
	Alternative AlternativeWithCriteria
	Value       Weight
}

type Alternative struct {
	Id string
}

type AlternativeWithCriteria struct {
	Alternative
	Criteria Weights
}
