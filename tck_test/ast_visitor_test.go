package tck_test

import (
	"github.com/mburbidg/cypher"
	"github.com/mburbidg/cypher/ast"
	"log"
	"strings"
)

type astVisitor struct {
	inPattern    bool
	inMatch      bool
	inCreate     bool
	inNode       bool
	inRel        bool
	inExpr       bool
	creatingNode bool
	hasProps     bool
	hasLabels    bool
	symbolTable  map[string]any
}

func newASTVisitor() ast.Visitor {
	return &astVisitor{symbolTable: map[string]any{}}
}

func (visitor *astVisitor) VisitSinglePartQueryEnter(query *ast.SinglePartQuery) error {
	log.Printf("Enter SinglePartQuery\n")
	return nil
}

func (visitor *astVisitor) VisitSinglePartQueryLeave(query *ast.SinglePartQuery) error {
	log.Printf("Leave SinglePartQuery\n")
	return nil
}

func (visitor *astVisitor) VisitReadingClauseEnter(clause []ast.ReadingClause) error {
	log.Printf("Enter ReadingClause\n")
	return nil
}

func (visitor *astVisitor) VisitReadingClauseLeave(clause []ast.ReadingClause) error {
	log.Printf("Leave ReadingClause\n")
	return nil
}

func (visitor *astVisitor) VisitUpdatingClauseEnter(clause []ast.UpdatingClause) error {
	log.Printf("Enter UpdatingClause\n")
	return nil
}

func (visitor *astVisitor) VisitUpdatingClauseLeave(clause []ast.UpdatingClause) error {
	log.Printf("Leave UpdatingClause\n")
	return nil
}

func (visitor *astVisitor) VisitCreateEnter(clause *ast.CreateClause) error {
	log.Printf("Enter Create\n")
	visitor.inCreate = true
	return nil
}

func (visitor *astVisitor) VisitCreateLeave(clause *ast.CreateClause) error {
	log.Printf("Leave Create\n")
	visitor.inCreate = false
	return nil
}

func (visitor *astVisitor) VisitMatchEnter(clause *ast.MatchClause) error {
	log.Printf("Enter Match\n")
	visitor.inMatch = true
	return nil
}

func (visitor *astVisitor) VisitMatchLeave(clause *ast.MatchClause) error {
	log.Printf("Leave Match\n")
	visitor.inMatch = false
	return nil
}

func (visitor *astVisitor) VisitPatternEnter(pattern *ast.Pattern) error {
	log.Printf("Enter Pattern\n")
	visitor.inPattern = true
	return nil
}

func (visitor *astVisitor) VisitPatternLeave(pattern *ast.Pattern) error {
	log.Printf("Leave Pattern\n")
	visitor.inPattern = false
	return nil
}

func (visitor *astVisitor) VisitPatternPartEnter(part *ast.PatternPart) error {
	log.Printf("Enter PatternPart\n")
	return nil
}

func (visitor *astVisitor) VisitPatternPartLeave(part *ast.PatternPart) error {
	log.Printf("Leave PatternPart\n")
	return nil
}

func (visitor *astVisitor) VisitPatternElementNestedEnter(part *ast.PatternElementNested) error {
	log.Printf("Enter PatternElementNested\n")
	return nil
}

func (visitor *astVisitor) VisitPatternElementNestedLeave(part *ast.PatternElementNested) error {
	log.Printf("Leave PatternElementNested\n")
	return nil
}

func (visitor *astVisitor) VisitPatternElementPatternEnter(part *ast.PatternElementPattern) error {
	log.Printf("Enter PatternElementPattern\n")
	if visitor.inCreate && len(part.Chain) == 0 {
		visitor.creatingNode = true
	}
	return nil
}

func (visitor *astVisitor) VisitPatternElementPatternLeave(part *ast.PatternElementPattern) error {
	log.Printf("Leave PatternElementPattern\n")
	visitor.creatingNode = false
	return nil
}

func (visitor *astVisitor) VisitProjectionEnter(projection *ast.Projection) error {
	log.Printf("Enter Projection\n")
	return nil
}

func (visitor *astVisitor) VisitProjectionLeave(projection *ast.Projection) error {
	log.Printf("Leave Projection\n")
	return nil
}

func (visitor *astVisitor) VisitSortOrderEnter(order *ast.SortOrder) error {
	log.Printf("Enter SortOrder\n")
	return nil
}

func (visitor *astVisitor) VisitSortOrderLeave(order *ast.SortOrder) error {
	log.Printf("Leave SortOrder\n")
	return nil
}

func (visitor *astVisitor) VisitProjectionItemsEnter(items *ast.ProjectionItems) error {
	log.Printf("Enter ProjectionItems\n")
	return nil
}

