package main

// TODO: parser method naming conventions are irregular

import (
	"fmt"
	"unicode"
)

const (
	_ = iota
	precBool
	precCmp
	precAdd
	precMul
)

var opPrecs = map[tokenKind]int{
	tokenPlus:  precAdd,
	tokenMinus: precAdd,
	tokenStar:  precMul,
	tokenSlash: precMul,

	tokenLess: precCmp,
	tokenMore: precCmp,

	tokenLessEq: precCmp,
	tokenMoreEq: precCmp,

	tokenFPlus:  precAdd,
	tokenFMinus: precAdd,
	tokenFStar:  precMul,
	tokenFSlash: precMul,

	tokenFMore: precCmp,
	tokenFLess: precCmp,

	tokenFLessEq: precCmp,
	tokenFMoreEq: precCmp,

	tokenAndAnd: precBool,
	tokenOrOr:   precBool,
}

// Simple type (no polymorphism). Maybe there's a more
// efficient way to encode all that.
type Type interface {
	aType()
	// NOTE/TODO: the string representation of types
	// is indirectly tested in typing_test.go and might
	// need adjustments / dedicated tests.
	String() string
}

type typ struct{}

func (t *typ) aType()         {}
func (t *typ) String() string { return "" }

type MissingType struct {
	typ
}

type ArrowType struct {
	typ
	left, right Type
}

type ProductType struct {
	typ
	left, right Type
}

type UnitType struct {
	typ
}

type BoolType struct {
	typ
}

type IntType struct {
	typ
}

type FloatType struct {
	typ
}

// type variable
type VarType struct {
	typ
	name string
}

func (t *MissingType) String() string {
	return "<missing>"
}

func (t *ArrowType) String() string {
	return fmt.Sprintf("%s → %s", t.left, t.right)
}

func (t *ProductType) String() string {
	var l, r string

	switch t.left.(type) {
	case *ArrowType:
		l = fmt.Sprintf("(%s)", t.left)
	default:
		l = fmt.Sprintf("%s", t.left)
	}

	switch t.right.(type) {
	case *ArrowType:
		r = fmt.Sprintf("(%s)", t.right)
	default:
		r = fmt.Sprintf("%s", t.right)
	}

	return fmt.Sprintf("%s × %s", l, r)
}

func (t *UnitType) String() string {
	return "*"
}

func (t *BoolType) String() string {
	return "bool"
}

func (t *IntType) String() string {
	return "int"
}

func (t *FloatType) String() string {
	return "float"
}

func (t *VarType) String() string {
	return t.name
}

// NOTE: I'm not sure we can implement a recursive union type
// easily with a generic: the compiler complains about recursivity,
// and we need our sub-types depending on Expr (e.g. AbsExpr) to be
// parametrized as well. But, I haven't digged too deep either.
//
// NOTE: the dummy aExpr() feels now useless because of get/setType().
//
// NOTE: this feels clumsy anyway.
type Expr interface {
	aExpr()
	getType() Type
	setType(Type)

	String() string
}

type expr struct {
	typ Type
}

func (e *expr) aExpr()           {}
func (e *expr) getType() Type    { return e.typ }
func (e *expr) setType(typ Type) { e.typ = typ }
func (e *expr) String() string   { return "" }

type IntExpr struct {
	expr
	v int64
}

type UnitExpr struct {
	expr
}

type FloatExpr struct {
	expr
	v float64
}

type BoolExpr struct {
	expr
	v bool
}

type VarExpr struct {
	expr
	name string
	// fresh bool
}

type AbsExpr struct {
	expr
	// The only type information we parse optional,
	// and pertaining to an abstraction's bounded variable.
	//
	// That type however is merely the left part of
	// an ArrowType{} which'll make the type of the AbsExpr,
	// so we can't fit it in expr.typ
	typ   Type
	name  string
	right Expr
}

