package main

import (
	"encoding/json"
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

func TestScannerScanAll(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"empty input",
			scanAll,
			[]any{"", ""},
			[]any{[]token{token{tokenEOF, 1, 1, ""}}, nil},
		},
		{
			"spaces",
			scanAll,
			[]any{"  \t\t\r\n", ""},
			[]any{[]token{token{tokenEOF, 2, 1, ""}}, nil},
		},
		{
			"single byte tokens",
			scanAll,
			[]any{"  \t\t\r\n().  :<× >", ""},
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
			[]any{"hello world", ""},
			[]any{[]token{
				token{tokenName, 1, 1, "hello"},
				token{tokenName, 1, 7, "world"},
				token{tokenEOF, 1, 12, ""},
			}, nil},
		},
		{
			"ifelse",
			scanAll,
			[]any{"\n(λp. λx. λy. p x y)", ""},
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
			[]any{"(λp:bool -> bool . M)", ""},
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
			[]any{"(λp:bool->  bool . M)", ""},
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
			[]any{"(λp:bool →  bool→. M)", ""},
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
			[]any{"0123", ""},
			[]any{[]token{
				token{tokenInt, 1, 1, "0123"},
				token{tokenEOF, 1, 5, ""},
			}, nil},
		},
		{
			"allow unusual number parsing terminator",
			scanAll,
			[]any{"0123aaa", ""},
			[]any{[]token{
				token{tokenInt, 1, 1, "0123"},
				token{tokenName, 1, 5, "aaa"},
				token{tokenEOF, 1, 8, ""},
			}, nil},
		},
		{
			"numbers",
			scanAll,
			[]any{"0123 123 123.46 .10 .", ""},
			[]any{[]token{
				token{tokenInt, 1, 1, "0123"},
				token{tokenInt, 1, 6, "123"},
				token{tokenFloat, 1, 10, "123.46"},
				token{tokenFloat, 1, 17, ".10"},
				token{tokenDot, 1, 21, "."},
				token{tokenEOF, 1, 22, ""},
			}, nil},
		},
		{
			"two-bytes operators, slashes",
			scanAll,
			[]any{"+. + . //.", ""},
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
			[]any{"false let truer\ttrue", ""},
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
			[]any{"foo||bar&&", ""},
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

func TestScannerProduct(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"〈〉",
			scanAll,
			[]any{"〈〉", ""},
			[]any{[]token{
				token{tokenLBracket, 1, 1, "〈"},
				token{tokenRBracket, 1, 2, "〉"},
				token{tokenEOF, 1, 3, ""},
			}, nil},
		},
		{
			"〈X〉",
			scanAll,
			[]any{"〈X〉", ""},
			[]any{[]token{
				token{tokenLBracket, 1, 1, "〈"},
				token{tokenName, 1, 2, "X"},
				token{tokenRBracket, 1, 3, "〉"},
				token{tokenEOF, 1, 4, ""},
			}, nil},
		},
		{
			"〈X,   Y〉",
			scanAll,
			[]any{"〈X,   Y〉", ""},
			[]any{[]token{
				token{tokenLBracket, 1, 1, "〈"},
				token{tokenName, 1, 2, "X"},
				token{tokenComa, 1, 3, ","},
				token{tokenName, 1, 7, "Y"},
				token{tokenRBracket, 1, 8, "〉"},
				token{tokenEOF, 1, 9, ""},
			}, nil},
		},
	})
}

func TestScannerExcl(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"!true",
			scanAll,
			[]any{"!true", ""},
			[]any{[]token{
				token{tokenExcl, 1, 1, "!"},
				token{tokenBool, 1, 2, "true"},
				token{tokenEOF, 1, 6, ""},
			}, nil},
		},
	})
}

func TestScannerFCmpOp(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"3.≤.5.",
			scanAll,
			[]any{"3.≤.5.", ""},
			[]any{[]token{
				token{tokenFloat, 1, 1, "3."},
				token{tokenFLessEq, 1, 3, "≤."},
				token{tokenFloat, 1, 5, "5."},
				token{tokenEOF, 1, 7, ""},
			}, nil},
		},
		{
			"≤x",
			scanAll,
			[]any{"≤x", ""},
			[]any{[]token{
				token{tokenLessEq, 1, 1, "≤"},
				token{tokenName, 1, 2, "x"},
				token{tokenEOF, 1, 3, ""},
			}, nil},
		},
	})
}

func TestScannerIdentifier(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"x0 y01",
			scanAll,
			[]any{"x0 y01λ", ""},
			[]any{[]token{
				token{tokenName, 1, 1, "x0"},
				token{tokenName, 1, 4, "y01"},
				token{tokenLambda, 1, 7, "λ"},
				token{tokenEOF, 1, 8, ""},
			}, nil},
		},
	})
}
