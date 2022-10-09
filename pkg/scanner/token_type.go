package scanner

type TokenType int

const (
	// Special tokens
	Illegal TokenType = iota
	EndOfInput
	None

	// Punctuation
	Period
	Dotdot
	OpenParen
	CloseParen
	OpenBrace
	CloseBrace
	OpenBracket
	CloseBracket
	Plus
	Minus
	Star
	ForwardSlash
	Percent
	Caret
	Equal
	NotEqual
	LessThan
	GreaterThan
	LessThanOrEqual
	GreaterThanOrEqual
	DollarSign
	Colon

	Identifier
	Double
	Integer
	DecimalInteger
	HexInteger
	OctInteger
	String

	// Reserved words

	All
	Asc
	Ascending
	By
	Create
	Delete
	Desc
	Descending
	Detach
	Exists
	Limit
	Match
	Merge
	On
	Optional
	Order
	Remove
	Return
	Set
	Skip
	Where
	With
	Union
	Unwind
	And
	As
	Contains
	Distinct
	Ends
	In
	Is
	Not
	Or
	Starts
	Xor
	False
	True
	Null
	Constraint
	Do
	For
	Require
	Unique
	Case
	When
	Then
	Else
	End
	Mandatory
	Scalar
	Of
	Add
	Drop
)
