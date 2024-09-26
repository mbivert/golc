package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/mbivert/ftests"
)

// unexported fields are unavailable to the json
// package, and thus aren't visible in tests...
func (e *IntExpr) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		T string
		V int64
	}{
		T: "int",
		V: e.v,
	})
}

func (e *FloatExpr) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		T string
		V float64
	}{
		T: "float",
		V: e.v,
	})
}

func (e *BoolExpr) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		T string
		V bool
	}{
		T: "bool",
		V: e.v,
	})
}

func (e *AbsExpr) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		T     string
		Typ   Type
		Name  string
		Right Expr
	}{
		T:     "abs",
		Typ:   e.typ,
		Name:  e.name,
		Right: e.right,
	})
}

func (e *AppExpr) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		T     string
		Left  Expr
		Right Expr
	}{
		T:     "app",
		Left:  e.left,
		Right: e.right,
	})
}

func (e *VarExpr) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		T    string
		Name string
	}{
		T:    "var",
		Name: e.name,
	})
}

func (e *BinaryExpr) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		T     string
		Op    string
		Left  Expr
		Right Expr
	}{
		T:     "bin",
		Op:    e.op.String(),
		Left:  e.left,
		Right: e.right,
	})
}

/*
 * Many of those tests are recycled from
 *	https://github.com/mbivert/nix-series-code/blob/master/lambda/parse_test.nix
 *	https://github.com/mbivert/nix-series-code/blob/master/exprs_test.nix
 *
 * More tests to import from:
 *	https://github.com/mbivert/nix-series-code/blob/master/lambda_test.nix
 *
 * We're focusing here on basic lambda calculus extended with some
 * scalar types (eg. bool, int, float) and basic arithmetic operations.
 */

