package parser

import (
	"github.com/mburbidg/cypher/pkg/ast"
	"github.com/mburbidg/cypher/pkg/scanner"
	"github.com/mburbidg/cypher/pkg/utils"
	"math"
	"unicode/utf8"
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

// Todo: Put must at the start of required expressions and nodes.

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
	if len(tokenTypes) > 0 {
		if t := p.scanner.NextToken(); t.T == tokenTypes[0] {
			if ok := p.matchPhrase(tokenTypes[1:]...); ok {
				return true
			} else {
				p.scanner.ReturnToken(t)
				return false
			}
		} else {
			p.scanner.ReturnToken(t)
			return false
		}
	} else {
		return true
	}
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
	for property, _ := p.propertyLookup(); property != nil; property, _ = p.propertyLookup() {
		if property != nil {
			properties = append(properties, property)
		}
	}
	labels, _ := p.NodeLabels()
	return &ast.PropertyLabelsExpr{atom, properties, labels}
}

func (p *Parser) propertyLookup() (ast.SchemaName, error) {
	if s, err := p.schemaName(); err == nil {
		return s, nil
	}
	return nil, nil
}

func (p *Parser) schemaName() (ast.SchemaName, error) {
	t := p.scanner.NextToken()
	if _, ok := scanner.ReservedWordTokens[t.T]; ok {
		return &ast.ReservedWordSchemaName{t.T}, nil
	}
	p.scanner.ReturnToken(t)
	if symbolicName, err := p.symbolicName(); err == nil {
		return &ast.SymbolicNameSchemaName{symbolicName}, nil
	} else {
		return nil, err
	}
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
	if _, ok := p.match(scanner.Colon); ok {
		s, err := p.schemaName()
		if err != nil {
			return nil, err
		}
		if s == nil {
			p.reporter.Error(p.scanner.Line(), "expecting schema name following ':'")
		}
		return s, nil
	}
	return nil, nil
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
	if expr, _ := p.parameter(); expr != nil {
		return expr
	}
	if expr, _ := p.caseExpr(); expr != nil {
		return expr
	}
	if ok := p.matchPhrase(scanner.Identifier, scanner.OpenParen, scanner.Star, scanner.CloseParen); ok {
		return &ast.OpExpr{ast.CountAll}
	}
	if expr, _ := p.listComprehensionExpr(); expr != nil {
		return expr
	}
	if expr, _ := p.patternComprehensionExpr(); expr != nil {
		return expr
	}
	if expr, _ := p.builtInFunction(); expr != nil {
		return expr
	}
	if expr, _ := p.relationshipsPattern(); expr != nil {
		return expr
	}
	if expr, _ := p.parenthesizedExpr(); expr != nil {
		return expr
	}
	if expr, _ := p.variable(); expr != nil {
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

func (p *Parser) parameter() (ast.Expr, error) {
	if _, ok := p.match(scanner.DollarSign); ok {
		s, err := p.symbolicName()
		if err != nil {
			return nil, err
		}
		if s != nil {
			return &ast.Parameter{SymbolicName: s}, nil
		}
		if t, ok := p.match(scanner.DecimalInteger); ok {
			return &ast.Parameter{N: &t}, nil
		} else {
			return nil, p.reporter.Error(t.Line, "expecting symbolic name or integer")
		}
	}
	return nil, nil
}

func (p *Parser) caseExpr() (ast.Expr, error) {
	caseExpr := &ast.CaseExpr{}
	if t, ok := p.match(scanner.Case); ok {
		caseExpr.Init = p.expr()
		caseAlt, err := p.caseAlt()
		for err == nil && caseAlt != nil {
			caseExpr.Alternatives = append(caseExpr.Alternatives, caseAlt)
		}
		if len(caseExpr.Alternatives) == 0 {
			return nil, p.reporter.Error(t.Line, "expecting WHEN after CASE or CASE initialization expression")
		}
		if t, ok := p.match(scanner.Else); ok {
			caseExpr.Else = p.expr()
			if caseExpr.Else == nil {
				return nil, p.reporter.Error(t.Line, "expecting expression after CASE ELSE")
			}
		}
		if t, ok := p.match(scanner.End); ok {
			return caseExpr, nil
		} else {
			return nil, p.reporter.Error(t.Line, "expecting CASE END")
		}
	}
	return nil, nil
}

func (p *Parser) caseAlt() (*ast.CaseAltNode, error) {
	if _, ok := p.match(scanner.When); ok {
		whenExpr := p.expr()
		if t, ok := p.match(scanner.Then); ok {
			thenExpr := p.expr()
			return &ast.CaseAltNode{whenExpr, thenExpr}, nil
		} else {
			return nil, p.reporter.Error(t.Line, "expecting symbolic name or integer")
		}
	}
	return nil, nil
}

func (p *Parser) listComprehensionExpr() (ast.Expr, error) {
	if _, ok := p.match(scanner.OpenBracket); !ok {
		return nil, nil
	}
	var err error
	listCompExpr := &ast.ListComprehensionExpr{}
	listCompExpr.FilterExpr, err = p.filterExpr()
	if err != nil {
		return nil, err
	}
	if t, ok := p.match(scanner.Pipe); !ok {
		return nil, p.reporter.Error(t.Line, "expecting '|' in list expression")
	}
	listCompExpr.Expr = p.expr()
	if listCompExpr.Expr == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting expression after '|'")
	}
	if t, ok := p.match(scanner.CloseBracket); !ok {
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
	if t, ok := p.match(scanner.In); !ok {
		return nil, p.reporter.Error(t.Line, "expecting 'IN'")
	}
	if filterExpr.InExpr = p.expr(); filterExpr.InExpr == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting 'IN' expression")
	}
	if _, ok := p.match(scanner.Where); ok {
		if filterExpr.WhereExpr = p.expr(); filterExpr.WhereExpr == nil {
			return nil, p.reporter.Error(p.scanner.Line(), "expecting 'WHERE' expression")
		}
	}
	return filterExpr, nil
}

var builtInNames = map[string]ast.Operator{
	"ALL":    ast.AllOp,
	"ANY":    ast.AnyOp,
	"NONE":   ast.NoneOp,
	"SINGLE": ast.SingleOp,
}

func (p *Parser) builtInFunction() (ast.Expr, error) {
	if t, ok := p.match(scanner.Identifier); ok {
		if op, ok := builtInNames[t.Lexeme]; ok {
			if _, ok := p.match(scanner.OpenParen); !ok {
				return nil, p.reporter.Error(p.scanner.Line(), "expecting '('")
			}
			expr, err := p.filterExpr()
			if err != nil {
				return nil, err
			}
			if expr == nil {
				return nil, p.reporter.Error(p.scanner.Line(), "expecting filter expression")
			}
			if _, ok := p.match(scanner.OpenParen); !ok {
				return nil, p.reporter.Error(p.scanner.Line(), "expecting ')'")
			}
			return &ast.BuiltInExpr{op, expr}, nil
		} else {
			p.scanner.ReturnToken(t)
		}
	}
	return nil, nil
}

func (p *Parser) variable() (ast.Expr, error) {
	s, err := p.symbolicName()
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting variable")
	}
	return &ast.VariableExpr{s}, nil
}

