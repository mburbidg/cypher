package ast

import "github.com/mburbidg/cypher/pkg/scanner"

type Node interface {
}

type Expr interface {
	Node
	exprNode()
}

type OpExpr struct {
	Op Operator
}

type UnaryExpr struct {
	Op   Operator
	Expr Expr
}

type BinaryExpr struct {
	Left  Expr
	Op    Operator
	Right Expr
}

type TernaryExpr struct {
	E1 Expr
	Op Operator
	E2 Expr
	E3 Expr
}

type ListExpr struct {
	List []Expr
}

type ListComprehensionExpr struct {
	FilterExpr Expr
	Expr       Expr
}

type PropertyLabelsExpr struct {
	Atom         Expr
	PropertyKeys []SchemaName
	Labels       []SchemaName
}

type SchemaName interface {
	schemaNameNode()
}

type SymbolicNameSchemaName struct {
	SymbolicName SymbolicName
}

type ReservedWordSchemaName struct {
	TokenType scanner.TokenType
}

type SymbolicName interface {
	symbolicNameNode()
}

type SymbolicNameIdentifier struct {
	Identifier scanner.Token
	Type       SymbolType
}

type SymbolicNameHexLetter struct {
	Letter rune
}

type ReservedWord struct {
	Token scanner.Token
}

type Label struct {
}

type Literal struct {
	Kind  scanner.TokenType
	Value interface{}
}

type Parameter struct {
	SymbolicName SymbolicName
	N            *scanner.Token
}

type CaseExpr struct {
	Init         Expr
	Alternatives []*CaseAltNode
	Else         Expr
}

type CaseAltNode struct {
	When Expr
	Then Expr
}

type BuiltInExpr struct {
	Op   Operator
	Expr Expr
}

type FilterExpr struct {
	Variable  Expr
	InExpr    Expr
	WhereExpr Expr
}

type VariableExpr struct {
	SymbolicName SymbolicName
}

type PatternComprehensionExpr struct {
	Variable            Expr
	ReltionshipsPattern *ReltionshipsPattern
	WhereExpr           Expr
	PipeExpr            Expr
}

type ReltionshipsPattern struct {
}

func (e *OpExpr) exprNode()                   {}
func (e *UnaryExpr) exprNode()                {}
func (e *BinaryExpr) exprNode()               {}
func (e *TernaryExpr) exprNode()              {}
func (e *ListExpr) exprNode()                 {}
func (e *Literal) exprNode()                  {}
func (e *PropertyLabelsExpr) exprNode()       {}
func (e *Parameter) exprNode()                {}
func (e *CaseExpr) exprNode()                 {}
func (e *ListComprehensionExpr) exprNode()    {}
func (e *FilterExpr) exprNode()               {}
func (e *BuiltInExpr) exprNode()              {}
func (e *VariableExpr) exprNode()             {}
func (e *PatternComprehensionExpr) exprNode() {}

func (s *SymbolicNameIdentifier) symbolicNameNode() {}
func (s *SymbolicNameHexLetter) symbolicNameNode()  {}

func (r *SymbolicNameSchemaName) schemaNameNode() {}
func (r *ReservedWordSchemaName) schemaNameNode() {}
