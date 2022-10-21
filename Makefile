.DEFAULT_GOAL := test

test:
	go test -v github.com/mburbidg/cypher/pkg/scanner
	go test -v github.com/mburbidg/cypher/pkg/parser
.PHONY:test