type AppExpr struct {
	expr
	left, right Expr
}

// TODO: have a specifc Operator type instead of tokenKind?
type UnaryExpr struct {
	expr
	op    tokenKind
	right Expr
}

type BinaryExpr struct {
	expr
	op          tokenKind
	left, right Expr
}

// NOTE/TODO: probably better with a []Expr, len ≥ 2
type ProductExpr struct {
	expr
	left, right Expr
}

func (e *IntExpr) String() string {
	return fmt.Sprintf("%d", e.v)
}

func (e *UnitExpr) String() string {
	return "*"
}

func (e *FloatExpr) String() string {
	return fmt.Sprintf("%f", e.v)
}

func (e *BoolExpr) String() string {
	return fmt.Sprintf("%t", e.v)
}

func (e *VarExpr) String() string {
	return fmt.Sprintf("%s", e.name)
}

func (e *AbsExpr) String() string {
	return fmt.Sprintf("λ%s:%s.%s", e.name, e.typ, e.right)
}

func (e *AppExpr) String() string {
	return fmt.Sprintf("((%s) %s)", e.left, e.right)
}

func (e *UnaryExpr) String() string {
	return fmt.Sprintf("(%s %s)", e.op, e.right)
}

func (e *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", e.left, e.op, e.right)
}

func (e *ProductExpr) String() string {
	return fmt.Sprintf("〈%s, %s〉", e.left, e.right)
}

type parser struct {
	scanner
	tok  token
	errf func(string, ...interface{})
}

func (p *parser) errHeref(m string, args ...interface{}) error {
	return fmt.Errorf("%s:%d:%d: %s", p.fn,
		p.tok.ln, p.tok.cn,
		fmt.Sprintf(m, args...))
}

func (p *parser) init(src string, fn string) {
	p.scanner.init([]byte(src), fn)
	p.errf = func(m string, args ...interface{}) {
		panic(p.errHeref(m, args...))
	}
}

func (p *parser) next() token {
	p.tok = p.scanner.scan()
	return p.tok
}

// shortcut; trying to avoid the parsing code to dig through
// p.tok directly.
func (p *parser) has(t tokenKind) bool {
	return p.tok.kind == t
}

func (p *parser) PrimitiveType() Type {
	switch k := p.tok.kind; k {
	case tokenTBool:
		p.next()
		return &BoolType{}
	case tokenTInt:
		p.next()
		return &IntType{}
	case tokenTFloat:
		p.next()
		return &FloatType{}
	case tokenTUnit:
		p.next()
		return &UnitType{}
	case tokenLParen:
		p.next()
		t := p.Type()
		if !p.has(tokenRParen) {
			p.errf("Expecting left paren, got: %s", k.String())
		}
		p.next()
		return t
	default:
		p.errf("Unexpected token: %s", k.String())
	}
	return nil
}

// NOTE: in qlambdabook.pdf, <M1, M2, ... > := <M1, <M2, ...>>,
// hence it's only natural for × to be right associative as well
// (I didn't saw such a shortcut being articulated in the λ-calculus
// notes)
func (p *parser) ProductType() Type {
	l := p.PrimitiveType()

	for p.has(tokenProduct) {
		p.next()
		r := p.ProductType()
		l = &ProductType{typ{}, l, r}
	}

	return l
}

// product (×) binds stronger than arrows; arrow is right
// associative.
func (p *parser) ArrowType() Type {
	l := p.ProductType()

	for p.has(tokenArrow) {
		p.next()
		r := p.ArrowType()
		l = &ArrowType{typ{}, l, r}
	}

	return l
}

func (p *parser) Type() Type {
	return p.ArrowType()
}