func (p *Parser) symbolicName() (ast.SymbolicName, error) {
	if t, ok := p.match(scanner.Identifier); ok {
		if symbolType, ok := ast.SymbolNames[t.Lexeme]; ok {
			return &ast.SymbolicNameIdentifier{t, symbolType}, nil
		}
		if len(t.Lexeme) == 1 {
			switch ch, _ := utf8.DecodeLastRuneInString(t.Lexeme); ch {
			case 'a', 'b', 'c', 'd', 'e', 'f':
				return &ast.SymbolicNameHexLetter{ch}, nil
			}
			return &ast.SymbolicNameIdentifier{t, ast.Identifier}, nil
		}
	}
	return nil, nil
}

func (p *Parser) patternComprehensionExpr() (ast.Expr, error) {
	var err error
	patternExpr := &ast.PatternComprehensionExpr{}
	if _, ok := p.match(scanner.OpenBracket); !ok {
		return nil, nil
	}
	patternExpr.Variable, err = p.variable()
	if err != nil {
		return nil, err
	}
	if patternExpr.Variable != nil {
		if t, ok := p.match(scanner.Equal); !ok {
			return nil, p.reporter.Error(t.Line, "expecting '=' following variable")
		}
	}
	patternExpr.ReltionshipsPattern, err = p.relationshipsPattern()
	if patternExpr.ReltionshipsPattern == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting relationship pattern")
	}
	if _, ok := p.match(scanner.Where); ok {
		patternExpr.WhereExpr = p.expr()
		if p.expr() == nil {
			return nil, p.reporter.Error(p.scanner.Line(), "expecting expression following 'WHERE'")
		}
	}
	if t, ok := p.match(scanner.Pipe); !ok {
		return nil, p.reporter.Error(t.Line, "expecting '|'")
	}
	patternExpr.PipeExpr = p.expr()
	if p.expr() == nil {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting expression following '|'")
	}
	if t, ok := p.match(scanner.CloseBracket); !ok {
		return nil, p.reporter.Error(t.Line, "expecting ']'")
	}
	return nil, nil
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
	t, ok := p.match(scanner.LessThan)
	if ok {
		pattern.Left = ast.Directed
	}
	if _, ok := p.match(scanner.Dash); !ok {
		p.scanner.ReturnToken(t)
		return nil, nil
	}
	pattern.RelationshipDetail, err = p.relationshipDetail()
	if err != nil {
		return nil, err
	}
	if _, ok := p.match(scanner.Dash); !ok {
		return nil, p.reporter.Error(p.scanner.Line(), "expecting '-'")
	}
	if _, ok := p.match(scanner.GreaterThan); ok {
		pattern.Right = ast.Directed
	}
	return pattern, nil
}

