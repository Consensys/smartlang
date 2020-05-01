package ast

import (
	"reflect"
	"smart/token"
	"smart/types"
)

var typeMap map[interface{}]types.IType = make(map[interface{}]types.IType, 0)

func GetType(node interface{}) (typ types.IType, ok bool) {
	typ, ok = typeMap[node]
	return typ, ok
}

func (p *Program) checkType() {
	for _, v := range p.decls {
		v.checkType()
	}
}

func (p *Contract) checkType() {
	for _, v := range p.decls {
		v.checkType()
	}
}

func (p *Class) checkType() {
	for _, v := range p.decls {
		v.checkType()
	}
}

func (p *ImportDecl) checkType() {
	p.prog.checkType()
}

func (p *Interface) checkType() {
}

func (p *FallbackDecl) checkType() {
	p.codeblock.checkType()
}

func (p *Constructor) checkType() {
	p.code.checkType()
}

func (p *PoolDecl) checkType() {
}

func (p *StorageDecl) checkType() {
	// TODO: Reject recursive definitions if necessary?
	// Nothing else to do here? We'll reject things like interfaces and contracts when we do code generation (determine storage locations).
}

func (p *ConstDecl) checkType() {
	constType := p.sym.GetType()
	valueType := p.value.checkType()
	if valueType != constType {
		p.value.Err("initializer type mismatch")
	}
}

func (p *EventDecl) checkType() {
	// Nothing to do here? We'll check for legal types during code generation?
}

func (p *FuncDef) checkType() {
	p.body.checkType()
}

// Statements

func (p *CodeBlock) checkType() {
	for _, v := range p.statements {
		v.checkType()
	}
}

func (p *VarDecl) checkType() {
	varType := p.sym.GetType()
	if p.value != nil {
		valueType := p.value.checkType()
		// TODO: Need an actual compatibility check here, not just equivalence? E.g. assigning an array literal
		if valueType != varType {
			p.value.Err("initializer type mismatch")
		}
	}
}

func (p *AssertStmt) checkType() {
	typ := p.cond.checkType()
	if typ != types.Bool {
		p.cond.Err("expected bool type")
	}
}

func (p *ReturnStmt) checkType() {
	if len(p.args) != len(p.ctx.function.Returns) {
		p.Err("wrong return value count")
	}
	for i, v := range p.args {
		expected := p.ctx.function.Returns[i].GetType()
		if v.checkType() != expected {
			v.Err("return type mismatch")
		}
		id := p.ctx.function.Returns[i].GetId()
		if len(id) > 0 {
			if v.id == nil {
				v.Err("missing return value name")
			} else if v.id.name != id {
				v.Err("return name mismatch")
			}
		} else {
			if v.id != nil {
				v.Err("unexpected return value name")
			}
		}
	}
}

func (p *ExternalCall) checkType() {
	if p.dest.checkType() != types.Address {
		p.dest.Err("expected address destination")
	}
	if p.value != nil && p.value.checkType() != types.Uint256 {
		p.value.Err("expected uint256 value")
	}
	if p.gas != nil && p.gas.checkType() != types.Uint256 {
		p.gas.Err("expected uint256 gas amount")
	}
	if p.data != nil {
		arrayType, ok := p.data.checkType().(*types.ArrayType)
		if !ok || arrayType.GetElementType() != types.Byte {
			p.data.Err("expected byte array data")
		}
	}
}

func (p *EmitStmt) checkType() {
	f, ok := p.id.checkType().(*types.EventType)
	if !ok {
		p.Err("attempt to emit a non-event")
	}
	if len(p.args) != len(f.Params) {
		p.Err("wrong number of event arguments")
	}
	for i, v := range p.args {
		v.checkType()
		if v.id.name != f.Params[i].GetId() {
			v.Err("wrong argument name")
		}
	}
}

func (p *IfStmt) checkType() {
	typ := p.condition.checkType()
	if typ != types.Bool {
		p.condition.Err("expected bool type")
	}
	p.thenpart.checkType()
	if p.elsepart != nil {
		p.elsepart.checkType()
	}
}

func (p *ForStmt) checkType() {
	if p.init != nil {
		p.init.checkType()
	}
	typ := p.cond.checkType()
	if typ != types.Bool {
		p.cond.Err("expected bool type")
	}
	if p.loop != nil {
		p.loop.checkType()
	}
	p.body.checkType()
}

