package utils

import "fmt"

type StdReporter struct {
	Cnt int
}

func (r *StdReporter) Error(line int, msg string) {
	r.Cnt += 1
	fmt.Printf("Error: %s (line %d)\n", msg, line)
}
