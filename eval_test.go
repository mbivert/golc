package main

import (
	"testing"

	"github.com/mbivert/ftests"
)

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
			"matching variable is substituted",
			substituteExpr,
			[]any{mustParse("x"), mustParse("λx. λy. x y"), "x"},
			[]any{
				mustParse("λx. λy. x y"),
			},
		},
		{
			"un-matching variable name",
			substituteExpr,
			[]any{mustParse("y"), mustParse("λx. λy. x y"), "x"},
			[]any{
				mustParse("y"),
			},
		},
		{
			"variable substituted in both parts of an apply",
			substituteExpr,
			[]any{mustParse("(x (x y))"), mustParse("λx. λy. x y"), "x"},
			[]any{
				mustParse("((λx. λy. x y) ((λx. λy. x y) y))"),
			},
		},
		{
			"bound variable not substituted",
			substituteExpr,
			[]any{mustParse("λx. λz. x z"), mustParse("λx. λy. x y"), "x"},
			[]any{
				mustParse("λx. λz. x z"),
			},
		},
		{
			"deeper bound variable not substituted",
			substituteExpr,
			[]any{mustParse("λx. λz. x z"), mustParse("λx. λy. x y"), "z"},
			[]any{
				mustParse("λx. λz. x z"),
			},
		},
		{
			"replacing a free variable, no conflict",
			substituteExpr,
			[]any{mustParse("λx. λy. x z"), mustParse("λx. λy. x y"), "z"},
			[]any{
				mustParse("λx. λy. x (λx. λy. x y)"),
			},
		},
		{
			"replacing a free variable, renaming",
			substituteExpr,
			[]any{mustParse("λx. λy. x z y"), mustParse("λx. λz. x y z"), "z"},
			[]any{
				mustParse("λx. λx0. x (λx. λz. x y z) x0"),
			},
		},
		{
			"replacing a free variable, renaming (bis)",
			substituteExpr,
			[]any{mustParse("λx. λy. x z y"), mustParse("λx. λz. x y z"), "z"},
			[]any{
				mustParse("λx. λx0. x (λx. λz. x y z) x0"),
			},
		},
		{
			"Selinger's example",
			substituteExpr,
			[]any{mustParse("λx. y x"), mustParse("λz. x z"), "y"},
			[]any{
				mustParse("λx0. (λz. x z) x0"),
			},
		},
		{
			"replacing bound variable by the variable to rename",
			substituteExpr,
			[]any{
				mustParse(`
					(λf. n.
						((λy.
							(
								(λn. x. y. (n (λz. y) x))
								n
								(λf. x. (f x)) y))
						(λx0.
							(n (f (λf. x. (n (λg. h. (h (g f))) (λu. x) (λu. u))) x0)))))
				`),
				mustParse("f"),
				"x0",
			},
			[]any{
				mustParse(`
					(λx1. n.
						((λy.
							(
								(λn. x. y. (n (λz. y) x))
								n
								(λx1. x. (x1 x)) y))
						(λx0.
							(n (x1 (λx1. x. (n (λg. h. (h (g x1))) (λu. x) (λu. u))) x0)))))
				`),
			},
		},
		{
			"don't re-use a name already used below",
			substituteExpr,
			[]any{
				mustParse(`(λn. x0. y. (n (λz. y) x0))`),
				mustParse(`
					(λx0.
						(n (x1 (λx1. x0.
							(n (λg. h. (h (g x1))) (λu. x0) (λu. u))) x0)))
				`),
				"y",
			},
			[]any{
				mustParse(`(λx2. x0. y. (x2 (λz. y) x0))`),
			},
		},
		{
			"\"complex\" substitute",
			substituteExpr,
			[]any{
				mustParse(`
					(λy.
						(λp. λx. λy. p x y)
						x
						(λx. λy. x)
						(
							(λp. λx. λy. p x y)
							y
							(λx. λy. x)
							(λx. λy. y)))
				`),
				mustParse(`(λx. λy. x)`),
				"x",
			},
			[]any{
				mustParse(`
					(λy.
						(λp. λx. λy. p x y)
						(λx. λy. x)
						(λx. λy. x)
						(
							(λp. λx. λy. p x y)
							y
							(λx. λy. x)
							(λx. λy. y)))
				`),
			},
		},
		{
			"\"complex\" substitute (bis)",
			substituteExpr,
			[]any{
				mustParse(`
					(λp. λx. λy. p x y)
					(λx. λy. x)
					(λx. λy. x)
					(
						(λp. λx. λy. p x y)
						y
						(λx. λy. x)
						(λx. λy. y))
				`),
				mustParse(`(λx. λy. x)`),
				"y",
			},
			[]any{
				mustParse(`
					(λp. λx. λy. p x y)
					(λx. λy. x)
					(λx. λy. x)
					(
						(λp. λx. λy. p x y)
						(λx. λy. x)
						(λx. λy. x)
						(λx. λy. y))
				`),
			},
		},
		{
			"\"complex\" substitute (ter)",
			substituteExpr,
			[]any{
				mustParse(`((λx. λy. x) (λx. λy. x) y)`),
				mustParse(`((λx. λy. x) (λx. λy. x) (λx. λy. y))`),
				"y",
			},
			[]any{
				mustParse(`
				(λx. λy. x) (λx. λy. x)
					((λx. λy. x) (λx. λy. x) (λx. λy. y))
				`),
			},
		},
		{
			"\"complex\" substitute (xor, 1)",
			substituteExpr,
			[]any{
				mustParse(`
					((((λp. λx. λy. (p x y)) x)
						((((λp. λx. λy. (p x y)) y) (λx. λy. y)) (λx. λy. x)))
						((((λp. λx. λy. (p x y)) y) (λx. λy. x)) (λx. λy. y)))
				`),
				mustParse(`(λx. (λy. x))`),
				"y",
			},
			[]any{
				mustParse(`
					((((λp. λx. λy. (p x y)) x)
						((((λp. λx. λy. (p x y)) (λx. (λy. x))) (λx. λy. y)) (λx. λy. x)))
						((((λp. λx. λy. (p x y)) (λx. (λy. x))) (λx. λy. x)) (λx. λy. y)))
				`),
			},
		},
		{
			"\"complex\" substitute (xor, 2)",
			substituteExpr,
			[]any{
				mustParse(`
					(λy.
						((((λp. λx. λy. (p x y)) x)
							((((λp. λx. λy. (p x y)) y) (λx. λy. y)) (λx. λy. x)))
							((((λp. λx. λy. (p x y)) y) (λx. λy. x)) (λx. λy. y))))
						(λx. (λy. x))
				`),
				mustParse(`(λx. (λy. x))`),
				"x",
			},
			[]any{
				mustParse(`
					(λy.
						((((λp. λx. λy. (p x y)) (λx. (λy. x)))
							((((λp. λx. λy. (p x y)) y) (λx. λy. y)) (λx. λy. x)))
							((((λp. λx. λy. (p x y)) y) (λx. λy. x)) (λx. λy. y))))
						(λx. (λy. x))
				`),
			},
		},
	})
}

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
			},
		},
		{
			"3+4*2",
			evalExpr,
			[]any{mustSTypeParse("3+4*2")},
			[]any{
				&IntExpr{expr{&IntType{typ{}}}, 11},
			},
		},
		{
			"(3+4)*2",
			evalExpr,
			[]any{mustSTypeParse("(3+4)*2")},
			[]any{
				&IntExpr{expr{&IntType{typ{}}}, 14},
			},
		},
		{
			"(2<3) && !(true)",
			evalExpr,
			[]any{mustSTypeParse("(2<3) && !(true)")},
			[]any{
				&BoolExpr{expr{&BoolType{typ{}}}, false},
			},
		},
	})
}

func TestEvalLambdMaths(t *testing.T) {
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
			"(λx:int. x+3) 5",
			evalExpr,
			[]any{mustSTypeParse("(λx:int. x+3) 5")},
			[]any{
				&IntExpr{expr{&IntType{typ{}}}, 8},
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
