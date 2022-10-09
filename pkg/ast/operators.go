package ast

import "github.com/mburbidg/cypher/pkg/scanner"

type Operator int

const (
	Xor Operator = iota
	And
	Or
	Not
	Negate
	IsNull
	IsNotNull
	StringOrListOp
	StartsWith
	EndsWith
	Contains
	InList
	ListIndex
	ListRange
)

const (
	Equal              Operator = Operator(scanner.Equal)
	NotEqual           Operator = Operator(scanner.NotEqual)
	LessThan           Operator = Operator(scanner.LessThan)
	GreaterThan        Operator = Operator(scanner.GreaterThan)
	LessThanOrEqual    Operator = Operator(scanner.LessThanOrEqual)
	GreaterThanOrEqual Operator = Operator(scanner.GreaterThanOrEqual)
	Add                Operator = Operator(scanner.Plus)
	Subtract           Operator = Operator(scanner.Minus)
	Multiply           Operator = Operator(scanner.Star)
	Divide             Operator = Operator(scanner.ForwardSlash)
	Modulo             Operator = Operator(scanner.Percent)
	PowerOf            Operator = Operator(scanner.Caret)
)
