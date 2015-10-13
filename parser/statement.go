/*
 Copyright (R) 2015 Simon Schmidt
 All rights reserved. NO LICENSE!
 */

package parser

import "github.com/maxymania/dream-lang/lexer"
import "github.com/maxymania/dream-lang/tree"
import "io"
import "bytes"
import "fmt"

//import "text/scanner"

func (e *ExprParser) Statement(el *lexer.Element, dest io.Writer) (*lexer.Element){
	//buf := new(bytes.Buffer)
	switch el.Val.Kwrd {
	case '{',lexer.K_IF,lexer.K_UNLESS:
		return e.multiStatement(el,dest)
	}
	return e.singleStatement(el,dest)
}

func (e *ExprParser) multiStatement(el *lexer.Element, dest io.Writer) (*lexer.Element){
	if el2 := MatchK(el,'{'); el2!=nil {
		el = el2
		for{
			el2 = MatchK(el,'}')
			if el2!=nil { return el2 }
			el2 = e.Statement(el,dest)
			if el2==nil {
				e.err(ErrMsg{el.Val.Pos,"expected '}', got '"+el.Val.Text+"'"})
				return nil
			}
			el = el2
		}
	}
	if el2 := MatchK(el,lexer.K_IF,'('); el2!=nil {
		el = el2
		//buf := new(bytes.Buffer)
		number := 1
		x,el2 := e.Expression(el)
		if el2==nil { return nil }
		el = el2
		if MatchK(el,')','{')==nil {
			e.err(ErrMsg{el.Val.Pos,"expected ')' '{', got '"+el.Val.Text+"'"})
			return nil
		}
		el = el.Next()
		r := x.Load(dest)
		fmt.Fprintf(dest,"if %s then\n",r)
		el = e.multiStatement(el,dest)
		if el==nil { return nil }
		for {
			el2 = MatchK(el,lexer.K_ELSIF,'(')
			if el2==nil { break }
			x,el2 := e.Expression(el2)
			if el2==nil { return nil }
			el = el2
			if MatchK(el,')','{')==nil {
				e.err(ErrMsg{el.Val.Pos,"expected ')' '{', got '"+el.Val.Text+"'"})
				return nil
			}
			buf := new(bytes.Buffer)
			r := x.Load(buf)
			if buf.Len()==0 {
				fmt.Fprintf(dest,"elseif %s then\n",r)
			}else{
				fmt.Fprintf(dest,"else\n")
				buf.WriteTo(dest)
				fmt.Fprintf(dest,"if %s then\n",r)
				number++
			}
			el = e.multiStatement(el.Next(),dest)
		}
		el2 = MatchK(el,lexer.K_ELSE,'{')
		if el2!=nil {
			fmt.Fprintln(dest,"else")
			el = e.multiStatement(el.Next(),dest)
			if el==nil { return nil }
		}
		for i := 0 ; i<number ; i++ {
			fmt.Fprintln(dest,"end")
		}
		return el
	}
	if el2 := MatchK(el,lexer.K_UNLESS,'('); el2!=nil {
		el = el2
		x,el2 := e.Expression(el)
		if el2==nil { return nil }
		el = el2
		if MatchK(el,')','{')==nil {
			e.err(ErrMsg{el.Val.Pos,"expected '{', got '"+el.Val.Text+"'"})
			return nil
		}
		el = el.Next()
		r := x.Load(dest)
		fmt.Fprintf(dest,"if not %s then\n",r)
		el = e.multiStatement(el,dest)
		if el==nil { return nil }
		el2 = MatchK(el,lexer.K_ELSE,'{')
		if el2!=nil {
			fmt.Fprintln(dest,"else")
			el = e.multiStatement(el.Next(),dest)
			if el==nil { return nil }
		}
		fmt.Fprintln(dest,"end")
		return el
	}
	return nil
}

func (e *ExprParser) singleStatement(el *lexer.Element, dest io.Writer) (*lexer.Element){
	buf := new(bytes.Buffer)
	el2 := e.subStatement(el,buf)
	if el2==nil { return nil }
	if el3 := MatchK(el2,';'); el3!=nil {
		buf.WriteTo(dest)
		return el3
	}
	if el3 := MatchK(el2,lexer.K_IF); el3!=nil {
		x,el4 := e.Expression(el3)
		el5 := MatchK(el4,';')
		if el5==nil {
			e.err(ErrMsg{el.Val.Pos,"expected ';', got '"+el.Val.Text+"'"})
			return nil
		}
		r := x.Load(dest)
		fmt.Fprintf(dest,"if %s then\n",r)
		buf.WriteTo(dest)
		fmt.Fprintln(dest,"end")
		return el5
	}
	if el3 := MatchK(el2,lexer.K_UNLESS); el3!=nil {
		x,el4 := e.Expression(el3)
		el5 := MatchK(el4,';')
		if el5==nil {
			e.err(ErrMsg{el.Val.Pos,"expected ';', got '"+el.Val.Text+"'"})
			return nil
		}
		r := x.Load(dest)
		fmt.Fprintf(dest,"if not %s then\n",r)
		buf.WriteTo(dest)
		fmt.Fprintln(dest,"end")
		return el5
	}
	e.err(ErrMsg{el.Val.Pos,"expected ';', got '"+el.Val.Text+"'"})
	return nil
}

func (e *ExprParser) subStatement(el *lexer.Element, dest io.Writer) (*lexer.Element){
	x,el2 := e.Expression(el)
	if el2==nil { return nil }
	r := x.Load(dest)
	if _,ok := x.(*tree.Assign); !ok {
		fmt.Fprintln(dest,r)
	}
	return el2;
}

