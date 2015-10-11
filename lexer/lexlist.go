/*
 Copyright (R) 2015 Simon Schmidt
 All rights reserved. NO LICENSE!
 */

package lexer

import "text/scanner"

const(
	K_IF = iota+128
	K_ELSIF
	K_ELSE
	K_UNLESS
)

type Token struct{
	Type rune
	Kwrd rune
	Text string
	Pos scanner.Position
}
func (t *Token) scan(s *scanner.Scanner){
	t.Type = s.Scan()
	t.Text = s.TokenText()
	t.Pos = s.Pos()
	if t.Type == scanner.RawString { t.Type = scanner.String }
	t.Kwrd = t.Type
	if t.Type==scanner.Ident {
		switch t.Text {
		case "if": t.Kwrd = K_IF
		case "elsif": t.Kwrd = K_ELSIF
		case "else": t.Kwrd = K_ELSE
		case "unless": t.Kwrd = K_UNLESS
		}
	}
}

func Scan(s *scanner.Scanner) *Element {
	e := new(Element)
	e.s = s
	e.Val.scan(s)
	return e
}

type Element struct{
	Val  Token
	next *Element
	s *scanner.Scanner
}
func (e *Element) Next() *Element {
	if e.next == nil {
		ee := new(Element)
		ee.s = e.s
		ee.Val.scan(e.s)
		e.next = ee
	}
	return e.next
}

