package main

/*
 * TODO: comments support
 *
 * TODO: we're duplicating some effort in parsing scalars
 *	(integers, floats): we identify them here, but actually
 *	parse them in the parser. Perhaps we could parse them
 *	here already, which would imply wasting a few token{} bytes.
 *
 *	For now set aside to focus on more interesting things.
 *
 * TODO: we're only partially supporting extended tokens in
 *	preparation for qlambda.
 *
 * NOTE: the parser is purposefully sophisticated in that it
 *	relies on a bufio.Scanner instead of e.g. assuming the
 *	the input files fits in a []byte (which should be the
 *	expected use-case).
 *
 *	This means in particular that we need to handle cases
 *	like  a read yielding incomplete runes.
 */

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

type token struct {
	kind   tokenKind
	ln, cn uint // line/column numbers (1-based)
	raw    string
}

type scanner struct {
	scan *bufio.Scanner

	fn     string // filename
	ln, cn uint   // line/column numbers (1-based)

	// last emitted token's width (runes).
	// the difference with utf8.RuneCountInString(tok.raw)
	// is that, in case we emit a token, then go through
	// the input buffer only to have to request more data
	// to conclude, tw is reset.
	tw uint

	tok token // last token parsed
}

// single-rune tokens
var ones = map[rune]tokenKind{
	'(': tokenLParen,
	')': tokenRParen,
	'.': tokenDot,
	':': tokenColon,
	'+': tokenPlus,
	'-': tokenMinus,
	'*': tokenStar,
	'/': tokenSlash,
	'<': tokenLess,
	'>': tokenMore,
	',': tokenComa,
	'=': tokenEqual,
	'〈': tokenLBracket,
	'〉': tokenRBracket,
	// NOTE: tokenOr and tokenAnd were added only
	// for 'foo||' to be parsed correctly (see isSep() below)
	// (nothing wrong with implementing them either)
	'|': tokenOr,
	'&': tokenAnd,
	'≤': tokenLessEq,
	'≥': tokenMoreEq,
	'×': tokenProduct, // <P,Q> -> S <=> P×Q → S, as an ascii variant?
	'→': tokenArrow,
	'λ': tokenLambda,
	'π': tokenPi,
	// '⊸'  : tokenRMultiMap,
	// '⊗'  : tokenOMult,
	// '⊕'  : tokenOPlus,
	// '⊤'  : tokenTrue,
	// '!'  : tokenExclamation,
}

// two ascii characters tokens
var twos = map[string]tokenKind{
	"+.": tokenFPlus,
	"-.": tokenFMinus,
	"*.": tokenFStar,
	"/.": tokenFSlash,
	"<.": tokenFLess,
	">.": tokenFMore,
	"->": tokenArrow,
	"&&": tokenAndAnd,
	"||": tokenOrOr,
}

