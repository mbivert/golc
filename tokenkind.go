package main

type tokenKind uint

//go:generate go run golang.org/x/tools/cmd/stringer -type tokenKind -linecomment tokenkind.go

const (
	tokenEOF tokenKind = iota // EOF
	tokenError      // error

	// Standard stuff for an untyped λ-calculus
	// XXX Rename tokenName tokenVar
	tokenName       // name
	tokenLambda     // λ

	tokenLParen     // (
	tokenRParen     // )

	tokenDot        // .

	tokenFloat      // float64
	tokenInt        // int64
	tokenBool       // bool

	tokenExcl       // !
	tokenPlus       // +
	tokenMinus      // -
	tokenStar       // *
	tokenSlash      // /
	tokenOr         // |
	tokenAnd        // &
	tokenMoreEq     // ≥
	tokenLessEq     // ≤
	tokenFPlus      // +.
	tokenFMinus     // -.
	tokenFStar      // *.
	tokenFSlash     // /.
	tokenFMoreEq    // ≥.
	tokenFLessEq    // ≤.

	tokenLess       // <
	tokenMore       // >

	tokenFLess      // <.
	tokenFMore      // >.

	tokenAndAnd     // &&
	tokenOrOr       // ||

	// Extensions see [Selinger2009]
	tokenColon      // :
	tokenPi         // π

	tokenArrow      // →
	tokenProduct    // ×

	tokenLet        // let
	tokenIn         // in
	tokenRec        // rec

	tokenMatch      // match
	tokenWith       // with

	tokenIf         // if
	tokenThen       // then
	tokenElse       // else

	tokenNew        // new
	tokenMeas       // meas
)
