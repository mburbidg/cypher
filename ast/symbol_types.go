package ast

type SymbolType int

const (
	Count SymbolType = iota
	Filter
	Extract
	Any
	None
	Single
	Identifier
)

var SymbolNames = map[string]SymbolType{
	"count":   Count,
	"filter":  Filter,
	"extract": Extract,
	"any":     Any,
	"none":    None,
	"single":  Single,
}
