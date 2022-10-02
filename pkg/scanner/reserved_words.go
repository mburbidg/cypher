package scanner

import "strings"

type tokenMap map[string]tokenInfo

type tokenInfo struct {
	t      TokenType
	lexeme string
}

var reservedWords = tokenMap{
	"create":     tokenInfo{Create, "CREATE"},
	"delete":     tokenInfo{Delete, "DELETE"},
	"detach":     tokenInfo{Detach, "DETACH"},
	"exists":     tokenInfo{Exists, "EXISTS"},
	"match":      tokenInfo{Match, "MATCH"},
	"merge":      tokenInfo{Merge, "MERGE"},
	"optional":   tokenInfo{Optional, "OPTIONAL"},
	"remove":     tokenInfo{Remove, "REMOVE"},
	"return":     tokenInfo{Return, "RETURN"},
	"set":        tokenInfo{Set, "SET"},
	"union":      tokenInfo{Union, "UNION"},
	"unwind":     tokenInfo{Unwind, "UNWIND"},
	"with":       tokenInfo{With, "WITH"},
	"limit":      tokenInfo{Limit, "LIMIT"},
	"order":      tokenInfo{Order, "ORDER"},
	"skip":       tokenInfo{Skip, "SKIP"},
	"where":      tokenInfo{Where, "WHERE"},
	"asc":        tokenInfo{Asc, "ASC"},
	"ascending":  tokenInfo{Ascending, "ASCENDING"},
	"by":         tokenInfo{By, "BY"},
	"desc":       tokenInfo{Desc, "DECS"},
	"descending": tokenInfo{Descending, "DESCENDING"},
	"on":         tokenInfo{On, "ON"},
	"all":        tokenInfo{All, "ALL"},
	"case":       tokenInfo{Case, "CASE"},
	"else":       tokenInfo{Else, "ELSE"},
	"end":        tokenInfo{End, "END"},
	"then":       tokenInfo{Then, "THEN"},
	"when":       tokenInfo{When, "WHEN"},
	"and":        tokenInfo{And, "AND"},
	"as":         tokenInfo{As, "AS"},
	"contains":   tokenInfo{Contains, "CONTAINS"},
	"distinct":   tokenInfo{Distinct, "DISTINCT"},
	"ends":       tokenInfo{Ends, "ENDS"},
	"in":         tokenInfo{In, "IN"},
	"is":         tokenInfo{Is, "IS"},
	"not":        tokenInfo{Not, "NOT"},
	"or":         tokenInfo{Or, "OR"},
	"starts":     tokenInfo{Starts, "STARTS"},
	"xor":        tokenInfo{Xor, "XOR"},
	"false":      tokenInfo{False, "FALSE"},
	"null":       tokenInfo{Null, "NULL"},
	"true":       tokenInfo{True, "TRUE"},
}

func (m tokenMap) token(symbolName string, line int) (Token, bool) {
	if info, ok := reservedWords[strings.ToLower(symbolName)]; ok {
		return newKeywordToken(info.t, info.lexeme, line), true
	}
	return Token{}, false
}
