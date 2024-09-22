package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/mbivert/ftests"
)

// unexported fields are unavailable to the json
// package, and thus aren't visible in tests...
func (t token) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Kind   string
		LN, CN uint
		Raw    string
	}{
		Kind: (t.kind).String(),
		LN:   t.ln,
		CN:   t.cn,
		Raw:  t.raw,
	})
}

func TestScanAll(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"empty input",
			scanAll,
			[]any{strings.NewReader(""), ""},
			[]any{[]token{token{tokenEOF, 1, 1, ""}}, nil},
		},
		{
			"spaces",
			scanAll,
			[]any{strings.NewReader("  \t\t\r\n"), ""},
			[]any{[]token{token{tokenEOF, 2, 1, ""}}, nil},
		},
		{
			"single byte tokens",
			scanAll,
			[]any{strings.NewReader("  \t\t\r\n().  :<× >"), ""},
			[]any{[]token{
				token{tokenLParen, 2, 1, "("},
				token{tokenRParen, 2, 2, ")"},
				token{tokenDot, 2, 3, "."},
				token{tokenColon, 2, 6, ":"},
				token{tokenLess, 2, 7, "<"},
				token{tokenProduct, 2, 8, "×"},
				token{tokenMore, 2, 10, ">"},
				token{tokenEOF, 2, 11, ""},
			}, nil},
		},
		{
			"multi-bytes words",
			scanAll,
			[]any{strings.NewReader("hello world"), ""},
			[]any{[]token{
				token{tokenName, 1, 1, "hello"},
				token{tokenName, 1, 7, "world"},
				token{tokenEOF, 1, 12, ""},
			}, nil},
		},
		{
			"ifelse",
			scanAll,
			[]any{strings.NewReader("\n(λp. λx. λy. p x y)"), ""},
			[]any{[]token{
				token{tokenLParen, 2, 1, "("},
				token{tokenLambda, 2, 2, "λ"},
				token{tokenName, 2, 3, "p"},
				token{tokenDot, 2, 4, "."},

				token{tokenLambda, 2, 6, "λ"},
				token{tokenName, 2, 7, "x"},
				token{tokenDot, 2, 8, "."},

				token{tokenLambda, 2, 10, "λ"},
				token{tokenName, 2, 11, "y"},
				token{tokenDot, 2, 12, "."},

				token{tokenName, 2, 14, "p"},
				token{tokenName, 2, 16, "x"},
				token{tokenName, 2, 18, "y"},
				token{tokenRParen, 2, 19, ")"},

				token{tokenEOF, 2, 20, ""},
			}, nil},
		},
		{
			"arrow",
			scanAll,
			[]any{strings.NewReader("(λp:bool -> bool . M)"), ""},
			[]any{[]token{
				token{tokenLParen, 1, 1, "("},
				token{tokenLambda, 1, 2, "λ"},
				token{tokenName, 1, 3, "p"},
				token{tokenColon, 1, 4, ":"},
				token{tokenTBool, 1, 5, "bool"},
				token{tokenArrow, 1, 10, "->"},
				token{tokenTBool, 1, 13, "bool"},
				token{tokenDot, 1, 18, "."},
				token{tokenName, 1, 20, "M"},
				token{tokenRParen, 1, 21, ")"},
				token{tokenEOF, 1, 22, ""},
			}, nil},
		},
		{
			"arrow (bis)",
			scanAll,
			[]any{strings.NewReader("(λp:bool->  bool . M)"), ""},
			[]any{[]token{
				token{tokenLParen, 1, 1, "("},
				token{tokenLambda, 1, 2, "λ"},
				token{tokenName, 1, 3, "p"},
				token{tokenColon, 1, 4, ":"},
				token{tokenTBool, 1, 5, "bool"},
				token{tokenArrow, 1, 9, "->"},
				token{tokenTBool, 1, 13, "bool"},
				token{tokenDot, 1, 18, "."},
				token{tokenName, 1, 20, "M"},
				token{tokenRParen, 1, 21, ")"},
				token{tokenEOF, 1, 22, ""},
			}, nil},
		},
		{
			"arrow (bis)",
			scanAll,
			[]any{strings.NewReader("(λp:bool →  bool→. M)"), ""},
			[]any{[]token{
				token{tokenLParen, 1, 1, "("},
				token{tokenLambda, 1, 2, "λ"},
				token{tokenName, 1, 3, "p"},
				token{tokenColon, 1, 4, ":"},
				token{tokenTBool, 1, 5, "bool"},
				token{tokenArrow, 1, 10, "→"},
				token{tokenTBool, 1, 13, "bool"},
				token{tokenArrow, 1, 17, "→"},
				token{tokenDot, 1, 18, "."},
				token{tokenName, 1, 20, "M"},
				token{tokenRParen, 1, 21, ")"},
				token{tokenEOF, 1, 22, ""},
			}, nil},
		},
		{
			"isolated integer",
			scanAll,
			[]any{strings.NewReader("0123"), ""},
			[]any{[]token{
				token{tokenInt, 1, 1, "0123"},
				token{tokenEOF, 1, 5, ""},
			}, nil},
		},
		{
			"allow unusual number parsing terminator",
			scanAll,
			[]any{strings.NewReader("0123aaa"), ""},
			[]any{[]token{
				token{tokenInt, 1, 1, "0123"},
				token{tokenName, 1, 5, "aaa"},
				token{tokenEOF, 1, 8, ""},
			}, nil},
		},
		{
			"numbers",
			scanAll,
			[]any{strings.NewReader("0123 123 123.46 .10 ."), ""},
			[]any{[]token{
				token{tokenInt, 1, 1, "0123"},
				token{tokenInt, 1, 6, "123"},
				token{tokenFloat, 1, 10, "123.46"},
				token{tokenFloat, 1, 16, ".10"},
				token{tokenDot, 1, 19, "."},
				token{tokenEOF, 1, 20, ""},
			}, nil},
		},
		{
			"two-bytes operators, slashes",
			scanAll,
			[]any{strings.NewReader("+. + . //."), ""},
			[]any{[]token{
				token{tokenFPlus, 1, 1, "+."},
				token{tokenPlus, 1, 4, "+"},
				token{tokenDot, 1, 6, "."},
				token{tokenSlash, 1, 8, "/"},
				token{tokenFSlash, 1, 9, "/."},
				token{tokenEOF, 1, 11, ""},
			}, nil},
		},
		{
			"multi-byte known tokens",
			scanAll,
			[]any{strings.NewReader("false let truer\ttrue"), ""},
			[]any{[]token{
				token{tokenBool, 1, 1, "false"},
				token{tokenLet, 1, 7, "let"},
				token{tokenName, 1, 11, "truer"},
				token{tokenBool, 1, 17, "true"},
				token{tokenEOF, 1, 21, ""},
			}, nil},
		},
		{
			"'twos' are reckognized as a separator (was a bug)",
			scanAll,
			[]any{strings.NewReader("foo||bar&&"), ""},
			[]any{[]token{
				token{tokenName, 1, 1, "foo"},
				token{tokenOrOr, 1, 4, "||"},
				token{tokenName, 1, 6, "bar"},
				token{tokenAndAnd, 1, 9, "&&"},
				token{tokenEOF, 1, 11, ""},
			}, nil},
		},
	})
}
