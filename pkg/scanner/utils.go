package scanner

import "github.com/smasher164/xid"

var spaceMap = map[rune]bool{
	' ':      true,
	'\t':     true,
	'\n':     true,
	'\v':     true,
	'\f':     true,
	'\r':     true,
	'\u001C': true,
	'\u001D': true,
	'\u001E': true,
	'\u001F': true,
	'\u1680': true,
	'\u180E': true,
	'\u2000': true,
	'\u2001': true,
	'\u2002': true,
	'\u2003': true,
	'\u2004': true,
	'\u2005': true,
	'\u2006': true,
	'\u2008': true,
	'\u2009': true,
	'\u200a': true,
	'\u2028': true,
	'\u2029': true,
	'\u205F': true,
	'\u3000': true,
	'\u00A0': true,
	'\u2007': true,
	'\u202F': true,
}

var hexMap = map[rune]bool{
	'0': true,
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
	'8': true,
	'9': true,
	'a': true,
	'b': true,
	'c': true,
	'd': true,
	'e': true,
	'f': true,
	'A': true,
	'B': true,
	'C': true,
	'D': true,
	'E': true,
	'F': true,
}

var octMap = map[rune]bool{
	'0': true,
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
}

func isSpace(ch rune) bool {
	_, ok := spaceMap[ch]
	return ok
}

func isHexDigit(ch rune) bool {
	_, ok := hexMap[ch]
	return ok
}

func isOctDigit(ch rune) bool {
	_, ok := octMap[ch]
	return ok
}

func isInvalidTerminator(ch rune) bool {
	switch {
	case xid.Start(ch):
		return true
	case ch == '"', ch == '\'':
		return true
	default:
		return false
	}
}
