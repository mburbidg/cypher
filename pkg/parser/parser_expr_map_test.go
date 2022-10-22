package parser

import (
	"github.com/mburbidg/cypher/pkg/scanner"
	"github.com/stretchr/testify/assert"
	"testing"
)

// We test the parser against the tkl use cases, for expressions.
// https://github.com/opencypher/openCypher/tree/master/tck/features/expressions/map

// map/Map1.feature
func TestMapKeysFunction(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"keys({name: 'Alice', age: 38, address: {city: 'London', residential: true}})": {"keys({name: 'Alice', age: 38, address: {city: 'London', residential: true}})", true},
		"xxxx": {"xxxx", true},
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
