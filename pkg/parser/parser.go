package parser

import (
	"github.com/mburbidg/cypher/pkg/ast"
	"github.com/mburbidg/cypher/pkg/scanner"
	"github.com/mburbidg/cypher/pkg/utils"
	"math"
	"strings"
)

type Parser struct {
	scanner  *scanner.Scanner
	reporter utils.Reporter
}

func New(scanner *scanner.Scanner, reporter utils.Reporter) *Parser {
	return &Parser{
		scanner:  scanner,
		reporter: reporter,
	}
}

func (p *Parser) Parse() (ast.Query, error) {
	return p.singlePartQuery()
}

func (p *Parser) match(tokenTypes ...scanner.TokenType) (scanner.Token, bool, error) {
	pos := p.scanner.Position
	token := p.scanner.NextToken()
	for _, tokenType := range tokenTypes {
		switch token.T {
		case scanner.Illegal:
			return scanner.Token{}, false, p.reporter.Error(token.Line, "illegal character")
		case scanner.EndOfInput:
			return token, false, nil
		case tokenType:
			return token, true, nil
		}
	}
	p.scanner.Position = pos
	return scanner.Token{}, false, nil
}

func (p *Parser) matchPhrase(tokenTypes ...scanner.TokenType) ([]scanner.Token, bool, error) {
	pos := p.scanner.Position
	if len(tokenTypes) > 0 {
		if t := p.scanner.NextToken(); t.T == tokenTypes[0] {
			if t.T == scanner.Illegal {
				p.scanner.Position = pos
				return nil, false, p.reporter.Error(t.Line, "illegal character")
			}
			if tokens, ok, err := p.matchPhrase(tokenTypes[1:]...); ok && err == nil {
				return append([]scanner.Token{t}, tokens...), true, nil
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
		return []scanner.Token{}, true, nil
	}
}

func (p *Parser) singlePartQuery() (ast.Query, error) {
	reading := []ast.Query{}
	updating := []ast.Query{}

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
	return &ast.SinglePartQuery{ReadingClause: reading, UpdatingClause: updating, Projection: projection}, nil
}

func (p *Parser) readingClause() (ast.Query, error) {
	return p.matchQuery()
}

func (p *Parser) updatingClause() (ast.Query, error) {
	return p.createQuery()
}

func (p *Parser) parseReturn() (*ast.Projection, error) {
	if _, ok, err := p.match(scanner.Return); err != nil {
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

func (p *Parser) projectionBody() (*ast.Projection, error) {
	distinct := false
	if _, ok, err := p.match(scanner.Distinct); err != nil {
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
	return &ast.Projection{Distinct: distinct, Items: items, Order: order, Skip: skip, Limit: limit}, nil
}

func (p *Parser) projectionItems() (*ast.ProjectionItems, error) {
	all := false
	items := []*ast.ProjectionItem{}

	// Parse '*' or first project item.
	if _, ok, err := p.match(scanner.Star); err != nil {
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
		if _, ok, err := p.match(scanner.Comma); err != nil {
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
	return &ast.ProjectionItems{All: all, Items: items}, nil
}

func (p *Parser) projectionItem() (*ast.ProjectionItem, error) {
	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	if _, ok, err := p.match(scanner.As); err != nil {
		return nil, err
	} else if !ok {
		return &ast.ProjectionItem{Expr: expr}, nil
	}

	// We receive an 'AS' token, so we expect a variable to follow.
	variable, err := p.variable()
	if err != nil {
		return nil, err
	} else if variable == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting variable following 'AS'")
	}
	return &ast.ProjectionItem{Expr: expr, Variable: variable}, nil
}

func (p *Parser) order() ([]*ast.SortItem, error) {
	if _, ok, err := p.matchPhrase(scanner.Order, scanner.By); err != nil {
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
	items := []*ast.SortItem{item}
	for {
		if _, ok, err := p.match(scanner.Comma); err != nil {
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

func (p *Parser) sortItem() (*ast.SortItem, error) {
	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	if expr == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting sort order expression after ',' or 'ORDER BY'")
	}
	var order ast.Order
	if t, ok, err := p.match(scanner.Asc, scanner.Ascending, scanner.Desc, scanner.Descending); err != nil {
		return nil, err
	} else if ok {
		switch t.T {
		case scanner.Asc, scanner.Ascending:
			order = ast.Asc
		case scanner.Desc, scanner.Descending:
			order = ast.Desc
		}
	}
	return &ast.SortItem{Order: order, Expr: expr}, nil
}

func (p *Parser) skip() (ast.Expr, error) {
	if _, ok, err := p.match(scanner.Skip); err != nil {
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

func (p *Parser) limit() (ast.Expr, error) {
	if _, ok, err := p.match(scanner.Limit); err != nil {
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

func (p *Parser) matchQuery() (ast.Query, error) {
	optional := false
	if _, ok, err := p.match(scanner.Optional); err != nil {
		return nil, err
	} else if ok {
		optional = true
	}
	if _, ok, err := p.match(scanner.Match); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	pattern, err := p.pattern()
	if err != nil {
		return nil, err
	}
	if _, ok, err := p.match(scanner.Where); err != nil {
		return nil, err
	} else if !ok {
		return &ast.MatchQuery{Optional: optional, Pattern: pattern}, nil
	}
	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	return &ast.MatchQuery{Optional: optional, Pattern: pattern, WhereExpr: expr}, nil
}

func (p *Parser) createQuery() (ast.Query, error) {
	if _, ok, err := p.match(scanner.Create); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	pattern, err := p.pattern()
	if err != nil {
		return nil, err
	}
	return &ast.CreateQuery{Pattern: pattern}, nil
}

func (p *Parser) pattern() (*ast.Pattern, error) {
	part, err := p.patternPart()
	if err != nil {
		return nil, err
	}
	parts := []*ast.PatternPart{part}
	for {
		if _, ok, err := p.match(scanner.Comma); err != nil {
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
	return &ast.Pattern{Parts: parts}, nil
}

func (p *Parser) patternPart() (*ast.PatternPart, error) {
	v, err := p.variable()
	if err != nil {
		return nil, err
	}

	// Handle variable assignment pattern part
	if v != nil {
		if _, ok, err := p.match(scanner.Equal); err != nil {
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
	return &ast.PatternPart{Variable: v, Element: part}, nil
}

func (p *Parser) anonymousPatternPart() (ast.PatternElement, error) {
	return p.patternElement()
}

func (p *Parser) patternElement() (ast.PatternElement, error) {
	node, err := p.nodePattern()
	if err != nil {
		return nil, err
	}

	// First handle the nested PatternElement production
	if node == nil {
		if _, ok, err := p.match(scanner.OpenParen); err != nil {
			return nil, err
		} else if ok {
			element, err := p.patternElement()
			if err != nil {
				return nil, err
			} else if element == nil {
				return nil, p.reporter.Error(p.scanner.Line(), "expecting nested pattern element")
			}
			return &ast.PatternElementNested{Element: element}, nil
		} else {
			return nil, p.reporter.Error(p.scanner.Line(), "expecting pattern element")
		}
	}

	// Else handle the chained PatternElement production
	chainList := []*ast.PatternElementChain{}
	for {
		chain, err := p.patternElementChain()
		if err != nil {
			return nil, err
		} else if chain == nil {
			break
		}
		chainList = append(chainList, chain)
	}
	return &ast.PatternElementPattern{Left: node, Chain: chainList}, nil
}

func (p *Parser) expr() (ast.Expr, error) {
	return p.orExpr()
}

func (p *Parser) orExpr() (ast.Expr, error) {
	expr, err := p.xorExpr()
	if err != nil {
		return nil, err
	}
	for {
		_, ok, err := p.match(scanner.Or)
		switch {
		case err != nil:
			return nil, err
		case ok:
			right, err := p.xorExpr()
			if err != nil {
				return nil, err
			}
			expr = &ast.BinaryExpr{expr, ast.Or, right}
		default:
			return expr, nil
		}
	}
}

func (p *Parser) xorExpr() (ast.Expr, error) {
	expr, err := p.andExpr()
	if err != nil {
		return nil, err
	}
	for {
		_, ok, err := p.match(scanner.Xor)
		switch {
		case err != nil:
			return nil, err
		case ok:
			right, err := p.andExpr()
			if err != nil {
				return nil, err
			}
			expr = &ast.BinaryExpr{expr, ast.Xor, right}
		default:
			return expr, nil
		}
	}
}

func (p *Parser) andExpr() (ast.Expr, error) {
	expr, err := p.notExpr()
	if err != nil {
		return nil, err
	}
	for {
		_, ok, err := p.match(scanner.And)
		switch {
		case err != nil:
			return nil, err
		case ok:
			right, err := p.notExpr()
			if err != nil {
				return nil, err
			}
			expr = &ast.BinaryExpr{expr, ast.And, right}
		default:
			return expr, nil
		}
	}
}

func (p *Parser) notExpr() (ast.Expr, error) {
	not := false
	for {
		_, ok, err := p.match(scanner.Not)
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
		expr = &ast.UnaryExpr{ast.Not, expr}
	}
	return expr, nil
}

func (p *Parser) comparisonExpr() (ast.Expr, error) {
	tokenTypes := []scanner.TokenType{
		scanner.Equal,
		scanner.NotEqual,
		scanner.LessThan,
		scanner.GreaterThan,
		scanner.LessThanOrEqual,
		scanner.GreaterThanOrEqual,
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
			expr = &ast.BinaryExpr{expr, op, right}
		default:
			return expr, nil
		}
	}
}

func (p *Parser) addOrSubtractExpr() (ast.Expr, error) {
	expr, err := p.multiplyDivideModuloExpr()
	if err != nil {
		return nil, err
	}
	for {
		t, ok, err := p.match(scanner.Plus, scanner.Dash)
		switch {
		case err != nil:
			return nil, err
		case ok:
			right, err := p.multiplyDivideModuloExpr()
			if err != nil {
				return nil, err
			}
			op, _ := opForTokens[t.T]
			expr = &ast.BinaryExpr{expr, op, right}
		default:
			return expr, nil
		}
	}
}

func (p *Parser) multiplyDivideModuloExpr() (ast.Expr, error) {
	tokenTypes := []scanner.TokenType{
		scanner.Star,
		scanner.ForwardSlash,
		scanner.Percent,
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
			expr = &ast.BinaryExpr{expr, op, right}
		default:
			return expr, nil
		}
	}
}

func (p *Parser) powerExpr() (ast.Expr, error) {
	expr, err := p.unaryAddOrSubtract()
	if err != nil {
		return nil, err
	}
	for {
		t, ok, err := p.match(scanner.Caret)
		switch {
		case err != nil:
			return nil, err
		case ok:
			right, err := p.unaryAddOrSubtract()
			if err != nil {
				return nil, err
			}
			expr = &ast.BinaryExpr{expr, ast.Operator(t.T), right}
		default:
			return expr, nil
		}
	}
}

func (p *Parser) unaryAddOrSubtract() (ast.Expr, error) {
	tokenTypes := []scanner.TokenType{
		scanner.Plus,
		scanner.Dash,
	}
	negate := false
	for {
		t, ok, err := p.match(tokenTypes...)
		if err != nil {
			return nil, err
		} else if ok {
			if t.T == scanner.Dash {
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
		expr = &ast.UnaryExpr{ast.Negate, expr}
	}
	return expr, nil
}

func (p *Parser) stringListNullOperatorExpr() (ast.Expr, error) {
	expr, err := p.propertyOrLabelsExpr()
	if err != nil {
		return nil, err
	}
	list := []ast.Expr{}
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
		return &ast.BinaryExpr{expr, ast.StringOrListOp, &ast.ListExpr{list}}, nil
	}
	return expr, nil
}

func (p *Parser) propertyOrLabelsExpr() (ast.Expr, error) {
	atom, err := p.atom()
	if err != nil {
		return nil, err
	}
	properties := []ast.SchemaName{}
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
	return &ast.PropertyLabelsExpr{atom, properties, labels}, nil
}

func (p *Parser) propertyLookup() (ast.SchemaName, error) {
	if _, ok, err := p.match(scanner.Period); err != nil {
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

func (p *Parser) schemaName() (ast.SchemaName, error) {
	pos := p.scanner.Position
	t := p.scanner.NextToken()
	if _, ok := scanner.ReservedWordTokens[t.T]; ok {
		return &ast.ReservedWordSchemaName{t.T}, nil
	}
	p.scanner.Position = pos
	name, err := p.symbolicName()
	if err != nil {
		return nil, err
	}
	if name != nil {
		return &ast.SymbolicNameSchemaName{name}, nil
	}
	return nil, nil
}

func (p *Parser) NodeLabels() ([]ast.SchemaName, error) {
	labels := []ast.SchemaName{}
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

func (p *Parser) NodeLabel() (ast.SchemaName, error) {
	if _, ok, err := p.match(scanner.Colon); err != nil {
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

func (p *Parser) isNullExpr() (ast.Expr, error) {
	if _, ok, err := p.match(scanner.Is); err != nil {
		return nil, err
	} else if ok {
		if _, ok, err := p.match(scanner.Null); err != nil {
			return nil, err
		} else if ok {
			return &ast.OpExpr{ast.IsNull}, nil
		}
		if _, ok, err := p.match(scanner.Not); err != nil {
			return nil, err
		} else if ok {
			if _, ok, err := p.match(scanner.Null); err != nil {
				return nil, err
			} else if ok {
				return &ast.OpExpr{ast.IsNotNull}, nil
			}
		}
	}
	return nil, nil
}

func (p *Parser) stringOpExpr() (ast.Expr, error) {
	if t, ok, err := p.match(scanner.Starts, scanner.Ends); err != nil {
		return nil, err
	} else if ok {
		if _, ok, err := p.match(scanner.With); err != nil {
			return nil, err
		} else if ok {
			expr, err := p.propertyOrLabelsExpr()
			if err != nil {
				return nil, err
			}
			switch t.T {
			case scanner.Starts:
				return &ast.UnaryExpr{ast.StartsWith, expr}, nil
			case scanner.Ends:
				return &ast.UnaryExpr{ast.EndsWith, expr}, nil
			}
		}
		p.reporter.Error(t.Line, "expecting WITH")
	}
	if _, ok, err := p.match(scanner.Contains); err != nil {
		return nil, err
	} else if ok {
		expr, err := p.propertyOrLabelsExpr()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{ast.Contains, expr}, nil
	}
	return nil, nil
}

func (p *Parser) listOpExpr() (ast.Expr, error) {
	if _, ok, err := p.match(scanner.In); err != nil {
		return nil, err
	} else if ok {
		list, err := p.propertyOrLabelsExpr()
		if err != nil {
			return nil, err
		}
		return &ast.ListOperatorExpr{Op: ast.InList, Expr: list}, nil
	}

	if _, ok, err := p.match(scanner.OpenBracket); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	// First check to see if this is slice operator with no start slice expression. The p.expr() method expects
	// to parse an expression, so we need to check for the '..' terminal case first.
	if _, ok, err := p.match(scanner.Dotdot); err != nil {
		return nil, err
	} else if ok {
		// Found '..', meaning there was no start slice expression. The end slice expression is also optional,
		// so first check for the ']' terminal case.
		if _, ok, err := p.match(scanner.CloseBracket); err != nil {
			return nil, err
		} else if ok {
			return &ast.ListOperatorExpr{Op: ast.ListRange}, nil
		}
		endExpr, err := p.expr()
		if err != nil {
			return nil, err
		}
		if _, ok, err := p.match(scanner.CloseBracket); err != nil {
			return nil, err
		} else if ok {
			return &ast.ListOperatorExpr{Op: ast.ListRange, EndExpr: endExpr}, nil
		}
		return nil, p.reporter.Error(p.scanner.Line(), "expecting ']' to close a list operator")
	}

	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	if _, ok, err := p.match(scanner.CloseBracket); err != nil {
		return nil, err
	} else if ok {
		// Found ']' so this is a list index operation as opposed to a slice operation.
		return &ast.ListOperatorExpr{Op: ast.ListIndex, Expr: expr}, nil
	}

	if _, ok, err := p.match(scanner.Dotdot); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting '..' between begin and end slice expressions")
	}

	// This is a slice operation. We have parsed up to the token just past the '..'. Since the end expression
	// is optional, check for the terminal case first.
	if _, ok, err := p.match(scanner.CloseBracket); err != nil {
		return nil, err
	} else if ok {
		return &ast.ListOperatorExpr{Op: ast.ListRange, Expr: expr}, nil
	}
	endExpr, err := p.expr()
	if err != nil {
		return nil, err
	}
	if _, ok, err := p.match(scanner.CloseBracket); err != nil {
		return nil, err
	} else if ok {
		return &ast.ListOperatorExpr{Op: ast.ListRange, Expr: expr, EndExpr: endExpr}, nil
	}
	return nil, p.reporter.Error(p.scanner.Line(), "expecting ']' to close a list operator")
}

func (p *Parser) atom() (ast.Expr, error) {
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
	if _, ok, err := p.matchPhrase(scanner.Identifier, scanner.OpenParen, scanner.Star, scanner.CloseParen); err != nil {
		return nil, err
	} else if ok {
		return &ast.OpExpr{ast.CountAll}, nil
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

func (p *Parser) literal() (ast.Expr, error) {
	if t, ok, err := p.match(scanner.DecimalInteger, scanner.HexInteger, scanner.OctInteger, scanner.Double, scanner.String, scanner.Null, scanner.False, scanner.True); err != nil {
		return nil, err
	} else if ok {
		switch t.T {
		case scanner.DecimalInteger, scanner.HexInteger, scanner.OctInteger:
			return &ast.PrimitiveLiteral{scanner.Integer, t.Literal}, nil
		case scanner.Double, scanner.String:
			return &ast.PrimitiveLiteral{t.T, t.Literal}, nil
		case scanner.False:
			return &ast.PrimitiveLiteral{t.T, false}, nil
		case scanner.True:
			return &ast.PrimitiveLiteral{t.T, true}, nil
		case scanner.Null:
			return &ast.PrimitiveLiteral{Kind: t.T}, nil
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

func (p *Parser) parameter() (ast.Expr, error) {
	if _, ok, err := p.match(scanner.DollarSign); err != nil {
		return nil, err
	} else if ok {
		s, err := p.symbolicName()
		if err != nil {
			return nil, err
		}
		if s != nil {
			return &ast.Parameter{SymbolicName: s}, nil
		}
		if t, ok, err := p.match(scanner.DecimalInteger); err != nil {
			return nil, err
		} else if ok {
			return &ast.Parameter{N: &t}, nil
		} else {
			return nil, p.reporter.Error(t.Line, "expecting symbolic name or integer")
		}
	}
	return nil, nil
}

func (p *Parser) caseExpr() (ast.Expr, error) {
	if t, ok, err := p.match(scanner.Case); err != nil {
		return nil, err
	} else if ok {
		initExpr, err := p.expr()
		if err != nil {
			return nil, err
		}
		caseAlts := []*ast.CaseAltNode{}
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
		var elseExpr ast.Expr
		if t, ok, err := p.match(scanner.Else); err != nil {
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
		if t, ok, err := p.match(scanner.End); err != nil {
			return nil, err
		} else if ok {
			return &ast.CaseExpr{Init: initExpr, Alternatives: caseAlts, Else: elseExpr}, nil
		} else {
			return nil, p.reporter.Error(t.Line, "expecting CASE END")
		}
	}
	return nil, nil
}

func (p *Parser) caseAlt() (*ast.CaseAltNode, error) {
	if _, ok, err := p.match(scanner.When); err != nil {
		return nil, err
	} else if ok {
		whenExpr, err := p.expr()
		if err != nil {
			return nil, err
		}
		if t, ok, err := p.match(scanner.Then); err != nil {
			return nil, err
		} else if ok {
			thenExpr, err := p.expr()
			if err != nil {
				return nil, err
			}
			return &ast.CaseAltNode{whenExpr, thenExpr}, nil
		} else {
			return nil, p.reporter.Error(t.Line, "expecting symbolic name or integer")
		}
	}
	return nil, nil
}

func (p *Parser) listComprehensionExpr() (ast.Expr, error) {
	if _, ok, err := p.match(scanner.OpenBracket); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	var err error
	listCompExpr := &ast.ListComprehensionExpr{}
	listCompExpr.FilterExpr, err = p.filterExpr()
	if err != nil {
		return nil, err
	}
	if t, ok, err := p.match(scanner.Pipe); err != nil {
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
	if t, ok, err := p.match(scanner.CloseBracket); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(t.Line, "expecting ']'")
	}
	return listCompExpr, nil
}

func (p *Parser) filterExpr() (ast.Expr, error) {
	filterExpr := &ast.FilterExpr{}
	var err error
	filterExpr.Variable, err = p.variable()
	if err != nil {
		return nil, err
	}
	if filterExpr.Variable == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting variable")
	}
	if t, ok, err := p.match(scanner.In); err != nil {
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
	if _, ok, err := p.match(scanner.Where); err != nil {
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

var quantifierNames = map[string]ast.Operator{
	"ALL":    ast.AllOp,
	"ANY":    ast.AnyOp,
	"NONE":   ast.NoneOp,
	"SINGLE": ast.SingleOp,
}

func (p *Parser) quantifierFunction() (ast.Expr, error) {
	pos := p.scanner.Position
	// 'ALL' is a keyword, the other three are not, so they are handled separately. The following code figures
	// out the quantifier operation being invoked.
	var op ast.Operator
	if t, ok, err := p.match(scanner.Identifier, scanner.All); err != nil {
		return nil, err
	} else if ok {
		switch t.T {
		case scanner.Identifier:
			if op, ok = quantifierNames[strings.ToUpper(t.Lexeme)]; !ok {
				p.scanner.Position = pos
				return nil, nil
			}
		case scanner.All:
			op = ast.AllOp
		}
	} else {
		return nil, nil
	}

	// Now parse the rest of the quantifier invocation.
	if _, ok, err := p.match(scanner.OpenParen); err != nil {
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
	if _, ok, err := p.match(scanner.CloseParen); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting ')'")
	}
	return &ast.QuantifierExpr{op, expr}, nil
}

func (p *Parser) variable() (ast.Expr, error) {
	s, err := p.symbolicName()
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, nil
	}
	return &ast.VariableExpr{s}, nil
}

func (p *Parser) symbolicName() (ast.SymbolicName, error) {
	if t, ok, err := p.match(scanner.Identifier); err != nil {
		return nil, err
	} else if ok {
		if symbolType, ok := ast.SymbolNames[t.Lexeme]; ok {
			return &ast.SymbolicNameIdentifier{t, symbolType}, nil
		}
		return &ast.SymbolicNameIdentifier{t, ast.Identifier}, nil
	}
	return nil, nil
}

func (p *Parser) patternComprehensionExpr() (ast.Expr, error) {
	var err error
	patternExpr := &ast.PatternComprehensionExpr{}
	if _, ok, err := p.match(scanner.OpenBracket); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	patternExpr.Variable, err = p.variable()
	if err != nil {
		return nil, err
	}
	if patternExpr.Variable != nil {
		if t, ok, err := p.match(scanner.Equal); err != nil {
			return nil, err
		} else if !ok {
			return nil, p.reporter.Error(t.Line, "expecting '=' following variable")
		}
	}
	patternExpr.ReltionshipsPattern, err = p.relationshipsPattern()
	if patternExpr.ReltionshipsPattern == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting relationship pattern")
	}
	if _, ok, err := p.match(scanner.Where); err != nil {
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
	if t, ok, err := p.match(scanner.Pipe); err != nil {
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
	if t, ok, err := p.match(scanner.CloseBracket); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(t.Line, "expecting ']'")
	}
	return patternExpr, nil
}

func (p *Parser) relationshipsPattern() (*ast.RelationshipsPattern, error) {
	var err error
	rel := &ast.RelationshipsPattern{
		Chain: []*ast.PatternElementChain{},
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

func (p *Parser) patternElementChain() (*ast.PatternElementChain, error) {
	var err error
	chain := &ast.PatternElementChain{}
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

func (p *Parser) relationshipPattern() (*ast.RelationshipPattern, error) {
	var err error
	pattern := &ast.RelationshipPattern{Left: ast.Undirected, Right: ast.Undirected}
	pos := p.scanner.Position
	_, ok, err := p.match(scanner.LessThan)
	if err != nil {
		return nil, err
	} else if ok {
		pattern.Left = ast.Directed
	}
	if _, ok, err := p.match(scanner.Dash); err != nil {
		return nil, err
	} else if !ok {
		p.scanner.Position = pos
		return nil, nil
	}
	pattern.RelationshipDetail, err = p.relationshipDetail()
	if err != nil {
		return nil, err
	}
	if _, ok, err := p.match(scanner.Dash); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting '-'")
	}
	if _, ok, err := p.match(scanner.GreaterThan); err != nil {
		return nil, err
	} else if ok {
		pattern.Right = ast.Directed
	}
	return pattern, nil
}

func (p *Parser) relationshipDetail() (*ast.RelationshipDetail, error) {
	var err error
	if _, ok, err := p.match(scanner.OpenBracket); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	detail := &ast.RelationshipDetail{}
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
	if t, ok, err := p.match(scanner.CloseBracket); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(t.Line, "expecting ']'")
	}
	return detail, nil
}

func (p *Parser) properties() (*ast.Properties, error) {
	var err error
	properties := &ast.Properties{}
	expr, err := p.mapLiteral()
	if err != nil {
		return nil, err
	}
	if expr != nil {
		properties.MapLiteral = expr.(*ast.MapLiteral)
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

func (p *Parser) rangeLiteral() (*ast.RangeLiteral, error) {
	literal := &ast.RangeLiteral{Begin: math.MinInt64, End: math.MaxInt64}
	if _, ok, err := p.match(scanner.Star); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	if t, ok, err := p.match(scanner.DecimalInteger, scanner.HexInteger, scanner.OctInteger); err != nil {
		return nil, err
	} else if ok {
		literal.Begin = t.Literal.(int64)
	}
	if _, ok, err := p.match(scanner.Dotdot); err != nil {
		return nil, err
	} else if ok {
		if t, ok, err := p.match(scanner.DecimalInteger, scanner.HexInteger, scanner.OctInteger); err != nil {
			return nil, err
		} else if ok {
			literal.End = t.Literal.(int64)
		}
	}
	return literal, nil
}

func (p *Parser) relationshipTypes() ([]ast.SchemaName, error) {
	if _, ok, err := p.match(scanner.Colon); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	typeNames := []ast.SchemaName{}
	s, err := p.schemaName()
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting relationship type name")
	}
	typeNames = append(typeNames, s)
	for {
		if _, ok, err := p.match(scanner.Pipe); err != nil {
			return nil, err
		} else if !ok {
			return typeNames, nil
		}
		_, _, err := p.match(scanner.Colon)
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

func (p *Parser) nodePattern() (*ast.NodePattern, error) {
	if _, ok, err := p.match(scanner.OpenParen); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	var err error
	np := &ast.NodePattern{}
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
			np.Properties = &ast.Properties{Parameter: parameter}
		}
	} else {
		np.Properties = &ast.Properties{MapLiteral: expr.(*ast.MapLiteral)}
	}
	if t, ok, err := p.match(scanner.CloseParen); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(t.Line, "expecting ')' following node pattern")
	}
	return np, nil
}

func (p *Parser) mapLiteral() (ast.Expr, error) {
	if _, ok, err := p.match(scanner.OpenBrace); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	literal := &ast.MapLiteral{PropertyKeyNames: []*ast.PropertyKeyNames{}}
	for {
		var err error
		pkn := &ast.PropertyKeyNames{}
		pkn.Name, err = p.schemaName()
		if err != nil {
			return nil, err
		}
		if pkn.Name == nil {
			break
		}
		if t, ok, err := p.match(scanner.Colon); err != nil {
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
		if _, ok, err := p.match(scanner.Comma); err != nil {
			return nil, err
		} else if !ok {
			break
		}
	}
	if t, ok, err := p.match(scanner.CloseBrace); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(t.Line, "expecting '}' following map literal")
	}
	return literal, nil
}

func (p *Parser) parenthesizedExpr() (ast.Expr, error) {
	if _, ok, err := p.match(scanner.OpenParen); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	if t, ok, err := p.match(scanner.CloseParen); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(t.Line, "expecting ')' following expression")
	}
	return expr, nil
}

func (p *Parser) functionInvocation() (ast.Expr, error) {
	fn, err := p.functionName()
	if err != nil {
		return nil, err
	}
	if _, ok, err := p.match(scanner.OpenParen); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	_, distinct, err := p.match(scanner.Distinct)
	if err != nil {
		return nil, err
	}

	// Check first to see if we have a close paren, which indicates there were no arguments. The method
	// p.expr() expects to find an expression, so we check the no arg case by checking for a close paren first.
	if _, ok, err := p.match(scanner.CloseParen); err != nil {
		return nil, err
	} else if ok {
		return &ast.FunctionInvocation{FunctionName: fn, Distinct: distinct}, nil
	}

	// The following code expects at least one argument.
	args := []ast.Expr{}
	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	if expr != nil {
		args = append(args, expr)
		for {
			if _, ok, err := p.match(scanner.Comma); err != nil {
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
	if t, ok, err := p.match(scanner.CloseParen); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(t.Line, "expecting ')' function parameters")
	}
	return &ast.FunctionInvocation{FunctionName: fn, Distinct: distinct, Args: args}, nil
}

func (p *Parser) functionName() (ast.FunctionName, error) {
	if _, ok, err := p.match(scanner.Exists); err != nil {
		return nil, err
	} else if ok {
		return &ast.ExistsFunctionName{}, nil
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
	return &ast.SymbolicFunctionName{Namespace: ns, FunctionName: name}, nil
}

func (p *Parser) namespace() ([]ast.SymbolicName, error) {
	namespace := []ast.SymbolicName{}
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
		if _, ok, err := p.match(scanner.Period); err != nil {
			return nil, err
		} else if !ok {
			p.scanner.Position = pos
			break
		}
		namespace = append(namespace, sn)
	}
	return namespace, nil
}

func (p *Parser) listLiteral() (ast.Expr, error) {
	if _, ok, err := p.match(scanner.OpenBracket); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	items := []ast.Expr{}

	// We need to check first to see if the list is the empty list, since p.expr() expects to see an
	// expression. We do this by checking for the list terminal ']'.
	if _, ok, err := p.match(scanner.CloseBracket); err != nil {
		return nil, err
	} else if ok {
		return &ast.ListLiteral{Items: items}, nil
	}

	expr, err := p.expr()
	if err != nil {
		return nil, err
	}
	items = append(items, expr)
	for {
		if _, ok, err := p.match(scanner.Comma); err != nil {
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
	if t, ok, err := p.match(scanner.CloseBracket); err != nil {
		return nil, err
	} else if !ok {
		return nil, p.reporter.Error(t.Line, "expecting ']' to close a list")
	}
	return &ast.ListLiteral{Items: items}, nil
}
