# Introduction
We're using a hand-rolled parser modelled after Go's:

  - [src/go/token/token.go][src/go/token/token.go];
  - [src/go/scanner/scanner.go][src/go/scanner/scanner.go];
  - [src/go/parser/parser.go][src/go/parser/parser.go].

Our tokens are implemented as a basic enumeration:

  - [tokenkind.go][gh-mb-golc-tokenkind.go].

The embedded ``go:generate`` stringer provides us with a
``func (i tokenKind) String() string {...}``:

  - [tokenkind_string.go][gh-mb-golc-tokenkind.go].

The tokens are initially processed from an [``io.Reader``][godoc-io-reader]
by our ``scanner``, which itself relies on a [``bufio.Scanner``][godoc-bufio-scanner]:

  - [scanner.go][gh-mb-golc-scanner.go];
  - [scanner_tests.go][gh-mb-golc-scanner_tests.go];

There are a few tricky bits; the scanning is way too sophisticated for
a toy lambda calculus, but I wanted to have something operating on a random
``io.Reader`` and not just a static string.

The scanned tokens are then parsed by a ad-hoc parser:

  - [parser.go][gh-mb-golc-parser.go];
  - [parser_tests.go][gh-mb-golc-parser_tests.go];

Currently, the parsing data structures aren't perfectly determined.

[src/go/token/token.go]: https://github.com/golang/go/blob/master/src/go/token/token.go
[src/go/scanner/scanner.go]: https://github.com/golang/go/blob/master/src/go/scanner/scanner.go
[src/go/parser/parser.go]: https://github.com/golang/go/blob/master/src/go/parser/parser.go

[gh-mb-golc-tokenkind.go]: https://github.com/mbiver/golc/blob/master/tokenkind.go
[gh-mb-golc-tokenkind_string.go]: https://github.com/mbiver/golc/blob/master/tokenkind_string.go

[gh-mb-golc-scanner.go]: https://github.com/mbiver/golc/blob/master/scanner.go
[gh-mb-golc-scanner_tests.go]: https://github.com/mbiver/golc/blob/master/scanner_tests.go

[gh-mb-golc-parser.go]: https://github.com/mbiver/golc/blob/master/parser.go
[gh-mb-golc-parser_tests.go]: https://github.com/mbiver/golc/blob/master/parser_tests.go

[godoc-io-reader]: https://pkg.go.dev/io#Reader
[godoc-bufio-scanner]: https://pkg.go.dev/bufio#Scanner
