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

// Some of those aren't (yet?) properly typed
func TestEvalRenameExpr(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"z | y, x",
			renameExpr,
			[]any{mustParse("z"), "y", "x"},
			[]any{
				mustParse("z"),
			},
		},
		{
			"x | y, x",
			renameExpr,
			[]any{mustParse("x"), "y", "x"},
			[]any{
				mustParse("y"),
			},
		},
		{
			"(x y) (y x z)  | y, x",
			renameExpr,
			[]any{mustParse("(x y) (y x z)"), "y", "x"},
			[]any{
				mustParse("(y y) (y y z) "),
			},
		},
		{
			"λx. x z  | y, x",
			renameExpr,
			[]any{mustParse("λx. x z"), "y", "x"},
			[]any{
				mustParse("λy. y z"),
			},
		},
		{
			"λx. x z  | y, y",
			renameExpr,
			[]any{mustParse("λx. x z"), "y", "y"},
			[]any{
				mustParse("λx. x z"),
			},
		},
		{
			"λx. λy. y z foo bar  | z, x",
			renameExpr,
			[]any{mustParse("λx. λy. y z foo bar"), "z", "x"},
			[]any{
				mustParse("λz. λy. y z foo bar"),
			},
		},
		{
			"λx. λy. y z foo bar  | foo, y",
			renameExpr,
			[]any{mustParse("λx. λy. y z foo bar"), "foo", "y"},
			[]any{
				mustParse("λx. λfoo. foo z foo bar"),
			},
		},
		{
			"(2<3) && !(true) && 3. ≤. 5. (no changes expected) | y, x",
			renameExpr,
			[]any{mustSTypeParse("(2<3) && !(true) && (3. ≤. 5.)"), "y", "x"},
			[]any{
				mustSTypeParse("(2<3) && !(true) && (3. ≤. 5.)"),
			},
		},
		// NOTE: type checking doesn't like for x to be unbounded hence mustParse()
		// instead of mustSTypeParse().
		// TODO: there are plans to allow it, as x's type can be infered.
		{
			"(2<3) && !(true) && 3. ≤. x | y, x",
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

func TestEvalSubstituteExpr(t *testing.T) {
	ftests.Run(t, []ftests.Test{
		{
			"x | λx. λy. x y | x (var, match replaced)",
			substituteExpr,
			[]any{mustParse("x"), mustParse("λx. λy. x y"), "x"},
			[]any{
				mustParse("λx. λy. x y"),
			},
		},
		{
			"y | λx. λy. x y | x (var, no match)",
			substituteExpr,
			[]any{mustParse("y"), mustParse("λx. λy. x y"), "x"},
			[]any{
				mustParse("y"),
			},
		},
		{
			"y | λx. λy. x y | x (two app occurences)",
			substituteExpr,
			[]any{mustParse("(x (x y))"), mustParse("λx. λy. x y"), "x"},
			[]any{
				mustParse("((λx. λy. x y) ((λx. λy. x y) y))"),
			},
		},
		{
			"λx. λz. x z | λx. λy. x y | z (bound var not substituted)",
			substituteExpr,
			[]any{mustParse("λx. λz. x z"), mustParse("λx. λy. x y"), "z"},
			[]any{
				mustParse("λx. λz. x z"),
			},
		},
		{
			"λx. λz. x z | λx. λy. x y | z (free variable substituted)",
			substituteExpr,
			[]any{mustParse("λx. λy. x z"), mustParse("λx. λy. x y"), "z"},
			[]any{
				mustParse("λx. λy. x (λx. λy. x y)"),
			},
		},
		{
			"λx. λy. x z y | λx. λy. x y | z (free variable substituted, with renaming)",
			substituteExpr,
			[]any{mustParse("λx. λy. x z y"), mustParse("λx. λz. x y z"), "z"},
			[]any{
				mustParse("λx. λx0. x (λx. λz. x y z) x0"),
			},
		},
	})
}

/*
	D/diff i t

	∂/pdiff t x u

	⊙/tmult

	⊕/tadd // depends on whether we want to combine λ-calcs


	1 ∂ λ

 */