func (p *WhileStmt) checkType() {
	typ := p.cond.checkType()
	if typ != types.Bool {
		p.Err("expected bool type")
	}
	p.body.checkType()
}

func (p *Move) checkType() {
	lhsType := p.lhs.checkType()
	if p.amount != nil {
		if p.amount.checkType() != types.Uint256 {
			p.Err("amount to move must by uint256")
		}
	}
	rhsType := p.rhs.checkType()
	if lhsType != rhsType {
		p.Err("type mismatch for move")
	}
	checkLHS(p, []IDesignator{p.lhs}, p.IsStore())
}

// Helper

func (p *Argument) checkType() types.IType {
	typ := p.value.checkType()
	typeMap[p] = typ
	return typ
}

// Expressions

func (p *Binary) checkType() types.IType {
	left := p.left.checkType()
	right := p.right.checkType()
	switch p.op {
	default:
		panic("unimplemented")
	case token.LOGAND, token.LOGOR:
		if left != types.Bool {
			p.left.Err("expected bool type")
		}
		if right != types.Bool {
			p.right.Err("expected bool type")
		}
		typeMap[p] = types.Bool
		return types.Bool
	case token.EQEQ, token.NEQ:
		// We only have one integer type for now (uint256), so this strict test is okay
		// TODO: Handle Int256!
		if left != right {
			p.Err("type mismatch for equality")
		}
		if left != types.Address && left != types.Bool && left != types.Uint256 && left != types.Bytes32 {
			p.Err("only primitive types can be compared")
		}
		typeMap[p] = types.Bool
		return types.Bool
	case token.GEQ, token.LEQ, token.LT, token.GT:
		if left != types.Uint256 || right != types.Uint256 {
			p.Err("comparisons are only supported for uint256")
		}
		return types.Bool
	case token.MINUS, token.PLUS, token.DIV, token.STAR, token.EXP:
		if left != types.Uint256 || right != types.Uint256 {
			p.Err("both operands must be uint256")
		}
		typeMap[p] = types.Uint256
		return types.Uint256
	case token.LSHIFT, token.RSHIFT:
		if left != types.Bytes32 {
			p.Err("left operand for a shift must be bytes32")
		}
		if right != types.Uint256 {
			p.Err("right operand for a shift must be uint256")
		}
		typeMap[p] = types.Bytes32
		return types.Bytes32
	case token.BITAND, token.BITOR:
		if left != types.Bytes32 || right != types.Bytes32 {
			p.Err("bitwise and and or are only supported for bytes32")
		}
		typeMap[p] = types.Bytes32
		return types.Bytes32
	}
}

func (p *AsExpr) checkType() types.IType {
	_ = p.expr.checkType()
	targetType := p.typeAST.buildType()
	typeMap[p] = targetType
	return targetType
	// TODO: Do we have to check compatibility here? Or do we do this during code generation?
}

func (p *Unary) checkType() types.IType {
	typ := p.expr.checkType()
	switch p.op {
	case token.BANG:
		if typ != types.Bool {
			p.Err("operand for negation must be a boolean")
		}
	default:
		panic("unimplemented")
	}
	return typ
}

func (p *CallExpr) checkType() types.IType {
	typ := p.fn.checkType()
	for _, v := range p.args {
		v.checkType()
	}
	ft, ok := typ.(*types.FunctionType)
	if !ok {
		panic("non-function type in CallExpr")
	}
	if len(ft.Returns) == 1 {
		typ = ft.Returns[0].GetType()
	}
	typeMap[p] = typ
	return typ
}

func (p *NewExpr) checkType() types.IType {
	_ = p.qual.checkType()
	for _, v := range p.args {
		v.checkType()
	}
	panic("unimplemented")
}

func (p *FieldAccess) checkType() types.IType {
	v, ok := p.expr.(*Variable)
	if !ok {
		panic("unimplemented")
	}
	fielder, ok := v.id.sym.(types.Fielder)
	if !ok {
		p.Err("type does not support field access")
	}
	sym, found := fielder.GetField(p.field.name)
	if !found {
		p.Err("field not found")
	}
	typ := sym.GetType()
	typeMap[p] = typ
	return typ
}

func (p *ArrayAccess) checkType() types.IType {
	array := p.array.checkType()
	indexer, ok := array.(types.Indexer)
	if !ok {
		p.Err("attempting to index into an unindexable type")
	}
	expected := indexer.GetIndexType()
	index := p.index.checkType()
	if expected != index {
		p.Err("incorrect index type")
	}

	typ := indexer.GetElementType()
	typeMap[p] = typ
	return typ
}

