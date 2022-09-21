package scanner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOperatorToken(t *testing.T) {
	token := newOperatorToken(Plus, 10)
	assert.Equal(t, Plus, token.t)
	assert.Empty(t, token.lexeme)
	assert.Empty(t, token.literal)
	assert.Equal(t, 10, token.line)
}

func TestNewKeywordToken(t *testing.T) {
	token := newKeywordToken(Create, "create", 20)
	assert.Equal(t, Create, token.t)
	assert.Equal(t, "create", token.lexeme)
	assert.Empty(t, token.literal)
	assert.Empty(t, 20, token.line)
}
