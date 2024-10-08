/*
 * Scanner/lexer, modelled after Go's.
 * (https://github.com/golang/go/blob/master/src/go/scanner/scanner.go)
 */
package main

import (
	"unicode"
	"unicode/utf8"
)

const (
	eof = -1
)

var identifiers = map[string]tokenKind{
	"and": tokenAndAnd,
	"or":  tokenOrOr,

	"lambda": tokenLambda,
	"let":    tokenLet,
	"in":     tokenIn,
	"match":  tokenMatch,
	"with":   tokenWith,
	"rec":    tokenRec,
	"pi":     tokenPi,
	"true":   tokenBool,
	"false":  tokenBool,

	"bool":  tokenTBool,
	"int":   tokenTInt,
	"float": tokenTFloat,

	// NOTE: we could have used an integer 1 and better
	// categorize it during parsing, but this is just simpler.
	"unit": tokenTUnit,

	// Untested. We're also missing all our gates:
	//	H (Hadamard) N (not) Vtheta (phase shift)
	//	X (exchange) N_C (controlled not)
	//	two more Pauli besides (not)?
	//	"new"    : tokenNew,
	//	"meas"   : tokenMeas,
}

type token struct {
	kind   tokenKind
	ln, cn uint // line/column numbers (1-based)
	raw    string
}

type scanner struct {
	src []byte // input file, fully loaded

	fn     string // filename
	ln, cn uint   // line/column numbers (1-based; columns counted in rune)

	ch      rune // current character/rune; set to eof when done
	offset  int  // ch's offset
	nextOff int  // offset + "len(ch)"
}

func (s *scanner) init(src []byte, fn string) {
	s.src = src
	s.fn = fn
	s.ln = 1
	s.cn = 0
	s.ch = ' '
	s.offset = 0
	s.nextOff = 0

	// load first rune
	s.next()
}

// grab next rune
//
// XXX/TODO: the cn/ln handling is really clumsy.
func (s *scanner) next() {

	// we still have something to read
	if s.nextOff < len(s.src) {
		s.offset = s.nextOff

		// assume common case: ascii character
		r, w := rune(s.src[s.offset]), 1

		if r >= utf8.RuneSelf {
			r, w = utf8.DecodeRune(s.src[s.offset:])
		}

		if r == utf8.RuneError && w == 1 {
			panic("TODO")
		}

		s.cn++
		if r == '\n' {
			s.ln++
			s.cn = 0
		}

		s.ch = r
		s.nextOff += w

	} else if s.ch != eof {
		s.offset = len(s.src)
		s.ch = eof
		if s.cn == 0 {
			s.cn = 1
		} else {
			s.cn++
		}
	}
}

func (s *scanner) peek() byte {
	if s.nextOff < len(s.src) {
		return s.src[s.nextOff]
	}
	return 0
}

// skip whitespaces
func (s *scanner) skipWhites() {
	// NOTE: next() handles reseting ln/cn
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
	}
}

func (s *scanner) switch2(tok0 tokenKind, ch1 rune, tok1 tokenKind) tokenKind {
	if s.ch == ch1 {
		s.next()
		return tok1
	}
	return tok0
}

func (s *scanner) switch3(
	tok0 tokenKind, ch1 rune, tok1 tokenKind, ch2 rune, tok2 tokenKind,
) tokenKind {
	if s.ch == ch1 {
		s.next()
		return tok1
	}
	if s.ch == ch2 {
		s.next()
		return tok2
	}
	return tok0
}

func (s *scanner) switch4(tok0, tok1, tok2, tok3 tokenKind) tokenKind {
	b0, b1 := s.ch, s.peek()

	if b0 == '=' && b1 == '.' {
		s.next()
		s.next()
		return tok3
	}
	if b0 == '=' {
		s.next()
		return tok2
	}
	if b0 == '.' {
		s.next()
		return tok1
	}

	return tok0
}

// returns lower-case ch iff ch is ASCII letter
func lower(ch rune) rune { return ('a' - 'A') | ch }

func isDigit(ch rune) bool { return '0' <= ch && ch <= '9' }

