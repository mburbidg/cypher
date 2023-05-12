package parser

import (
	ast2 "github.com/mburbidg/cypher/ast"
	"github.com/mburbidg/cypher/scanner"
	"github.com/stretchr/testify/assert"
	"testing"
)

// We test the parser against the tkl use cases, for expressions.
// https://github.com/opencypher/openCypher/tree/master/tck/features/expressions/mathematical

// mathematical/Mathematical2.feature
func TestAdditionExpr(t *testing.T) {
	reporter := newTestReporter()
	s := scanner.New([]byte("g.id = 1337"), reporter)
	p := New(s, reporter)
	tree, err := p.expr()
	assert.NoError(t, err)
	assert.Equal(t, "g", tree.(*ast2.BinaryExpr).Left.(*ast2.PropertyLabelsExpr).Atom.(*ast2.VariableExpr).SymbolicName.(*ast2.SymbolicNameIdentifier).Identifier.Lexeme)
	assert.Equal(t, 1, len(tree.(*ast2.BinaryExpr).Left.(*ast2.PropertyLabelsExpr).PropertyKeys))
	assert.Equal(t, "id", tree.(*ast2.BinaryExpr).Left.(*ast2.PropertyLabelsExpr).PropertyKeys[0].(*ast2.SymbolicNameSchemaName).SymbolicName.(*ast2.SymbolicNameIdentifier).Identifier.Lexeme)
	assert.Equal(t, ast2.Equal, tree.(*ast2.BinaryExpr).Op)
	assert.Equal(t, int64(1337), tree.(*ast2.BinaryExpr).Right.(*ast2.PropertyLabelsExpr).Atom.(*ast2.PrimitiveLiteral).Value.(int64))

	s = scanner.New([]byte("g.version + 5"), reporter)
	p = New(s, reporter)
	tree, err = p.expr()
	assert.NoError(t, err)
	assert.Equal(t, "g", tree.(*ast2.BinaryExpr).Left.(*ast2.PropertyLabelsExpr).Atom.(*ast2.VariableExpr).SymbolicName.(*ast2.SymbolicNameIdentifier).Identifier.Lexeme)
	assert.Equal(t, 1, len(tree.(*ast2.BinaryExpr).Left.(*ast2.PropertyLabelsExpr).PropertyKeys))
	assert.Equal(t, "version", tree.(*ast2.BinaryExpr).Left.(*ast2.PropertyLabelsExpr).PropertyKeys[0].(*ast2.SymbolicNameSchemaName).SymbolicName.(*ast2.SymbolicNameIdentifier).Identifier.Lexeme)
	assert.Equal(t, ast2.Add, tree.(*ast2.BinaryExpr).Op)
	assert.Equal(t, int64(5), tree.(*ast2.BinaryExpr).Right.(*ast2.PropertyLabelsExpr).Atom.(*ast2.PrimitiveLiteral).Value.(int64))
}

// mathematical/Mathematical3.feature
func TestInvalidHyphenInSubtraction(t *testing.T) {
	reporter := newTestReporter()
	s := scanner.New([]byte("42 — 41"), reporter)
	p := New(s, reporter)
	_, err := p.Parse()
	assert.Error(t, err)
}