func TestParserMathExprs(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"empty input",
			parse,
			[]any{strings.NewReader(""), ""},
			[]any{nil, fmt.Errorf(":1:1: Unexpected token: EOF")},
		},
		{
			"single int",
			parse,
			[]any{strings.NewReader("  1234"), ""},
			[]any{&IntExpr{expr{&IntType{}}, 1234}, nil},
		},
		{
			"single (int)",
			parse,
			[]any{strings.NewReader("  (1234)"), ""},
			[]any{&IntExpr{expr{&IntType{}}, 1234}, nil},
		},
		{
			"single ((int))",
			parse,
			[]any{strings.NewReader("  ((1234))"), ""},
			[]any{&IntExpr{expr{&IntType{}}, 1234}, nil},
		},
		{
			"single float",
			parse,
			[]any{strings.NewReader("  1234.45 "), ""},
			[]any{&FloatExpr{expr{&FloatType{}}, 1234.45}, nil},
		},
		{
			"single boolean",
			parse,
			[]any{strings.NewReader("  true "), ""},
			[]any{&BoolExpr{expr{&BoolType{}}, true}, nil},
		},
		{
			"single boolean (bis)",
			parse,
			[]any{strings.NewReader("  false "), ""},
			[]any{&BoolExpr{expr{&BoolType{}}, false}, nil},
		},
		// NOTE: this will be rejected during the type inference/checking phase
		{
			"two consecutives ints: 'bad' function call, still parses OK",
			parse,
			[]any{strings.NewReader("  1234 12"), ""},
			[]any{
				&AppExpr{expr{}, &IntExpr{expr{&IntType{}}, 1234}, &IntExpr{expr{&IntType{}}, 12}},
				nil,
			},
		},
		{
			"unary expression: -12",
			parse,
			[]any{strings.NewReader("  - 12"), ""},
			[]any{
				&UnaryExpr{expr{}, tokenMinus, &IntExpr{expr{&IntType{}}, 12}},
				nil,
			},
		},
		{
			"unary expression: +.12",
			parse,
			[]any{strings.NewReader("  +. 12"), ""},
			[]any{
				&UnaryExpr{expr{}, tokenFPlus, &IntExpr{expr{&IntType{}}, 12}},
				nil,
			},
		},
		{
			"unary expressions: ++.12",
			parse,
			[]any{strings.NewReader("  ++. 12"), ""},
			[]any{
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
			[]any{strings.NewReader("  (1234.45) "), ""},
			[]any{&FloatExpr{expr{&FloatType{}}, 1234.45}, nil},
		},
		{
			"single float in two pairs of parentheses",
			parse,
			[]any{strings.NewReader("  (  (1234.45)\t) "), ""},
			[]any{&FloatExpr{expr{&FloatType{}}, 1234.45}, nil},
		},
		{
			"Missing parenthesis",
			parse,
			[]any{strings.NewReader("  (  (1234.45)\t "), ""},
			[]any{
				nil,
				fmt.Errorf(":1:16: Expecting left paren, got: EOF"),
			},
		},
		{
			"single float in two pairs of parentheses, many unary operators",
			parse,
			[]any{strings.NewReader("  +.(  -  (-.-1234.45)\t) "), ""},
			[]any{
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
			[]any{strings.NewReader("1+2+ 3 "), ""},
			[]any{
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
			[]any{strings.NewReader("1-42+12"), ""},
			[]any{
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
			[]any{strings.NewReader("1.*.2.+. 3. "), ""},
			[]any{
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
			[]any{strings.NewReader("1.*.(2.+. 3.)"), ""},
			[]any{
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
			[]any{strings.NewReader("x < 3"), ""},
			[]any{
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
			[]any{strings.NewReader("x < (3 +67)"), ""},
			[]any{
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
			[]any{strings.NewReader("1- 3 * 5 + (1 + 34  )/ 3."), ""},
			[]any{
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
			[]any{strings.NewReader("0	/ 78 * 12"), ""},
			[]any{
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
			"!(true)",
			parse,
			[]any{strings.NewReader("!(true)"), ""},
			[]any{
				&UnaryExpr{
					expr{},
					tokenExcl,
					&BoolExpr{expr{&BoolType{typ{}}}, true},
				},
				nil,
			},
		},
	})
}

func TestParserUntypedλCalc(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"basic abstraction",
			parse,
			[]any{strings.NewReader("λx.x"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&typ{},
					"x",
					&VarExpr{expr{}, "x"},
				},
				nil,
			},
		},
		{
			"basic abstraction (optional λ)",
			parse,
			[]any{strings.NewReader("x.x"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&typ{},
					"x",
					&VarExpr{expr{}, "x"},
				},
				nil,
			},
		},
		{
			"abstraction: syntax error (missing dot)",
			parse,
			[]any{strings.NewReader("\nλx x"), ""},
			[]any{
				nil,
				fmt.Errorf(":2:4: Expecting dot after lambda variable name, got: name"),
			},
		},
		{
			"abstraction: syntax error (missing variable name)",
			parse,
			[]any{strings.NewReader("\nλ.x x"), ""},
			[]any{
				nil,
				fmt.Errorf(":2:2: Expecting variable name after lambda, got: ."),
			},
		},
		{
			"abstraction: syntax error (missing variable name)",
			parse,
			[]any{strings.NewReader("λx. "), ""},
			[]any{
				nil,
				fmt.Errorf(":1:5: Unexpected token: EOF"),
			},
		},
		{
			"(λx. x x) (λx. x x)",
			parse,
			[]any{strings.NewReader("(λx. x x) (λx. x x)"), ""},
			[]any{
				&AppExpr{
					expr{},
					&AbsExpr{
						expr{},
						&typ{},
						"x",
						&AppExpr{
							expr{},
							&VarExpr{expr{}, "x"},
							&VarExpr{expr{}, "x"},
						},
					},
					&AbsExpr{
						expr{},
						&typ{},
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
			[]any{strings.NewReader("(λx. (λy. x))"), ""},
			[]any{T, nil},
		},
		{
			"(λx. λy. x)",
			parse,
			[]any{strings.NewReader("(λx. λy. x)"), ""},
			[]any{T, nil},
		},
		{
			"λx.λy.x",
			parse,
			[]any{strings.NewReader("λx.λy.x"), ""},
			[]any{T, nil},
		},
		{
			"x. (one (two (three (four five))))",
			parse,
			[]any{strings.NewReader("x. (one (two (three (four five))))"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&typ{},
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
			[]any{strings.NewReader("x. one two three four five"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&typ{},
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

func TestParserPrimitiveType(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"boolean",
			parse,
			[]any{strings.NewReader("λx : bool . x && y"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&BoolType{},
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
			[]any{strings.NewReader("λx : int . x + y"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&IntType{},
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
			[]any{strings.NewReader("λx : float . x +. y"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&FloatType{},
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
			[]any{strings.NewReader("λx : unit. *"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&UnitType{},
					"x",
					&UnitExpr{expr{&UnitType{}}},
				},
				nil,
			},
		},
	})
}

func TestParserArrowType(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"bool → bool",
			parse,
			[]any{strings.NewReader("λx : bool → bool . x y"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&ArrowType{typ{}, &BoolType{}, &BoolType{}},
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
			[]any{strings.NewReader("λx : bool → bool → bool . x y z"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&ArrowType{
						typ{}, &BoolType{}, &ArrowType{
							typ{}, &BoolType{}, &BoolType{},
						},
					},
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
			[]any{strings.NewReader("λx : bool → bool → bool → int . (x y z) + 3"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&ArrowType{
						typ{}, &BoolType{}, &ArrowType{
							typ{}, &BoolType{}, &ArrowType{
								typ{}, &BoolType{}, &IntType{},
							},
						},
					},
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
			[]any{strings.NewReader("λx : (bool → bool) → bool . x y z"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&ArrowType{
						typ{}, &ArrowType{
							typ{}, &BoolType{}, &BoolType{},
						},
						&BoolType{},
					},
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

// same as TestArrowType(), but systematically using the short
// form (x: $type . $expr instead of λx: $type . $expr
func TestParserArrowTypeShort(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"bool → bool",
			parse,
			[]any{strings.NewReader("x : bool → bool . x y"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&ArrowType{typ{}, &BoolType{}, &BoolType{}},
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
			[]any{strings.NewReader("x : bool → bool → bool . x y z"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&ArrowType{
						typ{}, &BoolType{}, &ArrowType{
							typ{}, &BoolType{}, &BoolType{},
						},
					},
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
			[]any{strings.NewReader("x : bool → bool → bool → int . (x y z) + 3"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&ArrowType{
						typ{}, &BoolType{}, &ArrowType{
							typ{}, &BoolType{}, &ArrowType{
								typ{}, &BoolType{}, &IntType{},
							},
						},
					},
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
			[]any{strings.NewReader("x : (bool → bool) → bool . x y z"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&ArrowType{
						typ{}, &ArrowType{
							typ{}, &BoolType{}, &BoolType{},
						},
						&BoolType{},
					},
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

// "we adopt the convention that × binds stronger than →"
func TestParserArrowProductType(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"bool × int → bool := (bool×int) → bool",
			parse,
			[]any{strings.NewReader("λx : bool×int → bool . x y"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&ArrowType{typ{}, &ProductType{
						typ{}, &BoolType{}, &IntType{},
					}, &BoolType{}},
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
			"bool × (int → bool)",
			parse,
			[]any{strings.NewReader("λx : bool×(int → bool) . x y"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&ProductType{typ{}, &BoolType{}, &ArrowType{
						typ{}, &IntType{}, &BoolType{},
					}},
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
	})
}

// again, given what's in qlambdabook.pdf, we assume ×
// to be right-associative.
func TestParserProductType(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"bool × int × bool := bool×(int×bool)",
			parse,
			[]any{strings.NewReader("λx : bool×int×bool . x y"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&ProductType{typ{}, &BoolType{}, &ProductType{
						typ{}, &IntType{}, &BoolType{},
					}},
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
			"(bool × int) × bool",
			parse,
			[]any{strings.NewReader("λx : (bool×int)×bool . x y"), ""},
			[]any{
				&AbsExpr{
					expr{},
					&ProductType{typ{}, &ProductType{
						typ{}, &BoolType{}, &IntType{},
					}, &BoolType{}},
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
	})
}

func TestParserProduct(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"<>",
			parse,
			[]any{strings.NewReader("〈〉"), ""},
			[]any{
				nil,
				fmt.Errorf(":1:2: Unexpected token: 〉"),
			},
		},
		{
			"<X>",
			parse,
			[]any{strings.NewReader("〈X〉"), ""},
			[]any{
				&VarExpr{expr{}, "X"},
				nil,
			},
		},
		{
			"<X, Y>",
			parse,
			[]any{strings.NewReader("〈X, Y〉"), ""},
			[]any{
				&ProductExpr{expr{},
					&VarExpr{expr{}, "X"},
					&VarExpr{expr{}, "Y"},
				},
				nil,
			},
		},
		{
			"<X, Y, Z>",
			parse,
			[]any{strings.NewReader("〈X, Y, Z〉"), ""},
			[]any{
				&ProductExpr{expr{},
					&VarExpr{expr{}, "X"},
					&ProductExpr{expr{},
						&VarExpr{expr{}, "Y"},
						&VarExpr{expr{}, "Z"},
					},
				},
				nil,
			},
		},
		{
			"<X, <Y, Z>>",
			parse,
			[]any{strings.NewReader("〈X, 〈Y, Z〉〉"), ""},
			[]any{
				&ProductExpr{expr{},
					&VarExpr{expr{}, "X"},
					&ProductExpr{expr{},
						&VarExpr{expr{}, "Y"},
						&VarExpr{expr{}, "Z"},
					},
				},
				nil,
			},
		},
	})
}

func TestParserBasicLetIn(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"let",
			parse,
			[]any{strings.NewReader("let"), ""},
			[]any{
				nil,
				fmt.Errorf(":1:4: Expecting variable name after let, got: EOF"),
			},
		},
		{
			"let 42",
			parse,
			[]any{strings.NewReader("let 42"), ""},
			[]any{
				nil,
				fmt.Errorf(":1:5: Expecting variable name after let, got: int64"),
			},
		},
		{
			"let x 42",
			parse,
			[]any{strings.NewReader("let x 42"), ""},
			[]any{
				nil,
				fmt.Errorf(":1:7: Expecting equal after let $x, got: int64"),
			},
		},
		{
			"let x = 42",
			parse,
			[]any{strings.NewReader("let x = 42"), ""},
			[]any{
				nil,
				fmt.Errorf(":1:11: Expecting 'in' after let $x = $M, got EOF"),
			},
		},
		{
			"let x = 42 in",
			parse,
			[]any{strings.NewReader("let x = 42 in"), ""},
			[]any{
				nil,
				fmt.Errorf(":1:14: Unexpected token: EOF"),
			},
		},
		{
			"let x = 42 in x",
			parse,
			[]any{strings.NewReader("let x = 42 in x"), ""},
			[]any{
				&AppExpr{expr{},
					&AbsExpr{expr{},
						&typ{},
						"x",
						&VarExpr{expr{}, "x"},
					},
					&IntExpr{expr{&IntType{typ{}}}, 42},
				},
				nil,
			},
		},
		{
			"let x = 42 in x + 3",
			parse,
			[]any{strings.NewReader("let x = 42 in x + 3"), ""},
			[]any{
				&AppExpr{expr{},
					&AbsExpr{expr{},
						&typ{},
						"x",
						&BinaryExpr{expr{},
							tokenPlus,
							&VarExpr{expr{}, "x"},
							&IntExpr{expr{&IntType{typ{}}}, 3},
						},
					},
					&IntExpr{expr{&IntType{typ{}}}, 42},
				},
				nil,
			},
		},
	})
}
