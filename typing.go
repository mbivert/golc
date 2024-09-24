package main

import "fmt"

type Subst map[string]Type

// NOTE/TODO: for now, applySubst operates on types, but we may
// have to act on typed expressions later on. If so, then perhaps
// the parsing should pre-fill the types with something reasonable.

// apply the substitution s to x in-place; returns updated x
func applySubst(t Type, σ Subst) Type {
	switch t.(type) {
	case *VarType:
		v, ok := σ[t.(*VarType).name]
		if ok {
			return v
		}

	case *ArrowType:
		return &ArrowType{
			typ{},
			applySubst(t.(*ArrowType).left, σ),
			applySubst(t.(*ArrowType).right, σ),
		}

	case *ProductType:
		return &ProductType{
			typ{},
			applySubst(t.(*ProductType).left, σ),
			applySubst(t.(*ProductType).right, σ),
		}

	// "iotas" (unit / primitive types)
	case *UnitType:
	case *BoolType:
	case *IntType:
	case *FloatType:

	default:
		panic("O__o")
	}

	return t
}

// returns true if n -- assumed to be a VarType's name --
// occurs in t
func occursIn(t Type, n string) bool {
	switch t.(type) {
	case *VarType:
		return t.(*VarType).name == n

	case *ArrowType:
		return occursIn(t.(*ArrowType).left, n) ||
			occursIn(t.(*ArrowType).right, n)

	case *ProductType:
		return occursIn(t.(*ProductType).left, n) ||
			occursIn(t.(*ProductType).right, n)

	// "iotas" (unit / primitive types)
	case *UnitType:
	case *BoolType:
	case *IntType:
	case *FloatType:

	default:
		panic("X_x")
	}

	return false
}

func mguVarType(t Type, n string) (Subst, error) {
	if !occursIn(t, n) {
		// case 2 / 4
		return Subst{n: t}, nil
	} else {
		// case 3 / 5
		return nil, fmt.Errorf("%s occurs in %s", n, t)
	}
}

// Most General Unifier; we're closely following the algorithm
// description, being verbose on purpose/to reflect it ("space
// shuttle style" / "DO NOT ATTEMPT TO SIMPLIFY THIS CODE")
func mgu(as, bs []Type) (Subst, error) {
	if len(as) != len(bs) {
		panic("assert")
	}

	// TODO: may become useless once the end of mgu()
	// is properly wired.
	if len(as) == 0 {
		return Subst{}, nil
	}

	if len(as) == 1 {
		a, b := as[0], bs[0]

		switch a.(type) {
		case *VarType:
			n := a.(*VarType).name
			switch b.(type) {
			case *VarType:
				if n == b.(*VarType).name {
					// case 1.
					// no entry: id() assumed
					return Subst{}, nil
				}
			default:
				// case 2 / 3
				return mguVarType(b, n)
			}
		}

		switch b.(type) {
		case *VarType:
			// case 4 / 5
			return mguVarType(a, b.(*VarType).name)
		}

		// case 6
		switch a.(type) {
		case *BoolType:
			switch b.(type) {
			case *BoolType:
				return Subst{}, nil
			}
		case *IntType:
			switch b.(type) {
			case *IntType:
				return Subst{}, nil
			}
		case *FloatType:
			switch b.(type) {
			case *FloatType:
				return Subst{}, nil
			}
		}

		switch a.(type) {
		case *ArrowType:
			switch b.(type) {
			case *ArrowType:
				return mgu(
					[]Type{
						a.(*ArrowType).left,
						a.(*ArrowType).right,
					},
					[]Type{
						b.(*ArrowType).left,
						b.(*ArrowType).right,
					},
				)
			}
		case *ProductType:
			switch b.(type) {
			case *ProductType:
				return mgu(
					[]Type{
						a.(*ProductType).left,
						a.(*ProductType).right,
					},
					[]Type{
						b.(*ProductType).left,
						b.(*ProductType).right,
					},
				)
			}
		}

		// case 9
		switch a.(type) {
		case *UnitType:
			switch b.(type) {
			case *UnitType:
				return Subst{}, nil
			}
		}

		return nil, fmt.Errorf("Cannot unify %s with %s", a, b)
	}

	ρ, err := mgu(as[1:], bs[1:])
	if err != nil {
		return nil, err
	}

	τ, err := mgu(
		[]Type{applySubst(as[0], ρ)},
		[]Type{applySubst(bs[0], ρ)},
	)
	if err != nil {
		return nil, err
	}

	// TODO: build τ ο ρ; it's likely that we can't just merge
	// τ into ρ
	τ = τ

	return nil, nil
}

func inferType(x, y Expr) (Subst, error) {
	return nil, nil
}

// x                : 'a
// x+3              : int
// (λx. x+3)        : int   → int
// (λx:int. x+3)    : int   → int
// (λx:float. x+3)  : float → float
// (λx:bool. x+3)   : error "bool+int": undefined

/*
	So add int64, float64 and booleans so that
	we can have type inference rules.
*/
