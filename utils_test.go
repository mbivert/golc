package main

import (
	"testing"

	"github.com/mbivert/ftests"
)

func TestUtilsFreeVars(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"scalar expression: int",
			freeVars,
			[]any{mustParse("123")},
			[]any{map[string]bool{}},
		},
		{
			"scalar expression: bool",
			freeVars,
			[]any{mustParse("true")},
			[]any{map[string]bool{}},
		},
		{
			"scalar expression: float",
			freeVars,
			[]any{mustParse("42.42")},
			[]any{map[string]bool{}},
		},
		{
			"free variable",
			freeVars,
			[]any{mustParse("x")},
			[]any{map[string]bool{"x": true}},
		},
		{
			"simple abstraction, no free variable",
			freeVars,
			[]any{mustParse("x. x")},
			[]any{map[string]bool{}},
		},
		{
			"simple abstraction, one free variable",
			freeVars,
			[]any{mustParse("x. y")},
			[]any{map[string]bool{
				"y": true,
			}},
		},
		{
			"abstraction + applications",
			freeVars,
			[]any{mustParse("x. x y z")},
			[]any{map[string]bool{
				"y": true,
				"z": true,
			}},
		},
	})
}

func TestUtilsAllVars(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"scalar expression: int",
			allVars,
			[]any{mustParse("123")},
			[]any{map[string]bool{}},
		},
		{
			"scalar expression: bool",
			allVars,
			[]any{mustParse("true")},
			[]any{map[string]bool{}},
		},
		{
			"scalar expression: float",
			allVars,
			[]any{mustParse("42.42")},
			[]any{map[string]bool{}},
		},
		{
			"free variable",
			allVars,
			[]any{mustParse("x")},
			[]any{map[string]bool{"x": true}},
		},
		{
			"simple abstraction, no free variable, one bound",
			allVars,
			[]any{mustParse("x. x")},
			[]any{map[string]bool{"x": true}},
		},
		{
			"simple abstraction, one free variable, one bound",
			allVars,
			[]any{mustParse("x. y")},
			[]any{map[string]bool{
				"x": true,
				"y": true,
			}},
		},
		{
			"abstraction + applications",
			allVars,
			[]any{mustParse("x. x y z")},
			[]any{map[string]bool{
				"x": true,
				"y": true,
				"z": true,
			}},
		},
	})
}

// TODO: a bit light
func TestUtilsPrettyPrint(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"bare int",
			prettyPrint,
			[]any{mustParse("123")},
			[]any{"123"},
		},
		{
			"bare float",
			prettyPrint,
			[]any{mustParse("123.42")},
			[]any{"123.42"},
		},
		{
			"bare bool",
			prettyPrint,
			[]any{mustParse("true")},
			[]any{"true"},
		},
		{
			"bare variable",
			prettyPrint,
			[]any{mustParse("someVar")},
			[]any{"someVar"},
		},
		{
			"simple abstraction (id)",
			prettyPrint,
			[]any{mustParse("λ x. x")},
			[]any{"(λx. x)"},
		},
		{
			"arithmetic",
			prettyPrint,
			[]any{mustParse("(2+2)*3")},
			[]any{"((2 + 2) * 3)"},
		},
		{
			"imbricated abstraction + application",
			prettyPrint,
			[]any{mustParse("λx. y. x y")},
			[]any{"(λx. y. (x y))"},
		},
		{
			"and",
			prettyPrint,
			[]any{and},
			[]any{"(λx. y. (x y (λx. y. y)))"},
		},
		{
			"(((x y) q) (z p))",
			prettyPrint,
			[]any{mustParse("(((x y) q) (z p))")},
			[]any{"(x y q (z p))"},
		},
	})
}

func TestUtilsGetFresh(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"empty map",
			getFresh,
			[]any{map[string]bool{}},
			[]any{"x0"},
		},
		{
			"not empty, but x0 still free",
			getFresh,
			[]any{map[string]bool{"x": true, "y": true}},
			[]any{"x0"},
		},
		{
			"x0 already used",
			getFresh,
			[]any{map[string]bool{"x0": true}},
			[]any{"x1"},
		},
	})
}
