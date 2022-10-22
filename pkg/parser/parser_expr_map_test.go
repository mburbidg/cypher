package parser

import (
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
			runTest(t, reporter, tc)
		})
	}
}
