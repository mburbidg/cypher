package parser

import (
	"github.com/mburbidg/cypher/ast"
	"github.com/mburbidg/cypher/scanner"
)

var opForTokens = map[scanner.TokenType]ast.Operator{
	scanner.Equal:              ast.Equal,
	scanner.NotEqual:           ast.NotEqual,
	scanner.LessThan:           ast.LessThan,
	scanner.GreaterThan:        ast.GreaterThan,
	scanner.LessThanOrEqual:    ast.LessThanOrEqual,
	scanner.GreaterThanOrEqual: ast.GreaterThanOrEqual,
	scanner.Plus:               ast.Add,
	scanner.Dash:               ast.Subtract,
	scanner.Star:               ast.Multiply,
	scanner.ForwardSlash:       ast.Divide,
	scanner.Percent:            ast.Modulo,
	scanner.Caret:              ast.PowerOf,
}
