package parser

import (
	"github.com/mburbidg/cypher/pkg/ast"
	"github.com/mburbidg/cypher/pkg/scanner"
	"github.com/mburbidg/cypher/pkg/utils"
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

func (p *Parser) Parse() ast.Expr {
	return p.expr()
}

func (p *Parser) match(tokenTypes ...scanner.TokenType) (scanner.Token, bool) {
	token := p.scanner.NextToken()
	for _, tokenType := range tokenTypes {
		switch token.T {
		case scanner.EndOfInput:
			return token, false
		case tokenType:
			return token, true
		}
	}
	p.scanner.ReturnToken(token)
	return scanner.Token{}, false
}

func (p *Parser) matchPhrase(tokenTypes ...scanner.TokenType) bool {
	for i := 0; i < len(tokenTypes); i++ {
		if t, ok := p.match(tokenTypes[i]); !ok {
			p.reporter.Error(t.Line, "expecting phrase")
			return false
		}
	}
	return true
}

func (p *Parser) expr() ast.Expr {
	return p.orExpr()
}

func (p *Parser) orExpr() ast.Expr {
	expr := p.xorExpr()
	for _, ok := p.match(scanner.Or); ok; _, ok = p.match(scanner.Or) {
		expr = &ast.BinaryExpr{expr, ast.Or, p.xorExpr()}
	}
	return expr
}

func (p *Parser) xorExpr() ast.Expr {
	expr := p.andExpr()
	for _, ok := p.match(scanner.Xor); ok; _, ok = p.match(scanner.Xor) {
		expr = &ast.BinaryExpr{expr, ast.Xor, p.andExpr()}
	}
	return expr
}

func (p *Parser) andExpr() ast.Expr {
	expr := p.notExpr()
	for _, ok := p.match(scanner.And); ok; _, ok = p.match(scanner.And) {
		expr = &ast.BinaryExpr{expr, ast.And, p.notExpr()}
	}
	return expr
}

func (p *Parser) notExpr() ast.Expr {
	not := false
	for _, ok := p.match(scanner.Not); ok; _, ok = p.match(scanner.And) {
		not = !not
	}
	expr := p.comparisonExpr()
	if not {
		expr = &ast.UnaryExpr{ast.Not, expr}
	}
	return expr
}

func (p *Parser) comparisonExpr() ast.Expr {
	tokenTypes := []scanner.TokenType{
		scanner.Equal,
		scanner.NotEqual,
		scanner.LessThan,
		scanner.GreaterThan,
		scanner.LessThanOrEqual,
		scanner.GreaterThanOrEqual,
	}
	expr := p.addOrSubtractExpr()
	for t, ok := p.match(tokenTypes...); ok; _, ok = p.match(tokenTypes...) {
		expr = &ast.BinaryExpr{expr, ast.Operator(t.T), p.addOrSubtractExpr()}
	}
	return expr
}

func (p *Parser) addOrSubtractExpr() ast.Expr {
	tokenTypes := []scanner.TokenType{
		scanner.Plus,
		scanner.Minus,
	}
	expr := p.multiplyDivideModuloExpr()
	for t, ok := p.match(tokenTypes...); ok; _, ok = p.match(tokenTypes...) {
		expr = &ast.BinaryExpr{expr, ast.Operator(t.T), p.multiplyDivideModuloExpr()}
	}
	return expr
}

func (p *Parser) multiplyDivideModuloExpr() ast.Expr {
	tokenTypes := []scanner.TokenType{
		scanner.Star,
		scanner.ForwardSlash,
		scanner.Percent,
	}
	expr := p.powerExpr()
	for t, ok := p.match(tokenTypes...); ok; _, ok = p.match(tokenTypes...) {
		expr = &ast.BinaryExpr{expr, ast.Operator(t.T), p.powerExpr()}
	}
	return expr
}

func (p *Parser) powerExpr() ast.Expr {
	tokenTypes := []scanner.TokenType{
		scanner.Caret,
	}
	expr := p.unaryAddOrSubtract()
	for t, ok := p.match(tokenTypes...); ok; _, ok = p.match(tokenTypes...) {
		expr = &ast.BinaryExpr{expr, ast.Operator(t.T), p.unaryAddOrSubtract()}
	}
	return expr
}

func (p *Parser) unaryAddOrSubtract() ast.Expr {
	tokenTypes := []scanner.TokenType{
		scanner.Plus,
		scanner.Minus,
	}
	negate := false
	for t, ok := p.match(tokenTypes...); ok; _, ok = p.match(tokenTypes...) {
		if t.T == scanner.Minus {
			negate = !negate
		}
	}
	expr := p.stringListNullOperatorExpr()
	if negate {
		expr = &ast.UnaryExpr{ast.Negate, expr}
	}
	return expr
}

func (p *Parser) stringListNullOperatorExpr() ast.Expr {
	expr := p.propertyOrLabelsExpr()
	list := []ast.Expr{}
	for {
		if expr := p.stringOpExpr(); expr != nil {
			list = append(list, expr)
			continue
		}
		if expr := p.listOpExpr(); expr != nil {
			list = append(list, expr)
			continue
		}
		if expr := p.isNullExpr(); expr != nil {
			list = append(list, expr)
			continue
		}
		break
	}
	if len(list) > 0 {
		return &ast.BinaryExpr{expr, ast.StringOrListOp, &ast.ListExpr{list}}
	}
	return expr
}

func (p *Parser) propertyOrLabelsExpr() ast.Expr {
	atom := p.atom()
	properties := []ast.SchemaName{}
	for property := p.propertyLookup(); property != nil; property = p.propertyLookup() {
		if property != nil {
			properties = append(properties, property)
		}
	}
	labels := p.NodeLabels()
	return &ast.PropertyLabelsExpr{atom, properties, labels}
}

func (p *Parser) propertyLookup() ast.SchemaName {
	if _, ok := p.match(scanner.Period); ok {
		if t, ok := p.match(scanner.Identifier); ok {
			if symbolType, ok := ast.SymbolNames[t.Lexeme]; ok {
				return &ast.SymbolName{t, symbolType}
			} else {
				return &ast.SymbolName{t, ast.Identifier}
			}
		}
		t := p.scanner.NextToken()
		if _, ok := scanner.ReservedWordTokens[t.T]; ok {
			return &ast.ReservedWord{t}
		}
		p.scanner.ReturnToken(t)
		p.reporter.Error(t.Line, "expecting symbol name or reserved word")
	}
	return nil
}

func (p *Parser) NodeLabels() []ast.SchemaName {
	labels := []ast.SchemaName{}
	for label := p.NodeLabel(); label != nil; label = p.NodeLabel() {
		if label != nil {
			labels = append(labels, label)
		}
	}
	return labels
}

func (p *Parser) NodeLabel() ast.SchemaName {
	if _, ok := p.match(scanner.Colon); ok {
		if t, ok := p.match(scanner.Identifier); ok {
			if symbolType, ok := ast.SymbolNames[t.Lexeme]; ok {
				return &ast.SymbolName{t, symbolType}
			} else {
				return &ast.SymbolName{t, ast.Identifier}
			}
		}
		t := p.scanner.NextToken()
		if _, ok := scanner.ReservedWordTokens[t.T]; ok {
			return &ast.ReservedWord{t}
		}
		p.scanner.ReturnToken(t)
		p.reporter.Error(t.Line, "expecting symbol name or reserved word")
	}
	return nil
}

func (p *Parser) isNullExpr() ast.Expr {
	if _, ok := p.match(scanner.Is); ok {
		if _, ok := p.match(scanner.Null); ok {
			return &ast.OpExpr{ast.IsNull}
		}
		if _, ok := p.match(scanner.Not); ok {
			if _, ok := p.match(scanner.Null); ok {
				return &ast.OpExpr{ast.IsNotNull}
			}
		}
	}
	return nil
}

func (p *Parser) stringOpExpr() ast.Expr {
	if t, ok := p.match(scanner.Starts, scanner.Ends); ok {
		if _, ok := p.match(scanner.With); ok {
			expr := p.propertyOrLabelsExpr()
			switch t.T {
			case scanner.Starts:
				return &ast.UnaryExpr{ast.StartsWith, expr}
			case scanner.Ends:
				return &ast.UnaryExpr{ast.EndsWith, expr}
			}
		}
		p.reporter.Error(t.Line, "expecting WITH")
	}
	if _, ok := p.match(scanner.Contains); ok {
		expr := p.propertyOrLabelsExpr()
		return &ast.UnaryExpr{ast.Contains, expr}
	}
	return nil
}

func (p *Parser) listOpExpr() ast.Expr {
	if _, ok := p.match(scanner.In); ok {
		list := p.propertyOrLabelsExpr()
		return &ast.UnaryExpr{ast.InList, list}
	}
	if _, ok := p.match(scanner.OpenBracket); ok {
		start := p.expr()
		if t, ok := p.match(scanner.CloseBracket, scanner.Dotdot); ok {
			switch t.T {
			case scanner.CloseBracket:
				if start == nil {
					p.reporter.Error(t.Line, "expecting index expression")
					return nil
				}
				return &ast.UnaryExpr{ast.ListIndex, start}
			case scanner.Dotdot:
				end := p.expr()
				return &ast.BinaryExpr{start, ast.ListRange, end}
			}
		}
	}
	return nil
}

func (p *Parser) atom() ast.Expr {
	if expr := p.literal(); expr != nil {
		return expr
	}
	if expr := p.parameter(); expr != nil {
		return expr
	}
	if expr, _ := p.caseExpr(); expr != nil {
		return expr
	}
	return nil
}

func (p *Parser) literal() ast.Expr {
	if t, ok := p.match(scanner.DecimalInteger, scanner.HexInteger, scanner.OctInteger, scanner.Double, scanner.String, scanner.Null, scanner.False, scanner.True); ok {
		switch t.T {
		case scanner.DecimalInteger, scanner.HexInteger, scanner.OctInteger:
			return &ast.Literal{scanner.Integer, t.Literal}
		case scanner.Double, scanner.String:
			return &ast.Literal{t.T, t.Literal}
		case scanner.False:
			return &ast.Literal{t.T, false}
		case scanner.True:
			return &ast.Literal{t.T, true}
		case scanner.Null:
			return &ast.Literal{Kind: t.T}
		}
	}
	return nil
}

func (p *Parser) parameter() ast.Expr {
	if _, ok := p.match(scanner.DollarSign); ok {
		if t, ok := p.match(scanner.Identifier); ok {
			if symbolType, ok := ast.SymbolNames[t.Lexeme]; ok {
				return &ast.Parameter{SymbolName: &ast.SymbolName{t, symbolType}}
			} else {
				return &ast.Parameter{SymbolName: &ast.SymbolName{t, ast.Identifier}}
			}
		}
		if t, ok := p.match(scanner.DecimalInteger); ok {
			return &ast.Parameter{N: &t}
		} else {
			p.reporter.Error(t.Line, "expecting symbolic name or integer")
		}
	}
	return nil
}

func (p *Parser) caseExpr() (ast.Expr, error) {
	if _, ok := p.match(scanner.Case); ok {
		caseAltExpr, err := p.caseAltExpr()
		for err == nil && caseAltExpr != nil {

		}
	}
	return nil, nil
}

func (p *Parser) caseAltExpr() (ast.Expr, error) {
	if _, ok := p.match(scanner.When); ok {
		whenExpr := p.expr()
		if t, ok := p.match(scanner.Then); ok {
			thenExpr := p.expr()
			return &ast.CaseAltExpr{whenExpr, thenExpr}, nil
		} else {
			return nil, p.reporter.Error(t.Line, "expecting symbolic name or integer")
		}
	}
	return nil, nil
}

func (p *Parser) filterExpr() ast.Expr {
	return nil
}

func (p *Parser) functionInvocation() ast.Expr {
	return nil
}

func (p *Parser) variable() ast.Expr {
	return nil
}

func (p *Parser) patternComprehension() ast.Expr {
	return nil
}
