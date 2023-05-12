package parser

import (
	"testing"
)

// We test the parser against the tkl use cases, for expressions.
// https://github.com/opencypher/openCypher/tree/master/tck/features/expressions/graph

// graph/Graph3.feature
func TestNodeLabels(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"labels(node)":    {"labels(node)", true},
		"labels(list[0])": {"labels(list[0])", true},
		"labels(null)":    {"labels(null)", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runExprTest(t, reporter, tc)
		})
	}
}

// graph/Graph5.feature
func TestNodeEdgeLabels(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"a:B":   {"a:B", true},
		"a:A:B": {"a:A:B", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runExprTest(t, reporter, tc)
		})
	}
}

// graph/Graph6.feature
func TestStaticPropertyAccess(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"(list[1]).missing": {"(list[1]).missing", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runExprTest(t, reporter, tc)
		})
	}
}

// graph/Graph7.feature
func TestDynamicPropertyAccess(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"n['nam' + 'e']": {"n['nam' + 'e']", true},
		"n[$idx]":        {"n[$idx]", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runExprTest(t, reporter, tc)
		})
	}
}

// graph/Graph8.feature
func TestPropertyKeyFunction(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"'exists' IN keys(n)": {"'exists' IN keys(n)", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runExprTest(t, reporter, tc)
		})
	}
}

// graph/Graph9.feature
func TestRetrieveAsPropertyMap(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"properties({name: 'Popeye', level: 9001})": {"properties({name: 'Popeye', level: 9001})", true},
		"properties(1)":             {"properties(1)", true},
		"properties('Cypher')":      {"properties('Cypher')", true},
		"properties([true, false])": {"properties([true, false])", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runExprTest(t, reporter, tc)
		})
	}
}
