package ast

import (
	"fmt"
	"math/big"
	asm "smart/assembler"
	"smart/token"
	"smart/types"

	"github.com/ethereum/go-ethereum/crypto"
)

var labelCount int = 0

func labelGen(name string) string {
	ret := fmt.Sprintf("%s-%d", name, labelCount)
	labelCount += 1
	return ret
}

var functionLabels map[*types.FunctionSym]string

func CodeGen(p *Program, a *asm.Assembler) {
	functionLabels = make(map[*types.FunctionSym]string)
	p.codegen(a)
}

func (p *Program) codegen(a *asm.Assembler) {
	for _, v := range p.decls {
		v.codegen(a)
	}
}

func (p *ImportDecl) codegen(a *asm.Assembler) {
	// TODO
}

func (p *Interface) codegen(a *asm.Assembler) {
	// TODO
}

func (p *Class) codegen(a *asm.Assembler) {
	// TODO
}

func (p *Contract) codegen(a *asm.Assembler) {
	for _, v := range p.decls {
		v.codegen(a)
	}
}

func (p *FallbackDecl) codegen(a *asm.Assembler) {
	// TODO
}

func (p *Constructor) codegen(a *asm.Assembler) {
	// TODO handle incoming parameters
	// prologue
	a.Comment("initial FP = 64")
	a.Push1(64)
	a.Push1(0)
	a.Emit(asm.MSTORE)
	a.Comment("initial SP = 64")
	a.Push1(64)
	a.Push1(32)
	a.Emit(asm.MSTORE)

	p.code.codegen(a)

	a.Push1(0)
	a.Push1(0)
	a.Emit(asm.RETURN)
}

func (p *PoolDecl) codegen(a *asm.Assembler) {
	// TODO
}

func (p *StorageDecl) codegen(a *asm.Assembler) {
	// TODO
}

func (p *ConstDecl) codegen(a *asm.Assembler) {
	// TODO
}

func (p *FuncDef) codegen(a *asm.Assembler) {
	lab, ok := functionLabels[p.sym]
	if !ok {
		lab = labelGen(p.sig.id.name)
		functionLabels[p.sym] = lab
	}
	a.Label(lab)
	// desired layout:
	//   FP    -> OLDFP
	//            return location
	//            param1
	//            param2
	//            ...
	//   SP    ->

	// fetch old FP
	a.Comment("fetch old FP")
	a.Push1(0)
	a.Emit(asm.MLOAD)

	// store at SP
	a.Comment("store at SP")
	a.Push1(32)
	a.Emit(asm.MLOAD)
	a.Emit(asm.MSTORE)

	// update FP to point to SP
	a.Comment("update FP to point to SP")
	a.Push1(32)
	a.Emit(asm.MLOAD)
	a.Push1(0)
	a.Emit(asm.MSTORE)

	a.Comment("move parameters from EVM stack to memory")
	// stack looks like: p3, p2, p1
	for i := len(p.sig.params) - 1; i >= 0; i -= 1 {
		// TODO: this pointer
		// store relative to frame pointer
		a.Push1(0)
		a.Emit(asm.MLOAD)
		off := (i + 2) * 32 // + 2 is for old FP and return location
		if off >= 256 {
			panic("fix me by using PUSH2 or something else")
		}
		p.sig.params[i].sym.SetOffset(off)
		a.Push1(uint8(off))
		a.Emit(asm.ADD)
		a.Emit(asm.MSTORE)
	}

	// store return location right after old FP (so 32 bytes after FP)
	a.Comment("store return location right after old FP")
	a.Push1(0)
	a.Emit(asm.MLOAD)
	a.Push1(32)
	a.Emit(asm.ADD)
	a.Emit(asm.MSTORE)

	// compute new SP = FP + (2 + nParams)*32 // 2 is for old FP and return location
	a.Comment("move SP to account for old FP, return location, and all parameters")
	a.Push1(0)
	a.Emit(asm.MLOAD)
	off := 32 * (2 + len(p.sig.params))
	if off >= 256 {
		panic("fix me by using PUSH2 or something else")
	}
	a.Push1(uint8(off)) // TODO: this pointer
	a.Emit(asm.ADD)
	a.Push1(32)
	a.Emit(asm.MSTORE)

	a.Comment("function body")
	p.body.codegen(a)
}

