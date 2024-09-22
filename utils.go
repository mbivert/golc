package main

import (
	"fmt"
	"strconv"
)

// Compute all free variables within a given expression
// NOTE: this is more efficient than the previous version,
// but perhaps the previous version would still be preferable,
// were we to store the actual free variables at each nodes where
// we'll need it.
func freeVars(x Expr) map[string]bool {
	var aux func(Expr, map[string]bool) map[string]bool

	aux = func(x Expr, m map[string]bool) map[string]bool {
		switch x.(type) {
		case *VarExpr:
			m[x.(*VarExpr).name] = true
		case *AbsExpr:
			// If the variable being bound by the current abstraction
			// has already been declared higher up, then we don't
			// want to remove it. Think:
			//	(λx. y (λy. x y z))
			// Here, the second y is bound, but the first one is free.
			_, hasBefore := m[x.(*AbsExpr).bound]
			aux(x.(*AbsExpr).right, m)
			if !hasBefore {
				delete(m, x.(*AbsExpr).bound)
			}
		case *AppExpr:
			aux(x.(*AppExpr).left, m)
			aux(x.(*AppExpr).right, m)
		case *UnaryExpr:
			aux(x.(*UnaryExpr).right, m)
		case *BinaryExpr:
			aux(x.(*BinaryExpr).left, m)
			aux(x.(*BinaryExpr).right, m)

		// *IntExpr
		// *FloatExpr
		// *BoolExpr
		default:
		}
		return m
	}

	return aux(x, map[string]bool{})
}

func allVars(x Expr) map[string]bool {
	var aux func(Expr, map[string]bool) map[string]bool

	aux = func(x Expr, m map[string]bool) map[string]bool {
		switch x.(type) {
		case *VarExpr:
			m[x.(*VarExpr).name] = true
		case *AbsExpr:
			aux(x.(*AbsExpr).right, m)
		case *AppExpr:
			aux(x.(*AppExpr).left, m)
			aux(x.(*AppExpr).right, m)
		case *UnaryExpr:
			aux(x.(*UnaryExpr).right, m)
		case *BinaryExpr:
			aux(x.(*BinaryExpr).left, m)
			aux(x.(*BinaryExpr).right, m)

		// *IntExpr
		// *FloatExpr
		// *BoolExpr
		default:
		}

		return m
	}

	return aux(x, map[string]bool{})
}

func prettyPrint(x Expr) string {
	var aux func(Expr, bool, bool) string

	aux = func(x Expr, inAbs, inApp bool) string {
		switch x.(type) {
		case *VarExpr:
			return x.(*VarExpr).name
		case *AbsExpr:
			if inAbs {
				return fmt.Sprintf("%s. %s",
					x.(*AbsExpr).bound,
					aux(x.(*AbsExpr).right, true, false))
			} else {
				return fmt.Sprintf("(λ%s. %s)",
					x.(*AbsExpr).bound,
					aux(x.(*AbsExpr).right, true, false))
			}
		case *AppExpr:
			if inApp {
				return fmt.Sprintf("%s %s",
					aux(x.(*AppExpr).left, false, true),
					aux(x.(*AppExpr).right, false, false))
			} else {
				return fmt.Sprintf("(%s %s)",
					aux(x.(*AppExpr).left, false, true),
					aux(x.(*AppExpr).right, false, false))
			}

		// TODO: I'm sure we can do better for those two
		case *UnaryExpr:
			return fmt.Sprintf("%s (%s)",
				x.(*UnaryExpr).op,
				aux(x.(*UnaryExpr).right, false, false))
		case *BinaryExpr:
			return fmt.Sprintf("(%s %s %s)",
				prettyPrint(x.(*BinaryExpr).left),
				x.(*BinaryExpr).op,
				aux(x.(*BinaryExpr).right, false, false))

		case *IntExpr:
			return strconv.FormatInt(x.(*IntExpr).v, 10)
		case *FloatExpr:
			return strconv.FormatFloat(x.(*FloatExpr).v, 'g', -1, 64)
		case *BoolExpr:
			return strconv.FormatBool(x.(*BoolExpr).v)
		default:
			panic("O__o") // TODO
		}
	}

	return aux(x, false, false)
}

func getFresh(m map[string]bool) string {
	for n := 0;; n++ {
		s := fmt.Sprintf("x%d", n)
		if _, ok := m[s]; !ok {
			return s
		}
	}
}

/*
type DeBruijnBVarExpr struct {
	expr
	n int
}

type DeBruijnAbsExpr struct {
	expr
	right Expr
}

// For now just a toy
// https://plfa.github.io/DeBruijn/ (TODO: read)
func toDeBruijn(x Expr) Expr {
}
*/
