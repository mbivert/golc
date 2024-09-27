package main

import (
	"testing"

	"github.com/mbivert/ftests"
)

func TestEvalArithmetic(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		/*
			{
				"empty input",
				evalExpr,
				[]any{strings.NewReader(""), ""},
				[]any{nil, fmt.Errorf(":1:1: Unexpected token: EOF")},
			},
		*/
		{
			"3+4",
			evalExpr,
			[]any{mustSTypeParse("3+4")},
			[]any{
				&IntExpr{expr{&IntType{typ{}}}, 7},
				nil,
			},
		},
		{
			"3+4*2",
			evalExpr,
			[]any{mustSTypeParse("3+4*2")},
			[]any{
				&IntExpr{expr{&IntType{typ{}}}, 11},
				nil,
			},
		},
		{
			"(3+4)*2",
			evalExpr,
			[]any{mustSTypeParse("(3+4)*2")},
			[]any{
				&IntExpr{expr{&IntType{typ{}}}, 14},
				nil,
			},
		},
		{
			"(2<3) && !(true)",
			evalExpr,
			[]any{mustSTypeParse("(2<3) && !(true)")},
			[]any{
				&BoolExpr{expr{&BoolType{typ{}}}, false},
				nil,
			},
		},
	})
}

func TestEvalRenameExpr(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"(2<3) && !(true) && 3. ≤. 5. (no changes expected)",
			renameExpr,
			[]any{mustSTypeParse("(2<3) && !(true) && (3. ≤. 5.)"), "x", "y"},
			[]any{
				mustSTypeParse("(2<3) && !(true) && (3. ≤. 5.)"),
			},
		},
		// NOTE: type checking doesn't like for x to be unbounded hence mustParse()
		// instead of mustSTypeParse().
		// TODO: there are plans to allow it, as x's type can be infered.
		{
			"(2<3) && !(true) && 3. ≤. x",
			renameExpr,
			[]any{mustParse("(2<3) && !(true) && 3. ≤. x"), "y", "x"},
			[]any{
				mustParse("(2<3) && !(true) && 3. ≤. y"),
			},
		},
		{
			"λx:int. x+3",
			renameExpr,
			[]any{mustSTypeParse("λx:int. x+3"), "y", "x"},
			[]any{
				mustSTypeParse("λy:int. y+3"),
			},
		},
		{
			"λf:int→int.x:int. f (x+3)",
			renameExpr,
			[]any{mustSTypeParse("λf:int→int.x:int. f (x+3)"), "g", "f"},
			[]any{
				mustSTypeParse("λg:int→int.x:int. g (x+3)"),
			},
		},
	})
}
