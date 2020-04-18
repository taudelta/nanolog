.PHONY: lint test clear run-simple-example run-file-example

all: lint test
run-simple: lint test run-simple-example
run-file: lint test run-file-example

clear:
	rm -rf ./*.log

lint:
	golangci-lint run

test:
	go test -v .

run-simple-example:
	go run ./examples/simple/main.go

run-file-example:
	go run ./examples/file/main.go
