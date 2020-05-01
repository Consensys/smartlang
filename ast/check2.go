package ast

import (
	"fmt"
	"smart/types"
)

func (p *Program) pass2of2() {
	for _, v := range p.decls {
		v.pass2of2()
	}
}

func (p *Contract) pass2of2() {
	for _, v := range p.decls {
		v.pass2of2(p.ctx)
	}
}

func (p *Class) pass2of2() {
	for _, v := range p.decls {
		v.pass2of2(p.ctx)
	}
}

func (p *ImportDecl) pass2of2() {
	p.prog.pass2of2()
}

func (p *Interface) pass2of2() {
	for _, v := range p.decls {
		v.pass2of2(p.ctx)
	}
}

func (p *FallbackDecl) pass2of2(enclosing *Context) {
	p.codeblock.pass2of2(enclosing)
}

func returnsToTypes(ins []*ReturnType, enclosing *Context) []*types.ReturnSym {
	ts := make([]*types.ReturnSym, 0)
	for _, v := range ins {
		var name string
		if v.id != nil {
			name = v.id.name
		}
		v.typeAST.pass2of2(enclosing)
		typ := v.typeAST.buildType()
		sym := types.NewReturnSym(name, typ)
		if v.id != nil {
			InsertUnique(enclosing, v.id, sym)
		}
		ts = append(ts, sym)
	}
	return ts
}

func paramsToTypes(ins []*Param, enclosing *Context) []*types.ParamSym {
	ts := make([]*types.ParamSym, 0)
	for _, v := range ins {
		name := v.id.name
		v.typeAST.pass2of2(enclosing)
		typ := v.typeAST.buildType()
		sym := types.NewParamSym(name, typ)
		v.sym = sym
		InsertUnique(enclosing, v.id, sym)
		ts = append(ts, sym)
	}
	return ts
}

func (p *Constructor) pass2of2(enclosing *Context) {
	pts := paramsToTypes(p.params, p.ctx)
	c := types.NewFunctionType() // TODO add this somewhere?
	c.Params = pts
	p.code.pass2of2(p.ctx)
}

func (p *InterfaceSig) pass2of2(enclosing *Context) {
	p.sig.pass2of2(enclosing)
}

func (p *FuncSig) pass2of2(enclosing *Context) {
	if p == nil {
		return // for now...
	}
	pts := paramsToTypes(p.params, p.ctx)
	rts := returnsToTypes(p.ret, p.ctx)
	p.sym.Typ.Params = pts
	p.sym.Typ.Returns = rts
}

func (p *EventDecl) pass2of2(enclosing *Context) {
	pts := paramsToTypes(p.params, p.ctx)
	// f := types.NewEventType()
	// f.Params = pts
	p.sym.Typ.Params = pts
}

func (p *PoolDecl) pass2of2(enclosing *Context) {
	// nothing in pass 2
}

func (p *StorageDecl) pass2of2(enclosing *Context) {
	p.typeAST.pass2of2(enclosing)
	p.sym.Typ = p.typeAST.buildType()
	enclosing.symtab.Insert(p.id.name, p.sym)
}

func (p *ConstDecl) pass2of2(enclosing *Context) {
	if p.sym == nil {
		// must be in code block
		p.sym = types.NewConstSym(p.id.name)
		enclosing.symtab.Insert(p.id.name, p.sym)
	}
	p.typeAST.pass2of2(enclosing)
	p.value.pass2of2(enclosing)
	p.sym.Typ = p.typeAST.buildType()
}

func (p *FuncDef) pass2of2(enclosing *Context) {
	if p.sig == nil {
		return // for now...
	}
	p.sig.pass2of2(enclosing)
	p.sig.ctx.function = p.sig.sym.Typ
	p.body.pass2of2(p.sig.ctx) // ugly, but needed
}

func (p *CodeBlock) pass2of2(enclosing *Context) {
	p.ctx = enclosing.NestedContext(fmt.Sprintf("CodeBlock[%v]", p))
	for _, v := range p.statements {
		v.pass2of2(p.ctx)
	}
}

func (p *VarDecl) pass2of2(enclosing *Context) {
	p.typeAST.pass2of2(enclosing)
	typ := p.typeAST.buildType()
	p.sym = types.NewLocalSym(p.id.name, typ)
	enclosing.symtab.Insert(p.id.name, p.sym)
	if p.value != nil {
		p.value.pass2of2(enclosing)
	}
}

func (p *AssertStmt) pass2of2(enclosing *Context) {
	p.cond.pass2of2(enclosing)
}

func (p *ReturnStmt) pass2of2(enclosing *Context) {
	p.ctx = enclosing
	for _, v := range p.args {
		v.pass2of2(enclosing)
	}
}

func (p *ExternalCall) pass2of2(enclosing *Context) {
	p.dest.pass2of2(enclosing)
	if p.data != nil {
		p.data.pass2of2(enclosing)
	}
	if p.gas != nil {
		p.gas.pass2of2(enclosing)
	}
	if p.value != nil {
		p.value.pass2of2(enclosing)
	}
}

func (p *EmitStmt) pass2of2(enclosing *Context) {
	p.id.pass2of2(enclosing)
	for _, v := range p.args {
		v.pass2of2(enclosing)
	}
}

