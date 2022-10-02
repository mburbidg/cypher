package scanner

type TokenType int

const (
	// Special tokens
	Illegal TokenType = iota
	EndOfInput

	// Punctuation
	Period
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
	String

	// Reserved words

	// Clauses
	Create
	Delete
	Detach
	Exists
	Match
	Merge
	Optional
	Remove
	Return
	Set
	Union
	Unwind
	With

	// Subclauses
	Limit
	Order
	Skip
	Where

	// Modifiers
	Asc
	Ascending
	By
	Desc
	Descending
	On

	// Expressions
	All
	Case
	Else
	End
	Then
	When

	// Operators
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

	// Literals
	False
	Null
	True
)
