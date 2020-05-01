package ast

import (
	"smart/types"
)

func (p *ArrayType) buildType() types.IType {
	element := p.typeAST.buildType()
	return types.NewArrayType(element)
}

func (p *QualIdent) buildType() types.IType {
	typ := p.sym.GetType()
	p.typ = typ
	return typ
}

func (p *MapType) buildType() types.IType {
	k := p.keyType.buildType()
	v := p.valType.buildType()
	return types.NewMapType(k, v)
}