func (p *BoolLiteral) checkType() types.IType {
	typeMap[p] = types.Bool
	return types.Bool
}

func (p *ArrayLiteral) checkType() types.IType {
	expected := p.typeAST.buildType()
	for _, v := range p.values {
		if v.checkType() != expected {
			v.Err("incorrect type for array literal")
		}
	}
	typ := types.NewArrayType(expected)
	typeMap[p] = typ
	return typ
}

func (p *BytesLiteral) checkType() types.IType {
	typeMap[p] = types.Bytes32
	return types.Bytes32
}

func (p *HexLiteral) checkType() types.IType {
	typeMap[p] = types.Uint256
	return types.Uint256
}

func (p *DecLiteral) checkType() types.IType {
	typeMap[p] = types.Uint256
	return types.Uint256
}

func (p *StringLiteral) checkType() types.IType {
	typeMap[p] = types.String
	return types.String
}

func (p *Variable) checkType() types.IType {
	typ := p.id.checkType()
	typeMap[p] = typ
	return typ
}

func checkLHS(p IStatement, lhs []IDesignator, isStore bool) {
	anyStorage := false
	for _, v := range lhs {
		lvalue, ok := v.(LValue)
		if !ok {
			v.Err("non-lvalue on left-hand side of an assignment or move")
		}
		if lvalue.IsStorage() {
			anyStorage = true
		}
	}

	if anyStorage && !isStore {
		p.Err("attempting to write to a storage variable outside of a store!")
	}
	if !anyStorage && isStore {
		p.Err("no storage variables being assigned to in store!")
	}
}

func (p *Assignment) checkType() {
	if len(p.lhs) != len(p.rhs) {
		p.Err("mismatched assignment lengths")
	}
	checkLHS(p, p.lhs, p.IsStore())
	lhs := make([]types.IType, len(p.lhs))
	for i, v := range p.lhs {
		lhs[i] = v.checkType()
	}
	rhs := make([]types.IType, len(p.rhs))
	if p.from != nil {
		typ := p.from.checkType()
		f, ok := typ.(*types.FunctionType)
		if !ok {
			p.Err("non-function call in from expression")
		}
		if len(p.rhs) != len(f.Returns) {
			p.Err("wrong return value count")
		}
		for i, v := range p.rhs {
			variable, ok := v.(*Variable)
			if !ok {
				v.Err("non-identifier in from assignment")
			}
			if variable.id.name != f.Returns[i].GetId() {
				v.Err("mismatched return value name")
			}
		}
	} else {
		for i, v := range p.rhs {
			rhs[i] = v.checkType()
			if lhs[i] != rhs[i] {
				v.Err("assignment type mismatch: " + reflect.TypeOf(lhs[i]).String() + " vs. " + reflect.TypeOf(rhs[i]).String())
			}
		}
	}
}

func (p *FunctionCall) checkType() {
	f, ok := p.target.checkType().(*types.FunctionType)
	if !ok {
		p.Err("attempt to call non-function")
	}
	if len(p.arguments) != len(f.Params) {
		p.Err("wrong number of arguments")
	}
	for i, v := range p.arguments {
		v.checkType()
		if v.id.name != f.Params[i].GetId() {
			v.Err("wrong argument name")
		}
	}
}

func (p *AugmentedAssignment) checkType() {
	lhs := p.lhs.checkType()
	rhs := p.rhs.checkType()
	if lhs != types.Uint256 || rhs != types.Uint256 {
		p.Err("augmented assignments are only valid for uint256")
	}
	checkLHS(p, []IDesignator{p.lhs}, p.IsStore())
	// TODO need to check that assignment allowed
	// TODO need to check the operator and that the types work for that operator
}

func (p *QualIdent) checkType() types.IType {
	typ := p.sym.GetType()
	typeMap[p] = typ
	return typ
}

func (p *Identifier) checkType() types.IType {
	typ := p.sym.GetType()
	typeMap[p] = typ
	return typ
}

func (p *Mint) checkType() types.IType {
	if p.expr != nil {
		if p.expr.checkType() != types.Uint256 {
			p.Err("amount to move must be uint256")
		}
	}
	typ := p.typeAST.buildType()
	_, ok := typ.(*types.PoolType)
	if !ok {
		p.typeAST.Err("trying to mint non-pool type")
	}
	typeMap[p] = typ
	return typ
}