func (visitor *astVisitor) VisitProjectionItemsLeave(items *ast.ProjectionItems) error {
	log.Printf("Leave ProjectionItems\n")
	return nil
}

func (visitor *astVisitor) VisitProjectionItemEnter(item *ast.ProjectionItem) error {
	log.Printf("Enter ProjectionItem\n")
	return nil
}

func (visitor *astVisitor) VisitProjectionItemLeave(item *ast.ProjectionItem) error {
	log.Printf("Leave ProjectionItem\n")
	return nil
}

func (visitor *astVisitor) VisitSortItemEnter(item *ast.SortItem) error {
	log.Printf("Enter SortItem\n")
	return nil
}

func (visitor *astVisitor) VisitSortItemLeave(item *ast.SortItem) error {
	log.Printf("Leave SortItem\n")
	return nil
}

func (visitor *astVisitor) VisitOpExpr(expr *ast.OpExpr) error {
	log.Printf("OpExpr\n")
	return nil
}

func (visitor *astVisitor) VisitUnaryExprEnter(expr *ast.UnaryExpr) error {
	log.Printf("Enter UnaryExprEnter\n")
	return nil
}

func (visitor *astVisitor) VisitUnaryExprLeave(expr *ast.UnaryExpr) error {
	log.Printf("Leave UnaryExprEnter\n")
	return nil
}

func (visitor *astVisitor) VisitBinaryExprEnter(expr *ast.BinaryExpr) error {
	log.Printf("Enter BinaryExpr\n")
	return nil
}

func (visitor *astVisitor) VisitBinaryExprLeave(expr *ast.BinaryExpr) error {
	log.Printf("Leave BinaryExpr\n")
	return nil
}

func (visitor *astVisitor) VisitTernaryExprEnter(expr *ast.TernaryExpr) error {
	log.Printf("Enter TernaryExpr\n")
	return nil
}

func (visitor *astVisitor) VisitTernaryExprLeave(expr *ast.TernaryExpr) error {
	log.Printf("Leave TernaryExpr\n")
	return nil
}

func (visitor *astVisitor) VisitListExprEnter(expr *ast.ListExpr) error {
	log.Printf("Enter ListExpr\n")
	return nil
}

func (visitor *astVisitor) VisitListExprLeave(expr *ast.ListExpr) error {
	log.Printf("Leave ListExpr\n")
	return nil
}

func (visitor *astVisitor) VisitListComprehensionExprEnter(expr *ast.ListComprehensionExpr) error {
	log.Printf("Enter ListComprehensionExpr\n")
	return nil
}

func (visitor *astVisitor) VisitListComprehensionExprLeave(expr *ast.ListComprehensionExpr) error {
	log.Printf("Leave ListComprehensionExpr\n")
	return nil
}

func (visitor *astVisitor) VisitPropertyLabelsExprEnter(expr *ast.PropertyLabelsExpr) error {
	log.Printf("Enter PropertyLabelsExpr\n")
	return nil
}

func (visitor *astVisitor) VisitPropertyLabelsExprLeave(expr *ast.PropertyLabelsExpr) error {
	log.Printf("Leave PropertyLabelsExpr\n")
	return nil
}

func (visitor *astVisitor) VisitSymbolicNameSchemaNameEnter(name *ast.SymbolicNameSchemaName) error {
	log.Printf("Enter SymbolicNameSchemaName\n")
	return nil
}

func (visitor *astVisitor) VisitSymbolicNameSchemaNameLeave(name *ast.SymbolicNameSchemaName) error {
	log.Printf("Leave SymbolicNameSchemaName\n")
	return nil
}

func (visitor *astVisitor) VisitReservedWordSchemaName(name *ast.ReservedWordSchemaName) error {
	log.Printf("ReservedWordSchemaName\n")
	return nil
}

func (visitor *astVisitor) VisitSymbolicNameIdentifier(name *ast.SymbolicNameIdentifier) error {
	log.Printf("SymbolicNameIdentifier\n")
	return visitor.visitSymbolicName(name.Identifier.Lexeme)
}

func (visitor *astVisitor) VisitSymbolicNameHexLetter(name *ast.SymbolicNameHexLetter) error {
	log.Printf("SymbolicNameHexLetter\n")
	builder := strings.Builder{}
	builder.WriteRune(name.Letter)
	return visitor.visitSymbolicName(builder.String())
}

