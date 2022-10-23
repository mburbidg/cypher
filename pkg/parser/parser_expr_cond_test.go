package parser

import (
	"testing"
)

// We test the parser against the tkl use cases, for expressions.
// https://github.com/opencypher/openCypher/tree/master/tck/features/expressions/conditional

// mathematical/conditional1.feature
func TestCoalesce(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"coalesce(a.title, a.name)": {"coalesce(a.title, a.name)", true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runExprTest(t, reporter, tc)
		})
	}
}

// mathematical/conditional2.feature
func TestCase(t *testing.T) {
	tests := map[string]struct {
		src   string
		valid bool
	}{
		"case": {`CASE a
						WHEN -10 THEN 'minus ten'
						WHEN 0 THEN 'zero'
						WHEN 1 THEN 'one'
						WHEN 5 THEN 'five'
						WHEN 10 THEN 'ten'
						WHEN 3000 THEN 'three thousand'
						ELSE 'something else'
					   END`, true},
	}
	reporter := newTestReporter()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runExprTest(t, reporter, tc)
		})
	}
}
