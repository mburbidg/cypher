package cypher

import "fmt"

const (
	VariableAlreadyBound         = "VariableAlreadyBound"
	UndefinedVariable            = "UndefinedVariable"
	NoSingleRelationshipType     = "NoSingleRelationshipType"
	RequiresDirectedRelationship = "RequiresDirectedRelationship"
	CreatingVarLength            = "CreatingVarLength"
	InvalidParameterUse          = "InvalidParameterUse"
)

type CypherErr struct {
	Msg  string
	Code string
}

func (err CypherErr) Error() string {
	return err.Msg
}

func NewVariableAlreadyBoundErr(name string) error {
	return &CypherErr{
		Msg:  fmt.Sprintf("variable already bound: '%s'", name),
		Code: VariableAlreadyBound,
	}
}

func NewUndefinedVariableErr(name string) error {
	return &CypherErr{
		Msg:  fmt.Sprintf("undefined variable: '%s'", name),
		Code: UndefinedVariable,
	}
}

func NewNoSingleRelationshipType() error {
	return &CypherErr{
		Msg:  fmt.Sprintf("no single relationship type"),
		Code: NoSingleRelationshipType,
	}
}

func NewRequiresDirectedRelationship() error {
	return &CypherErr{
		Msg:  fmt.Sprintf("required directed relationship"),
		Code: RequiresDirectedRelationship,
	}
}

func NewCreatingVarLength() error {
	return &CypherErr{
		Msg:  fmt.Sprintf("creating variable-length relationship"),
		Code: CreatingVarLength,
	}
}

func NewInvalidParameterUse() error {
	return &CypherErr{
		Msg:  fmt.Sprintf("invalid parameter use"),
		Code: InvalidParameterUse,
	}
}
