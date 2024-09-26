package main

import (
	"fmt"
	"testing"

	"github.com/mbivert/ftests"
)

func TestSTypingSInferType(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"x",
			sInferType,
			[]any{mustParse("x")},
			[]any{
				nil,
				fmt.Errorf("'x' isn't bounded!"),
			},
		},
		{
			"42",
			sInferType,
			[]any{mustParse("42")},
			[]any{
				&IntExpr{expr{&IntType{typ{}}}, 42},
				nil,
			},
		},
		{
			"true",
			sInferType,
			[]any{mustParse("true")},
			[]any{
				&BoolExpr{expr{&BoolType{typ{}}}, true},
				nil,
			},
		},
		{
			"42.42",
			sInferType,
			[]any{mustParse("42.42")},
			[]any{
				&FloatExpr{expr{&FloatType{typ{}}}, 42.42},
				nil,
			},
		},
		{
			"λx:bool.*",
			sInferType,
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
			sInferType,
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
			sInferType,
			[]any{mustParse("λx:bool.y")},
			[]any{
				nil,
				fmt.Errorf("'y' isn't bounded!"),
			},
		},
		{
			"42 42",
			sInferType,
			[]any{mustParse("42 42")},
			[]any{
				nil,
				fmt.Errorf("Trying to apply to non-arrow: 'int'"),
			},
		},
		{
			"(λx:bool.x) true",
			sInferType,
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
			sInferType,
			[]any{mustParse("(λx:bool.x) 42")},
			[]any{
				nil,
				fmt.Errorf("Can't apply 'int' to 'bool → bool'"),
			},
		},
	})
}
