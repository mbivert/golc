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
		ir, ok := r.(*IntExpr)
		if !ok {
			panic("oooo")
		}
		return &IntExpr{expr{}, int64Ops[x.op](ir.v)}, nil

	case tokenFPlus:
		fallthrough
	case tokenFMinus:
		ir, ok := r.(*FloatExpr)
		if !ok {
			panic("oooo")
		}
		return &FloatExpr{expr{}, float64Ops[x.op](ir.v)}, nil
	default:
		panic("Unexpected unary operator " + x.op.String())
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

	/*
		int64CmpOps := map[tokenKind](func (int64,int64) bool) {
			tokenLess   : func (a, b int64) bool { return a<b  },
			tokenMore   : func (a, b int64) bool { return a>b  },
			tokenLessEq : func (a, b int64) bool { return a<=b },
			tokenMoreEq : func (a, b int64) bool { return a>=b },
		}
	*/

	float64Ops := map[tokenKind](func(float64, float64) float64){
		tokenPlus:  func(a, b float64) float64 { return a + b },
		tokenStar:  func(a, b float64) float64 { return a * b },
		tokenMinus: func(a, b float64) float64 { return a - b },
		tokenSlash: func(a, b float64) float64 { return a / b },
	}

	/*
		float64CmpOps := map[tokenKind](func (float64,float64) bool) {
			tokenFLess   : func (a, b float64) bool { return a<b  },
			tokenFMore   : func (a, b float64) bool { return a>b  },
			tokenFLessEq : func (a, b float64) bool { return a<=b },
			tokenFMoreEq : func (a, b float64) bool { return a>=b },
		}
	*/

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
		// TODO/XXX Theoretically, we should be running our typechecking
		// before starting the evaluation, because we're doing static
		// typechecking; so should we certain that this should hold?
		il, ok := l.(*IntExpr)
		if !ok {
			panic("ooo")
		}
		ir, ok := r.(*IntExpr)
		if !ok {
			panic("oooo")
		}
		return &IntExpr{expr{}, int64Ops[x.op](il.v, ir.v)}, nil

	case tokenFPlus:
		fallthrough
	case tokenFStar:
		fallthrough
	case tokenFMinus:
		fallthrough
	case tokenFSlash:
		// TODO/XXX Theoretically, we should be running our typechecking
		// before starting the evaluation, because we're doing static
		// typechecking; so should we certain that this should hold?
		il, ok := l.(*FloatExpr)
		if !ok {
			panic("ooo")
		}
		ir, ok := r.(*FloatExpr)
		if !ok {
			panic("oooo")
		}
		return &FloatExpr{expr{}, float64Ops[x.op](il.v, ir.v)}, nil

	default:
		panic("TODO")
	}
}

// α-renaming x{b,a}: renaming a as b in x.
//
// renaming is performed in-place (why not I guess?)
func renameExpr(x Expr, b, a string) Expr {
	switch x.(type) {

	// NOTE: "cannot fallthrough in type switch"
	case *UnitExpr:
		return x
	case *IntExpr:
		return x
	case *FloatExpr:
		return x
	case *BoolExpr:
		return x

	case *UnaryExpr:
		x.(*UnaryExpr).right = renameExpr(x.(*UnaryExpr).right, b, a)
		return x
	case *BinaryExpr:
		x.(*BinaryExpr).left = renameExpr(x.(*BinaryExpr).left, b, a)
		x.(*BinaryExpr).right = renameExpr(x.(*BinaryExpr).right, b, a)
		return x
	}

	return nil
}

// NOTE: we're returning an Expr here.
// This is because I expect computation to stop on irreducible
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
	}
	return nil, nil
}
