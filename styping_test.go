package main

import (
	"fmt"
	"testing"

	"github.com/mbivert/ftests"
)

func TestSTypingInferSTypeSingle(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"x",
			inferSType,
			[]any{mustParse("x")},
			[]any{
				nil,
				fmt.Errorf("'x' isn't bounded!"),
			},
		},
		{
			"42",
			inferSType,
			[]any{mustParse("42")},
			[]any{
				&IntExpr{expr{&IntType{typ{}}}, 42},
				nil,
			},
		},
		{
			"true",
			inferSType,
			[]any{mustParse("true")},
			[]any{
				&BoolExpr{expr{&BoolType{typ{}}}, true},
				nil,
			},
		},
		{
			"42.42",
			inferSType,
			[]any{mustParse("42.42")},
			[]any{
				&FloatExpr{expr{&FloatType{typ{}}}, 42.42},
				nil,
			},
		},
	})
}

func TestSTypingInferSTypeTypedLambdas(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"λx:bool.*",
			inferSType,
			[]any{mustParse("λx:bool.*")},
			[]any{
				&AbsExpr{expr{&ArrowType{typ{},
						&BoolType{typ{}},
						&UnitType{typ{}},
					}},
					&BoolType{typ{}},
					"x",
					&UnitExpr{expr{&UnitType{typ{}}}},
				},
				nil,
			},
		},
		{
			"λx:bool.x",
			inferSType,
			[]any{mustParse("λx:bool.x")},
			[]any{
				&AbsExpr{expr{&ArrowType{typ{},
						&BoolType{typ{}},
						&BoolType{typ{}},
					}},
					&BoolType{typ{}},
					"x",
					&VarExpr{expr{&BoolType{typ{}}}, "x"},
				},
				nil,
			},
		},
		{
			"λx:bool.y",
			inferSType,
			[]any{mustParse("λx:bool.y")},
			[]any{
				nil,
				fmt.Errorf("'y' isn't bounded!"),
			},
		},
	})
}

/*

// TODO: in the parsing, we'll want to allow types to
// be unspecified ("MissingType"). variable MissingType
// can be recovered e.g. by binary applications. Think:
//
//	λx. x+3 can be typed as int → int; the type of x
//	can be infered when typing x+3
//
func TestSTypingInferSTypeUntypedLambdas(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"λx.*",
			inferSType,
			[]any{mustParse("λx.*")},
			[]any{
				&AbsExpr{expr{&ArrowType{typ{},
					&MissingType{typ{}},
					&UnitType{typ{}},
				}},
					&BoolType{typ{}},
					"x",
					&UnitExpr{expr{&UnitType{typ{}}}},
				},
				nil,
			},
		},
	})
}
*/

func TestSTypingInferSTypeApps(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"42 42",
			inferSType,
			[]any{mustParse("42 42")},
			[]any{
				nil,
				fmt.Errorf("Trying to apply to non-arrow: 'int'"),
			},
		},
		{
			"(λx:bool.x) true",
			inferSType,
			[]any{mustParse("(λx:bool.x) true")},
			[]any{
				&AppExpr{expr{&BoolType{typ{}}},
					&AbsExpr{expr{&ArrowType{typ{},
						&BoolType{typ{}},
						&BoolType{typ{}},
					}},
						&BoolType{typ{}},
						"x",
						&VarExpr{expr{&BoolType{typ{}}}, "x"},
					},
					&BoolExpr{expr{&BoolType{typ{}}}, true},
				},
				nil,
			},
		},
		{
			"(λx:bool.x) 42",
			inferSType,
			[]any{mustParse("(λx:bool.x) 42")},
			[]any{
				nil,
				fmt.Errorf("Can't apply 'int' to 'bool → bool'"),
			},
		},
	})
}

