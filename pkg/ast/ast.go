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

type PropertyLabelsExpr struct {
	Atom         Expr
	PropertyKeys []SchemaName
	Labels       []SchemaName
}

type SchemaName interface {
	schemaNameNode()
}

type SymbolName struct {
	Identifier scanner.Token
	Type       SymbolType
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
	SymbolName *SymbolName
	N          *scanner.Token
}

type CaseExpr struct {
	Init         Expr
	Alternatives Expr
	Else         Expr
}

type CaseAltExpr struct {
	When Expr
	Then Expr
}

func (e *OpExpr) exprNode()             {}
func (e *UnaryExpr) exprNode()          {}
func (e *BinaryExpr) exprNode()         {}
func (e *TernaryExpr) exprNode()        {}
func (e *ListExpr) exprNode()           {}
func (e *Literal) exprNode()            {}
func (e *PropertyLabelsExpr) exprNode() {}
func (e *Parameter) exprNode()          {}
func (e *CaseExpr) exprNode()           {}
func (e *CaseAltExpr) exprNode()        {}

func (s *SymbolName) schemaNameNode()   {}
func (r *ReservedWord) schemaNameNode() {}
