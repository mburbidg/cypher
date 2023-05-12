package parser

import (
	ast2 "github.com/mburbidg/cypher/ast"
	scanner2 "github.com/mburbidg/cypher/scanner"
	"github.com/mburbidg/cypher/utils"
	"math"
	"strings"
)

type Parser struct {
	scanner  *scanner2.Scanner
	reporter utils.Reporter
}

type Statement struct {
	// The abstract syntax tree (AST) for the parser cypher query.
	AST ast2.Query

	// The cypher query.
	Cypher string
}

func New(scanner *scanner2.Scanner, reporter utils.Reporter) *Parser {
	return &Parser{
		scanner:  scanner,
		reporter: reporter,
	}
}

func (p *Parser) Parse() (Statement, error) {
	tree, err := p.singlePartQuery()
	if err != nil {
		return Statement{}, err
	}
	return Statement{
		AST:    tree,
		Cypher: p.scanner.String(),
	}, nil
}

func (p *Parser) match(tokenTypes ...scanner2.TokenType) (scanner2.Token, bool, error) {
	pos := p.scanner.Position
	token := p.scanner.NextToken()
	for _, tokenType := range tokenTypes {
		switch token.T {
		case scanner2.Illegal:
			return scanner2.Token{}, false, p.reporter.Error(token.Line, "illegal character")
		case scanner2.EndOfInput:
			return token, false, nil
		case tokenType:
			return token, true, nil
		}
	}
	p.scanner.Position = pos
	return scanner2.Token{}, false, nil
}

func (p *Parser) matchPhrase(tokenTypes ...scanner2.TokenType) ([]scanner2.Token, bool, error) {
	pos := p.scanner.Position
	if len(tokenTypes) > 0 {
		if t := p.scanner.NextToken(); t.T == tokenTypes[0] {
			if t.T == scanner2.Illegal {
				p.scanner.Position = pos
				return nil, false, p.reporter.Error(t.Line, "illegal character")
			}
			if tokens, ok, err := p.matchPhrase(tokenTypes[1:]...); ok && err == nil {
				return append([]scanner2.Token{t}, tokens...), true, nil
			} else if err != nil {
				p.scanner.Position = pos
				return nil, false, err
			} else {
				p.scanner.Position = pos
				return nil, false, nil
			}
		} else {
			p.scanner.Position = pos
			return nil, false, nil
		}
	} else {
		return []scanner2.Token{}, true, nil
	}
}

func (p *Parser) singlePartQuery() (ast2.Query, error) {
	reading := []ast2.Query{}
	updating := []ast2.Query{}

	// Parse productions for reading which includes MATCH, UNWIND and CALL
	for {
		query, err := p.readingClause()
		if err != nil {
			return nil, err
		}
		if query == nil {
			break
		}
		reading = append(reading, query)
	}

	// Parse productions for updating which includes CREATE, MERGE, DELETE, SET and REMOVE.
	for {
		query, err := p.updatingClause()
		if err != nil {
			return nil, err
		}
		if query == nil {
			break
		}
		updating = append(updating, query)
	}

	// Parse productions for RETURN. If there were no updating queries then RETURN is required.
	projection, err := p.parseReturn()
	if err != nil {
		return nil, err
	}
	if len(updating) == 0 && projection == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting 'RETURN' following MATCH clause")
	}
	return &ast2.SinglePartQuery{ReadingClause: reading, UpdatingClause: updating, Projection: projection}, nil
}

func (p *Parser) readingClause() (ast2.Query, error) {
	return p.matchQuery()
}

func (p *Parser) updatingClause() (ast2.Query, error) {
	return p.createQuery()
}

func (p *Parser) parseReturn() (*ast2.Projection, error) {
	if _, ok, err := p.match(scanner2.Return); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	projection, err := p.projectionBody()
	if err != nil {
		return nil, err
	}
	if projection == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting projection following 'RETURN'")
	}
	return projection, nil
}

func (p *Parser) projectionBody() (*ast2.Projection, error) {
	distinct := false
	if _, ok, err := p.match(scanner2.Distinct); err != nil {
		return nil, err
	} else if ok {
		distinct = true
	}
	items, err := p.projectionItems()
	if err != nil {
		return nil, err
	}
	order, err := p.order()
	if err != nil {
		return nil, err
	}
	skip, err := p.skip()
	if err != nil {
		return nil, err
	}
	limit, err := p.limit()
	if err != nil {
		return nil, err
	}
	return &ast2.Projection{Distinct: distinct, Items: items, Order: order, Skip: skip, Limit: limit}, nil
}

func (p *Parser) projectionItems() (*ast2.ProjectionItems, error) {
	all := false
	items := []*ast2.ProjectionItem{}

	// Parse '*' or first project item.
	if _, ok, err := p.match(scanner2.Star); err != nil {
		return nil, err
	} else if ok {
		all = true
	} else {
		item, err := p.projectionItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	// One or more, comma separated projection items can follow.
	for {
		if _, ok, err := p.match(scanner2.Comma); err != nil {
			return nil, err
		} else if !ok {
			break
		}
		item, err := p.projectionItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return &ast2.ProjectionItems{All: all, Items: items}, nil
}

func (p *Parser) projectionItem() (*ast2.ProjectionItem, error) {
	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	if _, ok, err := p.match(scanner2.As); err != nil {
		return nil, err
	} else if !ok {
		return &ast2.ProjectionItem{Expr: expr}, nil
	}

	// We receive an 'AS' token, so we expect a variable to follow.
	variable, err := p.variable()
	if err != nil {
		return nil, err
	} else if variable == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting variable following 'AS'")
	}
	return &ast2.ProjectionItem{Expr: expr, Variable: variable}, nil
}

func (p *Parser) order() ([]*ast2.SortItem, error) {
	if _, ok, err := p.matchPhrase(scanner2.Order, scanner2.By); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	item, err := p.sortItem()
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting sort expression following 'ORDER BY'")
	}
	items := []*ast2.SortItem{item}
	for {
		if _, ok, err := p.match(scanner2.Comma); err != nil {
			return nil, err
		} else if !ok {
			return items, nil
		}
		item, err := p.sortItem()
		if err != nil {
			return nil, err
		}
		if item == nil {
			return nil, p.reporter.Error(p.scanner.Line(), "expecting sort expression following ','")
		}
		items = append(items, item)
	}
}