// TODO: far from extensive
func TestSTypingInferSTypeExprs(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"3+true",
			inferSType,
			[]any{mustParse("3+true")},
			[]any{
				nil,
				fmt.Errorf("+ : (int×int) → int; got (int×bool)"),
			},
		},
		{
			"3+3",
			inferSType,
			[]any{mustParse("3+3")},
			[]any{
				&BinaryExpr{expr{&IntType{typ{}}},
					tokenPlus,
					&IntExpr{expr{&IntType{typ{}}}, 3},
					&IntExpr{expr{&IntType{typ{}}}, 3},
				},
				nil,
			},
		},
		{
			"3-.true",
			inferSType,
			[]any{mustParse("3-.true")},
			[]any{
				nil,
				fmt.Errorf("-. : (float×float) → float; got (int×bool)"),
			},
		},
		{
			"3.-.5.",
			inferSType,
			[]any{mustParse("3.-.5.")},
			[]any{
				&BinaryExpr{expr{&FloatType{typ{}}},
					tokenFMinus,
					&FloatExpr{expr{&FloatType{typ{}}}, 3.},
					&FloatExpr{expr{&FloatType{typ{}}}, 5.},
				},
				nil,
			},
		},
		{
			"3.&&5.",
			inferSType,
			[]any{mustParse("3.&&5.")},
			[]any{
				nil,
				fmt.Errorf("&& : (bool×bool) → bool; got (float×float)"),
			},
		},
		{
			"true&& false",
			inferSType,
			[]any{mustParse("true&& false")},
			[]any{
				&BinaryExpr{expr{&BoolType{typ{}}},
					tokenAndAnd,
					&BoolExpr{expr{&BoolType{typ{}}}, true},
					&BoolExpr{expr{&BoolType{typ{}}}, false},
				},
				nil,
			},
		},
		{
			"3<5",
			inferSType,
			[]any{mustParse("3<5")},
			[]any{
				&BinaryExpr{expr{&BoolType{typ{}}},
					tokenLess,
					&IntExpr{expr{&IntType{typ{}}}, 3},
					&IntExpr{expr{&IntType{typ{}}}, 5},
				},
				nil,
			},
		},
	})
}

func TestSTypingInferInferSTypeProduct(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"〈3, 3〉",
			inferSType,
			[]any{mustParse("〈3, 3〉")},
			[]any{
				&ProductExpr{expr{&ProductType{typ{},
						&IntType{typ{}},
						&IntType{typ{}},
					}},
					&IntExpr{expr{&IntType{typ{}}}, 3},
					&IntExpr{expr{&IntType{typ{}}}, 3},
				},
				nil,
			},
		},
		{
			"〈3, true〉",
			inferSType,
			[]any{mustParse("〈3, true〉")},
			[]any{
				&ProductExpr{expr{&ProductType{typ{},
						&IntType{typ{}},
						&BoolType{typ{}},
					}},
					&IntExpr{expr{&IntType{typ{}}}, 3},
					&BoolExpr{expr{&BoolType{typ{}}}, true},
				},
				nil,
			},
		},
		{
			"〈3, true, 5.〉",
			inferSType,
			[]any{mustParse("〈3, true, 5.〉")},
			[]any{
				&ProductExpr{expr{&ProductType{typ{},
						&IntType{typ{}},
						&ProductType{typ{},
							&BoolType{typ{}},
							&FloatType{typ{}},
						},
					}},
					&IntExpr{expr{&IntType{typ{}}}, 3},
					&ProductExpr{expr{&ProductType{typ{},
							&BoolType{typ{}},
							&FloatType{typ{}},
						}},
						&BoolExpr{expr{&BoolType{typ{}}}, true},
						&FloatExpr{expr{&FloatType{typ{}}}, 5.},
					},
				},
				nil,
			},
		},
	})
}

// TODO: again, not extensive
func TestSTypingInferInferSTypeUnary(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"+true",
			inferSType,
			[]any{mustParse("+true")},
			[]any{
				nil,
				fmt.Errorf("+ : int → int; got bool"),
			},
		},
		{
			"+.true",
			inferSType,
			[]any{mustParse("+.true")},
			[]any{
				nil,
				fmt.Errorf("+. : float → float; got bool"),
			},
		},
		{
			"+3",
			inferSType,
			[]any{mustParse("+3")},
			[]any{
				&UnaryExpr{expr{&IntType{typ{}}},
					tokenPlus,
					&IntExpr{expr{&IntType{typ{}}}, 3},
				},
				nil,
			},
		},
		{
			"-.3.",
			inferSType,
			[]any{mustParse("-.3.")},
			[]any{
				&UnaryExpr{expr{&FloatType{typ{}}},
					tokenFMinus,
					&FloatExpr{expr{&FloatType{typ{}}}, 3.},
				},
				nil,
			},
		},
		{
			"!3",
			inferSType,
			[]any{mustParse("!3")},
			[]any{
				nil,
				fmt.Errorf("! : bool → bool; got int"),
			},
		},
		{
			"!true",
			inferSType,
			[]any{mustParse("!true")},
			[]any{
				&UnaryExpr{expr{&BoolType{typ{}}},
					tokenExcl,
					&BoolExpr{expr{&BoolType{typ{}}}, true},
				},
				nil,
			},
		},
	})
}

