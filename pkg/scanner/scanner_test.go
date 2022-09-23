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
	assertTokens(t, []TokenType{Create, EndOfInput}, s)

	s = New(bytes.NewBufferString("+"))
	assertTokens(t, []TokenType{Plus, EndOfInput}, s)

	s = New(bytes.NewBufferString("MATCH (n) RETURN n WHERE n.foo = 1"))
	assertTokens(t, []TokenType{Match, OpenParen, SymbolName, CloseParen, Return, SymbolName, Where, SymbolName, Period, SymbolName, Equal, Integer, EndOfInput}, s)
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

func assertTokens(t *testing.T, expected []TokenType, scanner *Scanner) {
	for _, tokenType := range expected {
		token := scanner.NextToken()
		assert.Equal(t, tokenType, token.t)
	}
}

var query = "MATCH (n) RETURN n WHERE n.foo = 1"
var expected = []TokenType{Match, OpenParen, SymbolName, CloseParen, Return, SymbolName, Where, SymbolName, Period, SymbolName, Equal, Integer}
