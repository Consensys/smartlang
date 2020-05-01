package ast

import (
	"fmt"
	"smart/types"
)

func InsertUnique(context *Context, p *Identifier, sym types.ISymbol) {
	_, found := context.symtab.Lookup(p.name)
	if found {
		p.Err(fmt.Sprintf("duplicate name: %s", p.name))
	}
	context.symtab.Insert(p.name, sym)
}

func Check(p *Program) {
	s := types.RootSymTab("<builtins>", types.BuiltinSymbols)
	p.ctx = GlobalContext("<builtins>", s)
	p.pass1of2(p.ctx)
	p.pass2of2()
	p.checkType()
}

func (p *Program) pass1of2(enclosing *Context) {
	p.ctx = enclosing.NestedContext(fmt.Sprintf("Program[]"))
	for _, v := range p.decls {
		v.pass1of2(p.ctx)
	}
}

func (p *Contract) pass1of2(enclosing *Context) {
	p.ctx = enclosing.NestedContext(fmt.Sprintf("Contract[%s]", p.id.name))
	c := types.NewContractType()
	p.sym = types.NewContractSym(p.id.name, c, p.ctx.symtab)
	InsertUnique(enclosing, p.id, p.sym)
	for _, v := range p.decls {
		v.pass1of2(p.ctx)
	}
}

func (p *Class) pass1of2(enclosing *Context) {
	p.ctx = enclosing.NestedContext(fmt.Sprintf("Class[%s]", p.id.name))
	c := types.NewClassType()
	p.sym = types.NewClassSym(p.id.name, c, p.ctx.symtab)
	InsertUnique(enclosing, p.id, p.sym)
	for _, v := range p.decls {
		v.pass1of2(p.ctx)
	}
}

func (p *ImportDecl) pass1of2(enclosing *Context) {
	c := types.NewImportType()

	s := types.RootSymTab("<builtins>", types.BuiltinSymbols)
	p.ctx = GlobalContext("<builtins>", s)
	p.prog.pass1of2(p.ctx)

	p.sym = types.NewImportSym(p.alias.name, p.path, p.hash, c, p.prog.ctx.symtab)
	InsertUnique(enclosing, p.alias, p.sym)
}

func (p *Interface) pass1of2(enclosing *Context) {
	p.ctx = enclosing.NestedContext(fmt.Sprintf("Interface[%s]", p.id.name))
	c := types.NewInterfaceType()
	p.sym = types.NewInterfaceSym(p.id.name, c, p.ctx.symtab)
	InsertUnique(enclosing, p.id, p.sym)
	for _, v := range p.decls {
		v.pass1of2(p.ctx)
	}
}

func (p *InterfaceSig) pass1of2(enclosing *Context) {
	f := types.NewFunctionType()
	p.sym = types.NewFunctionSym(p.sig.id.name, f)
	p.sig.sym = p.sym
	InsertUnique(enclosing, p.sig.id, p.sym)
}

func (p *FuncSig) pass1of2(enclosing *Context) {
	p.ctx = enclosing.NestedContext(fmt.Sprintf("Func[%v]", p))
	// p.id was already added at outer scope.
	// TODO p.params
	// TODO p.ret
	// TODO p.decorators
}

func (p *FallbackDecl) pass1of2(enclosing *Context) {
	p.ctx = enclosing.NestedContext(fmt.Sprintf("Fallback[%v]", p))
	// TODO put something in the enclosing scope?
}

func (p *Constructor) pass1of2(enclosing *Context) {
	p.ctx = enclosing.NestedContext(fmt.Sprintf("Constructor[%v]", p))
	// TODO put something in the enclosing scope?
}

func (p *PoolDecl) pass1of2(enclosing *Context) {
	typ := types.NewPoolType(p.id.name)
	p.sym = types.NewPoolSym(p.id.name, typ)
	InsertUnique(enclosing, p.id, p.sym)
}

func (p *StorageDecl) pass1of2(enclosing *Context) {
	p.sym = types.NewStorageSym(p.id.name)
	InsertUnique(enclosing, p.id, p.sym)
}

func (p *ConstDecl) pass1of2(enclosing *Context) {
	p.sym = types.NewConstSym(p.id.name)
	InsertUnique(enclosing, p.id, p.sym)
}

func (p *FuncDef) pass1of2(enclosing *Context) {
	f := types.NewFunctionType()
	p.sym = types.NewFunctionSym(p.sig.id.name, f)
	p.sig.sym = p.sym
	InsertUnique(enclosing, p.sig.id, p.sym)
	p.sig.pass1of2(enclosing)
}

func (p *EventDecl) pass1of2(enclosing *Context) {
	p.ctx = enclosing.NestedContext(fmt.Sprintf("EventDecl[%v]", p))
	f := types.NewEventType()
	p.sym = types.NewEventSym(p.id.name, f)
	InsertUnique(enclosing, p.id, p.sym)
}
