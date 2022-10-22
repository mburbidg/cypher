package scanner

import (
	"fmt"
	"github.com/mburbidg/cypher/pkg/utils"
	"github.com/smasher164/xid"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Position struct {
	offset     int
	prevOffset int
	line       int
	eofRead    bool
}

type Scanner struct {
	src      []byte
	Position Position
	reporter utils.Reporter
}

const (
	eof = -1
)

func New(src []byte, reporter utils.Reporter) *Scanner {
	return &Scanner{
		src:      src,
		Position: Position{offset: 0, prevOffset: 0, line: 1, eofRead: false},
		reporter: reporter,
	}
}

func (s *Scanner) next() rune {
	if s.Position.offset < len(s.src) {
		r, w := utf8.DecodeRune(s.src[s.Position.offset:])
		s.Position.prevOffset = s.Position.offset
		s.Position.offset += w
		if r == '\n' {
			s.Position.line += 1
		}
		return r
	}
	s.Position.eofRead = true
	return eof
}

func (s *Scanner) peek() rune {
	if s.Position.offset < len(s.src) {
		r, _ := utf8.DecodeRune(s.src[s.Position.offset:])
		return r
	}
	return eof
}

func (s *Scanner) prev() {
	// Once we hit eof, no backing up.
	if !s.Position.eofRead {
		s.Position.offset = s.Position.prevOffset
	}
}

func (s *Scanner) Line() int {
	return s.Position.line
}

func (s *Scanner) NextToken() Token {
	if s.Position.offset >= len(s.src) {
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
			return newOperatorToken(Dotdot, s.Position.line)
		}
		s.prev()
		if unicode.IsDigit(ch) {
			return s.scanNumber('.')
		}
		return newOperatorToken(Period, s.Position.line)
	case ch == ',':
		return newOperatorToken(Comma, s.Position.line)
	case ch == '(':
		return newOperatorToken(OpenParen, s.Position.line)
	case ch == ')':
		return newOperatorToken(CloseParen, s.Position.line)
	case ch == '{':
		return newOperatorToken(OpenBrace, s.Position.line)
	case ch == '}':
		return newOperatorToken(CloseBrace, s.Position.line)
	case ch == '[':
		return newOperatorToken(OpenBracket, s.Position.line)
	case ch == ']':
		return newOperatorToken(CloseBracket, s.Position.line)
	case ch == '+':
		return newOperatorToken(Plus, s.Position.line)
	case ch == '-':
		return newOperatorToken(Dash, s.Position.line)
	case ch == '*':
		return newOperatorToken(Star, s.Position.line)
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
			return newOperatorToken(ForwardSlash, s.Position.line)
		}
	case ch == '%':
		return newOperatorToken(Percent, s.Position.line)
	case ch == '^':
		return newOperatorToken(Caret, s.Position.line)
	case ch == '=':
		return newOperatorToken(Equal, s.Position.line)
	case ch == '<':
		ch := s.next()
		switch ch {
		case '>':
			return newOperatorToken(NotEqual, s.Position.line)
		case '=':
			return newOperatorToken(LessThanOrEqual, s.Position.line)
		default:
			s.prev()
			return newOperatorToken(LessThan, s.Position.line)
		}
	case ch == '>':
		ch := s.next()
		if ch == '=' {
			return newOperatorToken(GreaterThanOrEqual, s.Position.line)
		}
		s.prev()
		return newOperatorToken(GreaterThan, s.Position.line)
	case ch == '$':
		return newOperatorToken(DollarSign, s.Position.line)
	case ch == ':':
		return newOperatorToken(Colon, s.Position.line)
	case ch == '|':
		return newOperatorToken(Pipe, s.Position.line)
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
	for {
		ch := s.next()
		if !isSpace(ch) {
			s.prev()
			break
		}
	}
}

func (s *Scanner) consumeMultilineComment() {
	for {
		ch := s.next()
		if ch == eof {
			s.reporter.Error(s.Position.line, "unterminated comment")
			return
		}
		if ch == '*' {
			ch := s.next()
			if ch == '/' {
				return
			}
		}
	}
}

