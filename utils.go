package main

import (
	"fmt"
	"strconv"
)

// Add n entries to m; returns m
func mergeMaps(m, n map[string]bool) map[string]bool {
	for k, v := range n {
		m[k] = v
	}
	return m
}

// Compute all free variables within a given expression
// TODO: we probably want to make this way more efficient.
func freeVars(x Expr) map[string]bool {
	switch x.(type) {
	case *VarExpr:
		return map[string]bool{x.(*VarExpr).name : true}
	case *AbsExpr:
		m := freeVars(x.(*AbsExpr).right)
		delete(m, x.(*AbsExpr).bound)
		return m
	case *AppExpr:
		return mergeMaps(freeVars(x.(*AppExpr).left), freeVars(x.(*AppExpr).right))
	case *UnaryExpr:
		return freeVars(x.(*UnaryExpr).right)
	case *BinaryExpr:
		return mergeMaps(freeVars(x.(*BinaryExpr).left), freeVars(x.(*BinaryExpr).right))

	// *IntExpr
	// *FloatExpr
	// *BoolExpr
	default:
		return make(map[string]bool)
	}
}

func allVars(x Expr) map[string]bool {
	switch x.(type) {
	case *VarExpr:
		return map[string]bool{x.(*VarExpr).name : true}
	case *AbsExpr:
		return freeVars(x.(*AbsExpr).right)
	case *AppExpr:
		return mergeMaps(freeVars(x.(*AppExpr).left), freeVars(x.(*AppExpr).right))
	case *UnaryExpr:
		return freeVars(x.(*UnaryExpr).right)
	case *BinaryExpr:
		return mergeMaps(freeVars(x.(*BinaryExpr).left), freeVars(x.(*BinaryExpr).right))

	// *IntExpr
	// *FloatExpr
	// *BoolExpr
	default:
		return make(map[string]bool)
	}
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
