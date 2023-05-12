package parser

import (
	scanner2 "github.com/mburbidg/cypher/scanner"
	"github.com/mburbidg/cypher/utils"
	"github.com/stretchr/testify/assert"
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
	s := scanner2.New([]byte("COUNT(*)"), reporter)
	p := New(s, reporter)
	_, ok, err := p.matchPhrase(scanner2.Identifier, scanner2.OpenParen, scanner2.Colon, scanner2.CloseParen, scanner2.EndOfInput)
	assert.NoError(t, err)
	assert.False(t, ok)
	_, ok, err = p.matchPhrase(scanner2.Identifier, scanner2.OpenParen, scanner2.Star, scanner2.CloseParen, scanner2.EndOfInput)
	assert.NoError(t, err)
	assert.True(t, ok)
}