func (p *Parser) sortItem() (*ast2.SortItem, error) {
	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting sort order expression after ',' or 'ORDER BY'")
	}
	var order ast2.Order
	if t, ok, err := p.match(scanner2.Asc, scanner2.Ascending, scanner2.Desc, scanner2.Descending); err != nil {
		return nil, err
	} else if ok {
		switch t.T {
		case scanner2.Asc, scanner2.Ascending:
			order = ast2.Asc
		case scanner2.Desc, scanner2.Descending:
			order = ast2.Desc
		}
	}
	return &ast2.SortItem{Order: order, Expr: expr}, nil
}

func (p *Parser) skip() (ast2.Expr, error) {
	if _, ok, err := p.match(scanner2.Skip); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (p *Parser) limit() (ast2.Expr, error) {
	if _, ok, err := p.match(scanner2.Limit); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (p *Parser) matchQuery() (ast2.Query, error) {
	optional := false
	if _, ok, err := p.match(scanner2.Optional); err != nil {
		return nil, err
	} else if ok {
		optional = true
	}
	if _, ok, err := p.match(scanner2.Match); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	pattern, err := p.pattern()
	if err != nil {
		return nil, err
	}
	if _, ok, err := p.match(scanner2.Where); err != nil {
		return nil, err
	} else if !ok {
		return &ast2.MatchQuery{Optional: optional, Pattern: pattern}, nil
	}
	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	return &ast2.MatchQuery{Optional: optional, Pattern: pattern, WhereExpr: expr}, nil
}

func (p *Parser) createQuery() (ast2.Query, error) {
	if _, ok, err := p.match(scanner2.Create); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	pattern, err := p.pattern()
	if err != nil {
		return nil, err
	}
	return &ast2.CreateQuery{Pattern: pattern}, nil
}

func (p *Parser) pattern() (*ast2.Pattern, error) {
	part, err := p.patternPart()
	if err != nil {
		return nil, err
	}
	parts := []*ast2.PatternPart{part}
	for {
		if _, ok, err := p.match(scanner2.Comma); err != nil {
			return nil, err
		} else if !ok {
			break
		}
		part, err := p.patternPart()
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}
	return &ast2.Pattern{Parts: parts}, nil
}

func (p *Parser) patternPart() (*ast2.PatternPart, error) {
	v, err := p.variable()
	if err != nil {
		return nil, err
	}

	// Handle variable assignment pattern part
	if v != nil {
		if _, ok, err := p.match(scanner2.Equal); err != nil {
			return nil, err
		} else if !ok {
			return nil, p.reporter.Error(p.scanner.Line(), "expecting '=' following variable")
		}
	}

	// Handle anonymous part
	part, err := p.anonymousPatternPart()
	if err != nil {
		return nil, err
	}
	return &ast2.PatternPart{Variable: v, Element: part}, nil
}

func (p *Parser) anonymousPatternPart() (ast2.PatternElement, error) {
	return p.patternElement()
}

func (p *Parser) patternElement() (ast2.PatternElement, error) {
	node, err := p.nodePattern()
	if err != nil {
		return nil, err
	}

	// First handle the nested PatternElement production
	if node == nil {
		if _, ok, err := p.match(scanner2.OpenParen); err != nil {
			return nil, err
		} else if ok {
			element, err := p.patternElement()
			if err != nil {
				return nil, err
			} else if element == nil {
				return nil, p.reporter.Error(p.scanner.Line(), "expecting nested pattern element")
			}
			return &ast2.PatternElementNested{Element: element}, nil
		} else {
			return nil, p.reporter.Error(p.scanner.Line(), "expecting pattern element")
		}
	}

	// Else handle the chained PatternElement production
	chainList := []*ast2.PatternElementChain{}
	for {
		chain, err := p.patternElementChain()
		if err != nil {
			return nil, err
		} else if chain == nil {
			break
		}
		chainList = append(chainList, chain)
	}
	return &ast2.PatternElementPattern{Left: node, Chain: chainList}, nil
}

func (p *Parser) expr() (ast2.Expr, error) {
	return p.orExpr()
}

func (p *Parser) orExpr() (ast2.Expr, error) {
	expr, err := p.xorExpr()
	if err != nil {
		return nil, err
	}
	for {
		_, ok, err := p.match(scanner2.Or)
		switch {
		case err != nil:
			return nil, err
		case ok:
			right, err := p.xorExpr()
			if err != nil {
				return nil, err
			}
			expr = &ast2.BinaryExpr{expr, ast2.Or, right}
		default:
			return expr, nil
		}
	}
}

func (p *Parser) xorExpr() (ast2.Expr, error) {
	expr, err := p.andExpr()
	if err != nil {
		return nil, err
	}
	for {
		_, ok, err := p.match(scanner2.Xor)
		switch {
		case err != nil:
			return nil, err
		case ok:
			right, err := p.andExpr()
			if err != nil {
				return nil, err
			}
			expr = &ast2.BinaryExpr{expr, ast2.Xor, right}
		default:
			return expr, nil
		}
	}
}

func (p *Parser) andExpr() (ast2.Expr, error) {
	expr, err := p.notExpr()
	if err != nil {
		return nil, err
	}
	for {
		_, ok, err := p.match(scanner2.And)
		switch {
		case err != nil:
			return nil, err
		case ok:
			right, err := p.notExpr()
			if err != nil {
				return nil, err
			}
			expr = &ast2.BinaryExpr{expr, ast2.And, right}
		default:
			return expr, nil
		}
	}
}

func (p *Parser) notExpr() (ast2.Expr, error) {
	not := false
	for {
		_, ok, err := p.match(scanner2.Not)
		if err != nil {
			return nil, err
		} else if ok {
			not = !not
		} else {
			break
		}
	}
	expr, err := p.comparisonExpr()
	if err != nil {
		return nil, err
	}
	if not {
		expr = &ast2.UnaryExpr{ast2.Not, expr}
	}
	return expr, nil
}

func (p *Parser) comparisonExpr() (ast2.Expr, error) {
	tokenTypes := []scanner2.TokenType{
		scanner2.Equal,
		scanner2.NotEqual,
		scanner2.LessThan,
		scanner2.GreaterThan,
		scanner2.LessThanOrEqual,
		scanner2.GreaterThanOrEqual,
	}
	expr, err := p.addOrSubtractExpr()
	if err != nil {
		return nil, err
	}
	for {
		t, ok, err := p.match(tokenTypes...)
		switch {
		case err != nil:
			return nil, err
		case ok:
			right, err := p.addOrSubtractExpr()
			if err != nil {
				return nil, err
			}
			op, _ := opForTokens[t.T]
			expr = &ast2.BinaryExpr{expr, op, right}
		default:
			return expr, nil
		}
	}
}

func (p *Parser) addOrSubtractExpr() (ast2.Expr, error) {
	expr, err := p.multiplyDivideModuloExpr()
	if err != nil {
		return nil, err
	}
	for {
		t, ok, err := p.match(scanner2.Plus, scanner2.Dash)
		switch {
		case err != nil:
			return nil, err
		case ok:
			right, err := p.multiplyDivideModuloExpr()
			if err != nil {
				return nil, err
			}
			op, _ := opForTokens[t.T]
			expr = &ast2.BinaryExpr{expr, op, right}
		default:
			return expr, nil
		}
	}
}

func (p *Parser) multiplyDivideModuloExpr() (ast2.Expr, error) {
	tokenTypes := []scanner2.TokenType{
		scanner2.Star,
		scanner2.ForwardSlash,
		scanner2.Percent,
	}
	expr, err := p.powerExpr()
	if err != nil {
		return nil, err
	}
	for {
		t, ok, err := p.match(tokenTypes...)
		switch {
		case err != nil:
			return nil, err
		case ok:
			right, err := p.powerExpr()
			if err != nil {
				return nil, err
			}
			op, _ := opForTokens[t.T]
			expr = &ast2.BinaryExpr{expr, op, right}
		default:
			return expr, nil
		}
	}
}

func (p *Parser) powerExpr() (ast2.Expr, error) {
	expr, err := p.unaryAddOrSubtract()
	if err != nil {
		return nil, err
	}
	for {
		t, ok, err := p.match(scanner2.Caret)
		switch {
		case err != nil:
			return nil, err
		case ok:
			right, err := p.unaryAddOrSubtract()
			if err != nil {
				return nil, err
			}
			expr = &ast2.BinaryExpr{expr, ast2.Operator(t.T), right}
		default:
			return expr, nil
		}
	}
}

func (p *Parser) unaryAddOrSubtract() (ast2.Expr, error) {
	tokenTypes := []scanner2.TokenType{
		scanner2.Plus,
		scanner2.Dash,
	}
	negate := false
	for {
		t, ok, err := p.match(tokenTypes...)
		if err != nil {
			return nil, err
		} else if ok {
			if t.T == scanner2.Dash {
				negate = !negate
			}
		} else {
			break
		}
	}
	expr, err := p.stringListNullOperatorExpr()
	if err != nil {
		return nil, err
	}
	if negate {
		expr = &ast2.UnaryExpr{ast2.Negate, expr}
	}
	return expr, nil
}

func (p *Parser) stringListNullOperatorExpr() (ast2.Expr, error) {
	expr, err := p.propertyOrLabelsExpr()
	if err != nil {
		return nil, err
	}
	list := []ast2.Expr{}
	for {
		if expr, err := p.stringOpExpr(); err != nil {
			return nil, err
		} else if expr != nil {
			list = append(list, expr)
			continue
		}
		if expr, err := p.listOpExpr(); err != nil {
			return nil, err
		} else if expr != nil {
			list = append(list, expr)
			continue
		}
		if expr, err := p.isNullExpr(); err != nil {
			return nil, err
		} else if expr != nil {
			list = append(list, expr)
			continue
		}
		break
	}
	if len(list) > 0 {
		return &ast2.BinaryExpr{expr, ast2.StringOrListOp, &ast2.ListExpr{list}}, nil
	}
	return expr, nil
}

func (p *Parser) propertyOrLabelsExpr() (ast2.Expr, error) {
	atom, err := p.atom()
	if err != nil {
		return nil, err
	}
	properties := []ast2.SchemaName{}
	for {
		property, err := p.propertyLookup()
		if err != nil {
			return nil, err
		} else if property != nil {
			properties = append(properties, property)
		} else {
			break
		}
	}
	if len(properties) == 0 {
		properties = nil
	}
	labels, _ := p.NodeLabels()
	return &ast2.PropertyLabelsExpr{atom, properties, labels}, nil
}

func (p *Parser) propertyLookup() (ast2.SchemaName, error) {
	if _, ok, err := p.match(scanner2.Period); err != nil {
		return nil, err
	} else if ok {
		sn, err := p.schemaName()
		if err != nil {
			return nil, err
		}
		if sn != nil {
			return sn, nil
		}
	}
	return nil, nil
}

func (p *Parser) schemaName() (ast2.SchemaName, error) {
	pos := p.scanner.Position
	t := p.scanner.NextToken()
	if _, ok := scanner2.ReservedWordTokens[t.T]; ok {
		return &ast2.ReservedWordSchemaName{t.T}, nil
	}
	p.scanner.Position = pos
	name, err := p.symbolicName()
	if err != nil {
		return nil, err
	}
	if name != nil {
		return &ast2.SymbolicNameSchemaName{name}, nil
	}
	return nil, nil
}

func (p *Parser) NodeLabels() ([]ast2.SchemaName, error) {
	labels := []ast2.SchemaName{}
	label, err := p.NodeLabel()
	for err == nil && label != nil {
		labels = append(labels, label)
		label, err = p.NodeLabel()
	}
	if err != nil {
		return nil, err
	}
	return labels, nil
}

func (p *Parser) NodeLabel() (ast2.SchemaName, error) {
	if _, ok, err := p.match(scanner2.Colon); err != nil {
		return nil, err
	} else if ok {
		s, err := p.schemaName()
		if err != nil {
			return nil, err
		}
		if s == nil {
			return nil, p.reporter.Error(p.scanner.Line(), "expecting schema name following ':'")
		}
		return s, nil
	}
	return nil, nil
}

func (p *Parser) isNullExpr() (ast2.Expr, error) {
	if _, ok, err := p.match(scanner2.Is); err != nil {
		return nil, err
	} else if ok {
		if _, ok, err := p.match(scanner2.Null); err != nil {
			return nil, err
		} else if ok {
			return &ast2.OpExpr{ast2.IsNull}, nil
		}
		if _, ok, err := p.match(scanner2.Not); err != nil {
			return nil, err
		} else if ok {
			if _, ok, err := p.match(scanner2.Null); err != nil {
				return nil, err
			} else if ok {
				return &ast2.OpExpr{ast2.IsNotNull}, nil
			}
		}
	}
	return nil, nil
}

func (p *Parser) stringOpExpr() (ast2.Expr, error) {
	if t, ok, err := p.match(scanner2.Starts, scanner2.Ends); err != nil {
		return nil, err
	} else if ok {
		if _, ok, err := p.match(scanner2.With); err != nil {
			return nil, err
		} else if ok {
			expr, err := p.propertyOrLabelsExpr()
			if err != nil {
				return nil, err
			}
			switch t.T {
			case scanner2.Starts:
				return &ast2.UnaryExpr{ast2.StartsWith, expr}, nil
			case scanner2.Ends:
				return &ast2.UnaryExpr{ast2.EndsWith, expr}, nil
			}
		}
		p.reporter.Error(t.Line, "expecting WITH")
	}
	if _, ok, err := p.match(scanner2.Contains); err != nil {
		return nil, err
	} else if ok {
		expr, err := p.propertyOrLabelsExpr()
		if err != nil {
			return nil, err
		}
		return &ast2.UnaryExpr{ast2.Contains, expr}, nil
	}
	return nil, nil
}

func (p *Parser) listOpExpr() (ast2.Expr, error) {
	if _, ok, err := p.match(scanner2.In); err != nil {
		return nil, err
	} else if ok {
		list, err := p.propertyOrLabelsExpr()
		if err != nil {
			return nil, err
		}
		return &ast2.ListOperatorExpr{Op: ast2.InList, Expr: list}, nil
	}

	if _, ok, err := p.match(scanner2.OpenBracket); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	// First check to see if this is slice operator with no start slice expression. The p.expr() method expects
	// to parse an expression, so we need to check for the '..' terminal case first.
	if _, ok, err := p.match(scanner2.Dotdot); err != nil {
		return nil, err
	} else if ok {
		// Found '..', meaning there was no start slice expression. The end slice expression is also optional,
		// so first check for the ']' terminal case.
		if _, ok, err := p.match(scanner2.CloseBracket); err != nil {
			return nil, err
		} else if ok {
			return &ast2.ListOperatorExpr{Op: ast2.ListRange}, nil
		}
		endExpr, err := p.expr()
		if err != nil {
			return nil, err
		}
		if _, ok, err := p.match(scanner2.CloseBracket); err != nil {
			return nil, err
		} else if ok {
			return &ast2.ListOperatorExpr{Op: ast2.ListRange, EndExpr: endExpr}, nil
		}
		return nil, p.reporter.Error(p.scanner.Line(), "expecting ']' to close a list operator")
	}

	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	if _, ok, err := p.match(scanner2.CloseBracket); err != nil {
		return nil, err
	} else if ok {
		// Found ']' so this is a list index operation as opposed to a slice operation.
		return &ast2.ListOperatorExpr{Op: ast2.ListIndex, Expr: expr}, nil
	}

	if _, ok, err := p.match(scanner2.Dotdot); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting '..' between begin and end slice expressions")
	}

	// This is a slice operation. We have parsed up to the token just past the '..'. Since the end expression
	// is optional, check for the terminal case first.
	if _, ok, err := p.match(scanner2.CloseBracket); err != nil {
		return nil, err
	} else if ok {
		return &ast2.ListOperatorExpr{Op: ast2.ListRange, Expr: expr}, nil
	}
	endExpr, err := p.expr()
	if err != nil {
		return nil, err
	}
	if _, ok, err := p.match(scanner2.CloseBracket); err != nil {
		return nil, err
	} else if ok {
		return &ast2.ListOperatorExpr{Op: ast2.ListRange, Expr: expr, EndExpr: endExpr}, nil
	}
	return nil, p.reporter.Error(p.scanner.Line(), "expecting ']' to close a list operator")
}

func (p *Parser) atom() (ast2.Expr, error) {
	pos := p.scanner.Position
	if expr, _ := p.patternComprehensionExpr(); expr != nil {
		return expr, nil
	}
	p.scanner.Position = pos
	if expr, _ := p.literal(); expr != nil {
		return expr, nil
	}
	p.scanner.Position = pos
	if expr, err := p.parameter(); err != nil {
		return nil, err
	} else if expr != nil {
		return expr, nil
	}
	if expr, err := p.caseExpr(); err != nil {
		return nil, err
	} else if expr != nil {
		return expr, nil
	}
	if _, ok, err := p.matchPhrase(scanner2.Identifier, scanner2.OpenParen, scanner2.Star, scanner2.CloseParen); err != nil {
		return nil, err
	} else if ok {
		return &ast2.OpExpr{ast2.CountAll}, nil
	}
	if expr, err := p.listComprehensionExpr(); err != nil {
		return nil, err
	} else if expr != nil {
		return expr, nil
	}
	if expr, err := p.quantifierFunction(); err != nil {
		return nil, err
	} else if expr != nil {
		return expr, nil
	}
	if expr, _ := p.relationshipsPattern(); expr != nil {
		return expr, nil
	}
	p.scanner.Position = pos
	if expr, err := p.parenthesizedExpr(); err != nil {
		return nil, err
	} else if expr != nil {
		return expr, nil
	}
	if expr, err := p.functionInvocation(); err != nil {
		return nil, err
	} else if expr != nil {
		return expr, nil
	}
	p.scanner.Position = pos
	if expr, err := p.variable(); err != nil {
		return nil, err
	} else if expr != nil {
		return expr, nil
	}
	return nil, p.reporter.Error(p.scanner.Line(), "expecting atom")
}

func (p *Parser) literal() (ast2.Expr, error) {
	if t, ok, err := p.match(scanner2.DecimalInteger, scanner2.HexInteger, scanner2.OctInteger, scanner2.Double, scanner2.String, scanner2.Null, scanner2.False, scanner2.True); err != nil {
		return nil, err
	} else if ok {
		switch t.T {
		case scanner2.DecimalInteger, scanner2.HexInteger, scanner2.OctInteger:
			return &ast2.PrimitiveLiteral{scanner2.Integer, t.Literal}, nil
		case scanner2.Double, scanner2.String:
			return &ast2.PrimitiveLiteral{t.T, t.Literal}, nil
		case scanner2.False:
			return &ast2.PrimitiveLiteral{t.T, false}, nil
		case scanner2.True:
			return &ast2.PrimitiveLiteral{t.T, true}, nil
		case scanner2.Null:
			return &ast2.PrimitiveLiteral{Kind: t.T}, nil
		}
	}
	if expr, err := p.mapLiteral(); err != nil {
		return nil, err
	} else if expr != nil {
		return expr, nil
	}
	if expr, err := p.listLiteral(); err != nil {
		return nil, err
	} else if expr != nil {
		return expr, nil
	}
	return nil, nil
}

func (p *Parser) parameter() (ast2.Expr, error) {
	if _, ok, err := p.match(scanner2.DollarSign); err != nil {
		return nil, err
	} else if ok {
		s, err := p.symbolicName()
		if err != nil {
			return nil, err
		}
		if s != nil {
			return &ast2.Parameter{SymbolicName: s}, nil
		}
		if t, ok, err := p.match(scanner2.DecimalInteger); err != nil {
			return nil, err
		} else if ok {
			return &ast2.Parameter{N: &t}, nil
		} else {
			return nil, p.reporter.Error(t.Line, "expecting symbolic name or integer")
		}
	}
	return nil, nil
}

func (p *Parser) caseExpr() (ast2.Expr, error) {
	if t, ok, err := p.match(scanner2.Case); err != nil {
		return nil, err
	} else if ok {
		initExpr, err := p.expr()
		if err != nil {
			return nil, err
		}
		caseAlts := []*ast2.CaseAltNode{}
		caseAlt, err := p.caseAlt()
		if err != nil {
			return nil, err
		}
		if caseAlt == nil {
			return nil, p.reporter.Error(p.scanner.Line(), "expecting case alternative")
		}
		for caseAlt != nil {
			caseAlts = append(caseAlts, caseAlt)
			caseAlt, err = p.caseAlt()
			if err != nil {
				return nil, err
			}
		}
		if len(caseAlts) == 0 {
			return nil, p.reporter.Error(t.Line, "expecting WHEN after CASE or CASE initialization expression")
		}
		var elseExpr ast2.Expr
		if t, ok, err := p.match(scanner2.Else); err != nil {
			return nil, err
		} else if ok {
			elseExpr, err = p.expr()
			if err != nil {
				return nil, err
			}
			if elseExpr == nil {
				return nil, p.reporter.Error(t.Line, "expecting expression after CASE ELSE")
			}
		}
		if t, ok, err := p.match(scanner2.End); err != nil {
			return nil, err
		} else if ok {
			return &ast2.CaseExpr{Init: initExpr, Alternatives: caseAlts, Else: elseExpr}, nil
		} else {
			return nil, p.reporter.Error(t.Line, "expecting CASE END")
		}
	}
	return nil, nil
}

func (p *Parser) caseAlt() (*ast2.CaseAltNode, error) {
	if _, ok, err := p.match(scanner2.When); err != nil {
		return nil, err
	} else if ok {
		whenExpr, err := p.expr()
		if err != nil {
			return nil, err
		}
		if t, ok, err := p.match(scanner2.Then); err != nil {
			return nil, err
		} else if ok {
			thenExpr, err := p.expr()
			if err != nil {
				return nil, err
			}
			return &ast2.CaseAltNode{whenExpr, thenExpr}, nil
		} else {
			return nil, p.reporter.Error(t.Line, "expecting symbolic name or integer")
		}
	}
	return nil, nil
}

func (p *Parser) listComprehensionExpr() (ast2.Expr, error) {
	if _, ok, err := p.match(scanner2.OpenBracket); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	var err error
	listCompExpr := &ast2.ListComprehensionExpr{}
	listCompExpr.FilterExpr, err = p.filterExpr()
	if err != nil {
		return nil, err
	}
	if t, ok, err := p.match(scanner2.Pipe); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(t.Line, "expecting '|' in list expression")
	}
	listCompExpr.Expr, err = p.expr()
	if err != nil {
		return nil, err
	}
	if listCompExpr.Expr == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting expression after '|'")
	}
	if t, ok, err := p.match(scanner2.CloseBracket); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(t.Line, "expecting ']'")
	}
	return listCompExpr, nil
}

