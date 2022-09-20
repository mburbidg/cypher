package scanner

import "strings"

type tokenMap map[string]tokenInfo

type tokenInfo struct {
	t      TokenType
	lexeme string
}

var reservedWords = tokenMap{
	"create":   tokenInfo{Create, "CREATE"},
	"delete":   tokenInfo{Create, "DELETE"},
	"detach":   tokenInfo{Create, "DETACH"},
	"exists":   tokenInfo{Create, "EXISTS"},
	"match":    tokenInfo{Create, "MATCH"},
	"merge":    tokenInfo{Create, "MERGE"},
	"optional": tokenInfo{Create, "OPTIONAL"},
	"remove":   tokenInfo{Create, "REMOVE"},
}

func (m tokenMap) token(symbolName string, line int) (Token, bool) {
	if info, ok := reservedWords[strings.ToLower(symbolName)]; ok {
		return newKeywordToken(info.t, info.lexeme, line), true
	}
	return Token{}, false
}
