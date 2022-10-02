.DEFAULT_GOAL := test

test:
	go test -v github.com/mburbidg/cypher/pkg/scanner
.PHONY:test