// mathematical/Mathematical8.feature
func TestMathematicalPrecedence(t *testing.T) {
	reporter := newTestReporter()
	s := scanner.New([]byte("12 / 4 * 3 - 2 * 4"), reporter)
	p := New(s, reporter)
	tree, err := p.expr()
	assert.NoError(t, err)
	assert.Equal(t, int64(12), tree.(*ast2.BinaryExpr).Left.(*ast2.BinaryExpr).Left.(*ast2.BinaryExpr).Left.(*ast2.PropertyLabelsExpr).Atom.(*ast2.PrimitiveLiteral).Value.(int64))
	assert.Equal(t, ast2.Divide, tree.(*ast2.BinaryExpr).Left.(*ast2.BinaryExpr).Left.(*ast2.BinaryExpr).Op)
	assert.Equal(t, int64(4), tree.(*ast2.BinaryExpr).Left.(*ast2.BinaryExpr).Left.(*ast2.BinaryExpr).Right.(*ast2.PropertyLabelsExpr).Atom.(*ast2.PrimitiveLiteral).Value.(int64))
	assert.Equal(t, ast2.Multiply, tree.(*ast2.BinaryExpr).Left.(*ast2.BinaryExpr).Op)
	assert.Equal(t, int64(3), tree.(*ast2.BinaryExpr).Left.(*ast2.BinaryExpr).Right.(*ast2.PropertyLabelsExpr).Atom.(*ast2.PrimitiveLiteral).Value.(int64))
	assert.Equal(t, ast2.Subtract, tree.(*ast2.BinaryExpr).Op)
	assert.Equal(t, int64(2), tree.(*ast2.BinaryExpr).Right.(*ast2.BinaryExpr).Left.(*ast2.PropertyLabelsExpr).Atom.(*ast2.PrimitiveLiteral).Value.(int64))
	assert.Equal(t, int64(4), tree.(*ast2.BinaryExpr).Right.(*ast2.BinaryExpr).Right.(*ast2.PropertyLabelsExpr).Atom.(*ast2.PrimitiveLiteral).Value.(int64))

	s = scanner.New([]byte("12 / 4 * (3 - 2 * 4)"), reporter)
	p = New(s, reporter)
	tree, err = p.expr()
	assert.NoError(t, err)
	assert.Equal(t, int64(12), tree.(*ast2.BinaryExpr).Left.(*ast2.BinaryExpr).Left.(*ast2.PropertyLabelsExpr).Atom.(*ast2.PrimitiveLiteral).Value.(int64))
	assert.Equal(t, ast2.Divide, tree.(*ast2.BinaryExpr).Left.(*ast2.BinaryExpr).Op)
	assert.Equal(t, int64(4), tree.(*ast2.BinaryExpr).Left.(*ast2.BinaryExpr).Right.(*ast2.PropertyLabelsExpr).Atom.(*ast2.PrimitiveLiteral).Value.(int64))
	assert.Equal(t, ast2.Multiply, tree.(*ast2.BinaryExpr).Op)
	assert.Equal(t, int64(3), tree.(*ast2.BinaryExpr).Right.(*ast2.PropertyLabelsExpr).Atom.(*ast2.BinaryExpr).Left.(*ast2.PropertyLabelsExpr).Atom.(*ast2.PrimitiveLiteral).Value.(int64))
	assert.Equal(t, ast2.Subtract, tree.(*ast2.BinaryExpr).Right.(*ast2.PropertyLabelsExpr).Atom.(*ast2.BinaryExpr).Op)
	assert.Equal(t, int64(2), tree.(*ast2.BinaryExpr).Right.(*ast2.PropertyLabelsExpr).Atom.(*ast2.BinaryExpr).Right.(*ast2.BinaryExpr).Left.(*ast2.PropertyLabelsExpr).Atom.(*ast2.PrimitiveLiteral).Value.(int64))
	assert.Equal(t, ast2.Multiply, tree.(*ast2.BinaryExpr).Right.(*ast2.PropertyLabelsExpr).Atom.(*ast2.BinaryExpr).Right.(*ast2.BinaryExpr).Op)
	assert.Equal(t, int64(4), tree.(*ast2.BinaryExpr).Right.(*ast2.PropertyLabelsExpr).Atom.(*ast2.BinaryExpr).Right.(*ast2.BinaryExpr).Right.(*ast2.PropertyLabelsExpr).Atom.(*ast2.PrimitiveLiteral).Value.(int64))
}

// mathematical/Mathematical11.feature
func TestAbsoluteFunction(t *testing.T) {
	reporter := newTestReporter()
	s := scanner.New([]byte("abs(-1)"), reporter)
	p := New(s, reporter)
	tree, err := p.expr()
	assert.NoError(t, err)
	name := tree.(*ast2.PropertyLabelsExpr).Atom.(*ast2.FunctionInvocation).FunctionName.(*ast2.SymbolicFunctionName).FunctionName.(*ast2.SymbolicNameIdentifier).Identifier.Lexeme
	assert.Equal(t, "abs", name)
	args := tree.(*ast2.PropertyLabelsExpr).Atom.(*ast2.FunctionInvocation).Args
	assert.Equal(t, 1, len(args))
	assert.Equal(t, ast2.Negate, args[0].(*ast2.UnaryExpr).Op)
	assert.Equal(t, int64(1), args[0].(*ast2.UnaryExpr).Expr.(*ast2.PropertyLabelsExpr).Atom.(*ast2.PrimitiveLiteral).Value.(int64))
}

// mathematical/Mathematical13.feature
func TestReturningFloatValues(t *testing.T) {
	reporter := newTestReporter()
	s := scanner.New([]byte("sqrt(12.96)"), reporter)
	p := New(s, reporter)
	tree, err := p.expr()
	assert.NoError(t, err)
	name := tree.(*ast2.PropertyLabelsExpr).Atom.(*ast2.FunctionInvocation).FunctionName.(*ast2.SymbolicFunctionName).FunctionName.(*ast2.SymbolicNameIdentifier).Identifier.Lexeme
	assert.Equal(t, "sqrt", name)
	args := tree.(*ast2.PropertyLabelsExpr).Atom.(*ast2.FunctionInvocation).Args
	assert.Equal(t, 1, len(args))
	assert.Equal(t, float64(12.96), args[0].(*ast2.PropertyLabelsExpr).Atom.(*ast2.PrimitiveLiteral).Value.(float64))
}