func (p *Parser) filterExpr() (ast2.Expr, error) {
	filterExpr := &ast2.FilterExpr{}
	var err error
	filterExpr.Variable, err = p.variable()
	if err != nil {
		return nil, err
	}
	if filterExpr.Variable == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting variable")
	}
	if t, ok, err := p.match(scanner2.In); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(t.Line, "expecting 'IN'")
	}
	filterExpr.InExpr, err = p.expr()
	if err != nil {
		return nil, err
	}
	if filterExpr.InExpr == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting 'IN' expression")
	}
	if _, ok, err := p.match(scanner2.Where); err != nil {
		return nil, err
	} else if ok {
		filterExpr.WhereExpr, err = p.expr()
		if err != nil {
			return nil, err
		}
		if filterExpr.WhereExpr == nil {
			return nil, p.reporter.Error(p.scanner.Line(), "expecting 'WHERE' expression")
		}
	}
	return filterExpr, nil
}

var quantifierNames = map[string]ast2.Operator{
	"ALL":    ast2.AllOp,
	"ANY":    ast2.AnyOp,
	"NONE":   ast2.NoneOp,
	"SINGLE": ast2.SingleOp,
}

func (p *Parser) quantifierFunction() (ast2.Expr, error) {
	pos := p.scanner.Position
	// 'ALL' is a keyword, the other three are not, so they are handled separately. The following code figures
	// out the quantifier operation being invoked.
	var op ast2.Operator
	if t, ok, err := p.match(scanner2.Identifier, scanner2.All); err != nil {
		return nil, err
	} else if ok {
		switch t.T {
		case scanner2.Identifier:
			if op, ok = quantifierNames[strings.ToUpper(t.Lexeme)]; !ok {
				p.scanner.Position = pos
				return nil, nil
			}
		case scanner2.All:
			op = ast2.AllOp
		}
	} else {
		return nil, nil
	}

	// Now parse the rest of the quantifier invocation.
	if _, ok, err := p.match(scanner2.OpenParen); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting '('")
	}
	expr, err := p.filterExpr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting filter expression")
	}
	if _, ok, err := p.match(scanner2.CloseParen); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting ')'")
	}
	return &ast2.QuantifierExpr{op, expr}, nil
}

