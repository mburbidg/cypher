package scanner

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

func isSpace(ch rune) bool {
	_, ok := spaceMap[ch]
	return ok
}

func isDigit(ch rune) bool {
	switch {
	case ch >= '0' && ch <= '9':
		return true
	default:
		return false
	}
}

func numberTerm(ch rune) bool {
	switch ch {
	case '.', '(', ')', '{', '}', '[', ']', '+', '-', '*', '/', '%', '^', '=', '<', '>', '$':
		return true
	default:
		return false
	}
}
