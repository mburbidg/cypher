package parser

import (
	"testing"
)

// We test the parser against the tkl use cases, for expressions.
// https://github.com/opencypher/openCypher/tree/master/tck/features/expressions/quantifier

// quantifier/Quantifier1.feature
func TestNoneQuantifier(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"none(x IN [] WHERE true)":                 {"none(x IN [] WHERE true)", true},
		"none(x IN nodes WHERE x.name = 'a')":      {"none(x IN nodes WHERE x.name = 'a')", true},
		"none(x IN [1, 2, 3] WHERE x IS NOT NULL)": {"none(x IN [1, 2, 3] WHERE x IS NOT NULL)", true},
		//"xxxx":                     {"xxxx", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runExprTest(t, reporter, tc)
		})
	}
}

// quantifier/Quantifier2.feature
func TestSingleQuantifier(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"single(x IN [] WHERE true)":                 {"single(x IN [] WHERE true)", true},
		"single(x IN nodes WHERE x.name = 'a')":      {"single(x IN nodes WHERE x.name = 'a')", true},
		"single(x IN [1, 2, 3] WHERE x IS NOT NULL)": {"single(x IN [1, 2, 3] WHERE x IS NOT NULL)", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runExprTest(t, reporter, tc)
		})
	}
}

// quantifier/Quantifier3.feature
func TestAnyQuantifier(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"any(x IN [] WHERE true)":                                  {"any(x IN [] WHERE true)", true},
		"any(x IN [] WHERE false)":                                 {"any(x IN [] WHERE false)", true},
		"any(x IN [] WHERE x)":                                     {"any(x IN [] WHERE x)", true},
		"any(x IN nodes WHERE x.name = 'a')":                       {"any(x IN nodes WHERE x.name = 'a')", true},
		"any(x IN [1, 2, 3] WHERE x IS NOT NULL)":                  {"any(x IN [1, 2, 3] WHERE x IS NOT NULL)", true},
		"any(x IN [1, 2, 3] WHERE x <> 3)":                         {"any(x IN [1, 2, 3] WHERE x <> 3)", true},
		"any(x IN [1, null, true, 4.5, 'abc', false] WHERE false)": {"any(x IN [1, null, true, 4.5, 'abc', false] WHERE false)", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runExprTest(t, reporter, tc)
		})
	}
}

// quantifier/Quantifier4.feature
func TestAllQuantifier(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"all(x IN [] WHERE true)":  {"all(x IN [] WHERE true)", true},
		"all(x IN [] WHERE false)": {"all(x IN [] WHERE false)", true},
		"all(x IN [] WHERE x)":     {"all(x IN [] WHERE x)", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runExprTest(t, reporter, tc)
		})
	}
}

// quantifier/Quantifier5.feature
func TestNoneQuantifierInterop(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"none(x IN [['abc'], ['abc', 'def']] WHERE true)":                                {"none(x IN [['abc'], ['abc', 'def']] WHERE true)", true},
		"none(x IN [1, 2, 3] WHERE x = y) = (NOT any(x IN [1, 2, 3] WHERE x <= z))":      {"none(x IN [1, 2, 3] WHERE x = y) = (NOT any(x IN [1, 2, 3] WHERE x <= z))", true},
		"none(x IN [1, 2, 3] WHERE x = y) = all(x IN [1, 2, 3] WHERE NOT (true))":        {"none(x IN [1, 2, 3] WHERE x = y) = all(x IN [1, 2, 3] WHERE NOT (true))", true},
		"none(x IN [1, 2, 3] WHERE true) = (size([x IN [1, 2, 3] WHERE false | x]) = 0)": {"none(x IN [1, 2, 3] WHERE true) = (size([x IN [1, 2, 3] WHERE false | x]) = 0)", true},
		//"xxxx":                     {"xxxx", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runExprTest(t, reporter, tc)
		})
	}
}
