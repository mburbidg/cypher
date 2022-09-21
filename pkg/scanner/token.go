package scanner

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
