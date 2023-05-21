package tck_test

import (
	"fmt"
	"github.com/mburbidg/cypher"
	"github.com/mburbidg/cypher/ast"
	"github.com/mburbidg/cypher/parser"
	"log"
	"strings"
)

const (
	none int = iota
	readingClause
	updatingClause
	matchClause
	createClause
	mergeClause
)

type scope struct {
	clauseType   int
	clauseDetail int
	creatingNode bool
	symbolTable  map[string]any
	prev         *scope
}

func newScope(prev *scope) *scope {
	return &scope{
		clauseType:   none,
		clauseDetail: none,
		creatingNode: false,
		symbolTable:  map[string]any{},
		prev:         prev,
	}
}

func (s *scope) pushScope() *scope {
	return newScope(s)
}

func (s *scope) popScope() (*scope, bool) {
	if s.prev == nil {
		return nil, false
	}
	return s.prev, true
}

func (s *scope) getClauseType() int {
	if s.clauseType != none {
		return s.clauseType
	}
	if s.prev != nil {
		return s.prev.getClauseType()
	}
	return none
}

func (s *scope) getClauseDetail() int {
	if s.clauseDetail != none {
		return s.clauseDetail
	}
	if s.prev != nil {
		return s.prev.getClauseDetail()
	}
	return none
}

type astRuntime struct {
}

func (runtime *astRuntime) eval(stmt parser.Statement) error {
	v := newASTVisitor()
	err := stmt.AST.Accept(v)
	return err
	//switch tree := stmt.AST.(type) {
	//case *ast.SinglePartQuery:
	//	return runtime.evalSinglePartQuery(newScope(nil), tree)
	//default:
	//	return fmt.Errorf("query not implemented")
	//}
	//return nil
}

func (runtime *astRuntime) evalSinglePartQuery(scope *scope, tree *ast.SinglePartQuery) error {
	log.Printf("evaluating single part query")
	err := runtime.evalReadingClause(scope, tree.ReadingClause)
	if err != nil {
		return err
	}
	err = runtime.evalUpdatingClause(scope, tree.UpdatingClause)
	if err != nil {
		return err
	}
	err = runtime.evalProjection(scope, tree.Projection)
	if err != nil {
		return err
	}
	return nil
}

