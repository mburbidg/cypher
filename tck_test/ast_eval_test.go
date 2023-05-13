package tck_test

import "github.com/mburbidg/cypher/parser"

type astRuntime struct {
}

func (runtime *astRuntime) eval(stmt parser.Statement) error {
	return nil
}
