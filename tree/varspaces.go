/*
 Copyright (R) 2015 Simon Schmidt
 All rights reserved. NO LICENSE!
 */

package tree

type TempReg interface{
	Reg(n int) string
	Add(n int)
}

type VarSpace interface{
	DefineScalar(s string)
	Scalar(s string) Expression
}

type VarSpaceImpl int
func (v VarSpaceImpl) DefineScalar(s string) {}
func (v VarSpaceImpl) Scalar(s string) Expression { return nil }

const EmptyVS = VarSpaceImpl(0)