func (runtime *astRuntime) evalReadingClause(scope *scope, tree []ast.ReadingClause) error {
	log.Printf("evaluating reading clause")
	scope.clauseType = readingClause
	defer func() { scope.clauseType = none }()
	for _, clause := range tree {
		switch c := clause.(type) {
		case *ast.MatchClause:
			err := runtime.evalMatchClause(scope, c)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (runtime *astRuntime) evalUpdatingClause(scope *scope, tree []ast.UpdatingClause) error {
	log.Printf("evaluating updating clause")
	scope.clauseType = updatingClause
	defer func() { scope.clauseType = none }()
	for _, clause := range tree {
		switch c := clause.(type) {
		case *ast.CreateClause:
			err := runtime.evalCreateClause(scope, c)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (runtime *astRuntime) evalMatchClause(scope *scope, tree *ast.MatchClause) error {
	log.Printf("evaluating match clause")
	scope.clauseDetail = matchClause
	defer func() { scope.clauseDetail = none }()
	for _, p := range tree.Pattern.Parts {
		err := runtime.evalPatternPart(scope, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (runtime *astRuntime) evalCreateClause(scope *scope, tree *ast.CreateClause) error {
	log.Printf("evaluating create clause")
	scope.clauseDetail = createClause
	defer func() { scope.clauseDetail = none }()
	for _, p := range tree.Pattern.Parts {
		err := runtime.evalPatternPart(scope, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (runtime *astRuntime) evalPatternPart(scope *scope, part *ast.PatternPart) error {
	log.Printf("evaluating pattern part")
	if err := runtime.bindVariable(scope, part.Variable); err != nil {
		return err
	}
	return runtime.evalPatternElement(scope, part.Element)
}

func (runtime *astRuntime) evalPatternElement(scope *scope, elem ast.PatternElement) error {
	switch elem := elem.(type) {
	case *ast.PatternElementPattern:
		if (scope.clauseDetail == createClause || scope.clauseDetail == mergeClause) && len(elem.Chain) == 0 {
			scope.creatingNode = true
		}
		defer func() { scope.creatingNode = false }()
		err := runtime.evalPatternElementPattern(scope, elem)
		if err != nil {
			return err
		}
	case *ast.PatternElementNested:
		return runtime.evalPatternElement(scope, elem.Element)
	default:
		return fmt.Errorf("pattern element type not implemented")
	}
	return nil
}

func (runtime *astRuntime) evalPatternElementPattern(scope *scope, elem *ast.PatternElementPattern) error {
	log.Printf("evaluating pattern element pattern")
	if err := runtime.evalNodePattern(scope, elem.Left); err != nil {
		return err
	}
	for _, c := range elem.Chain {
		if err := runtime.evalPatternElementChain(scope, c); err != nil {
			return err
		}
	}
	return nil
}

func (runtime *astRuntime) evalPatternElementChain(scope *scope, chain *ast.PatternElementChain) error {
	if err := runtime.evalRelationshipPattern(scope, chain.RelationshipPattern); err != nil {
		return err
	}
	if err := runtime.evalNodePattern(scope, chain.Right); err != nil {
		return err
	}
	return nil
}

func (runtime *astRuntime) evalNodePattern(scope *scope, pattern *ast.NodePattern) error {
	if pattern != nil {
		if err := runtime.bindVariable(scope, pattern.Variable); err != nil {
			return err
		}
		//if scope.creatingNode {
		//	if err := runtime.bindVariable(scope, pattern.Variable); err != nil {
		//		return err
		//	}
		//} else {
		//	runtime.getOrBindVariable(scope, pattern.Variable)
		//}
		//if err := runtime.evalProperties(scope, pattern.Properties); err != nil {
		//	return err
		//}
	}
	return nil
}

func (runtime *astRuntime) evalProperties(scope *scope, properties *ast.Properties) error {
	if properties != nil {
		if properties.Parameter != nil {
			if scope.getClauseDetail() == matchClause {
				return cypher.NewInvalidParameterUse()
			}
			if err := runtime.evalExpr(scope, properties.Parameter); err != nil {
				return err
			}
		}
		if properties.MapLiteral != nil {
			if err := runtime.evalExpr(scope, properties.MapLiteral); err != nil {
				return err
			}
		}
	}
	return nil
}

func (runtime *astRuntime) evalMapLiteral(scope *scope, literal *ast.MapLiteral) error {
	for _, p := range literal.PropertyKeyNames {
		if err := runtime.evalExpr(scope, p.Expr); err != nil {
			return err
		}
	}
	return nil
}

func (runtime *astRuntime) evalRelationshipPattern(scope *scope, pattern *ast.RelationshipPattern) error {
	if scope.getClauseType() == updatingClause {
		if pattern.Right == ast.Undirected && pattern.Left == ast.Undirected {
			return cypher.NewRequiresDirectedRelationship()
		}
		if pattern.Right == ast.Directed && pattern.Left == ast.Directed {
			return cypher.NewRequiresDirectedRelationship()
		}
	}
	if err := runtime.evalRelationshipDetail(scope, pattern.RelationshipDetail); err != nil {
		return err
	}
	return nil
}

func (runtime *astRuntime) evalRelationshipDetail(scope *scope, detail *ast.RelationshipDetail) error {
	if detail == nil {
		return cypher.NewNoSingleRelationshipType()
	}
	if len(detail.RelationshipTypes) > 1 {
		return cypher.NewNoSingleRelationshipType()
	}
	if detail.RangeLiteral != nil {
		return cypher.NewCreatingVarLength()
	}
	if err := runtime.bindVariable(scope, detail.Variable); err != nil {
		return err
	}
	return nil
}

func (runtime *astRuntime) evalProjection(scope *scope, tree *ast.Projection) error {
	log.Printf("evaluating reading clause")
	return nil
}

func (runtime *astRuntime) evalProjectionItems(scope *scope, items *ast.ProjectionItems) error {
	return nil
}

func (runtime *astRuntime) evalProjectionItem(scope *scope, clauseType int, item *ast.ProjectionItem) error {
	if err := runtime.bindVariable(scope, item.Variable); err != nil {
		return err
	}
	if err := runtime.evalExpr(scope, item.Expr); err != nil {
		return err
	}
	return nil
}

func (runtime *astRuntime) evalExpr(scope *scope, expr ast.Expr) error {
	switch e := expr.(type) {
	case *ast.OpExpr:
	case *ast.UnaryExpr:
		return runtime.evalExpr(scope, e.Expr)
	case *ast.BinaryExpr:
		if err := runtime.evalExpr(scope, e.Left); err != nil {
			return err
		}
		if err := runtime.evalExpr(scope, e.Right); err != nil {
			return err
		}
	case *ast.TernaryExpr:
		if err := runtime.evalExpr(scope, e.E1); err != nil {
			return err
		}
		if err := runtime.evalExpr(scope, e.E2); err != nil {
			return err
		}
		if err := runtime.evalExpr(scope, e.E3); err != nil {
			return err
		}
	case *ast.ListExpr:
		for _, e := range e.List {
			if err := runtime.evalExpr(scope, e); err != nil {
				return err
			}
		}
	case *ast.PrimitiveLiteral:
	case *ast.ListLiteral:
		for _, e := range e.Items {
			if err := runtime.evalExpr(scope, e); err != nil {
				return err
			}
		}
	case *ast.MapLiteral:
		for _, l := range e.PropertyKeyNames {
			if err := runtime.evalExpr(scope, l.Expr); err != nil {
				return err
			}
		}
	case *ast.PropertyLabelsExpr:
		if err := runtime.evalExpr(scope, e.Atom); err != nil {
			return err
		}
	case *ast.Parameter:
	case *ast.CaseExpr:
		if err := runtime.evalExpr(scope, e.Init); err != nil {
			return err
		}
		if err := runtime.evalExpr(scope, e.Else); err != nil {
			return err
		}
		for _, alt := range e.Alternatives {
			if err := runtime.evalExpr(scope, alt.Then); err != nil {
				return err
			}
			if err := runtime.evalExpr(scope, alt.When); err != nil {
				return err
			}
		}
	case *ast.ListComprehensionExpr:
		if err := runtime.evalExpr(scope, e.Expr); err != nil {
			return err
		}
		if err := runtime.evalExpr(scope, e.FilterExpr); err != nil {
			return err
		}
	case *ast.FilterExpr:
		if err := runtime.bindVariable(scope, e.Variable); err != nil {
			return err
		}
		if err := runtime.evalExpr(scope, e.InExpr); err != nil {
			return err
		}
		if err := runtime.evalExpr(scope, e.WhereExpr); err != nil {
			return err
		}
	case *ast.QuantifierExpr:
		if err := runtime.evalExpr(scope, e.Expr); err != nil {
			return err
		}
	case *ast.VariableExpr:
		return runtime.expectVariable(scope, e.SymbolicName)
	case *ast.PatternComprehensionExpr:
		if err := runtime.evalExpr(scope, e.WhereExpr); err != nil {
			return err
		}
		if err := runtime.evalExpr(scope, e.PipeExpr); err != nil {
			return err
		}
		if err := runtime.evalExpr(scope, e.ReltionshipsPattern); err != nil {
			return err
		}
	case *ast.RelationshipsPattern:
		if err := runtime.evalNodePattern(scope, e.Left); err != nil {
			return err
		}
		for _, chain := range e.Chain {
			if err := runtime.evalPatternElementChain(scope, chain); err != nil {
				return err
			}
		}
	case *ast.FunctionInvocation:
		for _, arg := range e.Args {
			if err := runtime.evalExpr(scope, arg); err != nil {
				return err
			}
		}
	case *ast.ListOperatorExpr:
		if err := runtime.evalExpr(scope, e.Expr); err != nil {
			return err
		}
		if err := runtime.evalExpr(scope, e.EndExpr); err != nil {
			return err
		}
	default:
		return fmt.Errorf("expression not implemented")
	}
	return nil
}

func (runtime *astRuntime) getIdentifier(symbolicName ast.SymbolicName) string {
	switch id := symbolicName.(type) {
	case *ast.SymbolicNameIdentifier:
		return id.Identifier.Lexeme
	case *ast.SymbolicNameHexLetter:
		builder := strings.Builder{}
		builder.WriteRune(id.Letter)
		return builder.String()
	}
	return ""
}

func (runtime *astRuntime) bindVariable(scope *scope, symbolicName ast.SymbolicName) error {
	if symbolicName != nil {
		id := runtime.getIdentifier(symbolicName)
		if _, ok := scope.symbolTable[id]; ok {
			return cypher.NewVariableAlreadyBoundErr(id)
		}
		scope.symbolTable[id] = true
	}
	return nil
}

func (runtime *astRuntime) getOrBindVariable(scope *scope, symbolicName ast.SymbolicName) bool {
	id := runtime.getIdentifier(symbolicName)
	if _, ok := scope.symbolTable[id]; ok {
		return true
	}
	_ = runtime.bindVariable(scope, symbolicName)
	return false
}

func (runtime *astRuntime) expectVariable(scope *scope, symbolicName ast.SymbolicName) error {
	id := runtime.getIdentifier(symbolicName)
	if _, ok := scope.symbolTable[id]; !ok {
		return cypher.NewUndefinedVariableErr(id)
	}
	return nil
}
