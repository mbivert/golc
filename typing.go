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

// τ ∘ ρ, at least for some substitutions
func composeSubst(τ, ρ Subst) Subst {
	σ := make(Subst)

	// Import substitution from ρ to σ. In case
	// of substitutions from ρ which will be altered
	// by τ, import that of τ directly.
	for n, t := range ρ {
		σ[n] = t
		if v, ok := t.(*VarType); ok {
			if w, ok := τ[v.name]; ok {
				σ[n] = w
			}
		}
	}

	for n, t := range τ {
		_, ok := ρ[n]

		// that substitution doesn't exist in ρ:
		if !ok {
			σ[n] = t
			continue
		}

		// tricky case: this substitution exists
		// in ρ and in τ.
		//
		// TODO/XXX: let's see how relevant that case
		// practically is
		panic("assert")
	}

	return σ
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

// space shuttle style / DO NOT ATTEMPT TO SIMPLIFY THIS CODE
func mgu1(a, b Type) (Subst, error) {
	if av, ok := a.(*VarType); ok {
		if bv, ok := b.(*VarType); ok {
			if av.name == bv.name {
				// case 1.
				// no entry: id() assumed
				return Subst{}, nil
			}
		} else {
			// case 2 / 3
			return mguVarType(b, av.name)
		}
	}

	if bv, ok := b.(*VarType); ok {
		// case 4 / 5
		return mguVarType(a, bv.name)
	}

	// case 6 (maybe)
	if _, ok := a.(*BoolType); ok {
		if _, ok := b.(*BoolType); ok {
			return Subst{}, nil
		}
	}
	if _, ok := a.(*IntType); ok {
		if _, ok := b.(*IntType); ok {
			return Subst{}, nil
		}
	}
	if _, ok := a.(*FloatType); ok {
		if _, ok := b.(*FloatType); ok {
			return Subst{}, nil
		}
	}

	if av, ok := a.(*ArrowType); ok {
		if bv, ok := b.(*ArrowType); ok {
			// case 7
			return mgu(
				[]Type{av.left, av.right},
				[]Type{bv.left, bv.right},
			)
		}
	}
	if av, ok := a.(*ProductType); ok {
		if bv, ok := b.(*ProductType); ok {
			// case 8
			return mgu(
				[]Type{av.left, av.right},
				[]Type{bv.left, bv.right},
			)
		}
	}

	// case 9
	if _, ok := a.(*UnitType); ok {
		if _, ok := b.(*UnitType); ok {
			return Subst{}, nil
		}
	}

	return nil, fmt.Errorf("Cannot unify '%s' with '%s'", a, b)
}

// Most General Unifier; we're closely following the algorithm
// description, being verbose on purpose/to reflect it.
//
// space shuttle style / DO NOT ATTEMPT TO SIMPLIFY THIS CODE
func mgu(as, bs []Type) (Subst, error) {
	if len(as) != len(bs) {
		panic("assert")
	}

	if len(as) == 0 {
		return Subst{}, nil
	} else if len(as) == 1 {
		return mgu1(as[0], bs[0])
	}

	ρ, err := mgu(as[1:], bs[1:])
	if err != nil {
		return nil, err
	}

	τ, err := mgu1(
		applySubst(as[0], ρ),
		applySubst(bs[0], ρ),
	)

	if err != nil {
		return nil, err
	}

	return composeSubst(τ, ρ), nil
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
