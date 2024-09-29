/*
 * Evaluation. We assume things have been typechecked.
 */
package main

import "fmt"

func evalUnaryExpr(x *UnaryExpr) Expr {
	r, _ := reduceExpr(x.right)

	int64Ops := map[tokenKind](func(int64) int64){
		tokenPlus:  func(a int64) int64 { return a },
		tokenMinus: func(a int64) int64 { return -a },
	}

	float64Ops := map[tokenKind](func(float64) float64){
		tokenPlus:  func(a float64) float64 { return a },
		tokenMinus: func(a float64) float64 { return -a },
	}

	switch x.op {
	case tokenPlus:
		fallthrough
	case tokenMinus:
		return &IntExpr{expr{&IntType{typ{}}}, int64Ops[x.op](r.(*IntExpr).v)}

	case tokenFPlus:
		fallthrough
	case tokenFMinus:
		return &FloatExpr{expr{&FloatType{typ{}}}, float64Ops[x.op](r.(*FloatExpr).v)}

	case tokenExcl:
		return &BoolExpr{expr{&BoolType{typ{}}}, !r.(*BoolExpr).v}

	default:
		panic("TODO: " + x.op.String())
	}

	return nil
}

func evalBinaryExpr(x *BinaryExpr) Expr {
	l, _ := reduceExpr(x.left)
	r, _ := reduceExpr(x.right)

	int64Ops := map[tokenKind](func(int64, int64) int64){
		tokenPlus:  func(a, b int64) int64 { return a + b },
		tokenStar:  func(a, b int64) int64 { return a * b },
		tokenMinus: func(a, b int64) int64 { return a - b },
		tokenSlash: func(a, b int64) int64 { return a / b },
	}

	int64CmpOps := map[tokenKind](func(int64, int64) bool){
		tokenLess:   func(a, b int64) bool { return a < b },
		tokenMore:   func(a, b int64) bool { return a > b },
		tokenLessEq: func(a, b int64) bool { return a <= b },
		tokenMoreEq: func(a, b int64) bool { return a >= b },
	}

	float64Ops := map[tokenKind](func(float64, float64) float64){
		tokenPlus:  func(a, b float64) float64 { return a + b },
		tokenStar:  func(a, b float64) float64 { return a * b },
		tokenMinus: func(a, b float64) float64 { return a - b },
		tokenSlash: func(a, b float64) float64 { return a / b },
	}

	float64CmpOps := map[tokenKind](func(float64, float64) bool){
		tokenFLess:   func(a, b float64) bool { return a < b },
		tokenFMore:   func(a, b float64) bool { return a > b },
		tokenFLessEq: func(a, b float64) bool { return a <= b },
		tokenFMoreEq: func(a, b float64) bool { return a >= b },
	}

	boolOps := map[tokenKind](func(bool, bool) bool){
		tokenAndAnd: func(a, b bool) bool { return a && b },
		tokenOrOr:   func(a, b bool) bool { return a || b },
	}

	switch x.op {
	// XXX/TODO: should we allow e.g. x + 3? where x
	// is undefined (why not I guess?)
	case tokenPlus:
		fallthrough
	case tokenStar:
		fallthrough
	case tokenMinus:
		fallthrough
	case tokenSlash:
		return &IntExpr{expr{&IntType{typ{}}},
			int64Ops[x.op](l.(*IntExpr).v, r.(*IntExpr).v),
		}

	case tokenLess:
		fallthrough
	case tokenMore:
		fallthrough
	case tokenLessEq:
		fallthrough
	case tokenMoreEq:
		return &BoolExpr{expr{&BoolType{typ{}}},
			int64CmpOps[x.op](l.(*IntExpr).v, r.(*IntExpr).v),
		}

	case tokenFPlus:
		fallthrough
	case tokenFStar:
		fallthrough
	case tokenFMinus:
		fallthrough
	case tokenFSlash:
		return &FloatExpr{expr{&FloatType{typ{}}},
			float64Ops[x.op](l.(*FloatExpr).v, r.(*FloatExpr).v),
		}

	case tokenFLess:
		fallthrough
	case tokenFMore:
		fallthrough
	case tokenFLessEq:
		fallthrough
	case tokenFMoreEq:
		return &BoolExpr{expr{&BoolType{typ{}}},
			float64CmpOps[x.op](l.(*FloatExpr).v, r.(*FloatExpr).v),
		}

	case tokenAndAnd:
		fallthrough
	case tokenOrOr:
		return &BoolExpr{expr{&BoolType{typ{}}},
			boolOps[x.op](l.(*BoolExpr).v, r.(*BoolExpr).v),
		}

	default:
		panic("TODO: " + x.op.String())
	}
}

