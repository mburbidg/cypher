package scanner

import "strconv"

type Token struct {
	t       TokenType
	lexeme  string
	literal any
	line    int
}

var endOfInputToken = Token{
	t: EndOfInput,
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

func newIntegerToken(lexeme string, base int, line int) Token {
	n, _ := strconv.ParseInt(lexeme, base, 64)
	return Token{
		t:       Integer,
		lexeme:  lexeme,
		literal: n,
		line:    line,
	}
}

func newDoubleToken(lexeme string, line int) Token {
	f, _ := strconv.ParseFloat(lexeme, 64)
	return Token{
		t:       Double,
		lexeme:  lexeme,
		literal: f,
		line:    line,
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

func newStringToken(lexeme string, literal string, line int) Token {
	return Token{
		t:       String,
		lexeme:  lexeme,
		literal: literal,
		line:    line,
	}
}

func newIdentifierToken(lexeme string, line int) Token {
	return Token{
		t:      Identifier,
		lexeme: lexeme,
		line:   line,
	}
}

func newIllegalToken(lexeme string) Token {
	return Token{
		lexeme: lexeme,
		t:      Illegal,
	}
}
