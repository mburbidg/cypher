package parser

import (
	"github.com/mburbidg/cypher/pkg/scanner"
	"github.com/stretchr/testify/assert"
	"testing"
)

// We test the parser against the tkl use cases, for expressions.
// https://github.com/opencypher/openCypher/tree/master/tck/features/expressions/aggregation

// mathematical/aggregation1.feature
func TestAggrecationCount(t *testing.T) {
	tests := map[string]struct {
		src string
	}{
		"n.name, count(n.num)": {"n.name, count(n.num)"},
		"count(r)":             {"count(r)"},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := scanner.New([]byte(tc.src), reporter)
			p := New(s, reporter)
			tree, err := p.Parse()
			assert.NoError(t, err)
			assert.NotNil(t, tree)
		})
	}
}

// mathematical/aggregation2.feature
func TestAggrecationMinMax(t *testing.T) {
	tests := map[string]struct {
		src string
	}{
		"max(x)":                               {"max(x)"},
		"min(x)":                               {"min(x)"},
		"[1, 2, 0, null, -1]":                  {"[1, 2, 0, null, -1]"},
		"[1.0, 2.0, 0.5, null]":                {"[1.0, 2.0, 0.5, null]"},
		"[1, 2.0, 5, null, 3.2, 0.1]":          {"[1, 2.0, 5, null, 3.2, 0.1]"},
		"['a', 'b', 'B', null, 'abc', 'abc1']": {"['a', 'b', 'B', null, 'abc', 'abc1']"},
		"[[1], [2], [2, 1]]":                   {"[[1], [2], [2, 1]]"},
		"[1, 'a', null, [1, 2], 0.2, 'b']":     {"[1, 'a', null, [1, 2], 0.2, 'b']"},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := scanner.New([]byte(tc.src), reporter)
			p := New(s, reporter)
			tree, err := p.Parse()
			assert.NoError(t, err)
			assert.NotNil(t, tree)
		})
	}
}

// mathematical/aggregation3.feature
func TestAggrecationSum(t *testing.T) {
	tests := map[string]struct {
		src string
	}{
		"n.name, sum(n.num)": {"n.name, sum(n.num)"},
		"sum(i)":             {"sum(i)"},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := scanner.New([]byte(tc.src), reporter)
			p := New(s, reporter)
			tree, err := p.Parse()
			assert.NoError(t, err)
			assert.NotNil(t, tree)
		})
	}
}

// mathematical/aggregation5.feature
func TestAggrecationCollect(t *testing.T) {
	tests := map[string]struct {
		src string
	}{
		"n, collect(x)":           {"n, collect(x)"},
		"collect(DISTINCT n.num)": {"collect(DISTINCT n.num)"},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := scanner.New([]byte(tc.src), reporter)
			p := New(s, reporter)
			tree, err := p.Parse()
			assert.NoError(t, err)
			assert.NotNil(t, tree)
		})
	}
}

// mathematical/aggregation6.feature
func TestAggrecationPercentiles(t *testing.T) {
	tests := map[string]struct {
		src string
	}{
		"percentileDisc(n.price, $percentile)": {"percentileDisc(n.price, $percentile)"},
		"percentileDisc(0.90, deg)":            {"percentileDisc(0.90, deg)"},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := scanner.New([]byte(tc.src), reporter)
			p := New(s, reporter)
			tree, err := p.Parse()
			assert.NoError(t, err)
			assert.NotNil(t, tree)
		})
	}
}

// mathematical/aggregation8.feature
func TestAggrecationDistinct(t *testing.T) {
	tests := map[string]struct {
		src string
	}{
		"count(DISTINCT a)":      {"count(DISTINCT a)"},
		"count(DISTINCT a.name)": {"count(DISTINCT a.name)"},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := scanner.New([]byte(tc.src), reporter)
			p := New(s, reporter)
			tree, err := p.Parse()
			assert.NoError(t, err)
			assert.NotNil(t, tree)
		})
	}
}
