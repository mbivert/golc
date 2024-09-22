package main

import (
	"fmt"
	"strings"
	"testing"
	// "encoding/json"
)

// True
var T = &AbsExpr{
	expr{},
	"x",
	&AbsExpr{
		expr{},
		"y",
		&VarExpr{expr{}, "x"},
	},
}

/*
 * Most of those tests are recycled from
 *	https://github.com/mbivert/nix-series-code/blob/master/lambda/parse_test.nix
 *	https://github.com/mbivert/nix-series-code/blob/master/exprs_test.nix
 *
 * Many more tests to import from:
 *	https://github.com/mbivert/nix-series-code/blob/master/lambda_test.nix
 *
 * We're focusing here on basic lambda calculus extended with some
 * scalar types (eg. bool, int, float) and basic arithmetic operations.
 *
 * TODO: split bare lambda calculus & scalar expressions.
 *
 * Type indications, quantum and syntactical extensions will be tested
 * separately.
 */
func TestTypelessParse(t *testing.T) {
	doTests(t, []test{
		{
			"empty input",
			parse,
			[]interface{}{strings.NewReader(""), ""},
			[]interface{}{nil, fmt.Errorf(":1:1: Unexpected token: EOF")},
		},
		{
			"single int",
			parse,
			[]interface{}{strings.NewReader("  1234"), ""},
			[]interface{}{&IntExpr{expr{&IntType{}}, 1234}, nil},
		},
		{
			"single (int)",
			parse,
			[]interface{}{strings.NewReader("  (1234)"), ""},
			[]interface{}{&IntExpr{expr{&IntType{}}, 1234}, nil},
		},
		{
			"single ((int))",
			parse,
			[]interface{}{strings.NewReader("  ((1234))"), ""},
			[]interface{}{&IntExpr{expr{&IntType{}}, 1234}, nil},
		},
		{
			"single float",
			parse,
			[]interface{}{strings.NewReader("  1234.45 "), ""},
			[]interface{}{&FloatExpr{expr{&FloatType{}}, 1234.45}, nil},
		},
		{
			"single boolean",
			parse,
			[]interface{}{strings.NewReader("  true "), ""},
			[]interface{}{&BoolExpr{expr{&BoolType{}}, true}, nil},
		},
		{
			"single boolean (bis)",
			parse,
			[]interface{}{strings.NewReader("  false "), ""},
			[]interface{}{&BoolExpr{expr{&BoolType{}}, false}, nil},
		},
		/*
			{
				"single int + garbage",
				parse,
				[]interface{}{strings.NewReader("  1234 12"), ""},
				[]interface{}{
					&IntExpr{expr{&IntType{}}, 1234},
					fmt.Errorf(":1:8: Unexpected token: int64"),
				},
			},
		*/
		{
			"two consecutives ints: 'bad' function call, still parses OK",
			parse,
			[]interface{}{strings.NewReader("  1234 12"), ""},
			[]interface{}{
				&AppExpr{expr{}, &IntExpr{expr{&IntType{}}, 1234}, &IntExpr{expr{&IntType{}}, 12}},
				nil,
			},
		},
		{
			"unary expression: -12",
			parse,
			[]interface{}{strings.NewReader("  - 12"), ""},
			[]interface{}{
				&UnaryExpr{expr{}, tokenMinus, &IntExpr{expr{&IntType{}}, 12}},
				nil,
			},
		},
		{
			"unary expression: +.12",
			parse,
			[]interface{}{strings.NewReader("  +. 12"), ""},
			[]interface{}{
				&UnaryExpr{expr{}, tokenFPlus, &IntExpr{expr{&IntType{}}, 12}},
				nil,
			},
		},
		{
			"unary expressions: ++.12",
			parse,
			[]interface{}{strings.NewReader("  ++. 12"), ""},
			[]interface{}{
				&UnaryExpr{
					expr{},
					tokenPlus,
					&UnaryExpr{expr{}, tokenFPlus, &IntExpr{expr{&IntType{}}, 12}},
				},
				nil,
			},
		},
		{
			"single float in parentheses",
			parse,
			[]interface{}{strings.NewReader("  (1234.45) "), ""},
			[]interface{}{&FloatExpr{expr{&FloatType{}}, 1234.45}, nil},
		},
		{
			"single float in two pairs of parentheses",
			parse,
			[]interface{}{strings.NewReader("  (  (1234.45)\t) "), ""},
			[]interface{}{&FloatExpr{expr{&FloatType{}}, 1234.45}, nil},
		},
		{
			"Missing parenthesis",
			parse,
			[]interface{}{strings.NewReader("  (  (1234.45)\t "), ""},
			[]interface{}{
				nil,
				fmt.Errorf(":1:16: Expecting left paren, got: EOF"),
			},
		},
		{
			"single float in two pairs of parentheses, many unary operators",
			parse,
			[]interface{}{strings.NewReader("  +.(  -  (-.-1234.45)\t) "), ""},
			[]interface{}{
				&UnaryExpr{
					expr{},
					tokenFPlus,
					&UnaryExpr{
						expr{},
						tokenMinus,
						&UnaryExpr{
							expr{},
							tokenFMinus,
							&UnaryExpr{
								expr{},
								tokenMinus,
								&FloatExpr{expr{&FloatType{}}, 1234.45},
							},
						},
					},
				},
				nil,
			},
		},
		{
			"left-associativy, addition",
			parse,
			[]interface{}{strings.NewReader("1+2+ 3 "), ""},
			[]interface{}{
				&BinaryExpr{
					expr{},
					tokenPlus,
					&BinaryExpr{
						expr{},
						tokenPlus,
						&IntExpr{expr{&IntType{}}, 1},
						&IntExpr{expr{&IntType{}}, 2},
					},
					&IntExpr{expr{&IntType{}}, 3},
				},
				nil,
			},
		},
		{
			"left-associativy, addition/substraction: 1-42+12 ≠ 1-(42+12)",
			parse,
			[]interface{}{strings.NewReader("1-42+12"), ""},
			[]interface{}{
				&BinaryExpr{
					expr{},
					tokenPlus,
					&BinaryExpr{
						expr{},
						tokenMinus,
						&IntExpr{expr{&IntType{}}, 1},
						&IntExpr{expr{&IntType{}}, 42},
					},
					&IntExpr{expr{&IntType{}}, 12},
				},
				nil,
			},
		},
		{
			"multiplication has precedence over addition",
			parse,
			[]interface{}{strings.NewReader("1.*.2.+. 3. "), ""},
			[]interface{}{
				&BinaryExpr{
					expr{},
					tokenFPlus,
					&BinaryExpr{
						expr{},
						tokenFStar,
						&FloatExpr{expr{&FloatType{}}, 1.},
						&FloatExpr{expr{&FloatType{}}, 2.},
					},
					&FloatExpr{expr{&FloatType{}}, 3.},
				},
				nil,
			},
		},
		{
			"random expression",
			parse,
			[]interface{}{strings.NewReader("1.*.(2.+. 3.)"), ""},
			[]interface{}{
				&BinaryExpr{
					expr{},
					tokenFStar,
					&FloatExpr{expr{&FloatType{}}, 1.},
					&BinaryExpr{
						expr{},
						tokenFPlus,
						&FloatExpr{expr{&FloatType{}}, 2.},
						&FloatExpr{expr{&FloatType{}}, 3.},
					},
				},
				nil,
			},
		},
		{
			"Comparison: x < 3",
			parse,
			[]interface{}{strings.NewReader("x < 3"), ""},
			[]interface{}{
				&BinaryExpr{
					expr{},
					tokenLess,
					&VarExpr{expr{}, "x"},
					&IntExpr{expr{&IntType{}}, 3},
				},
				nil,
			},
		},
		{
			"Comparison: x < (3+67)",
			parse,
			[]interface{}{strings.NewReader("x < (3 +67)"), ""},
			[]interface{}{
				&BinaryExpr{
					expr{},
					tokenLess,
					&VarExpr{expr{}, "x"},
					&BinaryExpr{
						expr{},
						tokenPlus,
						&IntExpr{expr{&IntType{}}, 3},
						&IntExpr{expr{&IntType{}}, 67},
					},
				},
				nil,
			},
		},
		{
			"1- 3 * 5 + (1 + 34  )/ 3.",
			parse,
			[]interface{}{strings.NewReader("1- 3 * 5 + (1 + 34  )/ 3."), ""},
			[]interface{}{
				&BinaryExpr{
					expr{},
					tokenPlus,
					&BinaryExpr{
						expr{},
						tokenMinus,
						&IntExpr{expr{&IntType{}}, 1},
						&BinaryExpr{
							expr{},
							tokenStar,
							&IntExpr{expr{&IntType{}}, 3},
							&IntExpr{expr{&IntType{}}, 5},
						},
					},
					&BinaryExpr{
						expr{},
						tokenSlash,
						&BinaryExpr{
							expr{},
							tokenPlus,
							&IntExpr{expr{&IntType{}}, 1},
							&IntExpr{expr{&IntType{}}, 34},
						},
						&FloatExpr{expr{&FloatType{}}, 3.},
					},
				},
				nil,
			},
		},
		{
			"0	/ 78 * 12",
			parse,
			[]interface{}{strings.NewReader("0	/ 78 * 12"), ""},
			[]interface{}{
				&BinaryExpr{
					expr{},
					tokenStar,
					&BinaryExpr{
						expr{},
						tokenSlash,
						&IntExpr{expr{&IntType{}}, 0},
						&IntExpr{expr{&IntType{}}, 78},
					},
					&IntExpr{expr{&IntType{}}, 12},
				},
				nil,
			},
		},
		{
			"basic abstraction",
			parse,
			[]interface{}{strings.NewReader("λx.x"), ""},
			[]interface{}{
				&AbsExpr{
					expr{},
					"x",
					&VarExpr{expr{}, "x"},
				},
				nil,
			},
		},
		{
			"basic abstraction (optional λ)",
			parse,
			[]interface{}{strings.NewReader("x.x"), ""},
			[]interface{}{
				&AbsExpr{
					expr{},
					"x",
					&VarExpr{expr{}, "x"},
				},
				nil,
			},
		},
		{
			"abstraction: syntax error (missing dot)",
			parse,
			[]interface{}{strings.NewReader("\nλx x"), ""},
			[]interface{}{
				nil,
				fmt.Errorf(":2:4: Expecting dot after lambda variable name, got: name"),
			},
		},
		{
			"abstraction: syntax error (missing variable name)",
			parse,
			[]interface{}{strings.NewReader("\nλ.x x"), ""},
			[]interface{}{
				nil,
				fmt.Errorf(":2:2: Expecting variable name after lambda, got: ."),
			},
		},
		{
			"abstraction: syntax error (missing variable name)",
			parse,
			[]interface{}{strings.NewReader("λx. "), ""},
			[]interface{}{
				nil,
				fmt.Errorf(":1:5: Unexpected token: EOF"),
			},
		},
		{
			"(λx. x x) (λx. x x)",
			parse,
			[]interface{}{strings.NewReader("(λx. x x) (λx. x x)"), ""},
			[]interface{}{
				&AppExpr{
					expr{},
					&AbsExpr{
						expr{},
						"x",
						&AppExpr{
							expr{},
							&VarExpr{expr{}, "x"},
							&VarExpr{expr{}, "x"},
						},
					},
					&AbsExpr{
						expr{},
						"x",
						&AppExpr{
							expr{},
							&VarExpr{expr{}, "x"},
							&VarExpr{expr{}, "x"},
						},
					},
				},
				nil,
			},
		},
		{
			"(λx. (λy. x))",
			parse,
			[]interface{}{strings.NewReader("(λx. (λy. x))"), ""},
			[]interface{}{T, nil},
		},
		{
			"(λx. λy. x)",
			parse,
			[]interface{}{strings.NewReader("(λx. λy. x)"), ""},
			[]interface{}{T, nil},
		},
		{
			"λx.λy.x",
			parse,
			[]interface{}{strings.NewReader("λx.λy.x"), ""},
			[]interface{}{T, nil},
		},
		{
			"x. (one (two (three (four five))))",
			parse,
			[]interface{}{strings.NewReader("x. (one (two (three (four five))))"), ""},
			[]interface{}{
				&AbsExpr{
					expr{},
					"x",
					&AppExpr{
						expr{},
						&VarExpr{expr{}, "one"},
						&AppExpr{
							expr{},
							&VarExpr{expr{}, "two"},
							&AppExpr{
								expr{},
								&VarExpr{expr{}, "three"},
								&AppExpr{
									expr{},
									&VarExpr{expr{}, "four"},
									&VarExpr{expr{}, "five"},
								},
							},
						},
					},
				},
				nil,
			},
		},
		{
			"x. one two three four five",
			parse,
			[]interface{}{strings.NewReader("x. one two three four five"), ""},
			[]interface{}{
				&AbsExpr{
					expr{},
					"x",
					&AppExpr{
						expr{},
						&AppExpr{
							expr{},
							&AppExpr{
								expr{},
								&AppExpr{
									expr{},
									&VarExpr{expr{}, "one"},
									&VarExpr{expr{}, "two"},
								},
								&VarExpr{expr{}, "three"},
							},
							&VarExpr{expr{}, "four"},
						},
						&VarExpr{expr{}, "five"},
					},
				},
				nil,
			},
		},
	})
}