func (p *EventDecl) codegen(a *asm.Assembler) {
	// TODO
}

func (p *CodeBlock) codegen(a *asm.Assembler) {
	// TODO prologue
	for _, v := range p.statements {
		v.codegen(a)
	}
	// TODO epilogue
}

func (p *VarDecl) codegen(a *asm.Assembler) {
	// TODO
}

func (p *AssertStmt) codegen(a *asm.Assembler) {
	// TODO check boolean type
	p.cond.rvalue(a)
	lab := labelGen("assert-out")
	a.PushLabel(lab)
	a.Emit(asm.JUMPI)
	a.Emit(asm.REVERT) // TODO need to set up data
	a.Label(lab)
}

func (p *ExternalCall) codegen(a *asm.Assembler) {
	// TODO
}

func (p *EmitStmt) codegen(a *asm.Assembler) {
	// TODO
}

func (p *IfStmt) codegen(a *asm.Assembler) {
	p.condition.rvalue(a)
	lab := labelGen("if-then")
	out := labelGen("if-out")
	a.PushLabel(lab)
	a.Emit(asm.JUMPI)
	p.elsepart.codegen(a)
	a.PushLabel(out)
	a.Emit(asm.JUMP)
	a.Label(lab)
	p.thenpart.codegen(a)
	a.Label(out)
}

func (p *ReturnStmt) codegen(a *asm.Assembler) {
	a.Comment("push return values onto stack")
	for _, arg := range p.args {
		arg.value.rvalue(a)
	}

	// put return location on the EVM stack
	a.Comment("fetch return location from FP + 32")
	a.Push1(0)
	a.Emit(asm.MLOAD)
	a.Push1(32)
	a.Emit(asm.ADD)
	a.Emit(asm.MLOAD)

	// restore old SP
	// - fetch FP
	a.Comment("fetch FP")
	a.Push1(0)
	a.Emit(asm.MLOAD)
	// - store into SP
	a.Comment("store into SP")
	a.Push1(32)
	a.Emit(asm.MSTORE)

	// restore old FP
	// - fetch FP
	a.Comment("fetch FP")
	a.Push1(0)
	a.Emit(asm.MLOAD)
	// - read value there (old FP)
	a.Comment("read value there (old FP)")
	a.Emit(asm.MLOAD)
	// - store into FP
	a.Comment("store into FP")
	a.Push1(0)
	a.Emit(asm.MSTORE)

	// jump to return location
	a.Comment("jump back to return location")
	a.Emit(asm.JUMP)
}

func (p *Assignment) codegen(a *asm.Assembler) {
	// TODO: does not work correctly for multiple assignment
	var locs []Location
	for _, lhs := range p.lhs {
		e, ok := lhs.(HasLocation)
		if !ok {
			panic("TODO: should have been checked earlier?")
		}
		loc := e.LocationOf(a)
		locs = append(locs, loc)
	}
	for _, rhs := range p.rhs {
		rhs.rvalue(a)
	}
	// TODO: this does assignments in the wrong order.  Not sure if there
	// is a trick to reorder them cleverly.
	for _, loc := range locs {
		loc.StoreTo(a)
	}
}

func (p *Move) codegen(a *asm.Assembler) {
	// TODO
}

func (p *FunctionCall) codegen(a *asm.Assembler) {
	// TODO
}

func (p *AugmentedAssignment) codegen(a *asm.Assembler) {
	// TODO
}