func (p *IfStmt) pass2of2(enclosing *Context) {
	p.condition.pass2of2(enclosing)
	p.thenpart.pass2of2(enclosing)
	if p.elsepart != nil {
		p.elsepart.pass2of2(enclosing)
	}
}

func (p *ForStmt) pass2of2(enclosing *Context) {
	p.ctx = enclosing.NestedContext(fmt.Sprintf("for[%v]", p))
	if p.init != nil {
		p.init.pass2of2(p.ctx)
	}
	p.cond.pass2of2(p.ctx)
	if p.loop != nil {
		p.loop.pass2of2(p.ctx)
	}
	p.body.pass2of2(p.ctx)
}

func (p *WhileStmt) pass2of2(enclosing *Context) {
	p.cond.pass2of2(enclosing)
	p.body.pass2of2(enclosing)
}

func (p *Move) pass2of2(enclosing *Context) {
	p.lhs.pass2of2(enclosing)
	if p.amount != nil {
		p.amount.pass2of2(enclosing)
	}
	p.rhs.pass2of2(enclosing)
}

func (p *NewExpr) pass2of2(enclosing *Context) {
	p.qual.pass2of2(enclosing)
	for _, v := range p.args {
		v.pass2of2(enclosing)
	}
}

func (p *CallExpr) pass2of2(enclosing *Context) {
	p.ctx = enclosing
	p.fn.pass2of2(enclosing)
	for _, v := range p.args {
		v.pass2of2(enclosing)
	}
}

func (p *FieldAccess) pass2of2(enclosing *Context) {
	p.expr.pass2of2(enclosing)
}

func (p *ArrayAccess) pass2of2(enclosing *Context) {
	p.array.pass2of2(enclosing)
	p.index.pass2of2(enclosing)
}

func (p *AsExpr) pass2of2(enclosing *Context) {
	p.expr.pass2of2(enclosing)
	p.typeAST.pass2of2(enclosing)
}

func (p *BoolLiteral) pass2of2(enclosing *Context) {
	// nothing
}

func (p *ArrayLiteral) pass2of2(enclosing *Context) {
	p.typeAST.pass2of2(enclosing)
	for _, v := range p.values {
		v.pass2of2(enclosing)
	}
}

func (p *BytesLiteral) pass2of2(enclosing *Context) {
	// nothing
}

func (p *HexLiteral) pass2of2(enclosing *Context) {
	// nothing
}

func (p *DecLiteral) pass2of2(enclosing *Context) {
	// nothing
}

func (p *StringLiteral) pass2of2(enclosing *Context) {
	// nothing
}

func (p *Variable) pass2of2(enclosing *Context) {
	// TODO do we need to copy symbol from Ident here?
	p.id.pass2of2(enclosing)
}

func (p *Unary) pass2of2(enclosing *Context) {
	p.expr.pass2of2(enclosing)
}

func (p *Binary) pass2of2(enclosing *Context) {
	p.left.pass2of2(enclosing)
	p.right.pass2of2(enclosing)
}

func (p *Assignment) pass2of2(enclosing *Context) {
	for _, v := range p.lhs {
		v.pass2of2(enclosing)
	}
	if p.from == nil {
		for _, v := range p.rhs {
			v.pass2of2(enclosing)
		}
	} else {
		p.from.pass2of2(enclosing)
	}
}

func (p *FunctionCall) pass2of2(enclosing *Context) {
	p.target.pass2of2(enclosing)
	for _, v := range p.arguments {
		v.pass2of2(enclosing)
	}
}

func (p *AugmentedAssignment) pass2of2(enclosing *Context) {
	p.lhs.pass2of2(enclosing)
	p.rhs.pass2of2(enclosing)
}

func (p *QualIdent) pass2of2(enclosing *Context) {
	p.ctx = enclosing
	var sym types.ISymbol
	if p.qualifier == nil {
		p.id.pass2of2(enclosing)
		sym = p.id.sym
	} else {
		p.id.pass2of2(enclosing)
		// TODO handle not found
		// Does GetField do Lookup?
		fielder, ok := p.id.sym.(types.Fielder)
		if !ok {
			p.Err("type does not support field access")
		}
		s, found := fielder.GetField(p.qualifier.name)
		if !found {
			p.qualifier.Err(fmt.Sprintf("unknown field: %v", p.qualifier.name))
		}
		sym = s
	}
	p.sym = sym
}

func (p *Argument) pass2of2(enclosing *Context) {
	p.ctx = enclosing
	// TODO need to look up p.id
	p.value.pass2of2(enclosing)
}

func (p *Identifier) pass2of2(enclosing *Context) {
	p.ctx = enclosing
	sym, found := enclosing.symtab.Lookup(p.name)
	if !found {
		p.Err("unknown identifier: " + p.name)
	}
	p.sym = sym
}

func (p *MapType) pass2of2(enclosing *Context) {
	p.keyType.pass2of2(enclosing)
	p.valType.pass2of2(enclosing)
}

func (p *ArrayType) pass2of2(enclosing *Context) {
	if p.length != nil {
		p.length.pass2of2(enclosing)
	}
	p.typeAST.pass2of2(enclosing)
}

func (p *Mint) pass2of2(enclosing *Context) {
	p.expr.pass2of2(enclosing)
	p.typeAST.pass2of2(enclosing)
}