func (p *Parser) variable() (ast2.Expr, error) {
	s, err := p.symbolicName()
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, nil
	}
	return &ast2.VariableExpr{s}, nil
}

func (p *Parser) symbolicName() (ast2.SymbolicName, error) {
	if t, ok, err := p.match(scanner2.Identifier); err != nil {
		return nil, err
	} else if ok {
		if symbolType, ok := ast2.SymbolNames[t.Lexeme]; ok {
			return &ast2.SymbolicNameIdentifier{t, symbolType}, nil
		}
		return &ast2.SymbolicNameIdentifier{t, ast2.Identifier}, nil
	}
	return nil, nil
}

func (p *Parser) patternComprehensionExpr() (ast2.Expr, error) {
	var err error
	patternExpr := &ast2.PatternComprehensionExpr{}
	if _, ok, err := p.match(scanner2.OpenBracket); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	patternExpr.Variable, err = p.variable()
	if err != nil {
		return nil, err
	}
	if patternExpr.Variable != nil {
		if t, ok, err := p.match(scanner2.Equal); err != nil {
			return nil, err
		} else if !ok {
			return nil, p.reporter.Error(t.Line, "expecting '=' following variable")
		}
	}
	patternExpr.ReltionshipsPattern, err = p.relationshipsPattern()
	if patternExpr.ReltionshipsPattern == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting relationship pattern")
	}
	if _, ok, err := p.match(scanner2.Where); err != nil {
		return nil, err
	} else if ok {
		patternExpr.WhereExpr, err = p.expr()
		if err != nil {
			return nil, err
		}
		if patternExpr.WhereExpr == nil {
			return nil, p.reporter.Error(p.scanner.Line(), "expecting expression following 'WHERE'")
		}
	}
	if t, ok, err := p.match(scanner2.Pipe); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(t.Line, "expecting '|'")
	}
	patternExpr.PipeExpr, err = p.expr()
	if err != nil {
		return nil, err
	}
	if patternExpr.PipeExpr == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting expression following '|'")
	}
	if t, ok, err := p.match(scanner2.CloseBracket); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(t.Line, "expecting ']'")
	}
	return patternExpr, nil
}

