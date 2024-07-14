package main

import (
	"testing"
//	"fmt"
//	"strings"
//	"encoding/json"
)

func TestEvalArithmetic(t *testing.T) {
	doTests(t, []test{
/*
		{
			"empty input",
			evalExpr,
			[]interface{}{strings.NewReader(""), ""},
			[]interface{}{nil, fmt.Errorf(":1:1: Unexpected token: EOF")},
		},
*/
		{
			"Basic addition",
			evalExpr,
			[]interface{}{&BinaryExpr{
				expr{},
				tokenPlus,
				&IntExpr{expr{}, 3},
				&IntExpr{expr{}, 4},
			}},
			[]interface{}{
				&IntExpr{expr{}, 7},
				nil,
			},
		},
	})
}