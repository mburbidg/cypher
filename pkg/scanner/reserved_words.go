package scanner

import "strings"

type tokenMap map[string]tokenInfo

type tokenInfo struct {
	t      TokenType
	lexeme string
}

func addReservedWord(t tokenInfo) tokenInfo {
	ReservedWordTokens[t.t] = true
	return t
}

var ReservedWordTokens = map[TokenType]bool{}

var reservedWords = tokenMap{
	"all":        addReservedWord(tokenInfo{All, "ALL"}),
	"asc":        addReservedWord(tokenInfo{Asc, "ASC"}),
	"ascending":  addReservedWord(tokenInfo{Ascending, "ASCENDING"}),
	"by":         addReservedWord(tokenInfo{By, "BY"}),
	"create":     addReservedWord(tokenInfo{Create, "CREATE"}),
	"delete":     addReservedWord(tokenInfo{Delete, "DELETE"}),
	"desc":       addReservedWord(tokenInfo{Desc, "DECS"}),
	"descending": addReservedWord(tokenInfo{Descending, "DESCENDING"}),
	"detach":     addReservedWord(tokenInfo{Detach, "DETACH"}),
	"exists":     addReservedWord(tokenInfo{Exists, "EXISTS"}),
	"limit":      addReservedWord(tokenInfo{Limit, "LIMIT"}),
	"match":      addReservedWord(tokenInfo{Match, "MATCH"}),
	"merge":      addReservedWord(tokenInfo{Merge, "MERGE"}),
	"on":         addReservedWord(tokenInfo{On, "ON"}),
	"optional":   addReservedWord(tokenInfo{Optional, "OPTIONAL"}),
	"order":      addReservedWord(tokenInfo{Order, "ORDER"}),
	"remove":     addReservedWord(tokenInfo{Remove, "REMOVE"}),
	"return":     addReservedWord(tokenInfo{Return, "RETURN"}),
	"set":        addReservedWord(tokenInfo{Set, "SET"}),
	"skip":       addReservedWord(tokenInfo{Skip, "SKIP"}),
	"where":      addReservedWord(tokenInfo{Where, "WHERE"}),
	"with":       addReservedWord(tokenInfo{With, "WITH"}),
	"union":      addReservedWord(tokenInfo{Union, "UNION"}),
	"unwind":     addReservedWord(tokenInfo{Unwind, "UNWIND"}),
	"and":        addReservedWord(tokenInfo{And, "AND"}),
	"as":         addReservedWord(tokenInfo{As, "AS"}),
	"contains":   addReservedWord(tokenInfo{Contains, "CONTAINS"}),
	"distinct":   addReservedWord(tokenInfo{Distinct, "DISTINCT"}),
	"ends":       addReservedWord(tokenInfo{Ends, "ENDS"}),
	"in":         addReservedWord(tokenInfo{In, "IN"}),
	"is":         addReservedWord(tokenInfo{Is, "IS"}),
	"not":        addReservedWord(tokenInfo{Not, "NOT"}),
	"or":         addReservedWord(tokenInfo{Or, "OR"}),
	"starts":     addReservedWord(tokenInfo{Starts, "STARTS"}),
	"xor":        addReservedWord(tokenInfo{Xor, "XOR"}),
	"false":      addReservedWord(tokenInfo{False, "FALSE"}),
	"true":       addReservedWord(tokenInfo{True, "TRUE"}),
	"null":       addReservedWord(tokenInfo{Null, "NULL"}),
	"constraint": addReservedWord(tokenInfo{Constraint, "CONSTRAINT"}),
	"do":         addReservedWord(tokenInfo{Do, "DO"}),
	"for":        addReservedWord(tokenInfo{For, "FOR"}),
	"require":    addReservedWord(tokenInfo{Require, "REQUIRE"}),
	"unique":     addReservedWord(tokenInfo{Unique, "UNIQUE"}),
	"case":       addReservedWord(tokenInfo{Case, "CASE"}),
	"when":       addReservedWord(tokenInfo{When, "WHEN"}),
	"then":       addReservedWord(tokenInfo{Then, "THEN"}),
	"else":       addReservedWord(tokenInfo{Else, "ELSE"}),
	"end":        addReservedWord(tokenInfo{End, "END"}),
	"mandatory":  addReservedWord(tokenInfo{Mandatory, "MANDATORY"}),
	"scalar":     addReservedWord(tokenInfo{Scalar, "SCALAR"}),
	"of":         addReservedWord(tokenInfo{Of, "OF"}),
	"add":        addReservedWord(tokenInfo{Add, "ADD"}),
	"drop":       addReservedWord(tokenInfo{Drop, "DROP"}),
}

func (m tokenMap) token(symbolName string, line int) (Token, bool) {
	if info, ok := reservedWords[strings.ToLower(symbolName)]; ok {
		return newKeywordToken(info.t, info.lexeme, line), true
	}
	return Token{}, false
}
