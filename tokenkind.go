package main

type tokenKind uint

//go:generate go run golang.org/x/tools/cmd/stringer -type tokenKind -linecomment tokenkind.go

const (
	tokenEOF   tokenKind = iota // EOF
	tokenError                  // error

	// Standard stuff for an untyped λ-calculus
	// XXX Rename tokenName tokenVar
	tokenName   // name
	tokenLambda // λ

	tokenLParen // (
	tokenRParen // )

	tokenDot // .

	tokenFloat // float64
	tokenInt   // int64
	tokenBool  // bool

	// XXX meh, potential confusion (stringers),
	// hopefully benign.
	tokenTBool  // bool
	tokenTInt   // int
	tokenTFloat // float
	tokenTUnit  // unit

	tokenExcl // !

	tokenPlus  // +
	tokenFPlus // +.

	tokenMinus  // -
	tokenFMinus // -.

	tokenStar   // *
	tokenFStar  // *.
	tokenSlash  // /
	tokenFSlash // /.

	tokenLess  // <
	tokenFLess // <.
	tokenMore  // >
	tokenFMore // >.

	tokenComa  // ,
	tokenEqual // =

	tokenLBracket // 〈
	tokenRBracket // 〉

	tokenOr     // |
	tokenOrOr   // ||
	tokenAnd    // &
	tokenAndAnd // &&

	tokenMoreEq  // ≥
	tokenFMoreEq // ≥.
	tokenLessEq  // ≤
	tokenFLessEq // ≤.

	tokenColon // :
	tokenPi    // π

	tokenArrow   // →
	tokenProduct // ×

	tokenLet // let
	tokenIn  // in
	tokenRec // rec

	tokenMatch // match
	tokenWith  // with

	tokenIf   // if
	tokenThen // then
	tokenElse // else

	tokenNew  // new
	tokenMeas // meas
)