func TestPrimitiveType(t *testing.T) {
	doTests(t, []test{
		{
			"boolean",
			parse,
			[]interface{}{strings.NewReader("λx : bool . x && y"), ""},
			[]interface{}{
				&AbsExpr{
					expr{&BoolType{}},
					"x",
					&BinaryExpr{
						expr{},
						tokenAndAnd,
						&VarExpr{expr{}, "x"},
						&VarExpr{expr{}, "y"},
					},
				},
				nil,
			},
		},
		{
			"int",
			parse,
			[]interface{}{strings.NewReader("λx : int . x + y"), ""},
			[]interface{}{
				&AbsExpr{
					expr{&IntType{}},
					"x",
					&BinaryExpr{
						expr{},
						tokenPlus,
						&VarExpr{expr{}, "x"},
						&VarExpr{expr{}, "y"},
					},
				},
				nil,
			},
		},
		{
			"float",
			parse,
			[]interface{}{strings.NewReader("λx : float . x +. y"), ""},
			[]interface{}{
				&AbsExpr{
					expr{&FloatType{}},
					"x",
					&BinaryExpr{
						expr{},
						tokenFPlus,
						&VarExpr{expr{}, "x"},
						&VarExpr{expr{}, "y"},
					},
				},
				nil,
			},
		},
		{
			"float",
			parse,
			[]interface{}{strings.NewReader("λx : unit. *"), ""},
			[]interface{}{
				&AbsExpr{
					expr{&UnitType{}},
					"x",
					&UnitExpr{expr{&UnitType{}}},
				},
				nil,
			},
		},
	})
}