func (p *Parser) relationshipsPattern() (*ast2.RelationshipsPattern, error) {
	var err error
	rel := &ast2.RelationshipsPattern{
		Chain: []*ast2.PatternElementChain{},
	}
	rel.Left, err = p.nodePattern()
	if err != nil {
		return nil, err
	}
	if rel.Left == nil {
		return nil, nil
	}
	chain, err := p.patternElementChain()
	if err != nil {
		return nil, err
	}
	if chain == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting pattern element chain")
	}
	rel.Chain = append(rel.Chain, chain)
	for {
		chain, err := p.patternElementChain()
		if err != nil {
			return nil, err
		}
		if chain == nil {
			break
		}
		rel.Chain = append(rel.Chain, chain)
	}
	return rel, nil
}

func (p *Parser) patternElementChain() (*ast2.PatternElementChain, error) {
	var err error
	chain := &ast2.PatternElementChain{}
	chain.RelationshipPattern, err = p.relationshipPattern()
	if err != nil {
		return nil, err
	}
	if chain.RelationshipPattern == nil {
		return nil, nil
	}
	chain.Right, err = p.nodePattern()
	if err != nil {
		return nil, err
	}
	if chain.Right == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting node pattern")
	}
	return chain, nil
}

func (p *Parser) relationshipPattern() (*ast2.RelationshipPattern, error) {
	var err error
	pattern := &ast2.RelationshipPattern{Left: ast2.Undirected, Right: ast2.Undirected}
	pos := p.scanner.Position
	_, ok, err := p.match(scanner2.LessThan)
	if err != nil {
		return nil, err
	} else if ok {
		pattern.Left = ast2.Directed
	}
	if _, ok, err := p.match(scanner2.Dash); err != nil {
		return nil, err
	} else if !ok {
		p.scanner.Position = pos
		return nil, nil
	}
	pattern.RelationshipDetail, err = p.relationshipDetail()
	if err != nil {
		return nil, err
	}
	if _, ok, err := p.match(scanner2.Dash); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting '-'")
	}
	if _, ok, err := p.match(scanner2.GreaterThan); err != nil {
		return nil, err
	} else if ok {
		pattern.Right = ast2.Directed
	}
	return pattern, nil
}

