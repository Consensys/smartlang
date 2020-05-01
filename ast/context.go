package ast

import (
	"smart/types"
)

type Context struct {
	name        string
	symtab      *types.SymTab
	enclosing   *Context
	function    *types.FunctionType
	insideStore bool
}

func (c *Context) NestedContext(name string) *Context {
	var newContext Context
	newContext = *c // copy
	newContext.enclosing = c
	newContext.name = name
	newContext.symtab = types.NewSymTab(name, c.symtab)
	return &newContext
}

func GlobalContext(name string, symtab *types.SymTab) *Context {
	return &Context{name: name, symtab: symtab}
}
