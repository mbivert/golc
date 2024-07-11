.PHONY: scanner-tests
scanner-tests: tokenkind_string.go
	@echo Runing scanner tests...
	@go test -v scanner_test.go scanner.go ftests.go tokenkind.go tokenkind_string.go

.PHONY: parser-tests
parser-tests: tokenkind_string.go
	@echo Running parser tests...
	@go test -v parser_test.go parser.go scanner.go ftests.go tokenkind.go tokenkind_string.go

.PHONY: tests
tests: scanner-tests parser-tests

tokenkind_string.go: tokenkind.go
	@echo Generating $@...
	@go generate $<