func (visitor *astVisitor) visitSymbolicName(id string) error {
	switch {
	case visitor.inExpr:
		if _, ok := visitor.symbolTable[id]; !ok {
			return cypher.NewUndefinedVariableErr(id)
		}
	case visitor.inRel:
		if _, ok := visitor.symbolTable[id]; ok {
			return cypher.NewVariableAlreadyBoundErr(id)
		}
		visitor.symbolTable[id] = true
	case visitor.inNode:
		if visitor.creatingNode {
			if _, ok := visitor.symbolTable[id]; ok {
				return cypher.NewVariableAlreadyBoundErr(id)
			}
			visitor.symbolTable[id] = true
		} else if visitor.inCreate {
			if visitor.hasProps || visitor.hasLabels {
				if _, ok := visitor.symbolTable[id]; ok {
					return cypher.NewVariableAlreadyBoundErr(id)
				}
			}
		}
		visitor.getOrBind(id)
	}
	return nil
}

func (visitor *astVisitor) getOrBind(id string) {
	if _, ok := visitor.symbolTable[id]; ok {
		return
	}
	visitor.symbolTable[id] = true
}

func (visitor *astVisitor) VisitReservedWord(word *ast.ReservedWord) error {
	log.Printf("ReservedWord\n")
	return nil
}

func (visitor *astVisitor) VisitLabel(label *ast.Label) error {
	log.Printf("Label\n")
	return nil
}

func (visitor *astVisitor) VisitPrimitiveLiteral(literal *ast.PrimitiveLiteral) error {
	log.Printf("PrimitiveLiteral\n")
	return nil
}

func (visitor *astVisitor) VisitListLiteralEnter(literal *ast.ListLiteral) error {
	log.Printf("Enter ListLiteral\n")
	return nil
}

func (visitor *astVisitor) VisitListLiteralLeave(literal *ast.ListLiteral) error {
	log.Printf("Leave ListLiteral\n")
	return nil
}

func (visitor *astVisitor) VisitParameterEnter(param *ast.Parameter) error {
	log.Printf("Enter Parameter\n")
	return nil
}

func (visitor *astVisitor) VisitParameterLeave(param *ast.Parameter) error {
	log.Printf("Leave Parameter\n")
	return nil
}

func (visitor *astVisitor) VisitCaseExprEnter(expr *ast.CaseExpr) error {
	log.Printf("Enter CaseExpr\n")
	return nil
}

func (visitor *astVisitor) VisitCaseExprLeave(expr *ast.CaseExpr) error {
	log.Printf("Leave SinglePartQuery\n")
	return nil
}

func (visitor *astVisitor) VisitCaseAltNodeEnter(alt *ast.CaseAltNode) error {
	log.Printf("Enter CaseAltNode\n")
	return nil
}

func (visitor *astVisitor) VisitCaseAltNodeLeave(alt *ast.CaseAltNode) error {
	log.Printf("Leave CaseAltNode\n")
	return nil
}

func (visitor *astVisitor) VisitQuantifierExprEnter(quantifier *ast.QuantifierExpr) error {
	log.Printf("Enter QuantifierExpr\n")
	return nil
}

func (visitor *astVisitor) VisitQuantifierExprLeave(quantifier *ast.QuantifierExpr) error {
	log.Printf("Leave QuantifierExpr\n")
	return nil
}

func (visitor *astVisitor) VisitFilterExprEnter(filter *ast.FilterExpr) error {
	log.Printf("Enter FilterExpr\n")
	return nil
}

func (visitor *astVisitor) VisitFilterExprLeave(filter *ast.FilterExpr) error {
	log.Printf("Leave FilterExpr\n")
	return nil
}

func (visitor *astVisitor) VisitVariableExprEnter(expr *ast.VariableExpr) error {
	log.Printf("Enter VariableExpr\n")
	visitor.inExpr = true
	return nil
}

func (visitor *astVisitor) VisitVariableExprLeave(expr *ast.VariableExpr) error {
	log.Printf("Leave VariableExpr\n")
	visitor.inExpr = false
	return nil
}

func (visitor *astVisitor) VisitPatternComprehensionExprEnter(expr *ast.PatternComprehensionExpr) error {
	log.Printf("Enter PatternComprehensionExpr\n")
	return nil
}

func (visitor *astVisitor) VisitPatternComprehensionExprLeave(expr *ast.PatternComprehensionExpr) error {
	log.Printf("Leave PatternComprehensionExpr\n")
	return nil
}

func (visitor *astVisitor) VisitNodePatternEnter(pattern *ast.NodePattern) error {
	log.Printf("Enter NodePattern\n")
	visitor.inNode = true
	if len(pattern.Labels) > 0 {
		visitor.hasLabels = true
	}
	return nil
}

func (visitor *astVisitor) VisitNodePatternLeave(pattern *ast.NodePattern) error {
	log.Printf("Leave NodePattern\n")
	visitor.inNode = false
	visitor.hasLabels = false
	return nil
}

func (visitor *astVisitor) VisitMapLiteralEnter(literal *ast.MapLiteral) error {
	log.Printf("Enter MapLiteral\n")
	return nil
}

