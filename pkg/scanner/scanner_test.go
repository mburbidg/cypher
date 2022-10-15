package scanner

import (
	"bytes"
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

func TestScanner(t *testing.T) {
}

func TestNumbers(t *testing.T) {
	tests := map[string]struct {
		src    string
		tokens []TokenType
	}{
		"zero":              {"0", []TokenType{DecimalInteger, EndOfInput}},
		"zero:na":           {"0a", []TokenType{DecimalInteger, Identifier, EndOfInput}},
		"zero:-n+":          {"-0+", []TokenType{Minus, DecimalInteger, Plus, EndOfInput}},
		"zero:a n b":        {"a 0 b", []TokenType{Identifier, DecimalInteger, Identifier, EndOfInput}},
		"integer":           {"240", []TokenType{DecimalInteger, EndOfInput}},
		"integer:-n+":       {"-240+", []TokenType{Minus, DecimalInteger, Plus, EndOfInput}},
		"integer:a n b":     {"a 240 b", []TokenType{Identifier, DecimalInteger, Identifier, EndOfInput}},
		"integer:xa":        {"10a", []TokenType{DecimalInteger, Identifier, EndOfInput}},
		"integer:0xa":       {"0x3ae1", []TokenType{HexInteger, EndOfInput}},
		"integer:0nn":       {"0371", []TokenType{OctInteger, EndOfInput}},
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
		"punctuation:1":        {".(){}[]..", []TokenType{Period, OpenParen, CloseParen, OpenBrace, CloseBrace, OpenBracket, CloseBracket, Dotdot, EndOfInput}},
		"punctuation:1/ws":     {". ( ) { } [ ] ..", []TokenType{Period, OpenParen, CloseParen, OpenBrace, CloseBrace, OpenBracket, CloseBracket, Dotdot, EndOfInput}},
		"punctuation:2":        {"+-*/%^$:", []TokenType{Plus, Minus, Star, ForwardSlash, Percent, Caret, DollarSign, Colon, EndOfInput}},
		"punctuation:2/ws":     {"+ - * / % ^ $ :", []TokenType{Plus, Minus, Star, ForwardSlash, Percent, Caret, DollarSign, Colon, EndOfInput}},
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

func TestKeywords(t *testing.T) {
	tests := map[string]struct {
		src    string
		tokens []TokenType
	}{
		"CREATE:create": {"create", []TokenType{Create, EndOfInput}},
		"CREATE:Create": {"Create", []TokenType{Create, EndOfInput}},
		"CREATE:CREATE": {"CREATE", []TokenType{Create, EndOfInput}},
		"DELETE":        {"delete", []TokenType{Delete, EndOfInput}},
		"DETACH":        {"detach", []TokenType{Detach, EndOfInput}},
		"EXISTS":        {"exists", []TokenType{Exists, EndOfInput}},
		"MATCH":         {"match", []TokenType{Match, EndOfInput}},
		"MERGE":         {"merge", []TokenType{Merge, EndOfInput}},
		"OPTIONAL":      {"optional", []TokenType{Optional, EndOfInput}},
		"REMOVE":        {"remove", []TokenType{Remove, EndOfInput}},
		"RETURN":        {"return", []TokenType{Return, EndOfInput}},
		"SET":           {"set", []TokenType{Set, EndOfInput}},
		"UNION":         {"union", []TokenType{Union, EndOfInput}},
		"UNWIND":        {"unwind", []TokenType{Unwind, EndOfInput}},
		"WITH":          {"with", []TokenType{With, EndOfInput}},
		"LIMIT":         {"limit", []TokenType{Limit, EndOfInput}},
		"ORDER":         {"order", []TokenType{Order, EndOfInput}},
		"SKIP":          {"Skip", []TokenType{Skip, EndOfInput}},
		"WHERE":         {"where", []TokenType{Where, EndOfInput}},
		"ASC":           {"asc", []TokenType{Asc, EndOfInput}},
		"ASCENDING":     {"ascending", []TokenType{Ascending, EndOfInput}},
		"BY":            {"by", []TokenType{By, EndOfInput}},
		"DESC":          {"desc", []TokenType{Desc, EndOfInput}},
		"DESCENDING":    {"descending", []TokenType{Descending, EndOfInput}},
		"ON":            {"on", []TokenType{On, EndOfInput}},
		"ALL":           {"all", []TokenType{All, EndOfInput}},
		"CASE":          {"case", []TokenType{Case, EndOfInput}},
		"ELSE":          {"else", []TokenType{Else, EndOfInput}},
		"END":           {"end", []TokenType{End, EndOfInput}},
		"THEN":          {"then", []TokenType{Then, EndOfInput}},
		"WHEN":          {"when", []TokenType{When, EndOfInput}},
		"AND":           {"and", []TokenType{And, EndOfInput}},
		"AS":            {"as", []TokenType{As, EndOfInput}},
		"CONTAINS":      {"contains", []TokenType{Contains, EndOfInput}},
		"DISTINCT":      {"distinct", []TokenType{Distinct, EndOfInput}},
		"ENDS":          {"ends", []TokenType{Ends, EndOfInput}},
		"IN":            {"in", []TokenType{In, EndOfInput}},
		"IS":            {"is", []TokenType{Is, EndOfInput}},
		"NOT":           {"not", []TokenType{Not, EndOfInput}},
		"OR":            {"or", []TokenType{Or, EndOfInput}},
		"STARTS":        {"starts", []TokenType{Starts, EndOfInput}},
		"XOR":           {"xor", []TokenType{Xor, EndOfInput}},
		"FALSE":         {"false", []TokenType{False, EndOfInput}},
		"NULL":          {"null", []TokenType{Null, EndOfInput}},
		"TRUE":          {"true", []TokenType{True, EndOfInput}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := New(bytes.NewBufferString(tc.src), newTestReporter())
			assertTokens(t, tc.tokens, s)
		})
	}
}

func TestStatement(t *testing.T) {
	tests := map[string]struct {
		src    string
		tokens []TokenType
	}{
		"statement:1": {"MATCH (e:Entity) /* with a comment */ RETURN e", []TokenType{Match, OpenParen, Identifier, Colon, Identifier, CloseParen, Return, Identifier, EndOfInput}},
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
		assert.Equal(t, tokenType, token.T)
	}
}

func assertLiteral(t *testing.T, expected []interface{}, scanner *Scanner) {
	for _, literal := range expected {
		token := scanner.NextToken()
		assert.Equal(t, literal, token.Literal)
	}
}
