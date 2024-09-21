package main

import (
	"fmt"
	"io"
	"strings"
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

// Simple type (no polymorphism)
type Type interface {
	aType()
}

type typ struct {}
func (t *typ) aType() {}

type arrowType struct {
	typ
	left, right *Type
}

type productType struct {
	typ
	left, right *Type
}

type unitType struct {
	typ
}

type boolType struct {
	typ
}

type intType struct {
	typ
}

type floatType struct {
	typ
}

// NOTE: I'm not sure we can implement a recursive union type
// easily with a generic: the compiler complains about recursivity,
// and we need our sub-types depending on Expr (e.g. AbsExpr) to be
// parametrized as well. But, I haven't digged too deep either.
type Expr interface {
	aExpr()
}

type expr struct{
	typ Type
}

func (e *expr) aExpr() {}

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
}

type AbsExpr struct {
	expr
	bound string
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

type parser struct {
	scanner
	errf func(string, ...interface{})
}

func (p *parser) errHeref(m string, args ...interface{}) error {
	return fmt.Errorf("%s:%d:%d: %s", p.fn,
		p.tok.ln, p.tok.cn,
		fmt.Sprintf(m, args...))
}

func (p *parser) init(in io.Reader, fn string) {
	p.scanner.init(in, fn)
	p.errf = func(m string, args ...interface{}) {
		panic(p.errHeref(m, args...))
	}
}

func (p *parser) next() token {
	// XXX refine/clarify
	if !p.scanner.next() && p.scan.Err() != nil {
		panic(p.scan.Err())
	}
	return p.tok
}

/*
func mapsTo() {
}

func product() {
}

func type() {
}

func pair() {
}
*/

/*

func let() {
}

func match() {
}

etc.

*/

func (p *parser) PrimitiveType() Type {
	switch k := p.tok.kind; k {
	case tokenTBool:
		p.next()
		return &boolType{}
	case tokenTInt:
		p.next()
		return &intType{}
	case tokenTFloat:
		p.next()
		return &floatType{}
	case tokenTUnit:
		p.next()
		return &unitType{}
	default:
		p.errf("Unexpected token: %s", k.String())
	}
	return nil
}

func (p *parser) ProductType() Type {
	return p.PrimitiveType()
}

// product (×) binds stronger than arrows; arrow is right
// associative.
func (p *parser) ArrowType() Type {
	return p.ProductType()
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
		return &FloatExpr{expr{&floatType{}}, (float64(a) + (b / c))}
	}
	return &IntExpr{expr{&intType{}}, a}
}

func (p *parser) bool() *BoolExpr {
	v := true
	if p.tok.raw == "false" {
		v = false
	}
	p.next()
	return &BoolExpr{expr{&boolType{}}, v}
}

func (p *parser) star() *UnitExpr {
	p.next()
	return &UnitExpr{expr{&unitType{}}}
}

func (p *parser) parenExpr() Expr {
	p.next()
	e := p.appExpr()
	if k := p.tok.kind; k != tokenRParen {
		p.errf("Expecting left paren, got: %s", k.String())
	}
	p.next()
	return e
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
	case tokenMinus, tokenPlus, tokenFMinus, tokenFPlus:
		return p.unaryOpExpr()
	case tokenName:
		return p.varExpr()
	default:
		p.errf("Unexpected token: %s", k.String())
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

func (p *parser) absExpr() Expr {
	if p.tok.kind != tokenLambda {
		x := p.binaryExprs()

		// is this the short form: "x. [...]" instead of "λx. [...]"
		y, ok := x.(*VarExpr)
		if !ok || p.tok.kind != tokenDot {
			return x
		}
		p.next()
		r := p.appExpr()
		return &AbsExpr{expr{}, y.name, r}
	}

	p.next()
	if p.tok.kind != tokenName {
		p.errf("Expecting variable name after lambda, got: %s", p.tok.kind.String())
	}
	n := p.tok.raw
	p.next()

	// a type information may be supplied
	e := expr{}
	if p.tok.kind == tokenColon {
		p.next()
		typ := p.Type()
		e = expr{typ}
	}

	if p.tok.kind != tokenDot {
		p.errf("Expecting dot after lambda variable name, got: %s", p.tok.kind.String())
	}
	p.next()
	return &AbsExpr{e, n, p.appExpr()}
}

func (p *parser) appExpr() Expr {
	l := p.absExpr()

	// XXX too fragile?
	for p.tok.kind != tokenEOF && p.tok.kind != tokenRParen {
		r := p.absExpr()
		l = &AppExpr{expr{}, l, r}
	}

	return l
}

// parsing entry point, only called once.
func (p *parser) parse() (e Expr, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			fmt.Println(err)
		}
	}()

	p.next()
	e = p.appExpr()
	return e, err
}

func parse(in io.Reader, fn string) (Expr, error) {
	var p parser
	p.init(in, fn)
	e, err := p.parse()
	// remaining input is unexpected
	if err == nil && p.tok.kind != tokenEOF {
		err = p.errHeref("Unexpected token: %s", p.tok.kind.String())
	}
	return e, err
}

// For eval_test.go and utils_test.go so far.
func mustParse(s string) Expr {
	e, err := parse(strings.NewReader(s), "")
	if err != nil {
		panic(err)
	}
	return e
}