func (visitor *astVisitor) VisitMapLiteralLeave(literal *ast.MapLiteral) error {
	log.Printf("Leave MapLiteral\n")
	return nil
}

func (visitor *astVisitor) VisitPropertyKeyNameEnter(name *ast.PropertyKeyName) error {
	log.Printf("Enter PropertyKeyName\n")
	return nil
}

func (visitor *astVisitor) VisitPropertyKeyNameLeave(name *ast.PropertyKeyName) error {
	log.Printf("Leave PropertyKeyName\n")
	return nil
}

func (visitor *astVisitor) VisitPropertiesEnter(props *ast.Properties) error {
	log.Printf("Enter Properties\n")
	if len(props.MapLiteral.PropertyKeyNames) > 0 {
		visitor.hasProps = true
	}
	return nil
}

func (visitor *astVisitor) VisitPropertiesLeave(props *ast.Properties) error {
	log.Printf("Leave Properties\n")
	visitor.hasProps = false
	return nil
}

func (visitor *astVisitor) VisitRelationshipsPatternEnter(pattern *ast.RelationshipsPattern) error {
	log.Printf("Enter RelationshipsPattern\n")
	return nil
}

func (visitor *astVisitor) VisitRelationshipsPatternLeave(pattern *ast.RelationshipsPattern) error {
	log.Printf("Leave RelationshipsPattern\n")
	return nil
}

func (visitor *astVisitor) VisitPatternElementChainEnter(chain *ast.PatternElementChain) error {
	log.Printf("Enter PatternElementChain\n")
	return nil
}

func (visitor *astVisitor) VisitPatternElementChainLeave(chain *ast.PatternElementChain) error {
	log.Printf("Leave PatternElementChain\n")
	return nil
}

func (visitor *astVisitor) VisitRelationshipPatternEnter(pattern *ast.RelationshipPattern) error {
	log.Printf("Enter RelationshipPattern\n")
	visitor.inRel = true
	if visitor.inCreate {
		if pattern.Right == ast.Undirected && pattern.Left == ast.Undirected {
			return cypher.NewRequiresDirectedRelationship()
		}
		if pattern.Right == ast.Directed && pattern.Left == ast.Directed {
			return cypher.NewRequiresDirectedRelationship()
		}
		if pattern.RelationshipDetail == nil {
			return cypher.NewNoSingleRelationshipType()
		}
	}
	return nil
}

func (visitor *astVisitor) VisitRelationshipPatternLeave(pattern *ast.RelationshipPattern) error {
	log.Printf("Leave RelationshipPattern\n")
	visitor.inRel = false
	return nil
}

func (visitor *astVisitor) VisitRelationshipDetailEnter(detail *ast.RelationshipDetail) error {
	log.Printf("Enter RelationshipDetail\n")
	visitor.inPattern = true
	if visitor.inCreate {
		if detail.RangeLiteral != nil {
			return cypher.NewCreatingVarLength()
		}
	}
	return nil
}

func (visitor *astVisitor) VisitRelationshipDetailLeave(detail *ast.RelationshipDetail) error {
	log.Printf("Leave RelationshipDetail\n")
	if visitor.inCreate {
		if len(detail.RelationshipTypes) != 1 {
			return cypher.NewNoSingleRelationshipType()
		}
	}
	visitor.inPattern = false
	return nil
}

func (visitor *astVisitor) VisitRangeLiteral(literal *ast.RangeLiteral) error {
	log.Printf("RangeLiteral\n")
	return nil
}

func (visitor *astVisitor) VisitFunctionInvocationEnter(fnc *ast.FunctionInvocation) error {
	log.Printf("Enter FunctionInvocation\n")
	return nil
}

func (visitor *astVisitor) VisitFunctionInvocationLeave(fnc *ast.FunctionInvocation) error {
	log.Printf("Leave FunctionInvocation\n")
	return nil
}

func (visitor *astVisitor) VisitSymbolicFunctionNameEnter(name *ast.SymbolicFunctionName) error {
	log.Printf("Enter SymbolicFunctionName\n")
	return nil
}

func (visitor *astVisitor) VisitSymbolicFunctionNameLeave(name *ast.SymbolicFunctionName) error {
	log.Printf("Leave SymbolicFunctionName\n")
	return nil
}

func (visitor *astVisitor) VisitListOperatorExprEnter(expr *ast.ListOperatorExpr) error {
	log.Printf("Enter ListOperatorExpr\n")
	return nil
}

func (visitor *astVisitor) VisitListOperatorExprLeave(expr *ast.ListOperatorExpr) error {
	log.Printf("Leave ListOperatorExpr\n")
	return nil
}

func (visitor *astVisitor) VisitExistsFunctionName(name *ast.ExistsFunctionName) error {
	log.Printf("ExistsFunctionName\n")
	return nil
}
