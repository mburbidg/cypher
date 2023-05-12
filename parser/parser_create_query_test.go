package parser

import (
	"testing"
)

// We test the parser against the tkl use cases, for expressions.
// https://github.com/opencypher/openCypher/tree/master/tck/features/clauses/create

// create/Create1.feature
func TestCreatingNodes(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"CREATE ()":                                   {"CREATE ()", true},
		"CREATE (), ()":                               {"CREATE (), ()", true},
		"CREATE (:Label)":                             {"CREATE (:Label)", true},
		"CREATE (:Label), (:Label)":                   {"CREATE (:Label), (:Label)", true},
		"CREATE (:A:B:C:D)":                           {"CREATE (:A:B:C:D)", true},
		"CREATE (:B:A:D), (:B:C), (:D:E:B)":           {"CREATE (:B:A:D), (:B:C), (:D:E:B)", true},
		"CREATE ({created: true})":                    {"CREATE ({created: true})", true},
		"CREATE (n {name: 'foo'})":                    {"CREATE (n {name: 'foo'})", true},
		"CREATE (n {id: 12, name: 'foo'})":            {"CREATE (n {id: 12, name: 'foo'})", true},
		"CREATE (n {id: 12, name: null})":             {"CREATE (n {id: 12, name: null})", true},
		"CREATE (n {name: 'foo'}) RETURN n.name AS p": {"CREATE (n {name: 'foo'}) RETURN n.name AS p", true},
		"CREATE (n {id: 12, name: 'foo'}) RETURN n.id AS id, n.name AS p": {"CREATE (n {id: 12, name: 'foo'}) RETURN n.id AS id, n.name AS p", true},
		"CREATE (n {id: 12, name: null}) RETURN n.id AS id, n.name AS p":  {"CREATE (n {id: 12, name: null}) RETURN n.id AS id, n.name AS p", true},
		"CREATE (p:TheLabel {id: 4611686018427387905}) RETURN p.id":       {"CREATE (p:TheLabel {id: 4611686018427387905}) RETURN p.id", true},
		"MATCH (a) CREATE (a)":                          {"MATCH (a) CREATE (a)", true},
		"MATCH (a) CREATE (a {name: 'foo'}) RETURN a":   {"MATCH (a) CREATE (a {name: 'foo'}) RETURN a", true},
		"CREATE (n:Foo)-[:T1]->(), (n:Bar)-[:T2]->()":   {"CREATE (n:Foo)-[:T1]->(), (n:Bar)-[:T2]->()", true},
		"CREATE ()<-[:T2]-(n:Foo), (n:Bar)<-[:T1]-()":   {"CREATE ()<-[:T2]-(n:Foo), (n:Bar)<-[:T1]-()", true},
		"CREATE (n:Foo) CREATE (n:Bar)-[:OWNS]->(:Dog)": {"CREATE (n:Foo) CREATE (n:Bar)-[:OWNS]->(:Dog)", true},
		"CREATE (n {}) CREATE (n:Bar)-[:OWNS]->(:Dog)":  {"CREATE (n {}) CREATE (n:Bar)-[:OWNS]->(:Dog)", true},
		"CREATE (n:Foo) CREATE (n {})-[:OWNS]->(:Dog)":  {"CREATE (n:Foo) CREATE (n {})-[:OWNS]->(:Dog)", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runQueryTest(t, reporter, tc)
		})
	}
}
