package utils

import "fmt"

type StdReporter struct {
}

func (r *StdReporter) Error(line int, msg string) {
	fmt.Printf("Error: %s (line %d)\n", msg, line)
}
