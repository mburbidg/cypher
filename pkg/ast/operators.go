package ast

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
	CountAll
	AllOp
	AnyOp
	NoneOp
	SingleOp
	Equal
	NotEqual
	LessThan
	GreaterThan
	LessThanOrEqual
	GreaterThanOrEqual
	Add
	Subtract
	Multiply
	Divide
	Modulo
	PowerOf
)

const ()
