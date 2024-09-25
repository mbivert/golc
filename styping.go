/*
 * Inference for a simply typed extended λ-calculus.
 *
 * By comparison with the incomplete typing.go / typing_test.go,
 * which targets a polymorphic (à la System F) λ-calculus.
 *
 * Essentially, the idea is that we don't have to deal with VarType.
 */
package main

import (
	"fmt"
	"reflect"
)

// map a bound variable name to its type
type Ctx map[string]Type

// Infer the types for a given expression and perform
// typechecking when relevant.
//
// We modify (and return) the expression in place
// so that it contains the relevant typing data.
func sInferType(x Expr) (Expr, error) {
	var aux func(Expr, Ctx) (Expr, error)

	aux = func(x Expr, ctx Ctx) (Expr, error) {
		switch x.(type) {

		// Those four cases have already been typed
		// during the parsing.
		case *IntExpr:
		case *FloatExpr:
		case *BoolExpr:
		case *UnitExpr:

		// We may need some typechecking here
		case *UnaryExpr:
			switch x.(*UnaryExpr).op {
			// Right must be int
			case tokenMinus:
			case tokenPlus:

			// Right must be float
			case tokenFMinus:
			case tokenFPlus:

			// Right must be bool
			// case tokenExclamation:
			}

		// Again, we may need to typecheck things here
		case *BinaryExpr:
			switch x.(*BinaryExpr).op {
			// Left right must be ints
			case tokenMinus:
			case tokenPlus:
			case tokenStar:
			case tokenSlash:
			case tokenLessEq:
			case tokenMoreEq:
			case tokenLess:
			case tokenMore:

			// Left right must be floats
			case tokenFMinus:
			case tokenFPlus:
			case tokenFStar:
			case tokenFSlash:
			case tokenFLessEq:
			case tokenFMoreEq:
			case tokenFLess:
			case tokenFMore:

			// Left/right must be bools
			case tokenOrOr:
			case tokenAndAnd:
			}

		case *AbsExpr:
			n := x.(*AbsExpr).name
			t := x.(*AbsExpr).typ
			r := x.(*AbsExpr).right

			// save previous ctx[n] if any
			t2, ok := ctx[n]

			// new var in env
			ctx[n] = t

			r, err := aux(r, ctx)
			if err != nil {
				return nil, err
			}
			x.setType(&ArrowType{typ{}, t, r.getType()})
			x.(*AbsExpr).right = r

			if ok {
				ctx[n] = t2
			} else {
				delete(ctx, n)
			}

		case *AppExpr:
			l := x.(*AppExpr).left
			r := x.(*AppExpr).right

			l, err := aux(l, ctx)
			if err != nil {
				return nil, err
			}

			r, err = aux(r, ctx)
			if err != nil {
				return nil, err
			}

			tl, ok := l.getType().(*ArrowType)
			if !ok {
				return nil, fmt.Errorf("Trying to apply to non-arrow: '%s'", l.getType())
			}
			if reflect.TypeOf(tl.left) != reflect.TypeOf(r.getType()) {
				return nil, fmt.Errorf("Can't apply '%s' to '%s'", r.getType(), l.getType())
			}

			x.setType(l.(*AbsExpr).right.getType())
			x.(*AppExpr).left = l
			x.(*AppExpr).right = r

		case *VarExpr:
			n := x.(*VarExpr).name
			t, ok := ctx[n]
			if !ok {
				return nil, fmt.Errorf("'%s' isn't bounded!", n)
			}
			x.setType(t)
			return x, nil

		// case *ProductExpr:
		}

		return x, nil
	}

	return aux(x, Ctx{})
}
