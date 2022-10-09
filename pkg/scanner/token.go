package scanner

import (
	"log"
	"strconv"
)

type Token struct {
	T       TokenType
	Lexeme  string
	Literal any
	Line    int
}

var endOfInputToken = Token{
	T: EndOfInput,
}

func newOperatorToken(t TokenType, line int) Token {
	return Token{
		T:    t,
		Line: line,
	}
}

func newKeywordToken(t TokenType, lexeme string, line int) Token {
	return Token{
		T:      t,
		Lexeme: lexeme,
		Line:   line,
	}
}

func newEndOfInputToken() Token {
	return Token{
		T: EndOfInput,
	}
}

func newIntegerToken(t TokenType, lexeme string, base int, line int) Token {
	n, err := strconv.ParseInt(lexeme, base, 64)
	log.Printf("err=%s\n", err)
	return Token{
		T:       t,
		Lexeme:  lexeme,
		Literal: n,
		Line:    line,
	}
}

func newDoubleToken(lexeme string, line int) Token {
	f, _ := strconv.ParseFloat(lexeme, 64)
	return Token{
		T:       Double,
		Lexeme:  lexeme,
		Literal: f,
		Line:    line,
	}
}

func newStringToken(lexeme string, literal string, line int) Token {
	return Token{
		T:       String,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func newIdentifierToken(lexeme string, line int) Token {
	return Token{
		T:      Identifier,
		Lexeme: lexeme,
		Line:   line,
	}
}

func newIllegalToken(lexeme string) Token {
	return Token{
		Lexeme: lexeme,
		T:      Illegal,
	}
}
