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

The tokens are scanned from a ``[]byte`` buffer:

  - [scanner.go][gh-mb-golc-scanner.go];
  - [scanner_test.go][gh-mb-golc-scanner_test.go];

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

Simple type inference (as in, simply-typed λ-calculus) can be
found in:

  - [styping.go][gh-mb-golc-styping.go];
  - [styping_test.go][gh-mb-golc-styping_test.go];

Budding polymorphic type inference in (not sure we'll need to
be this sophisticated?):

  - [typing.go][gh-mb-golc-typing.go];
  - [typing_test.go][gh-mb-golc-typing_test.go];


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

[gh-mb-golc-styping.go]: https://github.com/mbivert/golc/blob/master/styping.go
[gh-mb-golc-styping_test.go]: https://github.com/mbivert/golc/blob/master/styping_test.go

[gh-mb-golc-typing.go]: https://github.com/mbivert/golc/blob/master/typing.go
[gh-mb-golc-typing_test.go]: https://github.com/mbivert/golc/blob/master/typing_test.go


