package types

import (
	"fmt"
)

type SymTab struct {
	name      string
	T         map[string]ISymbol
	enclosing *SymTab
}

func RootSymTab(name string, syms map[string]ISymbol) *SymTab {
	return &SymTab{name: name, T: syms}
}
func NewSymTab(name string, enclosing *SymTab) *SymTab {
	T := make(map[string]ISymbol)
	return &SymTab{name: name, T: T, enclosing: enclosing}
}

func (p *SymTab) Lookup(name string) (sym ISymbol, found bool) {
	for ; p != nil; p = p.enclosing {
		if v, present := p.T[name]; present {
			return v, true
		}
	}
	return nil, false
}

func (p *SymTab) Fetch(name string) (sym ISymbol, found bool) {
	sym, found = p.T[name]
	return sym, found
}

func (p *SymTab) Insert(name string, sym ISymbol) {
	// TODO catch overwrites?
	p.T[name] = sym
}

func (p *SymTab) Dump() {
	if p != nil {
		fmt.Println(p.name, "= map:", p.T)
		p.enclosing.Dump()
	}
}