func (p *ForStmt) codegen(a *asm.Assembler) {
	// <init>
	// PUSH @loop-cond
	// JUMP
	// loop:
	//   <body>
	//   <loop>
	// loop-condition:
	//   <cond>
	//   PUSH @loop
	//   JUMPI

	loop := labelGen("loop")
	cond := labelGen("loop-condition")

	if p.init != nil {
		p.init.codegen(a)
	}
	a.PushLabel(cond)
	a.Emit(asm.JUMP)
	a.Label(loop)
	p.body.codegen(a)
	if p.loop != nil {
		p.loop.codegen(a)
	}
	a.Label(cond)
	p.cond.rvalue(a)
	a.PushLabel(loop)
	a.Emit(asm.JUMPI)
}

func (p *WhileStmt) codegen(a *asm.Assembler) {
	// PUSH @loop-cond
	// JUMP
	// loop:
	//   <body>
	// loop-condition:
	//   <cond>
	//   PUSH @loop
	//   JUMPI

	loop := labelGen("loop")
	cond := labelGen("loop-condition")

	a.PushLabel(cond)
	a.Emit(asm.JUMP)
	a.Label(loop)
	p.body.codegen(a)
	a.Label(cond)
	p.cond.rvalue(a)
	a.PushLabel(loop)
	a.Emit(asm.JUMPI)
}

func revertIfFalse(a *asm.Assembler) {
	out := labelGen("no-overflow")
	a.PushLabel(out)
	a.Emit(asm.JUMPI)
	a.Push1(0)
	a.Push1(0)
	a.Emit(asm.REVERT)
	a.Label(out)
}

func (p *Binary) rvalue(a *asm.Assembler) {
	// TODO: Reconsider order of evaluation (currently strict left-to-right)
	if p.op == token.LOGOR {
		t := labelGen("or-shortcircuit")
		done := labelGen("or-done")

		p.left.rvalue(a) // stack: a
		a.PushLabel(t)
		a.Emit(asm.JUMPI)

		p.right.rvalue(a) // stack: b
		a.PushLabel(done)
		a.Emit(asm.JUMP)

		a.Label(t)
		a.Push1(1) // stack: 1

		a.Label(done)

		return
	}
	if p.op == token.LOGAND {
		f := labelGen("and-shortcircuit")
		done := labelGen("and-done")

		p.left.rvalue(a)   // stack: a
		a.Emit(asm.ISZERO) // stack: !a
		a.PushLabel(f)
		a.Emit(asm.JUMPI)

		p.right.rvalue(a) // stack: b
		a.PushLabel(done)
		a.Emit(asm.JUMP)

		a.Label(f)
		a.Push1(0)

		a.Label(done)

		return
	}
	p.left.rvalue(a)
	p.right.rvalue(a)
	switch p.op {
	case token.PLUS:
		// stack: b, a
		a.Emit(asm.DUP2)   // stack : a, b, a
		a.Emit(asm.ADD)    // stack : a+b, a
		a.Emit(asm.SWAP1)  // stack: a, a+b
		a.Emit(asm.DUP2)   // stack: a+b, a, a+b
		a.Emit(asm.LT)     // stack: a+b<a, a+b
		a.Emit(asm.ISZERO) // stack: a+b>=a a+b
		revertIfFalse(a)
	case token.MINUS:
		// stack: b, a
		a.Emit(asm.SWAP1)  // stack: a, b
		a.Emit(asm.DUP2)   // stack: b, a, b
		a.Emit(asm.DUP2)   // stack: a, b, a, b
		a.Emit(asm.LT)     // stack: a < b, a, b
		a.Emit(asm.ISZERO) // stack: a >= b
		revertIfFalse(a)
		a.Emit(asm.SUB)
	case token.STAR:
		// stack: b, a
		a.Emit(asm.DUP2)   // stack: a, b, a
		a.Emit(asm.DUP2)   // stack: b, a, b, a
		a.Emit(asm.MUL)    // stack: b*a, b, a
		a.Emit(asm.DUP3)   // stack: a, b*a, b, a
		a.Emit(asm.ISZERO) // stack: a == 0, b*a, b, a
		out := labelGen("no-overflow")
		a.PushLabel(out)
		a.Emit(asm.JUMPI) // safe if a == 0
		// stack: b*a, b, a
		a.Emit(asm.SWAP2) // stack: a, b, b*a
		a.Emit(asm.DUP3)  // stack: b*a, a, b, b*a
		a.Emit(asm.DIV)   // stack: b*a/a, b, b*a
		a.Emit(asm.EQ)    // stack: b*a/a==b, b*a
		a.PushLabel(out)
		a.Emit(asm.JUMPI) // safe if b*a/a==b
		a.Push1(0)
		a.Push1(0)
		a.Emit(asm.REVERT)
		a.Label(out)
	case token.DIV:
		// stack: b, a
		a.Emit(asm.DUP1)   // stack: b, b, a
		a.Emit(asm.ISZERO) // stack: b==0, b, a
		a.Emit(asm.ISZERO) // stack: b!=0, b, a
		revertIfFalse(a)   // division by zero is not allowed
		a.Emit(asm.SWAP1)  // stack: a, b
		a.Emit(asm.DIV)    // stack: a/b
	case token.EXP:
		// TODO: How are we handling overflow in exponentiation?
		//       Do we have to write our own (expensive?) routine?
		a.Emit(asm.SWAP1)
		a.Emit(asm.EXP)
	case token.GEQ:
		a.Emit(asm.SWAP1)
		a.Emit(asm.LT)
		a.Emit(asm.ISZERO)
	case token.LEQ:
		a.Emit(asm.SWAP1)
		a.Emit(asm.GT)
		a.Emit(asm.ISZERO)
	case token.LT:
		a.Emit(asm.SWAP1)
		a.Emit(asm.LT)
	case token.GT:
		a.Emit(asm.SWAP1)
		a.Emit(asm.GT)
	case token.EQEQ:
		a.Emit(asm.EQ)
	case token.NEQ:
		a.Emit(asm.EQ)
		a.Emit(asm.ISZERO)
	case token.BITAND:
		a.Emit(asm.AND)
	case token.BITOR:
		a.Emit(asm.OR)
	case token.LSHIFT:
		a.Emit(asm.SHL)
	case token.RSHIFT:
		a.Emit(asm.SHR)
	default:
		// TODO
		panic("unimplemented")
	}
}

