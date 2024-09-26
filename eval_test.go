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
