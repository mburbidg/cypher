package parser

import (
	"github.com/mburbidg/cypher/pkg/scanner"
	"github.com/stretchr/testify/assert"
	"testing"
)

// We test the parser against the tkl use cases, for expressions.
// https://github.com/opencypher/openCypher/tree/master/tck/features/expressions/comparison

// mathematical/comparison1.feature
func TestComparisonEquality(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"toInteger(n.id) = expected": {"toInteger(n.id) = expected", true},
		"a = b":                      {"a = b", true},
		"1 = 1":                      {"1 = 1", true},
		"null = null":                {"null = null", true},
		"null <> null":               {"null <> null", true},
		"a <> b":                     {"a <> b", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := scanner.New([]byte(tc.src), reporter)
			p := New(s, reporter)
			tree, err := p.Parse()
			if tc.valid {
				assert.NoError(t, err)
				assert.NotNil(t, tree)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// mathematical/comparison2.feature
func TestHalfBoundedRange(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"i.var IS NOT NULL AND i.var > 'x'": {"i.var IS NOT NULL AND i.var > 'x'", true},
		"i.var IS NULL OR i.var > 'x'":      {"i.var IS NULL OR i.var > 'x'", true},
		"i <> j":                            {"i <> j", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := scanner.New([]byte(tc.src), reporter)
			p := New(s, reporter)
			tree, err := p.Parse()
			if tc.valid {
				assert.NoError(t, err)
				assert.NotNil(t, tree)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// mathematical/comparison3.feature
func TestFullBoundedRange(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"1 < n.num < 3":       {"1 < n.num < 3", true},
		"1 < n.num <= 3":      {"1 < n.num <= 3", true},
		"1 <= n.num < 3":      {"1 <= n.num < 3", true},
		"1 <= n.num <= 3":     {"1 <= n.num <= 3", true},
		"'a' < n.name < 'c'":  {"'a' < n.name < 'c'", true},
		"'a' < n.name <= 'c'": {"'a' < n.name <= 'c'", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := scanner.New([]byte(tc.src), reporter)
			p := New(s, reporter)
			tree, err := p.Parse()
			if tc.valid {
				assert.NoError(t, err)
				assert.NotNil(t, tree)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// mathematical/comparison4.feature
func TestComparisonComparison(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"n.prop1 < m.prop1 = n.prop2 <> m.prop2": {"n.prop1 < m.prop1 = n.prop2 <> m.prop2", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := scanner.New([]byte(tc.src), reporter)
			p := New(s, reporter)
			tree, err := p.Parse()
			if tc.valid {
				assert.NoError(t, err)
				assert.NotNil(t, tree)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
