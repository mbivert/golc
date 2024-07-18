package main

import (
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

func prettyPrint(x Expr) {
	x = x
}
