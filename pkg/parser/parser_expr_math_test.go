package parser

import (
	"github.com/mburbidg/cypher/pkg/ast"
	"github.com/mburbidg/cypher/pkg/scanner"
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
	assert.Equal(t, "g", tree.(*ast.BinaryExpr).Left.(*ast.PropertyLabelsExpr).Atom.(*ast.VariableExpr).SymbolicName.(*ast.SymbolicNameIdentifier).Identifier.Lexeme)
	assert.Equal(t, 1, len(tree.(*ast.BinaryExpr).Left.(*ast.PropertyLabelsExpr).PropertyKeys))
	assert.Equal(t, "id", tree.(*ast.BinaryExpr).Left.(*ast.PropertyLabelsExpr).PropertyKeys[0].(*ast.SymbolicNameSchemaName).SymbolicName.(*ast.SymbolicNameIdentifier).Identifier.Lexeme)
	assert.Equal(t, ast.Equal, tree.(*ast.BinaryExpr).Op)
	assert.Equal(t, int64(1337), tree.(*ast.BinaryExpr).Right.(*ast.PropertyLabelsExpr).Atom.(*ast.PrimitiveLiteral).Value.(int64))

	s = scanner.New([]byte("g.version + 5"), reporter)
	p = New(s, reporter)
	tree, err = p.expr()
	assert.NoError(t, err)
	assert.Equal(t, "g", tree.(*ast.BinaryExpr).Left.(*ast.PropertyLabelsExpr).Atom.(*ast.VariableExpr).SymbolicName.(*ast.SymbolicNameIdentifier).Identifier.Lexeme)
	assert.Equal(t, 1, len(tree.(*ast.BinaryExpr).Left.(*ast.PropertyLabelsExpr).PropertyKeys))
	assert.Equal(t, "version", tree.(*ast.BinaryExpr).Left.(*ast.PropertyLabelsExpr).PropertyKeys[0].(*ast.SymbolicNameSchemaName).SymbolicName.(*ast.SymbolicNameIdentifier).Identifier.Lexeme)
	assert.Equal(t, ast.Add, tree.(*ast.BinaryExpr).Op)
	assert.Equal(t, int64(5), tree.(*ast.BinaryExpr).Right.(*ast.PropertyLabelsExpr).Atom.(*ast.PrimitiveLiteral).Value.(int64))
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
	assert.Equal(t, int64(12), tree.(*ast.BinaryExpr).Left.(*ast.BinaryExpr).Left.(*ast.BinaryExpr).Left.(*ast.PropertyLabelsExpr).Atom.(*ast.PrimitiveLiteral).Value.(int64))
	assert.Equal(t, ast.Divide, tree.(*ast.BinaryExpr).Left.(*ast.BinaryExpr).Left.(*ast.BinaryExpr).Op)
	assert.Equal(t, int64(4), tree.(*ast.BinaryExpr).Left.(*ast.BinaryExpr).Left.(*ast.BinaryExpr).Right.(*ast.PropertyLabelsExpr).Atom.(*ast.PrimitiveLiteral).Value.(int64))
	assert.Equal(t, ast.Multiply, tree.(*ast.BinaryExpr).Left.(*ast.BinaryExpr).Op)
	assert.Equal(t, int64(3), tree.(*ast.BinaryExpr).Left.(*ast.BinaryExpr).Right.(*ast.PropertyLabelsExpr).Atom.(*ast.PrimitiveLiteral).Value.(int64))
	assert.Equal(t, ast.Subtract, tree.(*ast.BinaryExpr).Op)
	assert.Equal(t, int64(2), tree.(*ast.BinaryExpr).Right.(*ast.BinaryExpr).Left.(*ast.PropertyLabelsExpr).Atom.(*ast.PrimitiveLiteral).Value.(int64))
	assert.Equal(t, int64(4), tree.(*ast.BinaryExpr).Right.(*ast.BinaryExpr).Right.(*ast.PropertyLabelsExpr).Atom.(*ast.PrimitiveLiteral).Value.(int64))

	s = scanner.New([]byte("12 / 4 * (3 - 2 * 4)"), reporter)
	p = New(s, reporter)
	tree, err = p.expr()
	assert.NoError(t, err)
	assert.Equal(t, int64(12), tree.(*ast.BinaryExpr).Left.(*ast.BinaryExpr).Left.(*ast.PropertyLabelsExpr).Atom.(*ast.PrimitiveLiteral).Value.(int64))
	assert.Equal(t, ast.Divide, tree.(*ast.BinaryExpr).Left.(*ast.BinaryExpr).Op)
	assert.Equal(t, int64(4), tree.(*ast.BinaryExpr).Left.(*ast.BinaryExpr).Right.(*ast.PropertyLabelsExpr).Atom.(*ast.PrimitiveLiteral).Value.(int64))
	assert.Equal(t, ast.Multiply, tree.(*ast.BinaryExpr).Op)
	assert.Equal(t, int64(3), tree.(*ast.BinaryExpr).Right.(*ast.PropertyLabelsExpr).Atom.(*ast.BinaryExpr).Left.(*ast.PropertyLabelsExpr).Atom.(*ast.PrimitiveLiteral).Value.(int64))
	assert.Equal(t, ast.Subtract, tree.(*ast.BinaryExpr).Right.(*ast.PropertyLabelsExpr).Atom.(*ast.BinaryExpr).Op)
	assert.Equal(t, int64(2), tree.(*ast.BinaryExpr).Right.(*ast.PropertyLabelsExpr).Atom.(*ast.BinaryExpr).Right.(*ast.BinaryExpr).Left.(*ast.PropertyLabelsExpr).Atom.(*ast.PrimitiveLiteral).Value.(int64))
	assert.Equal(t, ast.Multiply, tree.(*ast.BinaryExpr).Right.(*ast.PropertyLabelsExpr).Atom.(*ast.BinaryExpr).Right.(*ast.BinaryExpr).Op)
	assert.Equal(t, int64(4), tree.(*ast.BinaryExpr).Right.(*ast.PropertyLabelsExpr).Atom.(*ast.BinaryExpr).Right.(*ast.BinaryExpr).Right.(*ast.PropertyLabelsExpr).Atom.(*ast.PrimitiveLiteral).Value.(int64))
}

