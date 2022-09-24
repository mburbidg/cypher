package scanner

import "strconv"

type Token struct {
	t       TokenType
	lexeme  string
	literal any
	line    int
}

func newOperatorToken(t TokenType, line int) Token {
	return Token{
		t:    t,
		line: line,
	}
}

func newKeywordToken(t TokenType, lexeme string, line int) Token {
	return Token{
		t:      t,
		lexeme: lexeme,
		line:   line,
	}
}

func newEndOfInputToken() Token {
	return Token{
		t: EndOfInput,
	}
}

func newNumberToken(t TokenType, lexeme string, line int) Token {
	if t == Double {
		f, _ := strconv.ParseFloat(lexeme, 64)
		return Token{
			t:       t,
			lexeme:  lexeme,
			literal: f,
			line:    line,
		}
	}
	if t == Integer {
		n, _ := strconv.ParseInt(lexeme, 10, 32)
		return Token{
			t:       t,
			lexeme:  lexeme,
			literal: n,
			line:    line,
		}
	}
	return Token{}
}

func newErrorToken() Token {
	return Token{
		t: Error,
	}
}
