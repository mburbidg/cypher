package scanner

import "strings"

type tokenMap map[string]tokenInfo

type tokenInfo struct {
	t      TokenType
	lexeme string
}

var reservedWords = tokenMap{
	"create":   tokenInfo{Create, "CREATE"},
	"delete":   tokenInfo{Delete, "DELETE"},
	"detach":   tokenInfo{Detach, "DETACH"},
	"exists":   tokenInfo{Exists, "EXISTS"},
	"match":    tokenInfo{Match, "MATCH"},
	"merge":    tokenInfo{Merge, "MERGE"},
	"optional": tokenInfo{Optional, "OPTIONAL"},
	"remove":   tokenInfo{Remove, "REMOVE"},
	"return":   tokenInfo{Return, "RETURN"},
}

func (m tokenMap) token(symbolName string, line int) (Token, bool) {
	if info, ok := reservedWords[strings.ToLower(symbolName)]; ok {
		return newKeywordToken(info.t, info.lexeme, line), true
	}
	return Token{}, false
}
