.DEFAULT_GOAL := test

test:
	go test -v github.com/mburbidg/cypher/scanner
	go test -v github.com/mburbidg/cypher/parser
.PHONY:test

tck-test:
	go test -v github.com/mburbidg/cypher/tck-test
.PHONY:tck-test
