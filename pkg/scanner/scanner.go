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
	tokens   []Token
}

const (
	eof = -1
)

func New(input io.Reader, reporter utils.Reporter) *Scanner {
	return &Scanner{
		input:    bufio.NewReader(input),
		reporter: reporter,
		line:     1,
		eof:      false,
		tokens:   make([]Token, 0, 10),
	}
}

func (s *Scanner) next() rune {
	if s.eof {
		return eof
	}
	ch, _, err := s.input.ReadRune()
	if err != nil {
		if err != io.EOF {
			s.reporter.Error(s.line, fmt.Sprintf("error reading input: err=%s\n", err))
		}
		s.eof = true
		return eof
	}
	return ch
}

func (s *Scanner) peek() rune {
	ch := s.next()
	if !s.eof {
		err := s.input.UnreadRune()
		if err != nil {
			s.reporter.Error(s.line, fmt.Sprintf("error un-reading input: err=%s\n", err))
		}
	}
	return ch
}

func (s *Scanner) prev() {
	if !s.eof {
		err := s.input.UnreadRune()
		if err != nil {
			s.reporter.Error(s.line, fmt.Sprintf("error un-reading input: err=%s\n", err))
		}
	}
}

func (s *Scanner) pushToken(token Token) {
	s.tokens = append(s.tokens, token)
}

func (s *Scanner) popToken() (Token, bool) {
	if len(s.tokens) > 0 {
		token := s.tokens[len(s.tokens)-1]
		s.tokens = s.tokens[:len(s.tokens)-1]
		return token, true
	}
	return Token{}, false
}

func (s *Scanner) Line() int {
	return s.line
}

func (s *Scanner) ReturnToken(token Token) {
	s.pushToken(token)
}

func (s *Scanner) NextToken() Token {
	if t, ok := s.popToken(); ok {
		return t
	}
	if s.eof == true {
		return endOfInputToken
	}
	ch := s.next()
	if ch == eof {
		return endOfInputToken
	}

	switch {
	case ch == '.':
		ch = s.next()
		if ch == eof {
			return endOfInputToken
		}
		if ch == '.' {
			return newOperatorToken(Dotdot, s.line)
		}
		s.prev()
		if unicode.IsDigit(ch) {
			return s.scanNumber('.')
		}
		return newOperatorToken(Period, s.line)
	case ch == ',':
		return newOperatorToken(Comma, s.line)
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
		ch = s.next()
		switch ch {
		case '*':
			s.consumeMultilineComment()
			return s.NextToken()
		case '/':
			s.consumeSingleLineComment()
			return s.NextToken()
		default:
			s.prev()
			return newOperatorToken(ForwardSlash, s.line)
		}
	case ch == '%':
		return newOperatorToken(Percent, s.line)
	case ch == '^':
		return newOperatorToken(Caret, s.line)
	case ch == '=':
		return newOperatorToken(Equal, s.line)
	case ch == '<':
		ch := s.next()
		switch ch {
		case '>':
			return newOperatorToken(NotEqual, s.line)
		case '=':
			return newOperatorToken(LessThanOrEqual, s.line)
		default:
			s.prev()
			return newOperatorToken(LessThan, s.line)
		}
	case ch == '>':
		ch := s.next()
		if ch == '=' {
			return newOperatorToken(GreaterThanOrEqual, s.line)
		}
		s.prev()
		return newOperatorToken(GreaterThan, s.line)
	case ch == '$':
		return newOperatorToken(DollarSign, s.line)
	case ch == ':':
		return newOperatorToken(Colon, s.line)
	case ch == '|':
		return newOperatorToken(Pipe, s.line)
	case ch == '-':
		return newOperatorToken(Dash, s.line)
	case unicode.IsDigit(ch):
		return s.scanNumber(ch)
	case ch == '"', ch == '\'':
		return s.scanString(ch)
	case xid.Start(ch):
		return s.scanIdentifier(ch)
	case isSpace(ch):
		s.consumeWhitespace(ch)
		return s.NextToken()
	}

	return Token{}
}

func (s *Scanner) consumeWhitespace(ch rune) {
	s.incLine(ch)
	for {
		ch := s.next()
		if !isSpace(ch) {
			s.prev()
			break
		}
		s.incLine(ch)
	}
}

func (s *Scanner) consumeMultilineComment() {
	for {
		ch := s.next()
		if ch == eof {
			s.reporter.Error(s.line, "unterminated comment")
			return
		}
		if ch == '*' {
			ch := s.next()
			if ch == '/' {
				return
			}
			s.incLine(ch)
		}
		s.incLine(ch)
	}
}

func (s *Scanner) consumeSingleLineComment() {
	for {
		ch := s.next()
		if ch == eof {
			return
		}
		if ch == '\n' {
			s.incLine(ch)
			return
		}
	}
}

