/*
 Copyright (R) 2015 Simon Schmidt
 All rights reserved. NO LICENSE!
 */

package parser

import "github.com/maxymania/dream-lang/lexer"
import "github.com/maxymania/dream-lang/tree"

import "text/scanner"

func (e *ExprParser) Expression(el *lexer.Element) (tree.Expression,*lexer.Element) {
	//panic("")
	return e.layerThreeExpression(el)
}
func (e *ExprParser) layerThreeExpression(el *lexer.Element) (tree.Expression,*lexer.Element) {
	var x,x2 tree.Expression
	x,el = e.layerTwoExpression(el)
	if x==nil { return nil,nil }
	for {
		if el2 := MatchK(el,'.'); el2!=nil {
			x2,el = e.layerTwoExpression(el2)
			if x2==nil { return nil,nil }
			x = &tree.BinOp{Arg1:x,Arg2:x2,Op:".."}
		} else { break }
	}
	return x,el
}
func (e *ExprParser) layerTwoExpression(el *lexer.Element) (tree.Expression,*lexer.Element) {
	var x,x2 tree.Expression
	x,el = e.layerOneExpression(el)
	if x==nil { return nil,nil }
	for {
		if el2 := MatchK(el,'+'); el2!=nil {
			x2,el = e.layerOneExpression(el2)
			if x2==nil { return nil,nil }
			x = &tree.BinOp{Arg1:x,Arg2:x2,Op:"+"}
		} else if el2 := MatchK(el,'-'); el2!=nil {
			x2,el = e.layerOneExpression(el2)
			if x2==nil { return nil,nil }
			x = &tree.BinOp{Arg1:x,Arg2:x2,Op:"-"}
		} else { break }
	}
	return x,el
}
func (e *ExprParser) layerOneExpression(el *lexer.Element) (tree.Expression,*lexer.Element) {
	var x,x2 tree.Expression
	x,el = e.innerExpression(el)
	if x==nil { return nil,nil }
	for {
		if el2 := MatchK(el,'*'); el2!=nil {
			x2,el = e.innerExpression(el2)
			if x2==nil { return nil,nil }
			x = &tree.BinOp{Arg1:x,Arg2:x2,Op:"*"}
		} else if el2 := MatchK(el,'/'); el2!=nil {
			x2,el = e.innerExpression(el2)
			if x2==nil { return nil,nil }
			x = &tree.BinOp{Arg1:x,Arg2:x2,Op:"/"}
		} else if el2 := MatchK(el,'%'); el2!=nil {
			x2,el = e.innerExpression(el2)
			if x2==nil { return nil,nil }
			x = &tree.BinOp{Arg1:x,Arg2:x2,Op:"%"}
		} else { break }
	}
	return x,el
}
func (e *ExprParser) innerExpression(el *lexer.Element) (tree.Expression,*lexer.Element) {
	var x tree.Expression
	switch el.Val.Kwrd {
	case '(':
		el = el.Next()
		x,el = e.Expression(el)
		if el==nil { return nil,nil }
		if el.Val.Kwrd!=')' {
			e.err(ErrMsg{el.Val.Pos,"expected ')', got '"+el.Val.Text+"'"})
			return nil,nil
		}
		el = el.Next()
		return x,el
	case scanner.Int,scanner.Float,scanner.String:
		x = &tree.ScalarLit{Data:el.Val.Text}
		el = el.Next()
		return x,el
	case '$':
		el = el.Next()
		if el.Val.Kwrd!=scanner.Ident {
			e.err(ErrMsg{el.Val.Pos,"expected Ident, got '"+el.Val.Text+"'"})
			return nil,nil
		}
		x = e.Vsp.Scalar(el.Val.Text) //&tree.ScalarVar{Name:el.Val.Text}
		if x==nil {
			e.err(ErrMsg{el.Val.Pos,"var not declared: "+el.Val.Text})
			return nil,nil
		}
		el = el.Next()
		return x,el
	}
	e.err(ErrMsg{el.Val.Pos,"not expression"})
	return nil,nil
}


