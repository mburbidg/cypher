package scanner

import (
	"bufio"
	"github.com/smasher164/xid"
	"io"
	"strings"
)

type Scanner struct {
	input *bufio.Reader
	line  int
	eof   bool
}

func New(input io.Reader) *Scanner {
	return &Scanner{
		input: bufio.NewReader(input),
		line:  1,
	}
}

func (s *Scanner) NextToken() (Token, error) {
	if s.eof == true {
		return newEndOfInputToken(), nil
	}

	ch, err := s.readRune()
	if err == io.EOF {
		return newEndOfInputToken(), nil
	}
	if err != nil {
		return Token{}, err
	}

	switch ch {
	case '.':
		return newOperatorToken(Period, s.line), nil
	case '(':
		return newOperatorToken(OpenParen, s.line), nil
	case ')':
		return newOperatorToken(CloseParen, s.line), nil
	case '{':
		return newOperatorToken(OpenBrace, s.line), nil
	case '}':
		return newOperatorToken(CloseBrace, s.line), nil
	case '[':
		return newOperatorToken(OpenBracket, s.line), nil
	case ']':
		return newOperatorToken(CloseBracket, s.line), nil
	case '+':
		return newOperatorToken(Plus, s.line), nil
	case '-':
		return newOperatorToken(Minus, s.line), nil
	case '*':
		return newOperatorToken(Star, s.line), nil
	case '/':
		return newOperatorToken(ForwardSlash, s.line), nil
	case '%':
		return newOperatorToken(Percent, s.line), nil
	case '^':
		return newOperatorToken(Caret, s.line), nil
	case '=':
		return newOperatorToken(Equal, s.line), nil
	case '<':
		match, err := s.matchNext('>')
		if err != nil {
			return Token{}, err
		}
		if match {
			return newOperatorToken(NotEqual, s.line), nil
		}
		match, err = s.matchNext('=')
		if err != nil {
			return Token{}, err
		}
		if match {
			return newOperatorToken(LessThanOrEqual, s.line), nil
		}
		return newOperatorToken(LessThan, s.line), nil
	case '>':
		match, err := s.matchNext('=')
		if err != nil {
			return Token{}, err
		}
		if match {
			return newOperatorToken(GreaterThanOrEqual, s.line), nil
		}
		return newOperatorToken(GreaterThan, s.line), nil
	case '$':
		return newOperatorToken(DollarSign, s.line), nil
	}

	if xid.Start(ch) {
		return s.scanSymbolicName(ch)
	}

	if space(ch) {
		return s.scanSpace(ch)
	}

	return Token{}, nil
}

func (s *Scanner) matchNext(r rune) (bool, error) {
	ch, _, err := s.input.ReadRune()
	if err != nil {
		return false, err
	}
	if ch == r {
		return true, nil
	}
	err = s.input.UnreadRune()
	if err != nil {
		return false, err
	}
	return false, nil
}

func (s *Scanner) scanSymbolicName(ch rune) (Token, error) {
	b := strings.Builder{}
	b.WriteRune(ch)
	for {
		ch, err := s.readRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return Token{}, err
		}
		if xid.Continue(ch) {
			b.WriteRune(ch)
		} else {
			s.input.UnreadRune()
			break
		}
	}
	if token, ok := reservedWords.token(b.String(), s.line); ok {
		return token, nil
	}
	return Token{
		t:       SymbolName,
		lexeme:  b.String(),
		literal: nil,
		line:    s.line,
	}, nil
}

func (s *Scanner) scanSpace(ch rune) (Token, error) {
	b := strings.Builder{}
	b.WriteRune(ch)
	if ch == '\n' {
		s.line = s.line + 1
	}
	for {
		ch, err := s.readRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return Token{}, err
		}
		if space(ch) {
			b.WriteRune(ch)
			if ch == '\n' {
				s.line = s.line + 1
			}
		} else {
			s.input.UnreadRune()
			break
		}
	}
	return Token{
		t:       WhiteSpace,
		lexeme:  b.String(),
		literal: nil,
		line:    s.line,
	}, nil
}

func (s *Scanner) readRune() (rune, error) {
	ch, _, err := s.input.ReadRune()
	if err == io.EOF {
		s.eof = true
		return ' ', err
	}
	if err != nil {
		return ' ', err
	}
	return ch, nil
}
