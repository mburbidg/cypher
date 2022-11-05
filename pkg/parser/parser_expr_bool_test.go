package parser

import (
	"github.com/mburbidg/cypher/pkg/scanner"
	"github.com/stretchr/testify/assert"
	"testing"
)

// We test the parser against the tkl use cases, for expressions.
// https://github.com/opencypher/openCypher/tree/master/tck/features/expressions/boolean

// mathematical/boolean1.feature
func TestAndLogical(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"true AND true":           {"true AND true", true},
		"true AND false":          {"true AND false", true},
		"false AND true":          {"false AND true", true},
		"false AND false":         {"false AND false", true},
		"false AND null":          {"false AND null", true},
		"null AND true":           {"null AND true", true},
		"null AND false":          {"null AND false", true},
		"null AND null":           {"null AND null", true},
		"true AND true AND true":  {"true AND true AND true", true},
		"true AND true AND false": {"true AND true AND false", true},
		"true AND false AND null": {"true AND false AND null", true},
		"null AND null AND null":  {"null AND null AND null", true},
		"null AND null AND null AND null AND false AND null AND null AND null AND null AND null AND null": {"null AND null AND null AND null AND false AND null AND null AND null AND null AND null AND null", true},
		"(a AND b) = (b AND a)":                                 {"(a AND b) = (b AND a)", true},
		"(a AND b) IS NULL = (b AND a) IS NULL":                 {"(a AND b) IS NULL = (b AND a) IS NULL", true},
		"(a AND (b AND c)) = ((a AND b) AND c)":                 {"(a AND (b AND c)) = ((a AND b) AND c)", true},
		"(a AND (b AND c)) IS NULL = ((a AND b) AND c) IS NULL": {"(a AND (b AND c)) IS NULL = ((a AND b) AND c) IS NULL", true},
		"<a> AND <b>": {"<a> AND <b>", false},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runExprTest(t, reporter, tc)
		})
	}
}

// mathematical/boolean2.feature
func TestOrLogical(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"true OR true":          {"true OR true", true},
		"true OR false":         {"true OR false", true},
		"false OR true":         {"false OR true", true},
		"false OR false":        {"false OR false", true},
		"false OR null":         {"false OR null", true},
		"null OR true":          {"null OR true", true},
		"null OR false":         {"null OR false", true},
		"null OR null":          {"null OR null", true},
		"true OR true OR true":  {"true OR true OR true", true},
		"true OR true OR false": {"true OR true OR false", true},
		"true OR false OR null": {"true OR false OR null", true},
		"null OR null OR null":  {"null OR null OR null", true},
		"null OR null OR null OR null OR false OR null OR null OR null OR null OR null OR null": {"null OR null OR null OR null OR false OR null OR null OR null OR null OR null OR null", true},
		"(a OR b) = (b OR a)":                               {"(a OR b) = (b OR a)", true},
		"(a OR b) IS NULL = (b OR a) IS NULL":               {"(a OR b) IS NULL = (b OR a) IS NULL", true},
		"(a OR (b OR c)) = ((a OR b) OR c)":                 {"(a OR (b OR c)) = ((a OR b) OR c)", true},
		"(a OR (b OR c)) IS NULL = ((a OR b) OR c) IS NULL": {"(a OR (b OR c)) IS NULL = ((a OR b) OR c) IS NULL", true},
		"<a> OR <b>": {"<a> OR <b>", false},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runExprTest(t, reporter, tc)
		})
	}
}

// mathematical/boolean3.feature
func TestXorLogical(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"true XOR true":           {"true XOR true", true},
		"true XOR false":          {"true XOR false", true},
		"false XOR true":          {"false XOR true", true},
		"false XOR false":         {"false XOR false", true},
		"false XOR null":          {"false XOR null", true},
		"null XOR true":           {"null XOR true", true},
		"null XOR false":          {"null XOR false", true},
		"null XOR null":           {"null XOR null", true},
		"true XOR true XOR true":  {"true XOR true XOR true", true},
		"true XOR true XOR false": {"true XOR true XOR false", true},
		"true XOR false XOR null": {"true XOR false XOR null", true},
		"null XOR null XOR null":  {"null XOR null XOR null", true},
		"null XOR null XOR null XOR null XOR false XOR null XOR null XOR null XOR null XOR null XOR null": {"null XOR null XOR null XOR null XOR false XOR null XOR null XOR null XOR null XOR null XOR null", true},
		"(a XOR b) = (b XOR a)":                                 {"(a XOR b) = (b XOR a)", true},
		"(a XOR b) IS NULL = (b XOR a) IS NULL":                 {"(a XOR b) IS NULL = (b XOR a) IS NULL", true},
		"(a XOR (b XOR c)) = ((a XOR b) XOR c)":                 {"(a XOR (b XOR c)) = ((a XOR b) XOR c)", true},
		"(a XOR (b XOR c)) IS NULL = ((a XOR b) XOR c) IS NULL": {"(a XOR (b XOR c)) IS NULL = ((a XOR b) XOR c) IS NULL", true},
		"<a> XOR <b>": {"<a> XOR <b>", false},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runExprTest(t, reporter, tc)
		})
	}
}

// mathematical/boolean4.feature
func TestNotLogical(t *testing.T) {
	tests := map[string]struct {
		src string
	}{
		"NOT true":                      {"NOT true"},
		"NOT false":                     {"NOT false"},
		"NOT null":                      {"NOT null"},
		"NOT NOT true":                  {"NOT NOT true"},
		"NOT NOT false":                 {"NOT NOT false"},
		"NOT NOT null":                  {"NOT NOT null"},
		"NOT(n.name = 'apa' AND false)": {"NOT(n.name = 'apa' AND false)"},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := scanner.New([]byte(tc.src), reporter)
			p := New(s, reporter)
			tree, err := p.expr()
			assert.NoError(t, err)
			assert.NotNil(t, tree)
		})
	}
}

// mathematical/boolean5.feature
func TestInteropLogical(t *testing.T) {
	tests := map[string]struct {
		src string
	}{
		"(a OR (b AND c)) = ((a OR b) AND (a OR c))":                    {"(a OR (b AND c)) = ((a OR b) AND (a OR c))"},
		"(a OR (b AND c)) IS NULL = ((a OR b) AND (a OR c)) IS NULL":    {"(a OR (b AND c)) IS NULL = ((a OR b) AND (a OR c)) IS NULL"},
		"(a AND (b OR c)) = ((a AND b) OR (a AND c))":                   {"(a AND (b OR c)) = ((a AND b) OR (a AND c))"},
		"(a AND (b OR c)) IS NULL = ((a AND b) OR (a AND c)) IS NULL":   {"(a AND (b OR c)) IS NULL = ((a AND b) OR (a AND c)) IS NULL"},
		"(a AND (b XOR c)) = ((a AND b) XOR (a AND c))":                 {"(a AND (b XOR c)) = ((a AND b) XOR (a AND c))"},
		"(a AND (b XOR c)) IS NULL = ((a AND b) XOR (a AND c)) IS NULL": {"(a AND (b XOR c)) IS NULL = ((a AND b) XOR (a AND c)) IS NULL"},
		"NOT (a OR b) = (NOT (a) AND NOT (b))":                          {"NOT (a OR b) = (NOT (a) AND NOT (b))"},
		"NOT (a AND b) = (NOT (a) OR NOT (b))":                          {"NOT (a AND b) = (NOT (a) OR NOT (b))"},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := scanner.New([]byte(tc.src), reporter)
			p := New(s, reporter)
			tree, err := p.expr()
			assert.NoError(t, err)
			assert.NotNil(t, tree)
		})
	}
}
