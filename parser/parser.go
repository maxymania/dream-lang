/*
 Copyright (R) 2015 Simon Schmidt
 All rights reserved. NO LICENSE!
 */

package parser

import "text/scanner"
import "github.com/maxymania/dream-lang/lexer"
import "github.com/maxymania/dream-lang/tree"
import "fmt"

type FuncVarSpace struct{
	Base tree.VarSpace
	C int
	N map[string]string
}
func NewFuncVarSpace(v tree.VarSpace) *FuncVarSpace {
	return &FuncVarSpace{
		v,
		0,
		make(map[string]string),
	}
}
func(f *FuncVarSpace) DefineScalar(s string){
	if _,ok := f.N[s]; ok { return }
	num := f.C; f.C++
	f.N[s] = fmt.Sprintf("var%d",num)
}
func(f *FuncVarSpace) Scalar(s string) tree.Expression {
	n,ok := f.N[s]
	if ok {
		return &tree.ScalarVar{Name:n}
	}
	return f.Base.Scalar(s)
}

type FuncTempReg struct{
	O int
	N map[int]string
}
func NewFuncTempReg() *FuncTempReg {
	return &FuncTempReg{0,make(map[int]string)}
}
func (f *FuncTempReg) Reg(n int) string {
	if d,ok := f.N[n]; ok { return d }
	s := fmt.Sprintf("reg%d",n)
	f.N[n] = s
	return s
}
func (f *FuncTempReg) Add(n int) {
	f.O += n
}


type ErrMsg struct{
	P scanner.Position
	E string
}

type ExprParser struct{
	Pos []ErrMsg
	Vsp tree.VarSpace
	Tmp tree.TempReg
}
func (e *ExprParser) err(em ErrMsg) {
	e.Pos = append(e.Pos,em)
}
func (e *ExprParser) backup() int{
	return len(e.Pos)
}
func (e *ExprParser) restore(i int){
	e.Pos = e.Pos[:i]
}

func MatchT(e *lexer.Element,tks ...rune) *lexer.Element{
	for _,t := range tks {
		if e.Val.Type != t { return nil }
		e = e.Next()
	}
	return e
}

func MatchK(e *lexer.Element,kwr ...rune) *lexer.Element{
	for _,k := range kwr {
		if e.Val.Kwrd != k { return nil }
		e = e.Next()
	}
	return e
}

