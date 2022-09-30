package scanner

import (
	"bufio"
	"fmt"
	"github.com/mburbidg/cypher/pkg/utils"
	"github.com/smasher164/xid"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Scanner struct {
	input    *bufio.Reader
	reporter utils.Reporter
	line     int
	eof      bool
}

func New(input io.Reader, reporter utils.Reporter) *Scanner {
	return &Scanner{
		input:    bufio.NewReader(input),
		reporter: reporter,
		line:     1,
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
		ch, _, err := s.input.ReadRune()
		if err != nil {
			return newEndOfInputToken()
		}
		if unicode.IsDigit(ch) {
			_ = s.input.UnreadRune()
			return s.scanNumber('.')
		}
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
	case ch == '"', ch == '\'':
		return s.scanString(ch)
	case xid.Start(ch):
		return s.scanIdentifier(ch)
	case isSpace(ch):
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

func (s *Scanner) peek() rune {
	ch, _, err := s.input.ReadRune()
	if err != nil {
		return 0
	}
	_ = s.input.UnreadRune()
	return ch
}

func (s *Scanner) consumeWhitespace(ch rune) error {
	s.incLine(ch)
	for {
		ch, _, err := s.input.ReadRune()
		if err != nil {
			return err
		}
		if !isSpace(ch) {
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

func (s *Scanner) scanIdentifier(ch rune) Token {
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
		t:       Identifier,
		lexeme:  b.String(),
		literal: nil,
		line:    s.line,
	}
}

func (s *Scanner) scanNumber(ch rune) Token {
	t := Integer
	b := strings.Builder{}
	b.WriteRune(ch)
	if ch == '0' {
		nextCh, _, err := s.input.ReadRune()
		switch {
		case err != nil:
			s.eof = true
			return newIntegerToken(b.String(), 10, s.line)
		case nextCh == 'x':
			b.WriteRune(nextCh)
			return s.scanHexInteger(b.String())
		case isOctDigit(nextCh):
			b.WriteRune(nextCh)
			return s.scanOctInteger(b.String())
		case nextCh == '.':
			b.WriteRune(nextCh)
			t = Double
			if !unicode.IsDigit(s.peek()) {
				return newErrorToken(b.String())
			}
		default:
			_ = s.input.UnreadRune()
			return newIntegerToken("0", 10, s.line)
		}
	}
	if ch == '.' {
		t = Double
		if !unicode.IsDigit(s.peek()) {
			return newErrorToken(b.String())
		}
	}
	for {
		ch, _, err := s.input.ReadRune()
		if err != nil {
			s.eof = true
			return newNumberToken(t, b.String(), s.line)
		}
		switch {
		case unicode.IsDigit(ch):
			b.WriteRune(ch)
		case ch == '.':
			b.WriteRune(ch)
			t = Double
			if !unicode.IsDigit(s.peek()) {
				return newErrorToken(b.String())
			}
		default:
			_ = s.input.UnreadRune()
			return newNumberToken(t, b.String(), s.line)
		}
	}
}

func (s *Scanner) scanHexInteger(prefix string) Token {
	b := strings.Builder{}
	b.WriteString(prefix)
	for {
		ch, _, err := s.input.ReadRune()
		if err != nil {
			s.eof = true
			return newIntegerToken(b.String(), 16, s.line)
		}
		if isHexDigit(ch) {
			b.WriteRune(ch)
		} else {
			_ = s.input.UnreadRune()
			return newIntegerToken(b.String(), 16, s.line)
		}
	}
}

func (s *Scanner) scanOctInteger(prefix string) Token {
	b := strings.Builder{}
	b.WriteString(prefix)
	for {
		ch, _, err := s.input.ReadRune()
		if err != nil {
			s.eof = true
			return newIntegerToken(b.String(), 8, s.line)
		}
		if isOctDigit(ch) {
			b.WriteRune(ch)
		} else {
			_ = s.input.UnreadRune()
			return newIntegerToken(b.String(), 8, s.line)
		}
	}
}

func (s *Scanner) scanString(ch rune) Token {
	startCh := ch
	lexeme := strings.Builder{}
	lexeme.WriteRune(startCh)
	literal := strings.Builder{}
	for {
		ch, _, err := s.input.ReadRune()
		if err != nil {
			s.reporter.Error(s.line, fmt.Sprintf("expecting string end character '%s'", string(startCh)))
			s.eof = true
			return newErrorToken(lexeme.String())
		}
		switch {
		case ch == startCh:
			lexeme.WriteRune(ch)
			return newStringToken(lexeme.String(), literal.String(), s.line)
		case ch == '\\':
			lexeme.WriteRune(ch)
			literal.WriteRune(s.scanEscapeCharacter())
		default:
			lexeme.WriteRune(ch)
			literal.WriteRune(ch)
		}
	}
}

func (s *Scanner) scanEscapeCharacter() rune {
	ch, _, err := s.input.ReadRune()
	if err != nil {
		s.reporter.Error(s.line, fmt.Sprintf("incomplete character escape"))
		s.eof = true
		return 0
	}
	switch ch {
	case '\'', '"':
		return ch
	case 'b':
		return '\b'
	case 'f':
		return '\f'
	case 'n':
		return '\n'
	case 'r':
		return '\r'
	case 't':
		return '\t'
	case 'u':
		b := strings.Builder{}
		b.WriteString("\\u")
		for i := 0; i < 4; i++ {
			ch, _, err := s.input.ReadRune()
			if err != nil {
				s.reporter.Error(s.line, fmt.Sprintf("incomplete character escape"))
				s.eof = true
				return 0
			}
			ok := unicode.In(ch, unicode.ASCII_Hex_Digit)
			if !ok {
				s.reporter.Error(s.line, fmt.Sprintf("invalid escaped unicode character"))
				return 0
			}
			b.WriteRune(ch)
		}
		ch, _ := utf8.DecodeLastRuneInString(b.String())
		return ch
	case 'U':
		b := strings.Builder{}
		b.WriteString("\\U")
		for i := 0; i < 8; i++ {
			ch, _, err := s.input.ReadRune()
			if err != nil {
				s.reporter.Error(s.line, fmt.Sprintf("incomplete character escape"))
				s.eof = true
				return 0
			}
			ok := unicode.In(ch, unicode.ASCII_Hex_Digit)
			if !ok {
				s.reporter.Error(s.line, fmt.Sprintf("invalid escaped unicode character"))
				return 0
			}
			b.WriteRune(ch)
		}
		ch, _ := utf8.DecodeLastRuneInString(b.String())
		return ch
	}
	return 0
}

func (s *Scanner) consumeNonWhitespace() {
	for {
		ch, _, err := s.input.ReadRune()
		if err != nil {
			s.eof = true
		}
		if !isSpace(ch) {
			_ = s.input.UnreadRune()
			break
		}
	}
}

func (s *Scanner) incLine(ch rune) {
	if ch == '\n' {
		s.line = s.line + 1
	}
}
