package scanner

import (
	"bufio"
	"bytes"
	"fmt"
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

func (r *testReporter) Error(line int, msg string) {
	r.errors = append(r.errors, errorMsg{
		line: line,
		msg:  msg,
	})
}

func TestScanner(t *testing.T) {
	s := New(bytes.NewBufferString("Create"), newTestReporter())
	assertTokens(t, []TokenType{Create, EndOfInput}, s)

	s = New(bytes.NewBufferString("+"), newTestReporter())
	assertTokens(t, []TokenType{Plus, EndOfInput}, s)

	s = New(bytes.NewBufferString("MATCH (n) RETURN n WHERE n.foo = 1"), newTestReporter())
	assertTokens(t, []TokenType{Match, OpenParen, SymbolName, CloseParen, Return, SymbolName, Where, SymbolName, Period, SymbolName, Equal, Integer, EndOfInput}, s)
}

func TestBogusNumber(t *testing.T) {
	reporter := newTestReporter()
	s := New(bytes.NewBufferString("1a 1.2"), reporter)
	assertTokens(t, []TokenType{Error, Double, EndOfInput}, s)
	assert.Equal(t, 1, len(reporter.errors))
}

func TestString(t *testing.T) {
	reporter := newTestReporter()
	s := New(bytes.NewBufferString("\"This is a \\u00fa \\U032bca08 string.\\n\""), reporter)
	assertTokens(t, []TokenType{String, EndOfInput}, s)
	assert.Equal(t, 0, len(reporter.errors))
	s = New(bytes.NewBufferString("'This is a \\u00fa \\U032bca08 string.\\n'"), reporter)
	assertTokens(t, []TokenType{String, EndOfInput}, s)
	assert.Equal(t, 0, len(reporter.errors))
}

func TestIt(t *testing.T) {
	buf := bytes.NewBufferString("Create")
	r := bufio.NewReader(buf)
	ch, _, err := r.ReadRune()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ch)
	ch, _, err = r.ReadRune()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ch)
	ch, _, err = r.ReadRune()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ch)
	ch, _, err = r.ReadRune()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ch)
	ch, _, err = r.ReadRune()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ch)
	ch, _, err = r.ReadRune()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ch)
	ch, _, err = r.ReadRune()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ch)
	err = r.UnreadRune()
	if err != nil {
		fmt.Println(err)
	}
}

func assertTokens(t *testing.T, expected []TokenType, scanner *Scanner) {
	for _, tokenType := range expected {
		token := scanner.NextToken()
		assert.Equal(t, tokenType, token.t)
	}
}
