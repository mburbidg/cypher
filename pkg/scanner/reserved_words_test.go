package scanner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReservedWordToken(t *testing.T) {
	token, ok := reservedWords.token("create", 10)
	assert.True(t, ok)
	assert.Equal(t, Create, token.T)
	assert.Equal(t, "CREATE", token.Lexeme)
	assert.Equal(t, 10, token.Line)

	token, ok = reservedWords.token("CrEate", 20)
	assert.True(t, ok)
	assert.Equal(t, Create, token.T)
	assert.Equal(t, "CREATE", token.Lexeme)
	assert.Equal(t, 20, token.Line)
}

func TestNotReservedWordToken(t *testing.T) {
	_, ok := reservedWords.token("creat", 10)
	assert.False(t, ok)
}
