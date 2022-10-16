package parser

import (
	"bytes"
	"github.com/mburbidg/cypher/pkg/scanner"
	"github.com/mburbidg/cypher/pkg/utils"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

type errorMsg struct {
	line int
	msg  string
}

type testReporter struct {
	errors []errorMsg
}

func newTestReporter() *testReporter {
	return &testReporter{
		errors: make([]errorMsg, 0, 10),
	}
}

func (r *testReporter) Error(line int, msg string) error {
	r.errors = append(r.errors, errorMsg{
		line: line,
		msg:  msg,
	})
	return utils.ParseError{line, msg}
}

func TestMatchPhrase(t *testing.T) {
	reporter := newTestReporter()
	s := scanner.New(bytes.NewBufferString("COUNT(*)"), reporter)
	p := New(s, reporter)
	assert.False(t, p.matchPhrase(scanner.Identifier, scanner.OpenParen, scanner.Colon, scanner.CloseParen, scanner.EndOfInput))
	assert.True(t, p.matchPhrase(scanner.Identifier, scanner.OpenParen, scanner.Star, scanner.CloseParen, scanner.EndOfInput))
}

func TestAdditionExpr(t *testing.T) {
	reporter := newTestReporter()
	s := scanner.New(bytes.NewBufferString("3 + 5"), reporter)
	p := New(s, reporter)
	ast, err := p.Parse()
	assert.NoError(t, err)
	log.Printf("ast=%v\n", ast)
}

func TestMathExpr(t *testing.T) {
	reporter := newTestReporter()
	s := scanner.New(bytes.NewBufferString("12 / 4 * (3 - 2 * 4)"), reporter)
	p := New(s, reporter)
	ast, err := p.Parse()
	assert.NoError(t, err)
	log.Printf("ast=%v\n", ast)
}
