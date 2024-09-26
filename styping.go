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
func inferSType(x Expr) (Expr, error) {
	var aux func(Expr, Ctx) (Expr, error)

	aux = func(x Expr, ctx Ctx) (Expr, error) {
		var err error

		switch x.(type) {

		// Those four cases have already been typed
		// during the parsing.
		case *IntExpr:
		case *FloatExpr:
		case *BoolExpr:
		case *UnitExpr:

		// We may need some typechecking here
		case *UnaryExpr:
			r := x.(*UnaryExpr).right

			if r, err = aux(r, ctx); err != nil {
				return nil, err
			}

			switch x.(*UnaryExpr).op {
			// Right must be int
			case tokenMinus:
				fallthrough
			case tokenPlus:
				if _, rok := r.getType().(*IntType); !rok {
					return nil, fmt.Errorf("%s : int → int; got %s",
						x.(*UnaryExpr).op, r.getType(),
					)
				}
				x.setType(&IntType{typ{}})

			// Right must be float
			case tokenFMinus:
				fallthrough
			case tokenFPlus:
				if _, rok := r.getType().(*FloatType); !rok {
					return nil, fmt.Errorf("%s : float → float; got %s",
						x.(*UnaryExpr).op, r.getType(),
					)
				}
				x.setType(&FloatType{typ{}})

			// Right must be bool
			case tokenExcl:
				if _, rok := r.getType().(*BoolType); !rok {
					return nil, fmt.Errorf("%s : bool → bool; got %s",
						x.(*UnaryExpr).op, r.getType(),
					)
				}
				x.setType(&BoolType{typ{}})

			default:
				panic("assert")
			}

		// Again, we may need to typecheck things here
		case *BinaryExpr:
			l := x.(*BinaryExpr).left
			r := x.(*BinaryExpr).right

			if l, err = aux(l, ctx); err != nil {
				return nil, err
			}
			if r, err = aux(r, ctx); err != nil {
				return nil, err
			}

			// NOTE/TODO: maybe generics can help here
			// (quick test yields an issue with the setType(T{typ{}}))

			switch x.(*BinaryExpr).op {
			// (int×int) → int
			case tokenMinus:
				fallthrough
			case tokenPlus:
				fallthrough
			case tokenStar:
				fallthrough
			case tokenSlash:
				_, lok := l.getType().(*IntType)
				_, rok := r.getType().(*IntType)
				if !lok || ! rok {
					return nil, fmt.Errorf("%s : (int×int) → int; got (%s×%s)",
						x.(*BinaryExpr).op, l.getType(), r.getType(),
					)
				}
				x.setType(&IntType{typ{}})

			// (int×int) → bool
			case tokenLessEq:
				fallthrough
			case tokenMoreEq:
				fallthrough
			case tokenLess:
				fallthrough
			case tokenMore:
				_, lok := l.getType().(*IntType)
				_, rok := r.getType().(*IntType)
				if !lok || ! rok {
					return nil, fmt.Errorf("%s : (int×int) → bool; got (%s×%s)",
						x.(*BinaryExpr).op, l.getType(), r.getType(),
					)
				}
				x.setType(&BoolType{typ{}})

			// (float×float) → float
			case tokenFMinus:
				fallthrough
			case tokenFPlus:
				fallthrough
			case tokenFStar:
				fallthrough
			case tokenFSlash:
				_, lok := l.getType().(*FloatType)
				_, rok := r.getType().(*FloatType)
				if !lok || ! rok {
					return nil, fmt.Errorf("%s : (float×float) → float; got (%s×%s)",
						x.(*BinaryExpr).op, l.getType(), r.getType(),
					)
				}
				x.setType(&FloatType{typ{}})

			// (float×float) → bool
			case tokenFLessEq:
				fallthrough
			case tokenFMoreEq:
				fallthrough
			case tokenFLess:
				fallthrough
			case tokenFMore:
				_, lok := l.getType().(*FloatType)
				_, rok := r.getType().(*FloatType)
				if !lok || ! rok {
					return nil, fmt.Errorf("%s : (float×float) → float; got (%s×%s)",
						x.(*BinaryExpr).op, l.getType(), r.getType(),
					)
				}
				x.setType(&BoolType{typ{}})

			// Left/right must be bools
			case tokenOrOr:
				fallthrough
			case tokenAndAnd:
				_, lok := l.getType().(*BoolType)
				_, rok := r.getType().(*BoolType)
				if !lok || ! rok {
					return nil, fmt.Errorf("%s : (bool×bool) → bool; got (%s×%s)",
						x.(*BinaryExpr).op, l.getType(), r.getType(),
					)
				}
				x.setType(&BoolType{typ{}})

			default:
				panic("assert")
			}

			x.(*BinaryExpr).left  = l
			x.(*BinaryExpr).right = r

		case *AbsExpr:
			n := x.(*AbsExpr).name
			t := x.(*AbsExpr).typ
			r := x.(*AbsExpr).right

			// save previous ctx[n] if any
			t2, ok := ctx[n]

			// new var in env
			ctx[n] = t

			if r, err = aux(r, ctx); err != nil {
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

			if l, err = aux(l, ctx); err != nil {
				return nil, err
			}

			if r, err = aux(r, ctx); err != nil {
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

		case *ProductExpr:
			l := x.(*ProductExpr).left
			r := x.(*ProductExpr).right

			if l, err = aux(l, ctx); err != nil {
				return nil, err
			}
			if r, err = aux(r, ctx); err != nil {
				return nil, err
			}

			x.setType(&ProductType{typ{}, l.getType(), r.getType()})
			x.(*ProductExpr).left = l
			x.(*ProductExpr).right = r
		default:
			panic("assert")
		}

		return x, nil
	}

	return aux(x, Ctx{})
}

// To ease tests so far
func mustSType(x Expr) Expr {
	y, err := inferSType(x)
	if err != nil {
		panic(err)
	}
	return y
}
