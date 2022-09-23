package scanner

import (
	"bufio"
	"github.com/smasher164/xid"
	"io"
	"strings"
	"unicode"
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

func (s *Scanner) NextToken() Token {
	if s.eof == true {
		return newEndOfInputToken()
	}

	ch, _, err := s.input.ReadRune()
	if err != nil {
		return newEndOfInputToken()
	}

	switch {
	case ch == '.':
		return newOperatorToken(Period, s.line)
	case ch == '(':
		return newOperatorToken(OpenParen, s.line)
	case ch == ')':
		return newOperatorToken(CloseParen, s.line)
	case ch == '{':
		return newOperatorToken(OpenBrace, s.line)
	case ch == '}':
		return newOperatorToken(CloseBrace, s.line)
	case ch == '[':
		return newOperatorToken(OpenBracket, s.line)
	case ch == ']':
		return newOperatorToken(CloseBracket, s.line)
	case ch == '+':
		return newOperatorToken(Plus, s.line)
	case ch == '-':
		return newOperatorToken(Minus, s.line)
	case ch == '*':
		return newOperatorToken(Star, s.line)
	case ch == '/':
		match, err := s.matchNext('*')
		if err != nil {
			return newOperatorToken(ForwardSlash, s.line)
		}
		if match {
			err := s.consumeMultilineComment()
			if err != nil {
				return s.endOfInputToken()
			}
			return s.NextToken()
		}
		match, err = s.matchNext('/')
		if err != nil {
			return newOperatorToken(ForwardSlash, s.line)
		}
		if match {
			err := s.consumeSingleLineComment()
			if err != nil {
				return s.endOfInputToken()
			}
			return s.NextToken()
		}
		return newOperatorToken(ForwardSlash, s.line)
	case ch == '%':
		return newOperatorToken(Percent, s.line)
	case ch == '^':
		return newOperatorToken(Caret, s.line)
	case ch == '=':
		return newOperatorToken(Equal, s.line)
	case ch == '<':
		match, err := s.matchNext('>')
		if err != nil {
			return s.endOfInputToken()
		}
		if match {
			return newOperatorToken(NotEqual, s.line)
		}
		match, err = s.matchNext('=')
		if err != nil {
			return s.endOfInputToken()
		}
		if match {
			return newOperatorToken(LessThanOrEqual, s.line)
		}
		return newOperatorToken(LessThan, s.line)
	case ch == '>':
		match, err := s.matchNext('=')
		if err != nil {
			return s.endOfInputToken()
		}
		if match {
			return newOperatorToken(GreaterThanOrEqual, s.line)
		}
		return newOperatorToken(GreaterThan, s.line)
	case ch == '$':
		return newOperatorToken(DollarSign, s.line)
	case unicode.IsDigit(ch):
		return s.scanNumber(ch)
	case xid.Start(ch):
		return s.scanSymbolicName(ch)
	case space(ch):
		err := s.consumeWhitespace(ch)
		if err != nil {
			return s.endOfInputToken()
		}
		return s.NextToken()
	}

	return Token{}
}

func (s *Scanner) endOfInputToken() Token {
	s.eof = true
	return newEndOfInputToken()
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

func (s *Scanner) consumeWhitespace(ch rune) error {
	s.incLine(ch)
	for {
		ch, _, err := s.input.ReadRune()
		if err != nil {
			return err
		}
		if !space(ch) {
			return s.input.UnreadRune()
		}
		s.incLine(ch)
	}
}

func (s *Scanner) consumeMultilineComment() error {
	for {
		ch, _, err := s.input.ReadRune()
		if err != nil {
			return err
		}
		if ch == '*' {
			ok, err := s.matchNext('/')
			if err != nil {
				return err
			}
			if ok {
				return nil
			}
		}
		s.incLine(ch)
	}
}

func (s *Scanner) consumeSingleLineComment() error {
	for {
		ch, _, err := s.input.ReadRune()
		if err != nil {
			return err
		}
		if ch == '\n' {
			s.incLine(ch)
			return nil
		}
	}
}

func (s *Scanner) scanSymbolicName(ch rune) Token {
	b := strings.Builder{}
	b.WriteRune(ch)
	for {
		ch, _, err := s.input.ReadRune()
		if err != nil {
			break
		}
		if xid.Continue(ch) {
			b.WriteRune(ch)
		} else {
			_ = s.input.UnreadRune()
			break
		}
	}
	if token, ok := reservedWords.token(b.String(), s.line); ok {
		return token
	}
	return Token{
		t:       SymbolName,
		lexeme:  b.String(),
		literal: nil,
		line:    s.line,
	}
}

func (s *Scanner) scanNumber(ch rune) Token {
	t := Integer
	b := strings.Builder{}
	b.WriteRune(ch)
	for {
		ch, _, err := s.input.ReadRune()
		if err != nil {
			return newNumberToken(t, b.String(), s.line)
		}
		switch {
		case unicode.IsDigit(ch):
			b.WriteRune(ch)
		case ch == '.':
			b.WriteRune(ch)
			t = Double
		default:
			_ = s.input.UnreadRune()
			return newNumberToken(t, b.String(), s.line)
		}
	}
}

func (s *Scanner) incLine(ch rune) {
	if ch == '\n' {
		s.line = s.line + 1
	}
}