func (p *Parser) relationshipDetail() (*ast2.RelationshipDetail, error) {
	var err error
	if _, ok, err := p.match(scanner2.OpenBracket); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	detail := &ast2.RelationshipDetail{}
	detail.Variable, err = p.variable()
	if err != nil {
		return nil, err
	}
	detail.RelationshipTypes, err = p.relationshipTypes()
	if err != nil {
		return nil, err
	}
	detail.RangeLiteral, err = p.rangeLiteral()
	if err != nil {
		return nil, err
	}
	detail.Properties, err = p.properties()
	if err != nil {
		return nil, err
	}
	if t, ok, err := p.match(scanner2.CloseBracket); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(t.Line, "expecting ']'")
	}
	return detail, nil
}

func (p *Parser) properties() (*ast2.Properties, error) {
	var err error
	properties := &ast2.Properties{}
	expr, err := p.mapLiteral()
	if err != nil {
		return nil, err
	}
	if expr != nil {
		properties.MapLiteral = expr.(*ast2.MapLiteral)
		return properties, nil
	}
	properties.Parameter, err = p.parameter()
	if err != nil {
		return nil, err
	}
	if properties.Parameter != nil {
		return properties, nil
	}
	return nil, nil
}

func (p *Parser) rangeLiteral() (*ast2.RangeLiteral, error) {
	literal := &ast2.RangeLiteral{Begin: math.MinInt64, End: math.MaxInt64}
	if _, ok, err := p.match(scanner2.Star); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	if t, ok, err := p.match(scanner2.DecimalInteger, scanner2.HexInteger, scanner2.OctInteger); err != nil {
		return nil, err
	} else if ok {
		literal.Begin = t.Literal.(int64)
	}
	if _, ok, err := p.match(scanner2.Dotdot); err != nil {
		return nil, err
	} else if ok {
		if t, ok, err := p.match(scanner2.DecimalInteger, scanner2.HexInteger, scanner2.OctInteger); err != nil {
			return nil, err
		} else if ok {
			literal.End = t.Literal.(int64)
		}
	}
	return literal, nil
}

