package parser

import (
	"testing"
)

// We test the parser against the tkl use cases, for expressions.
// https://github.com/opencypher/openCypher/tree/master/tck/features/expressions/pattern

// mathematical/pattern1.feature
func TestPatternPredicate(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"(n)-[]->()":                        {"(n)-[]->()", true},
		"(n)-[]-()":                         {"(n)-[]-()", true},
		"(n)<-[]-()":                        {"(n)<-[]-()", true},
		"(n)-[:REL1]->()":                   {"(n)-[:REL1]->()", true},
		"(n)-[:REL1]-()":                    {"(n)-[:REL1]-()", true},
		"(n)<-[:REL1]-()":                   {"(n)<-[:REL1]-()", true},
		"(n)-[:REL1*]->()":                  {"(n)-[:REL1*]->()", true},
		"(n)-[:REL1*]-()":                   {"(n)-[:REL1*]-()", true},
		"(n)<-[:REL1*]-()":                  {"(n)<-[:REL1*]-()", true},
		"(n)-[:REL1*2]-()":                  {"(n)-[:REL1*2]-()", true},
		"(n)-[]->(m)":                       {"(n)-[]->(m)", true},
		"(n)-[:REL1|REL2|REL3|REL4]-(m)":    {"(n)-[:REL1|REL2|REL3|REL4]-(m)", true},
		"(n)-[:REL1]->(m)":                  {"(n)-[:REL1]->(m)", true},
		"(n)-[:REL1]-(m)":                   {"(n)-[:REL1]-(m)", true},
		"(n)-[:REL1*]->(m)":                 {"(n)-[:REL1*]->(m)", true},
		"(n)-[:REL1*]-(m)":                  {"(n)-[:REL1*]-(m)", true},
		"(n)-[:REL1*2]-(m)":                 {"(n)-[:REL1*2]-(m)", true},
		"(n)-[:REL2]-()":                    {"(n)-[:REL2]-()", true},
		"(n)-[:REL1]-() AND (n)-[:REL3]-()": {"(n)-[:REL1]-() AND (n)-[:REL3]-()", true},
		"(n)-[:REL1]-() OR (n)-[:REL2]-()":  {"(n)-[:REL1]-() OR (n)-[:REL2]-()", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, reporter, tc)
		})
	}
}

// mathematical/pattern2.feature
func TestPatternComprehension(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"[p = (n)-->() | p]":                       {"[p = (n)-->() | p]", true},
		"[p = (n)-->(:B) | p]":                     {"[p = (n)-->(:B) | p]", true},
		"[p = (a)-->(b) | p]":                      {"[p = (a)-->(b) | p]", true},
		"[(n)-[:T]->(b) | b.name]":                 {"[(n)-[:T]->(b) | b.name]", true},
		"[(n)-[r:T]->() | r.name]":                 {"[(n)-[r:T]->() | r.name]", true},
		"count([p = (n)-[:HAS]->() | p])":          {"count([p = (n)-[:HAS]->() | p])", true},
		"[x IN nodes(p) | size([(x)-->(:Y) | 1])]": {"[x IN nodes(p) | size([(x)-->(:Y) | 1])]", true},
		"[p = (n)-[:HAS]->() | p]":                 {"[p = (n)-[:HAS]->() | p]", true},
		"[p = (a)-[*]->(b) | p]":                   {"[p = (a)-[*]->(b) | p]", true},
		"[p = (liker)--() | p]":                    {"[p = (liker)--() | p]", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runTest(t, reporter, tc)
		})
	}
}
