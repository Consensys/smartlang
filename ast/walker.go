package ast

func (p *Program) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(p)
	for _, d := range p.decls {
		d.visit(enter, exit)
	}
	exit(p)
}

func (c *Contract) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(c)
	for _, d := range c.decls {
		d.visit(enter, exit)
	}
	exit(c)
}

func (f *FuncSig) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(f)
	for _, p := range f.params {
		p.visit(enter, exit)
	}
	exit(f)
}

func (f *FallbackDecl) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(f)
	f.codeblock.visit(enter, exit)
	exit(f)
}

func (c *Constructor) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(c)
	for _, p := range c.params {
		p.visit(enter, exit)
	}
	c.code.visit(enter, exit)
	exit(c)
}

func (p *PoolDecl) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(p)
	exit(p)
}

func (p *Param) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(p)
	p.typeAST.visit(enter, exit)
	exit(p)
}

func (s *StorageDecl) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(s)
	s.typeAST.visit(enter, exit)
	exit(s)
}

func (c *ConstDecl) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(c)
	c.typeAST.visit(enter, exit)
	c.value.visit(enter, exit)
	exit(c)
}

func (f *FuncDef) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(f)
	f.sig.visit(enter, exit)
	f.body.visit(enter, exit)
	exit(f)
}

func (c *CodeBlock) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(c)
	for _, s := range c.statements {
		s.visit(enter, exit)
	}
	exit(c)
}

func (e *EventDecl) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(e)
	for _, p := range e.params {
		p.visit(enter, exit)
	}
	exit(e)
}

func (c *Class) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(c)
	for _, d := range c.decls {
		d.visit(enter, exit)
	}
	exit(c)
}

func (i *ImportDecl) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(i)
	i.prog.visit(enter, exit)
	exit(i)
}

func (i *Interface) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(i)
	exit(i)
}

func (v *VarDecl) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(v)
	v.typeAST.visit(enter, exit)
	if v.value != nil {
		v.value.visit(enter, exit)
	}
	exit(v)
}

func (a *AssertStmt) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(a)
	a.cond.visit(enter, exit)
	exit(a)
}

func (e *ExternalCall) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(e)
	e.dest.visit(enter, exit)
	if e.data != nil {
		e.data.visit(enter, exit)
	}
	if e.gas != nil {
		e.gas.visit(enter, exit)
	}
	if e.value != nil {
		e.value.visit(enter, exit)
	}
	exit(e)
}

func (e *EmitStmt) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(e)
	for _, a := range e.args {
		a.visit(enter, exit)
	}
	exit(e)
}

func (i *IfStmt) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(i)
	i.condition.visit(enter, exit)
	i.thenpart.visit(enter, exit)
	i.elsepart.visit(enter, exit)
	exit(i)
}

func (r *ReturnStmt) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(r)
	for _, a := range r.args {
		a.visit(enter, exit)
	}
	exit(r)
}

func (a *Assignment) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(a)
	for _, d := range a.lhs {
		d.visit(enter, exit)
	}
	for _, e := range a.rhs {
		e.visit(enter, exit)
	}
	if a.from != nil {
		a.from.visit(enter, exit)
	}
	exit(a)
}

func (m *Move) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(m)
	m.lhs.visit(enter, exit)
	if m.amount != nil {
		m.amount.visit(enter, exit)
	}
	m.rhs.visit(enter, exit)
	exit(m)
}

func (f *FunctionCall) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(f)
	f.target.visit(enter, exit)
	for _, a := range f.arguments {
		a.visit(enter, exit)
	}
	exit(f)
}

func (a *AugmentedAssignment) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(a)
	a.lhs.visit(enter, exit)
	a.rhs.visit(enter, exit)
	exit(a)
}

func (f *ForStmt) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(f)
	f.init.visit(enter, exit)
	f.cond.visit(enter, exit)
	f.loop.visit(enter, exit)
	f.body.visit(enter, exit)
	exit(f)
}

func (w *WhileStmt) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(w)
	w.cond.visit(enter, exit)
	w.body.visit(enter, exit)
	exit(w)
}

func (q *QualIdent) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(q)
	exit(q)
}

func (a *ArrayType) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(a)
	if a.length != nil {
		a.length.visit(enter, exit)
	}
	a.typeAST.visit(enter, exit)
	exit(a)
}

func (m *MapType) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(m)
	m.keyType.visit(enter, exit)
	m.valType.visit(enter, exit)
	exit(m)
}

func (b *Binary) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(b)
	b.left.visit(enter, exit)
	b.right.visit(enter, exit)
	exit(b)
}

func (a *AsExpr) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(a)
	a.expr.visit(enter, exit)
	a.typeAST.visit(enter, exit)
	exit(a)
}

func (c *CallExpr) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(c)
	c.fn.visit(enter, exit)
	for _, a := range c.args {
		a.visit(enter, exit)
	}
	exit(c)
}

func (u *Unary) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(u)
	u.expr.visit(enter, exit)
	exit(u)
}

func (v *Variable) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(v)
	exit(v)
}

func (f *FieldAccess) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(f)
	f.expr.visit(enter, exit)
	exit(f)
}

func (a *ArrayAccess) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(a)
	a.array.visit(enter, exit)
	a.index.visit(enter, exit)
	exit(a)
}

func (n *NewExpr) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(n)
	n.qual.visit(enter, exit)
	for _, a := range n.args {
		a.visit(enter, exit)
	}
	exit(n)
}

func (d *DecLiteral) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(d)
	exit(d)
}

func (h *HexLiteral) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(h)
	exit(h)
}

func (b *BytesLiteral) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(b)
	exit(b)
}

func (b *BoolLiteral) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(b)
	exit(b)
}

func (a *ArrayLiteral) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(a)
	exit(a)
}

func (s *StringLiteral) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(s)
	exit(s)
}

func (a *Argument) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(a)
	a.value.visit(enter, exit)
	exit(a)
}

func (m *Mint) visit(enter func(IVisitor), exit func(IVisitor)) {
	enter(m)
	m.expr.visit(enter, exit)
	m.typeAST.visit(enter, exit)
	exit(m)
}