func (p *Parser) relationshipDetail() (*ast.RelationshipDetail, error) {
	var err error
	if _, ok := p.match(scanner.OpenBracket); !ok {
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
	if t, ok := p.match(scanner.CloseBracket); !ok {
		return nil, p.reporter.Error(t.Line, "expecting ']'")
	}
	return detail, nil
}

func (p *Parser) properties() (*ast.Properties, error) {
	var err error
	properties := &ast.Properties{}
	properties.MapLiteral, err = p.mapLiteral()
	if err != nil {
		return nil, err
	}
	if properties.MapLiteral != nil {
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
	if _, ok := p.match(scanner.Star); !ok {
		return nil, nil
	}
	if t, ok := p.match(scanner.Integer); ok {
		literal.Begin = t.Literal.(int64)
	}
	if _, ok := p.match(scanner.Dotdot); ok {
		if t, ok := p.match(scanner.Integer); ok {
			literal.End = t.Literal.(int64)
		}
	}
	return literal, nil
}

func (p *Parser) relationshipTypes() ([]ast.SchemaName, error) {
	if _, ok := p.match(scanner.Colon); !ok {
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
		if _, ok := p.match(scanner.Pipe); !ok {
			return typeNames, nil
		}
		p.match(scanner.Colon)
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
	if _, ok := p.match(scanner.OpenParen); !ok {
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
	literal, err := p.mapLiteral()
	if err != nil {
		return nil, err
	}
	if literal == nil {
		parameter, err := p.parameter()
		if err != nil {
			return nil, err
		}
		if parameter != nil {
			np.Properties = &ast.Properties{Parameter: parameter}
		}
	} else {
		np.Properties = &ast.Properties{MapLiteral: literal}
	}
	if t, ok := p.match(scanner.CloseParen); !ok {
		return nil, p.reporter.Error(t.Line, "expecting ')' following node pattern")
	}
	return np, nil
}

func (p *Parser) mapLiteral() (*ast.MapLiteral, error) {
	if _, ok := p.match(scanner.OpenBrace); !ok {
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
		if t, ok := p.match(scanner.Colon); !ok {
			return nil, p.reporter.Error(t.Line, "expecting '}' following map literal")
		}
		if pkn.Expr = p.expr(); pkn.Expr == nil {
			return nil, p.reporter.Error(p.scanner.Line(), "expecting expression following ':'")
		}
		literal.PropertyKeyNames = append(literal.PropertyKeyNames, pkn)
	}
	if t, ok := p.match(scanner.CloseBrace); !ok {
		return nil, p.reporter.Error(t.Line, "expecting '}' following map literal")
	}
	return literal, nil
}

func (p *Parser) parenthesizedExpr() (ast.Expr, error) {
	if _, ok := p.match(scanner.OpenParen); !ok {
		return nil, nil
	}
	expr := p.expr()
	if t, ok := p.match(scanner.CloseParen); !ok {
		return nil, p.reporter.Error(t.Line, "expecting ')' following expression")
	}
	return expr, nil
}

func (p *Parser) functionInvocation() (ast.Expr, error) {
	fn, err := p.functionName()
	if err != nil {
		return nil, err
	}
	if t, ok := p.match(scanner.OpenParen); !ok {
		return nil, p.reporter.Error(t.Line, "expecting '(' function name")
	}
	_, distinct := p.match(scanner.Distinct)
	if t, ok := p.match(scanner.CloseParen); !ok {
		return nil, p.reporter.Error(t.Line, "expecting ')' function parameters")
	}
	args := []ast.Expr{}
	expr := p.expr()
	if expr != nil {
		args = append(args, expr)
		for {
			if _, ok := p.match(scanner.Comma); !ok {
				break
			}
			expr = p.expr()
			if expr == nil {
				return nil, p.reporter.Error(p.scanner.Line(), "expecting argument following ','")
			}
		}
	}
	return &ast.FunctionInvocation{FunctionName: fn, Distinct: distinct, Args: args}, nil
}

func (p *Parser) functionName() (ast.FunctionName, error) {
	if _, ok := p.match(scanner.Exists); ok {
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
		sn, err := p.symbolicName()
		if err != nil {
			return nil, err
		}
		if sn == nil {
			break
		}
		namespace = append(namespace, sn)
		if t, ok := p.match(scanner.Period); !ok {
			return nil, p.reporter.Error(t.Line, "expecting '.' after namespace identifier")
		}
	}
	return namespace, nil
}
