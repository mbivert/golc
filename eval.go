/*
 * Evaluation. We assume things have been typechecked.
 */
package main

func evalUnaryExpr(x *UnaryExpr) (Expr, error) {
	r, err := evalExpr(x.right)
	if err != nil {
		return nil, err
	}

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
		return &IntExpr{expr{&IntType{typ{}}}, int64Ops[x.op](r.(*IntExpr).v)}, nil

	case tokenFPlus:
		fallthrough
	case tokenFMinus:
		return &FloatExpr{expr{&FloatType{typ{}}}, float64Ops[x.op](r.(*FloatExpr).v)}, nil

	case tokenExcl:
		return &BoolExpr{expr{&BoolType{typ{}}}, !r.(*BoolExpr).v}, nil

	default:
		panic("TODO: " + x.op.String())
	}

	return nil, nil
}

func evalBinaryExpr(x *BinaryExpr) (Expr, error) {
	l, err := evalExpr(x.left)
	if err != nil {
		return nil, err
	}
	r, err := evalExpr(x.right)
	if err != nil {
		return nil, err
	}

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
		}, nil

	case tokenLess:
		fallthrough
	case tokenMore:
		fallthrough
	case tokenLessEq:
		fallthrough
	case tokenMoreEq:
		return &BoolExpr{expr{&BoolType{typ{}}},
			int64CmpOps[x.op](l.(*IntExpr).v, r.(*IntExpr).v),
		}, nil

	case tokenFPlus:
		fallthrough
	case tokenFStar:
		fallthrough
	case tokenFMinus:
		fallthrough
	case tokenFSlash:
		return &FloatExpr{expr{&FloatType{typ{}}},
			float64Ops[x.op](l.(*FloatExpr).v, r.(*FloatExpr).v),
		}, nil

	case tokenFLess:
		fallthrough
	case tokenFMore:
		fallthrough
	case tokenFLessEq:
		fallthrough
	case tokenFMoreEq:
		return &BoolExpr{expr{&BoolType{typ{}}},
			float64CmpOps[x.op](l.(*FloatExpr).v, r.(*FloatExpr).v),
		}, nil

	case tokenAndAnd:
		fallthrough
	case tokenOrOr:
		return &BoolExpr{expr{&BoolType{typ{}}},
			boolOps[x.op](l.(*BoolExpr).v, r.(*BoolExpr).v),
		}, nil

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
		x.(*AppExpr).right = renameExpr(x.(*AppExpr).right, b, a)
		return x

	default:
		panic("assert")
	}

	return nil
}

// NOTE: we're returning an Expr here.
//
// This is because computation is expected to stop on irreducible
// lambda expressions at some point.
func evalExpr(x Expr) (Expr, error) {
	switch x.(type) {
	// NOTE: "cannot fallthrough in type switch"
	case *UnitExpr:
		return x, nil
	case *IntExpr:
		return x, nil
	case *FloatExpr:
		return x, nil
	case *BoolExpr:
		return x, nil
	case *UnaryExpr:
		return evalUnaryExpr(x.(*UnaryExpr))
	case *BinaryExpr:
		return evalBinaryExpr(x.(*BinaryExpr))
	default:
		panic("assert")
	}
	return nil, nil
}

// To ease tests so far
func mustSTypeParse(s string) Expr {
	return mustSType(mustParse(s))
}