func TestArrowType(t *testing.T) {
	doTests(t, []test{
		{
			"bool → bool",
			parse,
			[]interface{}{strings.NewReader("λx : bool → bool . x y"), ""},
			[]interface{}{
				&AbsExpr{
					expr{&ArrowType{typ{}, &BoolType{}, &BoolType{}}},
					"x",
					&AppExpr{
						expr{},
						&VarExpr{expr{}, "x"},
						&VarExpr{expr{}, "y"},
					},
				},
				nil,
			},
		},
		{
			"bool → bool → bool (right associative: bool → (bool → bool))",
			parse,
			[]interface{}{strings.NewReader("λx : bool → bool → bool . x y z"), ""},
			[]interface{}{
				&AbsExpr{
					expr{&ArrowType{
						typ{}, &BoolType{}, &ArrowType{
							typ{}, &BoolType{}, &BoolType{},
						},
					}},
					"x",
					&AppExpr{
						expr{},
						&AppExpr{
							expr{},
							&VarExpr{expr{}, "x"},
							&VarExpr{expr{}, "y"},
						},
						&VarExpr{expr{}, "z"},
					},
				},
				nil,
			},
		},
		{
			"bool → bool → bool → int",
			parse,
			[]interface{}{strings.NewReader("λx : bool → bool → bool → int . (x y z) + 3"), ""},
			[]interface{}{
				&AbsExpr{
					expr{&ArrowType{
						typ{}, &BoolType{}, &ArrowType{
							typ{}, &BoolType{}, &ArrowType{
								typ{}, &BoolType{}, &IntType{},
							},
						},
					}},
					"x",
					&BinaryExpr{
						expr{},
						tokenPlus,
						&AppExpr{
							expr{},
							&AppExpr{
								expr{},
								&VarExpr{expr{}, "x"},
								&VarExpr{expr{}, "y"},
							},
							&VarExpr{expr{}, "z"},
						},
						&IntExpr{expr{&IntType{}}, 3},
					},
				},
				nil,
			},
		},
		{
			"(bool → bool) → bool (manually altered associativity)",
			parse,
			[]interface{}{strings.NewReader("λx : (bool → bool) → bool . x y z"), ""},
			[]interface{}{
				&AbsExpr{
					expr{&ArrowType{
						typ{}, &ArrowType{
							typ{}, &BoolType{}, &BoolType{},
						},
						&BoolType{},
					}},
					"x",
					&AppExpr{
						expr{},
						&AppExpr{
							expr{},
							&VarExpr{expr{}, "x"},
							&VarExpr{expr{}, "y"},
						},
						&VarExpr{expr{}, "z"},
					},
				},
				nil,
			},
		},
	})
}