func (p *Parser) relationshipTypes() ([]ast2.SchemaName, error) {
	if _, ok, err := p.match(scanner2.Colon); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	typeNames := []ast2.SchemaName{}
	s, err := p.schemaName()
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting relationship type name")
	}
	typeNames = append(typeNames, s)
	for {
		if _, ok, err := p.match(scanner2.Pipe); err != nil {
			return nil, err
		} else if !ok {
			return typeNames, nil
		}
		_, _, err := p.match(scanner2.Colon)
		if err != nil {
			return nil, err
		}
		s, err := p.schemaName()
		if err != nil {
			return nil, err
		}
		if s == nil {
			return nil, p.reporter.Error(p.scanner.Line(), "expecting relationship type name")
		}
	}
}

func (p *Parser) nodePattern() (*ast2.NodePattern, error) {
	if _, ok, err := p.match(scanner2.OpenParen); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	var err error
	np := &ast2.NodePattern{}
	np.Variable, err = p.variable()
	if err != nil {
		return nil, err
	}
	np.Labels, err = p.NodeLabels()
	if err != nil {
		return nil, err
	}
	expr, err := p.mapLiteral()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		parameter, err := p.parameter()
		if err != nil {
			return nil, err
		}
		if parameter != nil {
			np.Properties = &ast2.Properties{Parameter: parameter}
		}
	} else {
		np.Properties = &ast2.Properties{MapLiteral: expr.(*ast2.MapLiteral)}
	}
	if t, ok, err := p.match(scanner2.CloseParen); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(t.Line, "expecting ')' following node pattern")
	}
	return np, nil
}

