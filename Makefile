.PHONY: scanner-tests
scanner-tests: tokenkind_string.go
	@echo Runing scanner tests...
	@go test -v -run TestScanner

.PHONY: parser-tests
parser-tests: tokenkind_string.go
	@echo Running parser tests...
	@go test -v -run TestParser

.PHONY: eval-tests
eval-tests: tokenkind_string.go
	@echo Running eval tests...
	@go test -v -run TestEval

.PHONY: utils-tests
utils-tests: tokenkind_string.go
	@echo Running utils tests...
	@go test -v -run TestUtils

.PHONY: typing-tests
typing-tests: tokenkind_string.go
	@echo Running typing tests...
	@go test -v -run TestTyping

.PHONY: tests
tests:
	@echo Running tests...
	@go test -v .

tokenkind_string.go: tokenkind.go
	@echo Generating $@...
	@go generate $<

