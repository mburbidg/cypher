package utils

import "fmt"

type ParseError struct {
	Line int
	Msg  string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("Error: %s (line %d)", e.Msg, e.Line)
}

type Reporter interface {
	Error(line int, msg string) error
}
