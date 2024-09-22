# Introduction
Hand-rolled parser, inspired by Go's:

  - [src/go/token/token.go][src/go/token/token.go];
  - [src/go/scanner/scanner.go][src/go/scanner/scanner.go];
  - [src/go/parser/parser.go][src/go/parser/parser.go].

Tokens are implemented as a basic enumeration:

  - [tokenkind.go][gh-mb-golc-tokenkind.go].

The embedded ``go:generate`` stringer provides us with a
``func (i tokenKind) String() string {...}``:

  - [tokenkind_string.go][gh-mb-golc-tokenkind_string.go].

The tokens are initially processed from an [``io.Reader``][godoc-io-reader]
by our ``scanner``, which itself relies on a [``bufio.Scanner``][godoc-bufio-scanner]:

  - [scanner.go][gh-mb-golc-scanner.go];
  - [scanner_test.go][gh-mb-golc-scanner_test.go];

There are a few tricky bits (e.g. a single unicode rune may be split on two
consecutive read); the scanning is way too sophisticated for a toy λ-calculus,
but that's purposeful (having the input stored in a big string would have
been sufficient).

The scanned tokens are then parsed by a ad-hoc parser:

  - [parser.go][gh-mb-golc-parser.go];
  - [parser_test.go][gh-mb-golc-parser_test.go];

**<u>Note:</u>** Currently, the parsing data structures aren't perfectly determined.

The WIP evaluation is shared between some auxiliary utilities and the
actual "evaluation" (reduction):

  - [utils.go][gh-mb-golc-utils.go];
  - [utils_test.go][gh-mb-golc-utils_test.go];
  - [eval.go][gh-mb-golc-eval.go];
  - [eval_test.go][gh-mb-golc-eval_test.go];

[src/go/token/token.go]: https://github.com/golang/go/blob/master/src/go/token/token.go
[src/go/scanner/scanner.go]: https://github.com/golang/go/blob/master/src/go/scanner/scanner.go
[src/go/parser/parser.go]: https://github.com/golang/go/blob/master/src/go/parser/parser.go

[gh-mb-golc-tokenkind.go]: https://github.com/mbivert/golc/blob/master/tokenkind.go
[gh-mb-golc-tokenkind_string.go]: https://github.com/mbivert/golc/blob/master/tokenkind_string.go

[gh-mb-golc-scanner.go]: https://github.com/mbivert/golc/blob/master/scanner.go
[gh-mb-golc-scanner_test.go]: https://github.com/mbivert/golc/blob/master/scanner_test.go

[gh-mb-golc-parser.go]: https://github.com/mbivert/golc/blob/master/parser.go
[gh-mb-golc-parser_test.go]: https://github.com/mbivert/golc/blob/master/parser_test.go

[gh-mb-golc-utils.go]: https://github.com/mbivert/golc/blob/master/utils.go
[gh-mb-golc-utils_test.go]: https://github.com/mbivert/golc/blob/master/utils_test.go

[gh-mb-golc-eval.go]: https://github.com/mbivert/golc/blob/master/eval.go
[gh-mb-golc-eval_test.go]: https://github.com/mbivert/golc/blob/master/eval_test.go

[godoc-io-reader]: https://pkg.go.dev/io#Reader
[godoc-bufio-scanner]: https://pkg.go.dev/bufio#Scanner
