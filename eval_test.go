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
			"Basic addition",
			evalExpr,
			[]any{&BinaryExpr{
				expr{},
				tokenPlus,
				&IntExpr{expr{}, 3},
				&IntExpr{expr{}, 4},
			}},
			[]any{
				&IntExpr{expr{}, 7},
				nil,
			},
		},
	})
}
