# Introduction
This repository contains an implementation of a
[typed λ-calculus][wp-en-typed-lambda-calculus], extended to
support elementary arithmetic expressions, some syntactic sugar,
and more importantly, implements the [quantum λ-calculus][qlc][^0]
described by Peter Selinger. One might also be interested in
[Benoît Valiron's thesis][benoit-valiron-thesis].

Implemented in Go; the parsing/scanning stages mimic, for practice,
what has been done in the Go compiler itself.

**<u>Note:</u>** This is essentially building on a previous
[λ-calculus interpreter][tales-lambda], written in Nix's language.
In particular, I'm recylcing some of the tests from this earlier
implementation.

[^0]: Beware, there are multiple papers describing quantum
λ-calculus authored by Selinger, e.g. [this one][qlc].

[qlc]: https://www.mscs.dal.ca/~selinger/papers/qlambdabook.pdf
[benoit-valiron-thesis]: https://theses.hal.science/tel-00483944
[wp-en-typed-lambda-calculus]: https://en.wikipedia.org/wiki/Typed_lambda_calculus
[tales-lambda]: https://tales.mbivert.com/on-nix-language-lambda-calculus/

<!--

https://okmij.org/ftp/ML/generalization.html

https://nostarch.com/writing-c-compiler
https://github.com/rui314/chibicc

https://compilerbook.com/
https://interpreterbook.com/

-->