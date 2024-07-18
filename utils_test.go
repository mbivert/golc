package main

import (
	"testing"
//	"fmt"
//	"encoding/json"
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

func TestMergeMaps(t *testing.T) {
	doTests(t, []test{
		{
			"empty maps",
			mergeMaps,
			[]interface{}{map[string]bool{}, map[string]bool{}},
			[]interface{}{map[string]bool{}},
		},
		{
			"Both maps non-empty",
			mergeMaps,
			[]interface{}{map[string]bool{"a":true}, map[string]bool{"b": false}},
			[]interface{}{map[string]bool{
				"a" : true,
				"b" : false,
			}},
		},
		{
			"Overriding",
			mergeMaps,
			[]interface{}{map[string]bool{"a":true}, map[string]bool{"a": false}},
			[]interface{}{map[string]bool{
				"a" : false,
			}},
		},
	})
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
			[]interface{}{"(λx.x)"},
		},
		{
			"arithmetic",
			prettyPrint,
			[]interface{}{mustParse("(2+2)*3")},
			[]interface{}{"((2 + 2) * 3)"},
		},
		{
			"Imbricated abstraction (gets \"simplified\") + application",
			prettyPrint,
			[]interface{}{mustParse("λx.λy. x y")},
			[]interface{}{"(λx.y.(x y))"},
		},
	})
}
