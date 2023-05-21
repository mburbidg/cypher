package ast

import (
	scanner2 "github.com/mburbidg/cypher/scanner"
)

type Node interface {
}

type Query interface {
	Node
	Acceptor
	queryNode()
}

type ReadingClause interface {
	Node
	Acceptor
	readingClauseNode()
}

type UpdatingClause interface {
	Node
	Acceptor
	updatingClauseNode()
}

type SinglePartQuery struct {
	ReadingClause  []ReadingClause
	UpdatingClause []UpdatingClause
	*Projection
}

func (q *SinglePartQuery) Accept(visitor Visitor) error {
	if err := visitor.VisitSinglePartQueryEnter(q); err != nil {
		return err
	}
	for _, clause := range q.ReadingClause {
		if err := clause.Accept(visitor); err != nil {
			return err
		}
	}
	for _, clause := range q.UpdatingClause {
		if err := clause.Accept(visitor); err != nil {
			return err
		}
	}
	if q.Projection != nil {
		if err := q.Projection.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSinglePartQueryLeave(q)
}

type CreateClause struct {
	Pattern *Pattern
}

func (c *CreateClause) Accept(visitor Visitor) error {
	if err := visitor.VisitCreateEnter(c); err != nil {
		return err
	}
	if err := c.Pattern.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitCreateLeave(c)
}

type MatchClause struct {
	Optional  bool
	Pattern   *Pattern
	WhereExpr Expr
}

func (m *MatchClause) Accept(visitor Visitor) error {
	if err := visitor.VisitMatchEnter(m); err != nil {
		return err
	}
	if err := m.Pattern.Accept(visitor); err != nil {
		return err
	}
	if m.WhereExpr != nil {
		if err := m.WhereExpr.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitMatchLeave(m)
}

type Pattern struct {
	Parts []*PatternPart
}

func (p *Pattern) Accept(visitor Visitor) error {
	if err := visitor.VisitPatternEnter(p); err != nil {
		return err
	}
	for _, part := range p.Parts {
		if err := part.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitPatternLeave(p)
}

type PatternElement interface {
	Acceptor
	patternElementNode()
}

type PatternElementNested struct {
	Element PatternElement
}

func (p *PatternElementNested) Accept(visitor Visitor) error {
	if err := visitor.VisitPatternElementNestedEnter(p); err != nil {
		return err
	}
	if err := p.Element.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitPatternElementNestedLeave(p)
}

type PatternElementPattern struct {
	Left  *NodePattern
	Chain []*PatternElementChain
}

func (p *PatternElementPattern) Accept(visitor Visitor) error {
	if err := visitor.VisitPatternElementPatternEnter(p); err != nil {
		return err
	}
	if err := p.Left.Accept(visitor); err != nil {
		return err
	}
	for _, e := range p.Chain {
		if err := e.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitPatternElementPatternEnter(p)
}

type PatternPart struct {
	Variable SymbolicName
	Element  PatternElement
}

func (part *PatternPart) Accept(visitor Visitor) error {
	if err := visitor.VisitPatternPartEnter(part); err != nil {
		return err
	}
	if part.Variable != nil {
		if err := part.Variable.Accept(visitor); err != nil {
			return err
		}
	}
	if err := part.Element.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitPatternPartLeave(part)
}

type Projection struct {
	Distinct bool
	Items    *ProjectionItems
	Order    *SortOrder
	Skip     Expr
	Limit    Expr
}

func (p *Projection) Accept(visitor Visitor) error {
	if err := visitor.VisitProjectionEnter(p); err != nil {
		return err
	}
	if p.Items != nil {

	}
	if err := p.Items.Accept(visitor); err != nil {
		return err
	}
	if p.Order != nil {
		if err := p.Order.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitProjectionLeave(p)
}

type ProjectionItems struct {
	All   bool
	Items []*ProjectionItem
}

func (p *ProjectionItems) Accept(visitor Visitor) error {
	if err := visitor.VisitProjectionItemsEnter(p); err != nil {
		return err
	}
	for _, item := range p.Items {
		if err := item.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitProjectionItemsEnter(p)
}

type ProjectionItem struct {
	Expr     Expr
	Variable SymbolicName
}

func (p *ProjectionItem) Accept(visitor Visitor) error {
	if err := visitor.VisitProjectionItemEnter(p); err != nil {
		return err
	}
	if err := p.Expr.Accept(visitor); err != nil {
		return err
	}
	if p.Variable != nil {
		if err := p.Variable.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitProjectionItemEnter(p)
}

type SortOrder struct {
	Items []*SortItem
}

func (o *SortOrder) Accept(visitor Visitor) error {
	if err := visitor.VisitSortOrderEnter(o); err != nil {
		return err
	}
	for _, item := range o.Items {
		if err := item.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitSortOrderLeave(o)
}

type SortItem struct {
	Expr  Expr
	Order Order
}

func (item *SortItem) Accept(visitor Visitor) error {
	if err := visitor.VisitSortItemEnter(item); err != nil {
		return err
	}
	if err := item.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitSortItemLeave(item)
}

type Expr interface {
	Node
	Acceptor
	exprNode()
}

type OpExpr struct {
	Op Operator
}

func (expr *OpExpr) Accept(visitor Visitor) error {
	return visitor.VisitOpExpr(expr)
}

type UnaryExpr struct {
	Op   Operator
	Expr Expr
}

func (expr *UnaryExpr) Accept(visitor Visitor) error {
	if err := visitor.VisitUnaryExprEnter(expr); err != nil {
		return err
	}
	if err := expr.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitUnaryExprLeave(expr)
}

type BinaryExpr struct {
	Left  Expr
	Op    Operator
	Right Expr
}

func (expr *BinaryExpr) Accept(visitor Visitor) error {
	if err := visitor.VisitBinaryExprEnter(expr); err != nil {
		return err
	}
	if err := expr.Left.Accept(visitor); err != nil {
		return err
	}
	if err := expr.Right.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitBinaryExprLeave(expr)
}

type TernaryExpr struct {
	E1 Expr
	Op Operator
	E2 Expr
	E3 Expr
}

func (expr *TernaryExpr) Accept(visitor Visitor) error {
	if err := visitor.VisitTernaryExprEnter(expr); err != nil {
		return err
	}
	if err := expr.E1.Accept(visitor); err != nil {
		return err
	}
	if err := expr.E2.Accept(visitor); err != nil {
		return err
	}
	if err := expr.E3.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitTernaryExprEnter(expr)
}

type ListExpr struct {
	List []Expr
}

func (expr *ListExpr) Accept(visitor Visitor) error {
	if err := visitor.VisitListExprEnter(expr); err != nil {
		return err
	}
	for _, item := range expr.List {
		if err := item.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitListExprLeave(expr)
}

type ListComprehensionExpr struct {
	FilterExpr Expr
	Expr       Expr
}

func (expr *ListComprehensionExpr) Accept(visitor Visitor) error {
	if err := visitor.VisitListComprehensionExprEnter(expr); err != nil {
		return err
	}
	if err := expr.FilterExpr.Accept(visitor); err != nil {
		return err
	}
	if err := expr.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitListComprehensionExprLeave(expr)
}

type PropertyLabelsExpr struct {
	Atom         Expr
	PropertyKeys []SchemaName
	Labels       []SchemaName
}

func (expr *PropertyLabelsExpr) Accept(visitor Visitor) error {
	if err := visitor.VisitPropertyLabelsExprEnter(expr); err != nil {
		return err
	}
	if err := expr.Atom.Accept(visitor); err != nil {
		return err
	}
	for _, key := range expr.PropertyKeys {
		if err := key.Accept(visitor); err != nil {
			return err
		}
	}
	for _, label := range expr.Labels {
		if err := label.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitPropertyLabelsExprLeave(expr)
}

type SchemaName interface {
	Acceptor
	schemaNameNode()
}

type SymbolicNameSchemaName struct {
	SymbolicName SymbolicName
}

func (name *SymbolicNameSchemaName) Accept(visitor Visitor) error {
	if err := visitor.VisitSymbolicNameSchemaNameEnter(name); err != nil {
		return err
	}
	if err := name.SymbolicName.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitSymbolicNameSchemaNameLeave(name)
}

type ReservedWordSchemaName struct {
	TokenType scanner2.TokenType
}

func (name *ReservedWordSchemaName) Accept(visitor Visitor) error {
	return visitor.VisitReservedWordSchemaName(name)
}

type SymbolicName interface {
	Acceptor
	symbolicNameNode()
}

type SymbolicNameIdentifier struct {
	Identifier scanner2.Token
	Type       SymbolType
}

func (name *SymbolicNameIdentifier) Accept(visitor Visitor) error {
	return visitor.VisitSymbolicNameIdentifier(name)
}

type SymbolicNameHexLetter struct {
	Letter rune
}

func (name *SymbolicNameHexLetter) Accept(visitor Visitor) error {
	return visitor.VisitSymbolicNameHexLetter(name)
}

type ReservedWord struct {
	Token scanner2.Token
}

func (word *ReservedWord) Accept(visitor Visitor) error {
	return visitor.VisitReservedWord(word)
}

type Label struct {
}

func (label *Label) Accept(visitor Visitor) error {
	return visitor.VisitLabel(label)
}

type PrimitiveLiteral struct {
	Kind  scanner2.TokenType
	Value interface{}
}

func (literal *PrimitiveLiteral) Accept(visitor Visitor) error {
	return visitor.VisitPrimitiveLiteral(literal)
}

type ListLiteral struct {
	Items []Expr
}

func (literal *ListLiteral) Accept(visitor Visitor) error {
	if err := visitor.VisitListLiteralEnter(literal); err != nil {
		return err
	}
	for _, item := range literal.Items {
		if err := item.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitListLiteralLeave(literal)
}

type Parameter struct {
	SymbolicName SymbolicName
	N            *scanner2.Token
}

func (param *Parameter) Accept(visitor Visitor) error {
	if err := visitor.VisitParameterEnter(param); err != nil {
		return err
	}
	if err := param.SymbolicName.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitParameterLeave(param)
}

type CaseExpr struct {
	Init         Expr
	Alternatives []*CaseAltNode
	Else         Expr
}

func (expr *CaseExpr) Accept(visitor Visitor) error {
	if err := visitor.VisitCaseExprEnter(expr); err != nil {
		return err
	}
	if err := expr.Init.Accept(visitor); err != nil {
		return err
	}
	for _, alt := range expr.Alternatives {
		if err := alt.Accept(visitor); err != nil {
			return err
		}
	}
	if err := expr.Else.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitCaseExprEnter(expr)
}

type CaseAltNode struct {
	When Expr
	Then Expr
}

func (alt *CaseAltNode) Accept(visitor Visitor) error {
	if err := visitor.VisitCaseAltNodeEnter(alt); err != nil {
		return err
	}
	if err := alt.When.Accept(visitor); err != nil {
		return err
	}
	if err := alt.Then.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitCaseAltNodeLeave(alt)
}

type QuantifierExpr struct {
	Op   Operator
	Expr Expr
}

func (quantifier *QuantifierExpr) Accept(visitor Visitor) error {
	if err := visitor.VisitQuantifierExprEnter(quantifier); err != nil {
		return err
	}
	if err := quantifier.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitQuantifierExprEnter(quantifier)
}

type FilterExpr struct {
	Variable  SymbolicName
	InExpr    Expr
	WhereExpr Expr
}

func (filter *FilterExpr) Accept(visitor Visitor) error {
	if err := visitor.VisitFilterExprEnter(filter); err != nil {
		return err
	}
	if err := filter.Variable.Accept(visitor); err != nil {
		return err
	}
	if err := filter.InExpr.Accept(visitor); err != nil {
		return err
	}
	if err := filter.WhereExpr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitFilterExprLeave(filter)
}

type VariableExpr struct {
	SymbolicName SymbolicName
}

func (expr *VariableExpr) Accept(visitor Visitor) error {
	if err := visitor.VisitVariableExprEnter(expr); err != nil {
		return err
	}
	if err := expr.SymbolicName.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitVariableExprEnter(expr)
}

type PatternComprehensionExpr struct {
	Variable            SymbolicName
	ReltionshipsPattern Expr
	WhereExpr           Expr
	PipeExpr            Expr
}

func (expr *PatternComprehensionExpr) Accept(visitor Visitor) error {
	if err := visitor.VisitPatternComprehensionExprEnter(expr); err != nil {
		return err
	}
	if err := expr.Variable.Accept(visitor); err != nil {
		return err
	}
	if err := expr.ReltionshipsPattern.Accept(visitor); err != nil {
		return err
	}
	if err := expr.WhereExpr.Accept(visitor); err != nil {
		return err
	}
	if err := expr.PipeExpr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitPatternComprehensionExprLeave(expr)
}

type NodePattern struct {
	Variable   SymbolicName
	Labels     []SchemaName
	Properties *Properties
}

func (pattern *NodePattern) Accept(visitor Visitor) error {
	if err := visitor.VisitNodePatternEnter(pattern); err != nil {
		return err
	}
	if pattern.Variable != nil {
		if err := pattern.Variable.Accept(visitor); err != nil {
			return err
		}
	}
	for _, label := range pattern.Labels {
		if err := label.Accept(visitor); err != nil {
			return err
		}
	}
	if pattern.Properties != nil {
		if err := pattern.Properties.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitNodePatternEnter(pattern)
}

type MapLiteral struct {
	PropertyKeyNames []*PropertyKeyName
}

func (literal *MapLiteral) Accept(visitor Visitor) error {
	if err := visitor.VisitMapLiteralEnter(literal); err != nil {
		return err
	}
	for _, keyName := range literal.PropertyKeyNames {
		if err := keyName.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitMapLiteralLeave(literal)
}

type PropertyKeyName struct {
	Name SchemaName
	Expr Expr
}

func (name *PropertyKeyName) Accept(visitor Visitor) error {
	if err := visitor.VisitPropertyKeyNameEnter(name); err != nil {
		return err
	}
	if err := name.Name.Accept(visitor); err != nil {
		return err
	}
	if err := name.Expr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitPropertyKeyNameEnter(name)
}

type Properties struct {
	MapLiteral *MapLiteral
	Parameter  Expr
}

func (props *Properties) Accept(visitor Visitor) error {
	if err := visitor.VisitPropertiesEnter(props); err != nil {
		return err
	}
	if props.MapLiteral != nil {
		if err := props.MapLiteral.Accept(visitor); err != nil {
			return err
		}
	}
	if props.Parameter != nil {
		if err := props.Parameter.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitPropertiesEnter(props)
}

type RelationshipsPattern struct {
	Left  *NodePattern
	Chain []*PatternElementChain
}

func (pattern *RelationshipsPattern) Accept(visitor Visitor) error {
	if err := visitor.VisitRelationshipsPatternEnter(pattern); err != nil {
		return err
	}
	if err := pattern.Left.Accept(visitor); err != nil {
		return err
	}
	for _, elem := range pattern.Chain {
		if err := elem.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitRelationshipsPatternEnter(pattern)
}

type PatternElementChain struct {
	RelationshipPattern *RelationshipPattern
	Right               *NodePattern
}

func (chain *PatternElementChain) Accept(visitor Visitor) error {
	if err := visitor.VisitPatternElementChainEnter(chain); err != nil {
		return err
	}
	if err := chain.RelationshipPattern.Accept(visitor); err != nil {
		return err
	}
	if err := chain.Right.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitPatternElementChainEnter(chain)
}

type RelationshipPattern struct {
	Left               Relationship
	Right              Relationship
	RelationshipDetail *RelationshipDetail
}

func (pattern *RelationshipPattern) Accept(visitor Visitor) error {
	if err := visitor.VisitRelationshipPatternEnter(pattern); err != nil {
		return err
	}
	if pattern.RelationshipDetail != nil {
		if err := pattern.RelationshipDetail.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitRelationshipPatternLeave(pattern)
}

type RelationshipDetail struct {
	Variable          SymbolicName
	RelationshipTypes []SchemaName
	RangeLiteral      *RangeLiteral
	Properties        *Properties
}

func (detail *RelationshipDetail) Accept(visitor Visitor) error {
	if err := visitor.VisitRelationshipDetailEnter(detail); err != nil {
		return err
	}
	if detail.Variable != nil {
		if err := detail.Variable.Accept(visitor); err != nil {
			return err
		}
	}
	for _, rel := range detail.RelationshipTypes {
		if err := rel.Accept(visitor); err != nil {
			return err
		}
	}
	if err := detail.RangeLiteral.Accept(visitor); err != nil {
		return err
	}
	if detail.Properties != nil {
		if err := detail.Properties.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitRelationshipDetailLeave(detail)
}

type RangeLiteral struct {
	Begin int64
	End   int64
}

func (literal *RangeLiteral) Accept(visitor Visitor) error {
	return visitor.VisitRangeLiteral(literal)
}

type FunctionInvocation struct {
	FunctionName FunctionName
	Distinct     bool
	Args         []Expr
}

func (fnc *FunctionInvocation) Accept(visitor Visitor) error {
	if err := visitor.VisitFunctionInvocationEnter(fnc); err != nil {
		return err
	}
	if err := fnc.FunctionName.Accept(visitor); err != nil {
		return err
	}
	for _, arg := range fnc.Args {
		if err := arg.Accept(visitor); err != nil {
			return err
		}
	}
	return visitor.VisitFunctionInvocationLeave(fnc)
}

type FunctionName interface {
	Acceptor
	functionNameNode()
}

type SymbolicFunctionName struct {
	Namespace    []SymbolicName
	FunctionName SymbolicName
}

func (name *SymbolicFunctionName) Accept(visitor Visitor) error {
	if err := visitor.VisitSymbolicFunctionNameEnter(name); err != nil {
		return err
	}
	for _, n := range name.Namespace {
		if err := n.Accept(visitor); err != nil {
			return err
		}
	}
	if err := name.FunctionName.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitSymbolicFunctionNameLeave(name)
}

type ListOperatorExpr struct {
	Op      Operator
	Expr    Expr
	EndExpr Expr
}

func (expr *ListOperatorExpr) Accept(visitor Visitor) error {
	if err := visitor.VisitListOperatorExprEnter(expr); err != nil {
		return err
	}
	if err := expr.Expr.Accept(visitor); err != nil {
		return err
	}
	if err := expr.EndExpr.Accept(visitor); err != nil {
		return err
	}
	return visitor.VisitListOperatorExprEnter(expr)
}

type ExistsFunctionName struct{}

func (name *ExistsFunctionName) Accept(visitor Visitor) error {
	return visitor.VisitExistsFunctionName(name)
}

func (q *SinglePartQuery) queryNode()       {}
func (q *CreateClause) updatingClauseNode() {}
func (q *MatchClause) readingClauseNode()   {}

func (p *PatternElementPattern) patternElementNode() {}
func (p *PatternElementNested) patternElementNode()  {}

func (e *OpExpr) exprNode()                   {}
func (e *UnaryExpr) exprNode()                {}
func (e *BinaryExpr) exprNode()               {}
func (e *TernaryExpr) exprNode()              {}
func (e *ListExpr) exprNode()                 {}
func (e *PrimitiveLiteral) exprNode()         {}
func (e *ListLiteral) exprNode()              {}
func (e *MapLiteral) exprNode()               {}
func (e *PropertyLabelsExpr) exprNode()       {}
func (e *Parameter) exprNode()                {}
func (e *CaseExpr) exprNode()                 {}
func (e *ListComprehensionExpr) exprNode()    {}
func (e *FilterExpr) exprNode()               {}
func (e *QuantifierExpr) exprNode()           {}
func (e *VariableExpr) exprNode()             {}
func (e *PatternComprehensionExpr) exprNode() {}
func (e *RelationshipsPattern) exprNode()     {}
func (e *FunctionInvocation) exprNode()       {}
func (e *ListOperatorExpr) exprNode()         {}

func (s *SymbolicNameIdentifier) symbolicNameNode() {}
func (s *SymbolicNameHexLetter) symbolicNameNode()  {}

func (r *SymbolicNameSchemaName) schemaNameNode() {}
func (r *ReservedWordSchemaName) schemaNameNode() {}

func (fn *SymbolicFunctionName) functionNameNode() {}
func (fn *ExistsFunctionName) functionNameNode()   {}