func (p *Unary) rvalue(a *asm.Assembler) {
	switch p.op {
	case token.MINUS:
		panic("unimplemented")
		// For when we have signed integers:
		// p.expr.rvalue(a)
		// a.Push1(0)
		// a.Emit(asm.SUB)
	case token.BANG:
		p.expr.rvalue(a)
		a.Emit(asm.ISZERO)
	default:
		// TODO: token.BURN: subtract from implicit pool sum? Or should we get rid of unary burn?
		// TODO: token.LARR
		panic("unimplemented")
	}
}

func (p *AsExpr) rvalue(a *asm.Assembler) {
	// TODO
}

func (p *CallExpr) rvalue(a *asm.Assembler) {
	ret := labelGen("return")
	a.PushLabel(ret)
	for _, arg := range p.args {
		arg.value.rvalue(a)
	}
	v, ok := p.fn.(*Variable)
	if !ok {
		panic("calling non-variable")
	}
	sym, ok := p.ctx.symtab.Lookup(v.id.name)
	if !ok {
		panic("couldn't find symbol for function call")
	}
	fsym, ok := sym.(*types.FunctionSym)
	if !ok {
		panic("trying to call non-function")
	}
	lab, ok := functionLabels[fsym]
	if !ok {
		lab = labelGen(v.id.name)
		functionLabels[fsym] = lab
	}
	a.PushLabel(lab)
	a.Emit(asm.JUMP)
	a.Label(ret)
	// TODO
}

func (p *NewExpr) rvalue(a *asm.Assembler) {
	// TODO
}

func (p *ArrayLiteral) rvalue(a *asm.Assembler) {
	// TODO
}