// mathematical/Mathematical11.feature
func TestAbsoluteFunction(t *testing.T) {
	reporter := newTestReporter()
	s := scanner.New([]byte("abs(-1)"), reporter)
	p := New(s, reporter)
	tree, err := p.expr()
	assert.NoError(t, err)
	name := tree.(*ast.PropertyLabelsExpr).Atom.(*ast.FunctionInvocation).FunctionName.(*ast.SymbolicFunctionName).FunctionName.(*ast.SymbolicNameIdentifier).Identifier.Lexeme
	assert.Equal(t, "abs", name)
	args := tree.(*ast.PropertyLabelsExpr).Atom.(*ast.FunctionInvocation).Args
	assert.Equal(t, 1, len(args))
	assert.Equal(t, ast.Negate, args[0].(*ast.UnaryExpr).Op)
	assert.Equal(t, int64(1), args[0].(*ast.UnaryExpr).Expr.(*ast.PropertyLabelsExpr).Atom.(*ast.PrimitiveLiteral).Value.(int64))
}

// mathematical/Mathematical13.feature
func TestReturningFloatValues(t *testing.T) {
	reporter := newTestReporter()
	s := scanner.New([]byte("sqrt(12.96)"), reporter)
	p := New(s, reporter)
	tree, err := p.expr()
	assert.NoError(t, err)
	name := tree.(*ast.PropertyLabelsExpr).Atom.(*ast.FunctionInvocation).FunctionName.(*ast.SymbolicFunctionName).FunctionName.(*ast.SymbolicNameIdentifier).Identifier.Lexeme
	assert.Equal(t, "sqrt", name)
	args := tree.(*ast.PropertyLabelsExpr).Atom.(*ast.FunctionInvocation).Args
	assert.Equal(t, 1, len(args))
	assert.Equal(t, float64(12.96), args[0].(*ast.PropertyLabelsExpr).Atom.(*ast.PrimitiveLiteral).Value.(float64))
}
