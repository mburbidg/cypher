package scanner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOperatorToken(t *testing.T) {
	token := newOperatorToken(Plus, 10)
	assert.Equal(t, Plus, token.T)
	assert.Empty(t, token.Lexeme)
	assert.Empty(t, token.Literal)
	assert.Equal(t, 10, token.Line)
}

func TestNewKeywordToken(t *testing.T) {
	token := newKeywordToken(Create, "create", 20)
	assert.Equal(t, Create, token.T)
	assert.Equal(t, "create", token.Lexeme)
	assert.Empty(t, token.Literal)
	assert.Equal(t, 20, token.Line)
}
