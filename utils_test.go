package main

import (
	"testing"
//	"fmt"
//	"encoding/json"
)

// True
var T = &AbsExpr{
	expr{},
	&typ{},
	"x",
	&AbsExpr{
		expr{},
		&typ{},
		"y",
		&VarExpr{expr{}, "x"},
	},
}

// False
var F = &AbsExpr{
	expr{},
	&typ{},
	"x",
	&AbsExpr{
		expr{},
		&typ{},
		"y",
		&VarExpr{expr{}, "y"},
	},
}

var and = &AbsExpr{
	expr{},
	&typ{},
	"x",
	&AbsExpr{
		expr{},
		&typ{},
		"y",
		&AppExpr{
			expr{},
			&AppExpr{
				expr{},
				&VarExpr{expr{}, "x"},
				&VarExpr{expr{}, "y"},
			},
			F,
		},
	},
}

func TestFreeVars(t *testing.T) {
	doTests(t, []test{
		{
			"scalar expression: int",
			freeVars,
			[]interface{}{mustParse("123")},
			[]interface{}{map[string]bool{}},
		},
		{
			"scalar expression: bool",
			freeVars,
			[]interface{}{mustParse("true")},
			[]interface{}{map[string]bool{}},
		},
		{
			"scalar expression: float",
			freeVars,
			[]interface{}{mustParse("42.42")},
			[]interface{}{map[string]bool{}},
		},
		{
			"free variable",
			freeVars,
			[]interface{}{mustParse("x")},
			[]interface{}{map[string]bool{"x":true}},
		},
		{
			"simple abstraction, no free variable",
			freeVars,
			[]interface{}{mustParse("x. x")},
			[]interface{}{map[string]bool{}},
		},
		{
			"simple abstraction, one free variable",
			freeVars,
			[]interface{}{mustParse("x. y")},
			[]interface{}{map[string]bool{
				"y" : true,
			}},
		},
		{
			"abstraction + applications",
			freeVars,
			[]interface{}{mustParse("x. x y z")},
			[]interface{}{map[string]bool{
				"y" : true,
				"z" : true,
			}},
		},
	})
}

func TestallVars(t *testing.T) {
	doTests(t, []test{
		{
			"scalar expression: int",
			allVars,
			[]interface{}{mustParse("123")},
			[]interface{}{map[string]bool{}},
		},
		{
			"scalar expression: bool",
			allVars,
			[]interface{}{mustParse("true")},
			[]interface{}{map[string]bool{}},
		},
		{
			"scalar expression: float",
			allVars,
			[]interface{}{mustParse("42.42")},
			[]interface{}{map[string]bool{}},
		},
		{
			"free variable",
			allVars,
			[]interface{}{mustParse("x")},
			[]interface{}{map[string]bool{"x":true}},
		},
		{
			"simple abstraction, no free variable, one bound",
			allVars,
			[]interface{}{mustParse("x. x")},
			[]interface{}{map[string]bool{"x" : true}},
		},
		{
			"simple abstraction, one free variable, one bound",
			allVars,
			[]interface{}{mustParse("x. y")},
			[]interface{}{map[string]bool{
				"x" : true,
				"y" : true,
			}},
		},
		{
			"abstraction + applications",
			allVars,
			[]interface{}{mustParse("x. x y z")},
			[]interface{}{map[string]bool{
				"x" : true,
				"y" : true,
				"z" : true,
			}},
		},
	})
}

// TODO: a bit light
func TestPrettyPrint(t *testing.T) {
	doTests(t, []test{
		{
			"bare int",
			prettyPrint,
			[]interface{}{mustParse("123")},
			[]interface{}{"123"},
		},
		{
			"bare float",
			prettyPrint,
			[]interface{}{mustParse("123.42")},
			[]interface{}{"123.42"},
		},
		{
			"bare bool",
			prettyPrint,
			[]interface{}{mustParse("true")},
			[]interface{}{"true"},
		},
		{
			"bare variable",
			prettyPrint,
			[]interface{}{mustParse("someVar")},
			[]interface{}{"someVar"},
		},
		{
			"simple abstraction (id)",
			prettyPrint,
			[]interface{}{mustParse("λ x. x")},
			[]interface{}{"(λx. x)"},
		},
		{
			"arithmetic",
			prettyPrint,
			[]interface{}{mustParse("(2+2)*3")},
			[]interface{}{"((2 + 2) * 3)"},
		},
		{
			"imbricated abstraction + application",
			prettyPrint,
			[]interface{}{mustParse("λx. y. x y")},
			[]interface{}{"(λx. y. (x y))"},
		},
		{
			"and",
			prettyPrint,
			[]interface{}{and},
			[]interface{}{"(λx. y. (x y (λx. y. y)))"},
		},
		{
			"(((x y) q) (z p))",
			prettyPrint,
			[]interface{}{mustParse("(((x y) q) (z p))")},
			[]interface{}{"(x y q (z p))"},
		},
	})
}
