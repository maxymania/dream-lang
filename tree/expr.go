/*
 Copyright (R) 2015 Simon Schmidt
 All rights reserved. NO LICENSE!
 */

package tree

import "io"

const (
	E_LOAD = 1
	E_STORE = 2
)

const (
	E_TP_SCALAR = iota
)

type Expression interface{
	Cap() int
	Typ() int
	Load(dest io.Writer) (expr string)
	Store(dest io.Writer, expr string) (expr2 string)
	LSBegin(dest io.Writer)
	LSLoad() (expr string)
	LSStore(dest io.Writer, expr string) (expr2 string)
}
type noopExpr struct{}
func (n noopExpr) Cap() int { return 0 }
func (n noopExpr) Typ() int { return E_TP_SCALAR }
func (n noopExpr) Load(dest io.Writer) (expr string) { panic("Not Readable") }
func (n noopExpr) Store(dest io.Writer, expr string) (expr2 string) { panic("Not Writable") }
func (n noopExpr) LSBegin(dest io.Writer) { panic("Neighter Readable Nor Writable") }
func (n noopExpr) LSLoad() (expr string) { panic("Neighter Readable Nor Writable") }
func (n noopExpr) LSStore(dest io.Writer, expr string) (expr2 string) { panic("Neighter Readable Nor Writable") }

func CheckExpr(t int,e ...Expression) {
	for _,ee := range e {
		if t!=ee.Typ() {
			panic("Error Type Mismatch")
		}
	}
}
func CheckExprScalar(e ...Expression) {
	for _,ee := range e {
		t := ee.Typ()
		if t!=E_TP_SCALAR {
			panic("Error Type Mismatch")
		}
	}
}