// special names
var many = map[string]tokenKind{
	// XXX those were in twos, but aren't two bytes long.
	"≤.": tokenFLessEq,
	"≥.": tokenFMoreEq,
	// XXX those were missing (untested thus)
	"<=.": tokenFLessEq,
	">=.": tokenFMoreEq,
	// XXX and/or untested
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

func isSep(r rune) bool {
	if unicode.IsSpace(r) {
		return true
	}

	// single byte token
	// so far, this encompasses all punctuation used
	// (including twos' first byte); unused punctuation
	// symbols are considered part of a name for now
	_, ok := ones[r]
	return ok
}

// TODO no tests for error case
func (s *scanner) scanRune(xs []byte, p int, atEOF bool) (rune, int, bool, error) {
	r, w := utf8.DecodeRune(xs[p:])
	if r == utf8.RuneError {
		// either all input is in xs or we have enough
		// bytes in xs to parse a complete rune (4 bytes
		// per rune at most): that's an encoding error
		if atEOF || p+3 < len(xs) {
			s.tok.kind = tokenError
			return r, 0, false, fmt.Errorf("Encoding issue")
		}

		// ask for extra bytes to try to complete a rune
		return 0, 0, true, nil
	}

	return r, w, false, nil
}

func (s *scanner) setKind(k tokenKind) {
	s.tok.raw = ""
	s.tok.kind = k
}

func (s *scanner) skipWhites(xs []byte) int {
	p := 0
	for w := 0; p < len(xs); p += w {
		var r rune
		r, w = utf8.DecodeRune(xs[p:])
		if !unicode.IsSpace(r) {
			break
		}
		s.cn++
		if r == '\n' {
			s.ln++
			s.cn = 1
		}
	}
	return p
}

// Note that this is a bit redundant with what is later performed
// in 'parser.go:/\) number\(', but this actually follows what
// is done in Go's scanner/parser.
//
// As there are no union in Go, either we'd had to waste some bytes,
// for every token, or use a data structure mimicking what is done
// in the parsing (to "compensate" for the lack of union)
func (s *scanner) scanNumber(xs []byte, p int, atEOF bool) (int, []byte, error) {
	isFloat := false

	s.setKind(tokenInt)

	q := p
	for ; q < len(xs); q++ {
		if !unicode.IsDigit(rune(xs[q])) {
			break
		}
		s.tw++
	}

	if q < len(xs) && xs[q] == '.' {
		isFloat = true
		for q++; q < len(xs); q++ {
			if xs[q] < '0' || xs[q] > '9' {
				break
			}
			s.tw++
		}
	}

	// ask for more
	if q == len(xs) && !atEOF {
		s.tw = 0
		return p, nil, nil
	}

	if isFloat {
		s.setKind(tokenFloat)
	}

	return q, xs[p:q], nil
}

func (s *scanner) scanNameOrId(xs []byte, p, q int) (int, []byte, error) {
	s.setKind(tokenName)
	y := xs[p:q]
	if k, ok := many[string(y)]; ok {
		s.setKind(k)
	}
	return q, y, nil
}

func (s *scanner) init(in io.Reader, fn string) {
	s.scan = bufio.NewScanner(in)
	s.fn = fn
	s.ln = 1
	s.cn = 1
	s.tw = 0

	// fragile.
	//	s.tok.kind = tokenError

	s.scan.Split(func(xs []byte, atEOF bool) (int, []byte, error) {
		s.cn += uint(s.tw)
		s.tw = 0

		p := s.skipWhites(xs)

		if p == len(xs) && atEOF {
			s.setKind(tokenEOF)
			return p, nil, bufio.ErrFinalToken
		}

		// we want at least 2 bytes here so we can
		// decide whether we have a 2 ascii characters long token
		// (see twos)
		if p >= len(xs)-2 && !atEOF {
			// ask for more
			return p, nil, nil
		}

		r, w, more, err := s.scanRune(xs, p, atEOF)
		if more {
			return p, nil, nil
		} else if err != nil {
			return p, nil, err
		}

		// 2 ascii characters long tokens
		if w == 1 && p+w < len(xs) {
			if k, ok := twos[string(xs[p:p+2])]; ok {
				s.setKind(k)
				s.tw += 2
				return p + 2, xs[p : p+2], nil
			}
		}

		// special case
		if r == '.' && p+w < len(xs) && unicode.IsDigit(rune(xs[p+w])) {
			return s.scanNumber(xs, p, atEOF)
		}

		// single rune tokens
		if k, ok := ones[r]; ok {
			s.setKind(k)
			s.tw++
			return p + w, xs[p : p+w], nil
		}

		// .<digit> already managed earlier
		if unicode.IsDigit(r) {
			return s.scanNumber(xs, p, atEOF)
		}

		// (try to) read a name
		q := p
		for w := 0; q < len(xs); q += w {
			var r rune
			r, w, more, err = s.scanRune(xs, q, atEOF)
			if more {
				s.tw = 0
				return p, nil, nil
			} else if err != nil {
				return q, nil, err
			}

			if isSep(r) {
				return s.scanNameOrId(xs, p, q)
			}
			s.tw++
		}

		// This was the last word
		if atEOF {
			return s.scanNameOrId(xs, p, q)
		}

		// ask for more; start after spaces next time
		s.tw = 0
		return p, nil, nil
	})
}

func (s *scanner) next() bool {
	r := s.scan.Scan()
	// Not a big fan of this but:
	//	s.scan.Err() returns nil when we're at EOF.
	//
	//	which means we'll return a EOF token here
	//	AND next() will return false. Which means
	//	caller (see e.g. scanAll below) will have to
	//	call next() one more time after having received
	//	false to catch the EOF, which is bug prone.
	//
	//	The other option would involve keeping a "eof"
	//	flag somewhere.; maybe there's a more idiomatic
	//	way of handling the EOF token, so I'm leaving
	//	this is-is for now.
	if r || s.scan.Err() == nil {
		s.tok.ln = s.ln
		s.tok.cn = s.cn
		s.tok.raw = s.scan.Text()
	}
	return r
}

// to ease tests
func scanAll(in io.Reader, fn string) ([]token, error) {
	var s scanner
	s.init(in, fn)
	ts := []token{}
	for s.next() {
		ts = append(ts, s.tok)
	}
	ts = append(ts, s.tok)
	//	println("Error: ", s.scan.Err())
	return ts, s.scan.Err()
}
