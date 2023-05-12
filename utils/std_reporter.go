package utils

import "fmt"

type StdReporter struct {
	Cnt int
}

func (r *StdReporter) Error(line int, msg string) error {
	r.Cnt += 1
	err := &ParseError{line, msg}
	fmt.Printf("%s\n", err)
	return err
}