// α-renaming x{b,a}: renaming a as b in x.
//
// renaming is performed in-place (why not I guess?)
func renameExpr(x Expr, b, a string) Expr {
	switch x.(type) {
	case *UnitExpr:
		return x
	case *IntExpr:
		return x
	case *FloatExpr:
		return x
	case *BoolExpr:
		return x
	case *ProductExpr:
		x.(*ProductExpr).left = renameExpr(x.(*ProductExpr).left, b, a)
		x.(*ProductExpr).right = renameExpr(x.(*ProductExpr).right, b, a)
		return x

	case *UnaryExpr:
		x.(*UnaryExpr).right = renameExpr(x.(*UnaryExpr).right, b, a)
		return x

	case *BinaryExpr:
		x.(*BinaryExpr).left = renameExpr(x.(*BinaryExpr).left, b, a)
		x.(*BinaryExpr).right = renameExpr(x.(*BinaryExpr).right, b, a)
		return x

	case *AppExpr:
		x.(*AppExpr).left = renameExpr(x.(*AppExpr).left, b, a)
		x.(*AppExpr).right = renameExpr(x.(*AppExpr).right, b, a)
		return x

	case *VarExpr:
		if x.(*VarExpr).name == a {
			x.(*VarExpr).name = b
		}
		return x

	case *AbsExpr:
		if x.(*AbsExpr).name == a {
			x.(*AbsExpr).name = b
		}
		x.(*AbsExpr).right = renameExpr(x.(*AbsExpr).right, b, a)
		return x

	default:
		panic("assert")
	}

	return nil
}

// β-substitution: x[y/a]: substituing a for y in x
func substituteExpr(x, y Expr, a string) Expr {
	switch x.(type) {
	case *UnitExpr:
		return x
	case *IntExpr:
		return x
	case *FloatExpr:
		return x
	case *BoolExpr:
		return x
	case *ProductExpr:
		x.(*ProductExpr).left = substituteExpr(x.(*ProductExpr).left, y, a)
		x.(*ProductExpr).right = substituteExpr(x.(*ProductExpr).right, y, a)
		return x

	case *UnaryExpr:
		x.(*UnaryExpr).right = substituteExpr(x.(*UnaryExpr).right, y, a)
		return x

	case *BinaryExpr:
		x.(*BinaryExpr).left = substituteExpr(x.(*BinaryExpr).left, y, a)
		x.(*BinaryExpr).right = substituteExpr(x.(*BinaryExpr).right, y, a)
		return x

	case *AppExpr:
		x.(*AppExpr).left = substituteExpr(x.(*AppExpr).left, y, a)
		x.(*AppExpr).right = substituteExpr(x.(*AppExpr).right, y, a)
		return x

	case *VarExpr:
		if x.(*VarExpr).name == a {
			return y
		}
		return x

	case *AbsExpr:
		name := x.(*AbsExpr).name
		if name == a {
			return x
		}
		if !isFree(y, name) {
			x.(*AbsExpr).right = substituteExpr(x.(*AbsExpr).right, y, a)
			return x
		}
		// bounded variable name of x occurs freely in y:
		// if we're about so swap a for y below the current
		// abstraction, we need to make sure our name won't
		// conflict with what happens in y. Hence, we need
		// to get a name which would conflict with nothing in
		// x, y or a for that matter.
		b := getFresh(allVars(x.(*AbsExpr).right), allVars(y), map[string]bool{a:true})
		x.(*AbsExpr).name = b
		x.(*AbsExpr).right = substituteExpr(
			renameExpr(x.(*AbsExpr).right, b, name), y, a,
		)
		return x

	default:
		panic("assert")
	}

	return nil
}

func reduceExpr(x Expr) (Expr, bool) {
	switch x.(type) {
	// NOTE: "cannot fallthrough in type switch"
	case *UnitExpr:
		return x, false

	case *IntExpr:
		return x, false

	case *FloatExpr:
		return x, false

	case *BoolExpr:
		return x, false

	case *UnaryExpr:
		return evalUnaryExpr(x.(*UnaryExpr)), true

	case *BinaryExpr:
		return evalBinaryExpr(x.(*BinaryExpr)), true

	case *AbsExpr:
		var b bool
		x.(*AbsExpr).right, b = reduceExpr(x.(*AbsExpr).right)
		return x, b

	case *VarExpr:
		return x, false

	case *AppExpr:
		// XXX hmm, will this always be an AbsEexpr?
		if _, ok := x.(*AppExpr).left.(*AbsExpr); ok {
			return substituteExpr(
				x.(*AppExpr).left.(*AbsExpr).right,
				x.(*AppExpr).right,
				x.(*AppExpr).left.(*AbsExpr).name,
			), true
		}
		var bl, br bool

		x.(*AppExpr).left, bl = reduceExpr(x.(*AppExpr).left)
		x.(*AppExpr).right, br = reduceExpr(x.(*AppExpr).right)
		return x, bl || br

	default:
		panic("assert")
	}
}

// NOTE: we're returning an Expr here.
//
// This is because computation is expected to stop on irreducible
// lambda expressions at some point.
//
// TODO: add a configurable timeout here
// TODO: termination detection feels clumsy as hell; we can't
// compare x with y, as reduction will modify its input in-place.
func evalExpr(x Expr) Expr {
	var y Expr
	var b bool

	for {
		fmt.Printf("%s\n", x)
		y, b = reduceExpr(x)

		if !b {
			return x
		}
		x = y
	}
}

// To ease tests so far
func mustSTypeParse(s string) Expr {
	return mustSType(mustParse(s))
}
