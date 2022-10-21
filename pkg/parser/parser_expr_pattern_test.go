package parser

import (
	"github.com/mburbidg/cypher/pkg/scanner"
	"github.com/stretchr/testify/assert"
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