func (p *StringLiteral) rvalue(a *asm.Assembler) {
	// TODO
}

func (p *BytesLiteral) rvalue(a *asm.Assembler) {
	// TODO
}

func (p *BoolLiteral) rvalue(a *asm.Assembler) {
	if p.value {
		a.Push1(1)
	} else {
		a.Push1(0)
	}
}

func (p *HexLiteral) rvalue(a *asm.Assembler) {
	// TODO
}

func (p *DecLiteral) rvalue(a *asm.Assembler) {
	// TODO: Maybe we add the bigint to the AST node during semantic checking? These errors should come earlier.
	maxUint256 := big.NewInt(0)
	maxUint256.Exp(big.NewInt(2), big.NewInt(256), nil)

	i := big.NewInt(0)
	_, success := i.SetString(p.lit, 10)
	if !success {
		p.Err("failed to parse decimal literal")
	} else if i.Cmp(maxUint256) >= 0 {
		p.Err("decimal literal too large")
	} else if i.Cmp(big.NewInt(0)) < 0 {
		// TODO
		panic("negative numbers not implemented")
	}
	if i.BitLen() == 0 {
		a.Push1(0)
	} else {
		a.PushBytes(i.Bytes())
	}
}

func (p *ArrayAccess) rvalue(a *asm.Assembler) {
	loc := p.LocationOf(a)
	loc.Fetch(a)
}

func (p *FieldAccess) rvalue(a *asm.Assembler) {
	loc := p.LocationOf(a)
	loc.Fetch(a)
}

func (p *Mint) rvalue(a *asm.Assembler) {
	p.expr.rvalue(a)
}

func (p *Variable) rvalue(a *asm.Assembler) {
	// TODO:  not handling constants or other kinds of names
	loc := p.LocationOf(a)
	loc.Fetch(a)
}

func (p *ArrayAccess) LocationOf(a *asm.Assembler) Location {
	arr, _ := p.array.(HasLocation) // should already have been checked?
	loc := arr.LocationOf(a)
	p.index.rvalue(a)

	// TODO: how are these runtime values combined?  arrays vs. maps
	// I'm just adding them for now.
	a.Emit(asm.ADD)

	return loc
}

func (p *FieldAccess) LocationOf(a *asm.Assembler) Location {
	e, _ := p.expr.(HasLocation) // should already have been checked?
	loc := e.LocationOf(a)
	var off uint8 = 99 // TODO: what is the offest of this field?
	a.Push1(off)
	a.Emit(asm.ADD)
	return loc
}

func (p *Variable) LocationOf(a *asm.Assembler) Location {
	// TODO:  How do I get the offest, type, and storage class info here?
	// For now global storage and parameters are the only supported things.
	if p.IsStorage() {
		a.PushBytes(crypto.Keccak256([]byte(p.id.name)))
		s := &StorageLocation{}
		return s
	} else {
		ps, ok := p.id.sym.(*types.ParamSym)
		if !ok {
			panic("unimplemented")
		}
		off := ps.GetOffset()
		if off >= 256 {
			panic("fix me by using push2 or something")
		}
		a.Push1(0)
		a.Emit(asm.MLOAD)
		a.Push1(uint8(off))
		a.Emit(asm.ADD)
		m := &MemoryLocation{}
		return m
	}
}

type StorageLocation struct {
}

func (p *StorageLocation) Fetch(a *asm.Assembler) {
	a.Emit(asm.SLOAD)
}

func (p *StorageLocation) StoreTo(a *asm.Assembler) {
	a.Emit(asm.SWAP1)
	a.Emit(asm.SSTORE)
}

type MemoryLocation struct {
}

func (p *MemoryLocation) Fetch(a *asm.Assembler) {
	a.Emit(asm.MLOAD)
}
func (p *MemoryLocation) StoreTo(a *asm.Assembler) {
	a.Emit(asm.SWAP1)
	a.Emit(asm.MSTORE)
}
