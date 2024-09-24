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

// most general unifier
func mgu(as, bs []Type) (Subst, error) {
	if len(as) != len(bs) {
		panic("assert")
	}

	if len(as) == 1 {
		a, b := as[0], bs[0]

		switch a.(type) {
		case *VarType:
			switch b.(type) {
			case *VarType:
				if a.(*VarType).name == b.(*VarType).name {
					// no entry: id() assumed
					return Subst{}, nil
				}
			default:
			}
		case *ArrowType:
		case *ProductType:
		case *UnitType:
		case *BoolType:
		case *IntType:
		case *FloatType:
		default:
			panic("O_O")
		}
	}

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