// TODO: Rename IntExpr to IntLit & cie?
func (p *parser) number() Expr {
	xs := []byte(p.tok.raw)
	k := p.tok.kind

	// parsing x = a + b; b < 1
	var a int64
	var b float64
	var c float64

	q := 0
	for ; q < len(xs); q++ {
		if !unicode.IsDigit(rune(xs[q])) {
			break
		}
		a = 10*a + int64(xs[q]-'0')
	}

	if q < len(xs) && xs[q] == '.' {
		for c, q = 1, q+1; q < len(xs); q++ {
			if xs[q] < '0' || xs[q] > '9' {
				break
			}
			b = 10*b + float64(xs[q]-'0')
			c *= 10
		}
	}

	p.next()
	if k == tokenFloat {
		return &FloatExpr{expr{&FloatType{}}, (float64(a) + (b / c))}
	}
	return &IntExpr{expr{&IntType{}}, a}
}

func (p *parser) bool() *BoolExpr {
	v := true
	if p.tok.raw == "false" {
		v = false
	}
	p.next()
	return &BoolExpr{expr{&BoolType{}}, v}
}

func (p *parser) star() *UnitExpr {
	p.next()
	return &UnitExpr{expr{&UnitType{}}}
}

func (p *parser) parenExpr() Expr {
	p.next()
	x := p.appExpr()
	if !p.has(tokenRParen) {
		p.errf("Expecting left paren, got: %s", p.tok.kind.String())
	}
	p.next()
	return x
}

func (p *parser) unaryOpExpr() *UnaryExpr {
	o := p.tok.kind
	p.next()
	return &UnaryExpr{expr{}, o, p.binaryExprs()}
}

func (p *parser) varExpr() *VarExpr {
	n := p.tok.raw
	p.next()
	return &VarExpr{expr{}, n}
}

// NOTE: we're using 〈〉 over <> to avoid confusion with < as an operator
// (e.g. <x, 1> will mess things up: parseBinary will expects something after
// the 1, and not consider it the end of a product)
func (p *parser) productExpr() Expr {
	p.next()

	var ret *ProductExpr

	x := &ret

	for {
		y := p.appExpr()

		hasComa := p.has(tokenComa)
		hasRBracket := p.has(tokenRBracket)

		// <Y> parsed as Y
		if hasRBracket && *x == nil {
			p.next()
			return y
		}

		// first element of a pair
		if hasComa && *x == nil {
			p.next()
			*x = &ProductExpr{expr{}, y, nil}
			continue
		}

		if hasComa || hasRBracket {
			p.next()
			if (*x).right == nil {
				(*x).right = y
			} else {
				z := (*x).right
				t := &ProductExpr{expr{}, z, y}
				(*x).right = t
				x = &t
			}
		}

		if hasRBracket {
			return ret
		}
	}

	return ret
}

func (p *parser) unaryExpr() Expr {
	switch k := p.tok.kind; k {
	case tokenInt, tokenFloat:
		return p.number()
	case tokenStar:
		return p.star()
	case tokenBool:
		return p.bool()
	case tokenLParen:
		return p.parenExpr()
	case tokenMinus, tokenPlus, tokenFMinus, tokenFPlus, tokenExcl:
		return p.unaryOpExpr()
	case tokenName:
		return p.varExpr()
	case tokenLBracket:
		return p.productExpr()
	default:
		p.errf("Unexpected token: %s", k)
	}
	return nil
}

func (p *parser) hasOp() int {
	x, ok := opPrecs[p.tok.kind]
	if !ok {
		x = -1
	}
	return x
}

func (p *parser) binaryExpr(prec int) Expr {
	left := p.unaryExpr()

	// we start the recursive parsing with prec == 0, so we'll
	// have to get there again and slurp the whole expression
	// (all genuine operators have an precedence > 0)
	for x := p.hasOp(); x > prec; x = p.hasOp() {
		op := p.tok.kind
		p.next()
		right := p.binaryExpr(x)
		left = &BinaryExpr{expr{}, op, left, right}
	}

	return left
}

func (p *parser) binaryExprs() Expr {
	return p.binaryExpr(0)
}