func (s *Scanner) scanIdentifier(ch rune) Token {
	b := strings.Builder{}
	b.WriteRune(ch)
	for {
		ch := s.next()
		if xid.Continue(ch) {
			b.WriteRune(ch)
		} else {
			s.prev()
			break
		}
	}
	if token, ok := reservedWords.token(b.String(), s.line); ok {
		return token
	}
	return newIdentifierToken(b.String(), s.line)
}

func (s *Scanner) scanNumber(ch rune) Token {
	b := strings.Builder{}
	b.WriteRune(ch)
	switch ch {
	case '0':
		switch ch := s.next(); ch {
		case 'x':
			b.WriteRune(ch)
			return s.scanHexInteger(&b)
		case '1', '2', '3', '4', '5', '6', '7':
			b.WriteRune(ch)
			return s.scanOctInteger(&b)
		case 'E':
			b.WriteRune(ch)
			return s.scanExponent(&b)
		case '.':
			b.WriteRune(ch)
			return s.scanDouble(&b)
		default:
			s.prev()
			return newIntegerToken(DecimalInteger, "0", 10, s.line)
		}
	case '.':
		return s.scanDouble(&b)
	}
	for {
		ch := s.next()
		if ch == eof {
			s.prev()
			return newIntegerToken(DecimalInteger, b.String(), 10, s.line)
		}
		switch ch {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			b.WriteRune(ch)
		case '.':
			b.WriteRune(ch)
			return s.scanDouble(&b)
		case 'E':
			b.WriteRune(ch)
			return s.scanExponent(&b)
		default:
			s.prev()
			return newIntegerToken(DecimalInteger, b.String(), 10, s.line)
		}
	}
}

func (s *Scanner) scanDouble(b *strings.Builder) Token {
	switch ch := s.next(); ch {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		b.WriteRune(ch)
	default:
		s.prev()
		s.reporter.Error(s.line, "expecting fractional part of double")
		return newIllegalToken(b.String())
	}
	for {
		switch ch := s.next(); ch {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			b.WriteRune(ch)
		case 'E':
			b.WriteRune(ch)
			return s.scanExponent(b)
		default:
			s.prev()
			return newDoubleToken(b.String(), s.line)
		}
	}
}

func (s *Scanner) scanHexInteger(b *strings.Builder) Token {
	for {
		switch ch := s.next(); ch {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			b.WriteRune(ch)
		case 'a', 'b', 'c', 'd', 'e', 'f', 'A', 'B', 'C', 'D', 'E', 'F':
			b.WriteRune(ch)
		default:
			s.prev()
			if len(b.String()) < 3 {
				s.reporter.Error(s.line, "expecting hex digit following 'x'")
				return newIllegalToken(b.String())
			}
			return newIntegerToken(HexInteger, b.String(), 0, s.line)
		}
	}
}

func (s *Scanner) scanOctInteger(b *strings.Builder) Token {
	for {
		switch ch := s.next(); ch {
		case '0', '1', '2', '3', '4', '5', '6', '7':
			b.WriteRune(ch)
		default:
			s.prev()
			return newIntegerToken(OctInteger, b.String(), 0, s.line)
		}
	}
}

func (s *Scanner) scanExponent(b *strings.Builder) Token {
	for i := 0; ; i++ {
		switch ch := s.next(); ch {
		case '-':
			if i > 0 {
				s.prev()
				s.reporter.Error(s.line, "invalid exponent")
				return newIllegalToken(b.String())
			}
			if !isDigit(s.peek()) {
				s.reporter.Error(s.line, "invalid exponent")
				return newIllegalToken(b.String())
			}
			b.WriteRune(ch)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			b.WriteRune(ch)
		default:
			s.prev()
			switch i {
			case 0:
				s.reporter.Error(s.line, "expecting exponent")
				return newIllegalToken(b.String())
			case 1:
				if ch == '-' {
					s.reporter.Error(s.line, "expecting digit following '-' in exponent")
					return newIllegalToken(b.String())
				}
				return newDoubleToken(b.String(), s.line)
			default:
				return newDoubleToken(b.String(), s.line)
			}
		}
	}
}

func (s *Scanner) scanString(ch rune) Token {
	startCh := ch
	lexeme := strings.Builder{}
	lexeme.WriteRune(startCh)
	literal := strings.Builder{}
	for {
		ch := s.next()
		if ch == eof {
			s.reporter.Error(s.line, fmt.Sprintf("expecting string end character '%s'", string(startCh)))
			return newIllegalToken(lexeme.String())
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
	ch := s.next()
	if ch == eof {
		s.reporter.Error(s.line, fmt.Sprintf("incomplete character escape"))
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
			ch := s.next()
			if ch == eof {
				s.reporter.Error(s.line, fmt.Sprintf("incomplete character escape"))
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
			ch := s.next()
			if ch == eof {
				s.reporter.Error(s.line, fmt.Sprintf("incomplete character escape"))
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
