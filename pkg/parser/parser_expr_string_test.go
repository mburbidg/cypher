package parser

import (
	"testing"
)

// We test the parser against the tkl use cases, for expressions.
// https://github.com/opencypher/openCypher/tree/master/tck/features/expressions/string

// string/String8.feature
func TestStringPrefixMatch(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"a.name STARTS WITH 'ABCDEF'": {"a.name STARTS WITH 'ABCDEF'", true},
		"a.name STARTS WITH ''":       {"a.name STARTS WITH ''", true},
		"a.name STARTS WITH ' '":      {"a.name STARTS WITH ' '", true},
		"a.name STARTS WITH null":     {"a.name STARTS WITH null", true},
		"NOT a.name STARTS WITH null": {"NOT a.name STARTS WITH null", true},
		"op1 STARTS WITH op2":         {"op1 STARTS WITH op2", true},
		"NOT a.name STARTS WITH 'ab'": {"NOT a.name STARTS WITH 'ab'", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, reporter, tc)
		})
	}
}

// string/String9.feature
func TestStringSuffixMatch(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"a.name ENDS WITH 'AB'":      {"a.name ENDS WITH 'AB'", true},
		"a.name ENDS WITH ''":        {"a.name ENDS WITH ''", true},
		"a.name ENDS WITH ' '":       {"a.name ENDS WITH ' '", true},
		"a.name ENDS WITH null":      {"a.name ENDS WITH null", true},
		"op1 ENDS WITH op2":          {"op1 ENDS WITH op2", true},
		"NOT a.name ENDS WITH 'def'": {"NOT a.name ENDS WITH 'def'", true},
		//"xxxx":                        {"xxxx", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, reporter, tc)
		})
	}
}

// string/String10.feature
func TestSubstringMatch(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"a.name CONTAINS 'ABCDEF'": {"a.name CONTAINS 'ABCDEF'", true},
		"a.name CONTAINS ''":       {"a.name CONTAINS ''", true},
		"a.name CONTAINS ' '":      {"a.name CONTAINS ' '", true},
		"a.name CONTAINS null":     {"a.name CONTAINS null", true},
		"NOT a.name CONTAINS null": {"NOT a.name CONTAINS null", true},
		"op1 CONTAINS op2":         {"op1 CONTAINS op2", true},
		"NOT a.name CONTAINS 'b'":  {"NOT a.name CONTAINS 'b'", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, reporter, tc)
		})
	}
}

// string/String11.feature
func TestStringMatch(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"a.name STARTS WITH 'a' AND a.name ENDS WITH 'f'":                          {"a.name STARTS WITH 'a' AND a.name ENDS WITH 'f'", true},
		"a.name STARTS WITH 'A' AND a.name CONTAINS 'C' AND a.name ENDS WITH 'EF'": {"a.name STARTS WITH 'A' AND a.name CONTAINS 'C' AND a.name ENDS WITH 'EF'", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, reporter, tc)
		})
	}
}
