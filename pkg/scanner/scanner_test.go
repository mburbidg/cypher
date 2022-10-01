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
}

func TestNumbers(t *testing.T) {
	tests := map[string]struct {
		src    string
		tokens []TokenType
	}{
		"zero":        {"0", []TokenType{Integer, EndOfInput}},
		"integer":     {"240", []TokenType{Integer, EndOfInput}},
		"integer:xa":  {"10a", []TokenType{Integer, Identifier, EndOfInput}},
		"integer:0xa": {"0x3ae1", []TokenType{Integer, EndOfInput}},
		"integer:0nn": {"0371", []TokenType{Integer, EndOfInput}},
		"double:0.x":  {"0.1", []TokenType{Double, EndOfInput}},
		"double:x.x":  {"25.1", []TokenType{Double, EndOfInput}},
		"double:.x":   {".15", []TokenType{Double, EndOfInput}},
		"double:x.":   {"14.", []TokenType{Illegal, EndOfInput}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := New(bytes.NewBufferString(tc.src), newTestReporter())
			assertTokens(t, tc.tokens, s)
		})
	}
}

func TestBogusNumber(t *testing.T) {
	reporter := newTestReporter()
	s := New(bytes.NewBufferString("1a 1.2"), reporter)
	assertTokens(t, []TokenType{Illegal, Double, EndOfInput}, s)
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