func isLetter(ch rune) bool {
	return ('a' <= lower(ch) && lower(ch) <= 'z') ||
		ch == '_' || (ch >= utf8.RuneSelf && unicode.IsLetter(ch))
}

func (s *scanner) idOrName() tokenKind {
	off := s.offset

	// we know that the first s.ch is a letter ≠ λ already,
	// so we can look for numbers already
	for (isLetter(s.ch) || isDigit(s.ch)) && s.ch != 'λ' {
		s.next()
	}

	if kind, ok := identifiers[string(s.src[off:s.offset])]; ok {
		return kind
	}
	return tokenName
}

func (s *scanner) skipDigits() {
	for isDigit(s.ch) {
		s.next()
	}
}

func (s *scanner) number() tokenKind {
	var kind tokenKind

	if s.ch == '.' {
		s.next()
		kind = tokenFloat
		s.skipDigits()
		return kind
	}

	s.skipDigits()

	if s.ch == '.' {
		s.next()
		kind = tokenFloat
		s.skipDigits()
		return kind
	}

	kind = tokenInt
	return kind
}

// grab next token
func (s *scanner) scan() token {
	s.skipWhites()

	var kind tokenKind

	ln, cn, off := s.ln, s.cn, s.offset

	switch ch := s.ch; {

	case isLetter(ch) && ch != 'λ':
		kind = s.idOrName()
	case isDigit(ch) || (ch == '.' && isDigit(rune(s.peek()))):
		kind = s.number()

	case ch == eof:
		kind = tokenEOF

	default:
		s.next()

		switch ch {
		case 'λ':
			kind = tokenLambda
		case '(':
			kind = tokenLParen
		case ')':
			kind = tokenRParen
		case '.':
			// floats (e.g. ".3") managed by outer switch
			kind = tokenDot

		case '!':
			kind = tokenExcl

		case '+':
			kind = s.switch2(tokenPlus, '.', tokenFPlus)
		case '-':
			kind = s.switch3(tokenMinus, '.', tokenFMinus, '>', tokenArrow)
		case '*':
			kind = s.switch2(tokenStar, '.', tokenFStar)
		case '/':
			kind = s.switch2(tokenSlash, '.', tokenFSlash)

		// TODO: make sure all those are tested
		case '<':
			kind = s.switch4(
				tokenLess,    // <
				tokenFLess,   // <.
				tokenLessEq,  // <=
				tokenFLessEq, // <=.
			)
		case '>':
			kind = s.switch4(
				tokenMore,    // >
				tokenFMore,   // >.
				tokenMoreEq,  // >=
				tokenFMoreEq, // >=.
			)

		case ',':
			kind = tokenComa
		case '=':
			kind = tokenEqual

		case '〈':
			kind = tokenLBracket
		case '〉':
			kind = tokenRBracket

		case '|':
			kind = s.switch2(tokenOr, '|', tokenOrOr)
		case '&':
			kind = s.switch2(tokenAnd, '&', tokenAndAnd)

		case '≤':
			kind = s.switch2(tokenLessEq, '.', tokenFLessEq)
		case '≥':
			kind = s.switch2(tokenMoreEq, '.', tokenFMoreEq)

		case ':':
			kind = tokenColon

		case 'π':
			kind = tokenPi

		case '→':
			kind = tokenArrow

		case '×':
			kind = tokenProduct

		case eof:
			kind = tokenEOF

		// case '⊸': tokenRMultiMap
		// case '⊗': tokenOMult
		// case '⊕': tokenOPlus
		// case '⊤': tokenTrue

		default:
			panic("assert TODO")
		}
	}

	return token{kind, ln, cn, string(s.src[off:s.offset])}
}

// grab next token
func (s *scanner) scanAll() ([]token, error) {
	var toks []token

	for {
		tok := s.scan()
		toks = append(toks, tok)
		if tok.kind == tokenEOF {
			return toks, nil
		}
	}
}

// slurp all tokens
func scanAll(src string, fn string) ([]token, error) {
	var s scanner
	s.init([]byte(src), fn)
	return s.scanAll()
}
