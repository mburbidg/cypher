package scanner

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScanner(t *testing.T) {
	s := New(bytes.NewBufferString("Create"))
	token, err := s.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, Create, token.t)

	s = New(bytes.NewBufferString("+"))
	token, err = s.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, Plus, token.t)

	s = New(bytes.NewBufferString("MATCH (n) RETURN n"))
	token, err = s.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, Match, token.t)
	token, err = s.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, WhiteSpace, token.t)
	token, err = s.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, OpenParen, token.t)
	token, err = s.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, SymbolName, token.t)
	assert.Equal(t, "n", token.lexeme)
	token, err = s.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, CloseParen, token.t)
	token, err = s.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, WhiteSpace, token.t)
	token, err = s.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, Return, token.t)
	token, err = s.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, WhiteSpace, token.t)
	token, err = s.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, SymbolName, token.t)
	assert.Equal(t, "n", token.lexeme)
	token, err = s.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, EndOfInput, token.t)

}

func TestIt(t *testing.T) {
	buf := bytes.NewBufferString("Create")
	r := bufio.NewReader(buf)
	ch, _, err := r.ReadRune()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ch)
	ch, _, err = r.ReadRune()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ch)
	ch, _, err = r.ReadRune()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ch)
	ch, _, err = r.ReadRune()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ch)
	ch, _, err = r.ReadRune()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ch)
	ch, _, err = r.ReadRune()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ch)
	ch, _, err = r.ReadRune()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ch)
	err = r.UnreadRune()
	if err != nil {
		fmt.Println(err)
	}
}
