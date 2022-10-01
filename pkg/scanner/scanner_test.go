package scanner

import (
	"bytes"
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
		"zero":              {"0", []TokenType{Integer, EndOfInput}},
		"zero:na":           {"0a", []TokenType{Integer, Identifier, EndOfInput}},
		"zero:-n+":          {"-0+", []TokenType{Minus, Integer, Plus, EndOfInput}},
		"zero:a n b":        {"a 0 b", []TokenType{Identifier, Integer, Identifier, EndOfInput}},
		"integer":           {"240", []TokenType{Integer, EndOfInput}},
		"integer:-n+":       {"-240+", []TokenType{Minus, Integer, Plus, EndOfInput}},
		"integer:a n b":     {"a 240 b", []TokenType{Identifier, Integer, Identifier, EndOfInput}},
		"integer:xa":        {"10a", []TokenType{Integer, Identifier, EndOfInput}},
		"integer:0xa":       {"0x3ae1", []TokenType{Integer, EndOfInput}},
		"integer:0nn":       {"0371", []TokenType{Integer, EndOfInput}},
		"double:0.x":        {"0.1", []TokenType{Double, EndOfInput}},
		"double:-0.x+":      {"-0.1+", []TokenType{Minus, Double, Plus, EndOfInput}},
		"double:a 0.x b":    {"a 0.1 b", []TokenType{Identifier, Double, Identifier, EndOfInput}},
		"double:0.xa":       {"0.1a", []TokenType{Double, Identifier, EndOfInput}},
		"double:x.x":        {"25.1", []TokenType{Double, EndOfInput}},
		"double:(x.x)":      {"(25.1)", []TokenType{OpenParen, Double, CloseParen, EndOfInput}},
		"double:a x.x b":    {"a 25.1 b", []TokenType{Identifier, Double, Identifier, EndOfInput}},
		"double:x.xa":       {"25.1a", []TokenType{Double, Identifier, EndOfInput}},
		"double:.x":         {".15", []TokenType{Double, EndOfInput}},
		"double:-.x+":       {"-.15+", []TokenType{Minus, Double, Plus, EndOfInput}},
		"double:.a x b":     {"a .15 b", []TokenType{Identifier, Double, Identifier, EndOfInput}},
		"double:.xa":        {".15a", []TokenType{Double, Identifier, EndOfInput}},
		"double:xEx":        {"10E3", []TokenType{Double, EndOfInput}},
		"double:-xEx+":      {"-10E3+", []TokenType{Minus, Double, Plus, EndOfInput}},
		"double:a xEx b":    {"a 10E3 b", []TokenType{Identifier, Double, Identifier, EndOfInput}},
		"double:xExa":       {"10E3a", []TokenType{Double, Identifier, EndOfInput}},
		"double:x.xEx":      {"1.35E3", []TokenType{Double, EndOfInput}},
		"double:-x.xEx+":    {"-1.35E3+", []TokenType{Minus, Double, Plus, EndOfInput}},
		"double:a x.xEx b":  {"a 1.35E3b ", []TokenType{Identifier, Double, Identifier, EndOfInput}},
		"double:x.xExa":     {"1.35E3a", []TokenType{Double, Identifier, EndOfInput}},
		"double:x.xE-x":     {"1.35E-3", []TokenType{Double, EndOfInput}},
		"double:(x.xE-x)":   {"(1.35E-3)", []TokenType{OpenParen, Double, CloseParen, EndOfInput}},
		"double:a x.xE-x b": {"a 1.35E-3 b", []TokenType{Identifier, Double, Identifier, EndOfInput}},
		"double:x.xE-xa":    {"1.35E-3a", []TokenType{Double, Identifier, EndOfInput}},
		"double:x.xE":       {"1.35E", []TokenType{Illegal, EndOfInput}},
		"double:x.xE-":      {"1.35E-", []TokenType{Illegal, EndOfInput}},
		"double:x.":         {"14.", []TokenType{Illegal, EndOfInput}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := New(bytes.NewBufferString(tc.src), newTestReporter())
			assertTokens(t, tc.tokens, s)
		})
	}
}

