package parser

import (
	"github.com/mburbidg/cypher/pkg/scanner"
	"github.com/stretchr/testify/assert"
	"testing"
)

// We test the parser against the tkl use cases, for expressions.
// https://github.com/opencypher/openCypher/tree/master/tck/features/expressions/list

// list/List1.feature
func TestDynamicElementAccess(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"[1, 2, 3][0]":         {"[1, 2, 3][0]", true},
		"[[1]][0][0]":          {"[[1]][0][0]", true},
		"expr[idx]":            {"expr[idx]", true},
		"expr[$idx]":           {"expr[$idx]", true},
		"expr[toInteger(idx)]": {"expr[toInteger(idx)]", true},
		"xxxx":                 {"xxxx", true},
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

// list/List2.feature
func TestListSlicing(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"list[1..3]":       {"list[1..3]", true},
		"list[1..]":        {"list[1..]", true},
		"list[..2]":        {"list[..2]", true},
		"list[-3..-1]":     {"list[-3..-1]", true},
		"foo(1)":           {"foo(1)", true},
		"":                 {"", true},
		"list[$from..$to]": {"list[$from..$to]", true},
		"xxxx":             {"xxxx", true},
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