func (s *Scanner) consumeSingleLineComment() {
	for {
		ch := s.next()
		if ch == eof {
			return
		}
		if ch == '\n' {
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
	if token, ok := reservedWords.token(b.String(), s.Position.line); ok {
		return token
	}
	return newIdentifierToken(b.String(), s.Position.line)
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
			if s.peek() == '.' {
				s.prev()
				return newIntegerToken(DecimalInteger, b.String(), 10, s.Position.line)
			}
			b.WriteRune(ch)
			return s.scanDouble(&b)
		default:
			s.prev()
			return newIntegerToken(DecimalInteger, "0", 10, s.Position.line)
		}
	case '.':
		return s.scanDouble(&b)
	}
	for {
		ch := s.next()
		if ch == eof {
			s.prev()
			return newIntegerToken(DecimalInteger, b.String(), 10, s.Position.line)
		}
		switch ch {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			b.WriteRune(ch)
		case '.':
			if s.peek() == '.' {
				s.prev()
				return newIntegerToken(DecimalInteger, b.String(), 10, s.Position.line)
			}
			b.WriteRune(ch)
			return s.scanDouble(&b)
		case 'E':
			b.WriteRune(ch)
			return s.scanExponent(&b)
		default:
			s.prev()
			return newIntegerToken(DecimalInteger, b.String(), 10, s.Position.line)
		}
	}
}

func (s *Scanner) scanDouble(b *strings.Builder) Token {
	switch ch := s.next(); ch {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		b.WriteRune(ch)
	default:
		s.prev()
		s.reporter.Error(s.Position.line, "expecting fractional part of double")
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
			return newDoubleToken(b.String(), s.Position.line)
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
				s.reporter.Error(s.Position.line, "expecting hex digit following 'x'")
				return newIllegalToken(b.String())
			}
			return newIntegerToken(HexInteger, b.String(), 0, s.Position.line)
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
			return newIntegerToken(OctInteger, b.String(), 0, s.Position.line)
		}
	}
}

func (s *Scanner) scanExponent(b *strings.Builder) Token {
	for i := 0; ; i++ {
		switch ch := s.next(); ch {
		case '-':
			if i > 0 {
				s.prev()
				s.reporter.Error(s.Position.line, "invalid exponent")
				return newIllegalToken(b.String())
			}
			if !isDigit(s.peek()) {
				s.reporter.Error(s.Position.line, "invalid exponent")
				return newIllegalToken(b.String())
			}
			b.WriteRune(ch)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			b.WriteRune(ch)
		default:
			s.prev()
			switch i {
			case 0:
				s.reporter.Error(s.Position.line, "expecting exponent")
				return newIllegalToken(b.String())
			case 1:
				if ch == '-' {
					s.reporter.Error(s.Position.line, "expecting digit following '-' in exponent")
					return newIllegalToken(b.String())
				}
				return newDoubleToken(b.String(), s.Position.line)
			default:
				return newDoubleToken(b.String(), s.Position.line)
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
			s.reporter.Error(s.Position.line, fmt.Sprintf("expecting string end character '%s'", string(startCh)))
			return newIllegalToken(lexeme.String())
		}
		switch {
		case ch == startCh:
			lexeme.WriteRune(ch)
			return newStringToken(lexeme.String(), literal.String(), s.Position.line)
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
		s.reporter.Error(s.Position.line, fmt.Sprintf("incomplete character escape"))
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
				s.reporter.Error(s.Position.line, fmt.Sprintf("incomplete character escape"))
				return 0
			}
			ok := unicode.In(ch, unicode.ASCII_Hex_Digit)
			if !ok {
				s.reporter.Error(s.Position.line, fmt.Sprintf("invalid escaped unicode character"))
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
				s.reporter.Error(s.Position.line, fmt.Sprintf("incomplete character escape"))
				return 0
			}
			ok := unicode.In(ch, unicode.ASCII_Hex_Digit)
			if !ok {
				s.reporter.Error(s.Position.line, fmt.Sprintf("invalid escaped unicode character"))
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
		ch := s.next()
		if !isSpace(ch) {
			s.prev()
			break
		}
	}
}
