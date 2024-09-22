.PHONY: scanner-tests
scanner-tests: tokenkind_string.go
	@echo Runing scanner tests...
	@go test -v scanner_test.go scanner.go tokenkind.go tokenkind_string.go

.PHONY: parser-tests
parser-tests: tokenkind_string.go
	@echo Running parser tests...
	@go test -v parser_test.go parser.go scanner.go tokenkind.go tokenkind_string.go

.PHONY: eval-tests
eval-tests: tokenkind_string.go
	@echo Running eval tests...
	@go test -v eval_test.go eval.go parser.go scanner.go tokenkind.go tokenkind_string.go

.PHONY: utils-tests
utils-tests: tokenkind_string.go
	@echo Running utils tests...
	@go test -v utils_test.go utils.go parser.go scanner.go tokenkind.go tokenkind_string.go

.PHONY: tests
tests: scanner-tests parser-tests eval-tests utils-tests

tokenkind_string.go: tokenkind.go
	@echo Generating $@...
	@go generate $<

