package main

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
		l := applySubst(t.(*ArrowType).left, σ)
		r := applySubst(t.(*ArrowType).right, σ)

		return &ArrowType{typ{}, l, r}

	case *ProductType:
		l := applySubst(t.(*ProductType).left, σ)
		r := applySubst(t.(*ProductType).right, σ)

		return &ProductType{typ{}, l, r}

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

// most general unifier
func mgu() Subst {
	return nil
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