func (p *Parser) mapLiteral() (ast2.Expr, error) {
	if _, ok, err := p.match(scanner2.OpenBrace); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	literal := &ast2.MapLiteral{PropertyKeyNames: []*ast2.PropertyKeyNames{}}
	for {
		var err error
		pkn := &ast2.PropertyKeyNames{}
		pkn.Name, err = p.schemaName()
		if err != nil {
			return nil, err
		}
		if pkn.Name == nil {
			break
		}
		if t, ok, err := p.match(scanner2.Colon); err != nil {
			return nil, err
		} else if !ok {
			return nil, p.reporter.Error(t.Line, "expecting '}' following map literal")
		}
		pkn.Expr, err = p.expr()
		if err != nil {
			return nil, err
		}
		if pkn.Expr == nil {
			return nil, p.reporter.Error(p.scanner.Line(), "expecting expression following ':'")
		}
		literal.PropertyKeyNames = append(literal.PropertyKeyNames, pkn)
		if _, ok, err := p.match(scanner2.Comma); err != nil {
			return nil, err
		} else if !ok {
			break
		}
	}
	if t, ok, err := p.match(scanner2.CloseBrace); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(t.Line, "expecting '}' following map literal")
	}
	return literal, nil
}

func (p *Parser) parenthesizedExpr() (ast2.Expr, error) {
	if _, ok, err := p.match(scanner2.OpenParen); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	if t, ok, err := p.match(scanner2.CloseParen); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(t.Line, "expecting ')' following expression")
	}
	return expr, nil
}

func (p *Parser) functionInvocation() (ast2.Expr, error) {
	fn, err := p.functionName()
	if err != nil {
		return nil, err
	}
	if _, ok, err := p.match(scanner2.OpenParen); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	_, distinct, err := p.match(scanner2.Distinct)
	if err != nil {
		return nil, err
	}

	// Check first to see if we have a close paren, which indicates there were no arguments. The method
	// p.expr() expects to find an expression, so we check the no arg case by checking for a close paren first.
	if _, ok, err := p.match(scanner2.CloseParen); err != nil {
		return nil, err
	} else if ok {
		return &ast2.FunctionInvocation{FunctionName: fn, Distinct: distinct}, nil
	}

	// The following code expects at least one argument.
	args := []ast2.Expr{}
	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	if expr != nil {
		args = append(args, expr)
		for {
			if _, ok, err := p.match(scanner2.Comma); err != nil {
				return nil, err
			} else if !ok {
				break
			}
			expr, err := p.expr()
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, p.reporter.Error(p.scanner.Line(), "expecting argument following ','")
			}
		}
	}
	if t, ok, err := p.match(scanner2.CloseParen); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(t.Line, "expecting ')' function parameters")
	}
	return &ast2.FunctionInvocation{FunctionName: fn, Distinct: distinct, Args: args}, nil
}

func (p *Parser) functionName() (ast2.FunctionName, error) {
	if _, ok, err := p.match(scanner2.Exists); err != nil {
		return nil, err
	} else if ok {
		return &ast2.ExistsFunctionName{}, nil
	}
	ns, err := p.namespace()
	if err != nil {
		return nil, err
	}
	name, err := p.symbolicName()
	if err != nil {
		return nil, err
	}
	if name == nil {
		if len(ns) == 0 {
			return nil, nil
		}
		return nil, p.reporter.Error(p.scanner.Line(), "expecting function name")
	}
	return &ast2.SymbolicFunctionName{Namespace: ns, FunctionName: name}, nil
}

func (p *Parser) namespace() ([]ast2.SymbolicName, error) {
	namespace := []ast2.SymbolicName{}
	for {
		pos := p.scanner.Position
		sn, err := p.symbolicName()
		if err != nil {
			return nil, err
		}
		if sn == nil {
			p.scanner.Position = pos
			return nil, nil
		}
		if _, ok, err := p.match(scanner2.Period); err != nil {
			return nil, err
		} else if !ok {
			p.scanner.Position = pos
			break
		}
		namespace = append(namespace, sn)
	}
	return namespace, nil
}

func (p *Parser) listLiteral() (ast2.Expr, error) {
	if _, ok, err := p.match(scanner2.OpenBracket); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	items := []ast2.Expr{}

	// We need to check first to see if the list is the empty list, since p.expr() expects to see an
	// expression. We do this by checking for the list terminal ']'.
	if _, ok, err := p.match(scanner2.CloseBracket); err != nil {
		return nil, err
	} else if ok {
		return &ast2.ListLiteral{Items: items}, nil
	}

	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	items = append(items, expr)
	for {
		if _, ok, err := p.match(scanner2.Comma); err != nil {
			return nil, err
		} else if !ok {
			break
		}
		expr, err := p.expr()
		if err != nil {
			return nil, err
		}
		if expr != nil {
			items = append(items, expr)
		}
	}
	if t, ok, err := p.match(scanner2.CloseBracket); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(t.Line, "expecting ']' to close a list")
	}
	return &ast2.ListLiteral{Items: items}, nil
}
