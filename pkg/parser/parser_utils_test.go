package parser

import (
	"github.com/mburbidg/cypher/pkg/scanner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func runTest(t *testing.T, reporter *testReporter, tc struct {
	src   string
	valid bool
}) {
	s := scanner.New([]byte(tc.src), reporter)
	p := New(s, reporter)
	tree, err := p.Parse()
	if tc.valid {
		assert.NoError(t, err)
		assert.NotNil(t, tree)
	} else {
		assert.Error(t, err)
	}
}
