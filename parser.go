package main

import (
	"io"
	"fmt"
	"unicode"
//	"strings"
)

const (
	_ = iota
	precCmp
	precAdd
	precMul
)

var opPrecs = map[tokenKind]int {
	tokenPlus    : precAdd,
	tokenMinus   : precAdd,
	tokenStar    : precMul,
	tokenSlash   : precMul,

	tokenLess    : precCmp,
	tokenMore    : precCmp,

	tokenLessEq  : precCmp,
	tokenMoreEq  : precCmp,

	tokenFPlus   : precAdd,
	tokenFMinus  : precAdd,
	tokenFStar   : precMul,

	tokenFMore   : precCmp,
	tokenFLess   : precCmp,

	tokenFLessEq : precCmp,
	tokenFMoreEq : precCmp,
}

/*
type Type interface {
	aType()
}

type type struct {}
func (t *type) aType() {}

type arrowType struct {
	type
	left, right *Type
}

type productType struct {
	type
	left, right *Type
}

type unitType struct {
	type
}

type boolType struct {
	type
}

type int64Type struct {
	type
}

type float64Type struct {
	type
}

type typeKind uint

const (
	KBool typeKind = iota // Bool
	KNat                  // Natural
	KArrow                // Arrow
	KProduct              // Product
	KUnit                 // Unit
)
*/

/*
// Do we need? As everything is an expression here
type Node interface {
	aNode()
}
type node struct {}

*/

type Expr interface {
	aExpr()
}

type expr struct {}
func (e *expr) aExpr() {}

type IntExpr struct {
	expr
	v int64
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
type UnaryOpExpr struct {
	expr
	op    tokenKind
	right Expr
}

type BinaryExpr struct {
	expr
	op tokenKind
	left, right Expr
}

type parser struct {
	scanner
	errf func(string, ...interface{}) ()
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

// TODO: Rename IntExpr to IntLit & cie?
func (p *parser) number() Expr {
	xs := []byte(p.tok.raw)
	k  := p.tok.kind

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
		return &FloatExpr{expr{}, (float64(a) + (b/c))}
	}
	return &IntExpr{expr{}, a}
}

func (p *parser) bool() *BoolExpr {
	v := true
	if p.tok.raw == "false" {
		v = false
	}
	p.next()
	return &BoolExpr{expr{}, v}
}

func (p *parser) parenExpr() Expr {
	p.next()
	e := p.app()
	if k := p.tok.kind; k != tokenRParen {
		p.errf("Expecting left paren, got: %s", k.String())
	}
	p.next()
	return e
}

func (p *parser) unaryOpExpr() *UnaryOpExpr {
	o := p.tok.kind
	p.next()
	return &UnaryOpExpr{expr{}, o, p.expr()}
}

func (p *parser) xvar() *VarExpr {
	n := p.tok.raw
	p.next()
	return &VarExpr{expr{}, n}
}

func (p *parser) unaryExpr() Expr {
	fmt.Println(p.tok)
	fmt.Println(p.tok.kind)
	switch k := p.tok.kind; k {
	case tokenInt, tokenFloat:
		return p.number()
	case tokenBool:
		return p.bool()
	case tokenLParen:
		return p.parenExpr()
	case tokenMinus, tokenPlus, tokenFMinus, tokenFPlus:
		return p.unaryOpExpr()
	case tokenName:
		return p.xvar()
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

	for x := p.hasOp() ; x > prec; x = p.hasOp() {
		op := p.tok.kind
		p.next()
		right := p.binaryExpr(x)
		left = &BinaryExpr{expr{}, op, left, right}
	}

	return left
}

func (p *parser) expr() Expr {
	return p.binaryExpr(0)
}

func (p *parser) abs() Expr {
	if p.tok.kind != tokenLambda {
		x := p.expr()
		y, ok := x.(*VarExpr)
		if !ok ||  p.tok.kind != tokenDot {
			return x
		}
		p.next()
		r := p.app()
		return &AbsExpr{expr{}, y.name, r}
	}

	p.next()
	if p.tok.kind != tokenName {
		p.errf("Expecting variable name after lambda, got: %s", p.tok.kind.String())
	}
	n := p.tok.raw
	p.next()
	if p.tok.kind != tokenDot {
		p.errf("Expecting dot after lambda variable name, got: %s", p.tok.kind.String())
	}
	p.next()
	return &AbsExpr{expr{}, n, p.app()}
}

func (p *parser) app() Expr {
	l := p.abs()

	// maybe this is a bit too fragile?
	for p.tok.kind != tokenEOF && p.tok.kind != tokenRParen {
		r := p.abs()
		l = &AppExpr{expr{}, l, r}
	}

	return l
}

func (p *parser) parse() (e Expr, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
		}
	}()

	p.next()
	e = p.app()
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