// XXX naming convention is confusing
//
// TODO: no rec, no let 〈x,y,...〉, no let *
func (p *parser) letIn() Expr {
	p.next()

	if !p.has(tokenName) {
		p.errf("Expecting variable name after let, got: %s", p.tok.kind)
	}

	n := p.varExpr()

	if !p.has(tokenEqual) {
		p.errf("Expecting equal after let $x, got: %s", p.tok.kind)
	}

	p.next()

	x := p.appExpr()

	t := Type(&typ{})

	if p.has(tokenColon) {
		p.next()
		t = p.Type()
	}

	if !p.has(tokenIn) {
		p.errf("Expecting 'in' after let $x = $M, got %s", p.tok.kind)
	}

	p.next()

	y := p.appExpr()

	// Desugar now; perhaps we'd want to have a dedicated pass.
	// XXX meh, no typing annotation
	return &AppExpr{expr{},
		&AbsExpr{expr{},
			//			&MissingType{typ{}},
			t,
			n.name,
			y,
		},
		x,
	}
}

func (p *parser) absExpr() Expr {
	var n string

	// TODO: hopefully this is good enough to insert it here
	if p.has(tokenLet) {
		return p.letIn()
	}

	if !p.has(tokenLambda) {
		x := p.binaryExprs()

		// is this the short form: "x. [...]" instead of "λx. [...]"
		// (eventually with a type annotation)
		y, ok := x.(*VarExpr)

		// not a VarExpr: definitely not a short form
		// not followed by either a dot or a colon: not a short form either
		if !ok || (!p.has(tokenDot) && !p.has(tokenColon)) {
			return x
		}

		n = y.name
	} else {
		p.next()
		if !p.has(tokenName) {
			p.errf("Expecting variable name after lambda, got: %s", p.tok.kind.String())
		}
		n = p.tok.raw
		p.next()
	}

	// a type information may be supplied
	//	t := Type(&MissingType{typ{}})
	t := Type(&typ{})
	if p.has(tokenColon) {
		p.next()
		t = p.Type()
	}

	if !p.has(tokenDot) {
		p.errf("Expecting dot after lambda variable name, got: %s", p.tok.kind.String())
	}
	p.next()

	return &AbsExpr{expr{}, t, n, p.appExpr()}
}

// tokens marking the end of an application. parser.appExpr()
// is the parsing entry point: we get back there again in a few
// cases (parser.parenExpr(), parser.productExpr(), parser.letIn())
// and need to detect the end of such cases.
var endAppExpr = map[tokenKind]bool{
	// nothing else to parse
	tokenEOF: true,

	// we were parsing something between parenthesis
	tokenRParen: true,

	// we were parsing something between brackets (product)
	tokenRBracket: true,

	// we're parsing something between brackets (product)
	tokenComa: true,

	// we just parsed the expression $expr associated to a bound
	// name $x of a let/in construct (let $x = $expr in ...)
	tokenIn: true,

	tokenColon: true,
}

func (p *parser) appExpr() Expr {
	l := p.absExpr()

	for {
		if _, stop := endAppExpr[p.tok.kind]; stop {
			break
		}
		r := p.absExpr()
		l = &AppExpr{expr{}, l, r}
	}

	return l
}

// parsing entry point, only called once.
func (p *parser) parse() (x Expr, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			fmt.Println(err)
		}
	}()

	p.next()
	x = p.appExpr()
	return x, err
}

func parse(src string, fn string) (Expr, error) {
	var p parser
	p.init(src, fn)
	x, err := p.parse()
	// remaining input is unexpected
	if err == nil && !p.has(tokenEOF) {
		err = p.errHeref("Unexpected token: %s", p.tok.kind.String())
	}
	return x, err
}

// To ease tests so far
func mustParse(src string) Expr {
	x, err := parse(src, "")
	if err != nil {
		panic(err)
	}
	return x
}
