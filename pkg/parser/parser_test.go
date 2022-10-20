package parser

import (
	"github.com/mburbidg/cypher/pkg/scanner"
	"github.com/mburbidg/cypher/pkg/utils"
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
	s := scanner.New([]byte("COUNT(*)"), reporter)
	p := New(s, reporter)
	_, ok, err := p.matchPhrase(scanner.Identifier, scanner.OpenParen, scanner.Colon, scanner.CloseParen, scanner.EndOfInput)
	assert.NoError(t, err)
	assert.False(t, ok)
	_, ok, err = p.matchPhrase(scanner.Identifier, scanner.OpenParen, scanner.Star, scanner.CloseParen, scanner.EndOfInput)
	assert.NoError(t, err)
	assert.True(t, ok)
}