func TestNumberValues(t *testing.T) {
	tests := map[string]struct {
		src      string
		literals []interface{}
	}{
		"zero":        {"0", []interface{}{int64(0)}},
		"integer:145": {"145", []interface{}{int64(145)}},
		"integer:0xa": {"0xa", []interface{}{int64(10)}},
		"integer:077": {"077", []interface{}{int64(63)}},
		"double:.10":  {".10", []interface{}{float64(.10)}},
		"double:1.10": {"1.10", []interface{}{float64(1.10)}},
		"double:1E3":  {"1E3", []interface{}{float64(1000.0)}},
		"double:1E-3": {"1E-3", []interface{}{float64(1.0 / 1000.0)}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := New(bytes.NewBufferString(tc.src), newTestReporter())
			assertLiteral(t, tc.literals, s)
		})
	}
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

func TestPunctuation(t *testing.T) {
	tests := map[string]struct {
		src    string
		tokens []TokenType
	}{
		"punctuation:1":        {".(){}[]", []TokenType{Period, OpenParen, CloseParen, OpenBrace, CloseBrace, OpenBracket, CloseBracket, EndOfInput}},
		"punctuation:1/ws":     {". ( ) { } [ ]", []TokenType{Period, OpenParen, CloseParen, OpenBrace, CloseBrace, OpenBracket, CloseBracket, EndOfInput}},
		"punctuation:2":        {"+-*/%^$", []TokenType{Plus, Minus, Star, ForwardSlash, Percent, Caret, DollarSign, EndOfInput}},
		"punctuation:2/ws":     {"+ - * / % ^ $", []TokenType{Plus, Minus, Star, ForwardSlash, Percent, Caret, DollarSign, EndOfInput}},
		"equal":                {"a=b", []TokenType{Identifier, Equal, Identifier, EndOfInput}},
		"equal/ws":             {"a = b", []TokenType{Identifier, Equal, Identifier, EndOfInput}},
		"!equal":               {"a<>b", []TokenType{Identifier, NotEqual, Identifier, EndOfInput}},
		"!equal/ws":            {"a <> b", []TokenType{Identifier, NotEqual, Identifier, EndOfInput}},
		"lessthan":             {"a<b", []TokenType{Identifier, LessThan, Identifier, EndOfInput}},
		"lessthan/ws":          {"a < b", []TokenType{Identifier, LessThan, Identifier, EndOfInput}},
		"lessthanorequal":      {"a<=b", []TokenType{Identifier, LessThanOrEqual, Identifier, EndOfInput}},
		"lessthanorequal/ws":   {"a <= b", []TokenType{Identifier, LessThanOrEqual, Identifier, EndOfInput}},
		"greaterthan":          {"a>b", []TokenType{Identifier, GreaterThan, Identifier, EndOfInput}},
		"greaterthan/ws":       {"a > b", []TokenType{Identifier, GreaterThan, Identifier, EndOfInput}},
		"greaterthanorequal":   {"a>=b", []TokenType{Identifier, GreaterThanOrEqual, Identifier, EndOfInput}},
		"greaterhanorequal/ws": {"a >= b", []TokenType{Identifier, GreaterThanOrEqual, Identifier, EndOfInput}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := New(bytes.NewBufferString(tc.src), newTestReporter())
			assertTokens(t, tc.tokens, s)
		})
	}
}

func assertTokens(t *testing.T, expected []TokenType, scanner *Scanner) {
	for _, tokenType := range expected {
		token := scanner.NextToken()
		assert.Equal(t, tokenType, token.t)
	}
}

func assertLiteral(t *testing.T, expected []interface{}, scanner *Scanner) {
	for _, literal := range expected {
		token := scanner.NextToken()
		assert.Equal(t, literal, token.literal)
	}
}
