/*
 Copyright (R) 2015 Simon Schmidt
 All rights reserved. NO LICENSE!
 */

package tree

import "io"
import "fmt"

type ScalarLit struct{
	noopExpr
	Data string
}
func (s *ScalarLit) Cap() int { return E_LOAD }
func (s *ScalarLit) Load(dest io.Writer) (expr string) { return s.Data }

type ScalarVar struct{
	noopExpr
	Name string
}
func (s *ScalarVar) Cap() int { return E_LOAD|E_STORE }
func (s *ScalarVar) Load(dest io.Writer) (expr string) { return s.Name }
func (s *ScalarVar) Store(dest io.Writer, expr string) (expr2 string) {
	fmt.Fprintf(dest,"%s = %s\n",s.Name , expr)
	return s.Name
}
func (s *ScalarVar) LSBegin(dest io.Writer) {}
func (s *ScalarVar) LSLoad() (expr string) { return s.Name }
func (s *ScalarVar) LSStore(dest io.Writer, expr string) (expr2 string) {
	fmt.Fprintf(dest,"%s = %s\n",s.Name , expr)
	return s.Name
}

type BinOp struct{
	noopExpr
	Arg1,Arg2 Expression
	Op  string
}
func (b *BinOp) Cap() int { return E_LOAD }
func (b *BinOp) Load(dest io.Writer) (expr string) {
	CheckExprScalar(b.Arg1,b.Arg2)
	a1 := b.Arg1.Load(dest)
	a2 := b.Arg2.Load(dest)
	return "("+a1+" "+b.Op+" "+a2+")"
}

type ArrayLkup struct {
	noopExpr
	Arg1,Arg2 Expression
	Tptr TempReg
}
func (b *ArrayLkup) Cap() int { return E_LOAD|E_STORE }
func (b *ArrayLkup) Load(dest io.Writer) (expr string) {
	CheckExprScalar(b.Arg1,b.Arg2)
	a1 := b.Arg1.Load(dest)
	a2 := b.Arg2.Load(dest)
	return a1+"["+a2+"]"
}
func (b *ArrayLkup) Store(dest io.Writer, expr string) (expr2 string) {
	CheckExprScalar(b.Arg1,b.Arg2)
	a1 := b.Arg1.Load(dest)
	a2 := b.Arg2.Load(dest)
	fmt.Fprintf(dest,"%s[%s] = %s\n",a1,a2,expr)
	return a1+"["+a2+"]"
}
func (b *ArrayLkup) LSBegin(dest io.Writer) {
	CheckExprScalar(b.Arg1,b.Arg2)
	a1 := b.Tptr.Reg(0)
	a2 := b.Tptr.Reg(1)
	b.Tptr.Add(2)
	defer b.Tptr.Add(-2)
	fmt.Fprintf(dest,"%s = %s\n",a1,b.Arg1.Load(dest))
	fmt.Fprintf(dest,"%s = %s\n",a2,b.Arg2.Load(dest))
}
func (b *ArrayLkup) LSLoad() (expr string) {
	a1 := b.Tptr.Reg(0)
	a2 := b.Tptr.Reg(1)
	return a1+"["+a2+"]"
}
func (b *ArrayLkup) LSStore(dest io.Writer, expr string) (expr2 string) {
	a1 := b.Tptr.Reg(0)
	a2 := b.Tptr.Reg(1)
	fmt.Fprintf(dest,"%s[%s] = %s\n",a1,a2,expr)
	return a1+"["+a2+"]"
}

