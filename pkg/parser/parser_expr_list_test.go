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
		"foo()":            {"foo()", true},
		"list[$from..$to]": {"list[$from..$to]", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, reporter, tc)
		})
	}
}

// list/List3.feature
func TestListEquality(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"[1, 2] = 'foo'":                                     {"[1, 2] = 'foo'", true},
		"[1] = [1, null]":                                    {"[1] = [1, null]", true},
		"[1, 2] = [null, 'foo']":                             {"[1, 2] = [null, 'foo']", true},
		"[1, 2] = [null, 2]":                                 {"[1, 2] = [null, 2]", true},
		"[[1]] = [[1], [null]]":                              {"[[1]] = [[1], [null]]", true},
		"[[1, 2], [1, 3]] = [[1, 2], [null, 'foo']]":         {"[[1, 2], [1, 3]] = [[1, 2], [null, 'foo']]", true},
		"[[1, 2], ['foo', 'bar']] = [[1, 2], [null, 'bar']]": {"[[1, 2], ['foo', 'bar']] = [[1, 2], [null, 'bar']]", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, reporter, tc)
		})
	}
}

// list/List4.feature
func TestListConcatenation(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"[1, 10, 100] + [4, 5]": {"[1, 10, 100] + [4, 5]", true},
		"[false, true] + false": {"[false, true] + false", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, reporter, tc)
		})
	}
}

// list/List5.feature
func TestListMembershipValidation(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"3 IN list[0]":                                    {"3 IN list[0]", true},
		"3 IN [[1, 2, 3]][0]":                             {"3 IN [[1, 2, 3]][0]", true},
		"3 IN list[0..1]":                                 {"3 IN list[0..1]", true},
		"3 IN [1, 2, 3][0..1]":                            {"3 IN [1, 2, 3][0..1]", true},
		"1 IN ['1', 2]":                                   {"1 IN ['1', 2]", true},
		"[1, 2] IN [1, [1, '2']]":                         {"[1, 2] IN [1, [1, '2']]", true},
		"[1] IN [1, 2]":                                   {"[1] IN [1, 2]", true},
		"[1, 2] IN [1, 2]":                                {"[1, 2] IN [1, 2]", true},
		"[1] IN [1, 2, [1]]":                              {"[1] IN [1, 2, [1]]", true},
		"[1, 2] IN [1, [1, 2]]":                           {"[1, 2] IN [1, [1, 2]]", true},
		"[1, 2] IN [1, [2, 1]]":                           {"[1, 2] IN [1, [2, 1]]", true},
		"[1, 2] IN [1, [1, 2, 3]]":                        {"[1, 2] IN [1, [1, 2, 3]]", true},
		"[1, 2] IN [1, [[1, 2]]]":                         {"[1, 2] IN [1, [[1, 2]]]", true},
		"[[1, 2], [3, 4]] IN [5, [[1, 2], [3, 4]]]":       {"[[1, 2], [3, 4]] IN [5, [[1, 2], [3, 4]]]", true},
		"[[1, 2], 3] IN [1, [[1, 2], 3]]":                 {"[[1, 2], 3] IN [1, [[1, 2], 3]]", true},
		"[[1]] IN [2, [[1]]]":                             {"[[1]] IN [2, [[1]]]", true},
		"[[1, 3]] IN [2, [[1, 3]]]":                       {"[[1, 3]] IN [2, [[1, 3]]]", true},
		"[[1]] IN [2, [1]]":                               {"[[1]] IN [2, [1]]", true},
		"[[1, 3]] IN [2, [1, 3]]":                         {"[[1, 3]] IN [2, [1, 3]]", true},
		"null IN [null]":                                  {"null IN [null]", true},
		"[null] IN [[null]]":                              {"[null] IN [[null]]", true},
		"[null] IN [null]":                                {"[null] IN [null]", true},
		"[1] IN [[1, null]]":                              {"[1] IN [[1, null]]", true},
		"3 IN [1, null, 3]":                               {"3 IN [1, null, 3]", true},
		"4 IN [1, null, 3]":                               {"4 IN [1, null, 3]", true},
		"[1, 2] IN [[null, 'foo'], [1, 2]]":               {"[1, 2] IN [[null, 'foo'], [1, 2]]", true},
		"[1, 2] IN [1, [1, 2], null]":                     {"[1, 2] IN [1, [1, 2], null]", true},
		"[1, 2] IN [[null, 'foo']]":                       {"[1, 2] IN [[null, 'foo']]", true},
		"[1, 2] IN [[null, 2]]":                           {"[1, 2] IN [[null, 2]]", true},
		"[1, 2] IN [1, [1, 2, null]]":                     {"[1, 2] IN [1, [1, 2, null]]", true},
		"[1, 2, null] IN [1, [1, 2, null]]":               {"[1, 2, null] IN [1, [1, 2, null]]", true},
		"[1, 2] IN [[null, 2], [1, 2]]":                   {"[1, 2] IN [[null, 2], [1, 2]]", true},
		"[[1, 2], [3, 4]] IN [5, [[1, 2], [3, 4], null]]": {"[[1, 2], [3, 4]] IN [5, [[1, 2], [3, 4], null]]", true},
		"[1, 2] IN [[null, 2], [1, 3]]":                   {"[1, 2] IN [[null, 2], [1, 3]]", true},
		"[] IN [[]]":                                      {"[] IN [[]]", true},
		"[] IN []":                                        {"[] IN []", true},
		"[] IN [1, []]":                                   {"[] IN [1, []]", true},
		"[] IN [1, 2]":                                    {"[] IN [1, 2]", true},
		"[[]] IN [1, [[]]]":                               {"[[]] IN [1, [[]]]", true},
		"[] IN [1, 2, null]":                              {"[] IN [1, 2, null]", true},
		"[[], []] IN [1, [[], []]]":                       {"[[], []] IN [1, [[], []]]", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, reporter, tc)
		})
	}
}

// list/List6.feature
func TestListSize(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"size([1, 2, 3])":                {"size([1, 2, 3])", true},
		"size(n.numbers)":                {"size(n.numbers)", true},
		"size([[], []] + [[]])":          {"size([[], []] + [[]])", true},
		"size(null)":                     {"size(null)", true},
		"size([(n)--() | 1]) > 0":        {"size([(n)--() | 1]) > 0", true},
		"size([(a)-->() | 1])":           {"size([(a)-->() | 1])", true},
		"size([(a)-[:T]->() | 1])":       {"size([(a)-[:T]->() | 1])", true},
		"size([(a)-[:T|OTHER]->() | 1])": {"size([(a)-[:T|OTHER]->() | 1])", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, reporter, tc)
		})
	}
}
