package ast

import (
	asm "smart/assembler"
	"smart/token"
	"smart/types"
)

type Locator interface {
	IsStorage() bool
}

type Location interface {
	Fetch(a *asm.Assembler)
	StoreTo(a *asm.Assembler)
}

type HasLocation interface {
	LocationOf(a *asm.Assembler) Location
}

type LValue interface {
	Locator
	Lvalue()
}

type IDecorator interface {
}

type Storer interface {
	SetStore(bool)
}

type IStatement interface {
	IVisitor
	Positioner
	Err(msg string)
	pass2of2(enclosing *Context)
	checkType()
	codegen(asm *asm.Assembler)
}

type IDesignator interface {
	IExpr
}

type ITypeAST interface {
	IVisitor
	Err(msg string)
	pass2of2(enclosing *Context)
	buildType() types.IType
}

type IExpr interface {
	IVisitor
	Positioner
	Err(msg string)
	pass2of2(enclosing *Context)
	// getType() types.IType
	checkType() types.IType
	rvalue(asm *asm.Assembler)
}

type IClassBodyDecl interface {
	IVisitor
	pass1of2(enclosing *Context)
	pass2of2(enclosing *Context)
	checkType()
}

type IDecl interface {
	IVisitor
	pass1of2(enclosing *Context)
	pass2of2(enclosing *Context)
	checkType()
	codegen(asm *asm.Assembler)
}

type IProgramBodyDecl interface {
	IVisitor
	pass1of2(enclosing *Context)
	pass2of2()
	checkType()
	codegen(asm *asm.Assembler)
}

type IVisitor interface {
	visit(enter func(IVisitor), exit func(IVisitor))
}

type IDumper interface {
	DumpSymbols() string
}

type Positioner interface {
	SetCoords(pos token.Position)
}
