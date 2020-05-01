package ast

import (
	"fmt"
	"smart/token"
	"smart/types"
)

type ASTNode struct {
	IVisitor
	IDumper
	Positioner
	coords token.Position
	ctx    *Context
}

func (a *ASTNode) Err(msg string) {
	panic(fmt.Sprintf("%s line %d, column %d: %s", a.coords.Filename, a.coords.Line, a.coords.Column, msg))
}

func (a *ASTNode) SetCoords(pos token.Position) {
	a.coords = pos
}

type Identifier struct {
	ASTNode
	name string
	sym  types.ISymbol
}

func NewIdentifier(name string, pos token.Position) *Identifier {
	id := &Identifier{name: name}
	id.coords = pos
	return id
}

// prog = { contract | class | importdecl | interface } EOF

type Program struct {
	ASTNode
	decls []IProgramBodyDecl
}

func NewProgram(decls []IProgramBodyDecl) *Program {
	return &Program{decls: decls}
}

// importdecl = IMPORT [ STRINGLITERAL ] [ LBRACK HEXLITERAL RBRACK ] [ AS IDENT ] SEMI
type ImportDecl struct {
	ASTNode
	path  string
	hash  []byte
	alias *Identifier
	prog  *Program
	sym   *types.ImportSym
}

func NewImportDecl(path string, hash []byte, alias *Identifier, prog *Program) *ImportDecl {
	return &ImportDecl{path: path, hash: hash, alias: alias, prog: prog}
}

// contract = CONTRACT IDENT contractBody
type Contract struct {
	ASTNode
	id    *Identifier
	decls []IDecl
	sym   *types.ContractSym
}

func NewContract(id *Identifier, decls []IDecl) *Contract {
	return &Contract{id: id, decls: decls}
}

// contractBody = LCURL { classStuff | fallback } RCURL
// classStuff = constructor | pooldecl | storagedecl | constdecl SEMI | funcdef | eventDecl

type FuncDef struct {
	ASTNode
	sig  *FuncSig
	body *CodeBlock
	sym  *types.FunctionSym
}

func NewFuncDef(sig *FuncSig, body *CodeBlock) *FuncDef {
	return &FuncDef{sig: sig, body: body}
}

// fallback = FALLBACK codeblock

type FallbackDecl struct {
	ASTNode
	codeblock *CodeBlock
}

func NewFallbackDecl(c *CodeBlock) *FallbackDecl {
	return &FallbackDecl{codeblock: c}
}

// eventDecl = EVENT IDENT LPAREN [ IDENT type [ INDEXED ] { COMMA IDENT type [ INDEXED ] } ] RPAREN SEMI

type EventDecl struct {
	ASTNode
	id     *Identifier
	params []*Param
	sym    *types.EventSym
}

type Param struct {
	ASTNode
	id      *Identifier
	typeAST ITypeAST
	indexed bool
	sym     *types.ParamSym
}

func NewEventDecl(id *Identifier, params []*Param) *EventDecl {
	return &EventDecl{id: id, params: params}
}

func NewParam(id *Identifier, typeAST ITypeAST, indexed bool) *Param {
	return &Param{id: id, typeAST: typeAST, indexed: indexed}
}

// constructor = CONSTRUCTOR LPAREN [ paramList ] RPAREN codeblock
type Constructor struct {
	ASTNode
	params []*Param
	code   *CodeBlock
}

func NewConstructor(params []*Param, code *CodeBlock) *Constructor {
	return &Constructor{params: params, code: code}
}

// paramList = IDENT [MINT] type { COMMA IDENT [MINT] type }

// class = CLASS IDENT classBody

type Class struct {
	ASTNode
	id    *Identifier
	decls []IClassBodyDecl
	sym   *types.ClassSym
}

func NewClass(id *Identifier, decls []IClassBodyDecl) *Class {
	return &Class{id: id, decls: decls}
}

// classBody = LCURL { classStuff } RCURL

// pooldecl = POOL IDENT SEMI

type PoolDecl struct {
	ASTNode
	id  *Identifier
	sym *types.PoolSym
}

func NewPoolDecl(id *Identifier) *PoolDecl {
	return &PoolDecl{id: id}
}

// vardecl = VAR IDENT type [EQ expr]

type VarDecl struct {
	ASTNode
	id      *Identifier
	typeAST ITypeAST
	value   IExpr
	sym     *types.LocalSym
}

func NewVarDecl(id *Identifier, typeAST ITypeAST, value IExpr) *VarDecl {
	return &VarDecl{id: id, typeAST: typeAST, value: value}
}

// storagedecl = STORAGE IDENT type SEMI

type StorageDecl struct {
	ASTNode
	id      *Identifier
	typeAST ITypeAST
	sym     *types.StorageSym
}

func NewStorageDecl(id *Identifier, typeAST ITypeAST) *StorageDecl {
	return &StorageDecl{id: id, typeAST: typeAST}
}

// constdecl = CONST IDENT type EQ expr

type ConstDecl struct {
	ASTNode
	id      *Identifier
	typeAST ITypeAST
	value   IExpr
	sym     *types.ConstSym
}

func NewConstDecl(id *Identifier, typeAST ITypeAST, value IExpr) *ConstDecl {
	return &ConstDecl{id: id, typeAST: typeAST, value: value}
}

// type = qualident | arraytype | maptype
// qualident = IDENT [ DOT IDENT ]
type QualIdent struct {
	ASTNode
	qualifier *Identifier
	id        *Identifier
	sym       types.ISymbol
	typ       types.IType
}

func NewQualIdent(qualifier *Identifier, id *Identifier) *QualIdent {
	return &QualIdent{qualifier: qualifier, id: id}
}

// arraytype = LBRACK [ expr ] RBRACK type

type ArrayType struct {
	ASTNode
	length  IExpr
	typeAST ITypeAST
}

func NewArrayType(length IExpr, typeAST ITypeAST) *ArrayType {
	return &ArrayType{length: length, typeAST: typeAST}
}

// maptype = MAP LBRACK type RBRACK type

type MapType struct {
	ASTNode
	keyType ITypeAST
	valType ITypeAST
}

func NewMapType(keyType ITypeAST, valType ITypeAST) *MapType {
	return &MapType{keyType: keyType, valType: valType}
}

// expr  = expr00 { LOGOR expr00 }
// expr00 = expr0 { LOGAND expr0 }
// expr0 = expr1 [ relation expr1 ]
// expr1 = expr2 { addop expr2 }
// expr2 = expr3 { mulop expr3 }
// expr3 = expr4 { EXP expr4 }
// expr4 = expr5 [ AS type ]
// expr5 = designator [ LPAREN [ argList ] RPAREN ] | literal | LPAREN expr RPAREN | MINUS expr5 | LARR expr5 | new | MINT expr5 | BURN expr5 | BANG expr5
// designator = IDENT { DOT IDENT | LBRACK expr RBRACK }
// arrayLiteral = LBRACK [ exprList ] RBRACK
// literal = DECLITERAL | HEXLITERAL | BYTESLITERAL | TRUE | FALSE | arrayLiteral | STRINGLITERAL
// relation = EQEQ | NEQ | LT | LEQ | GT | GEQ
// addop = PLUS | MINUS | BITAND | BITOR | LSHIFT | RSHIFT
// mulop = STAR | DIV
// new = NEW IDENT LPAREN [ argList ] RPAREN

type FieldAccess struct {
	ASTNode
	expr  IExpr
	field *Identifier
}

func NewFieldAccess(expr IExpr, field *Identifier) *FieldAccess {
	return &FieldAccess{expr: expr, field: field}
}
func (f *FieldAccess) IsStorage() bool {
	l, ok := f.expr.(LValue)
	if !ok {
		panic("unimplemented")
	}
	return l.IsStorage()
}

type ArrayAccess struct {
	ASTNode
	array IExpr
	index IExpr
}

func NewArrayAccess(array IExpr, index IExpr) *ArrayAccess {
	return &ArrayAccess{array: array, index: index}
}
func (a *ArrayAccess) IsStorage() bool {
	l, ok := a.array.(LValue)
	if !ok {
		panic("unimplemented")
	}
	return l.IsStorage()
}
func (ArrayAccess) Lvalue() {}

type DecLiteral struct {
	ASTNode
	lit string
}

func NewDecLiteral(lit string) *DecLiteral {
	return &DecLiteral{lit: lit}
}

type HexLiteral struct {
	ASTNode
	lit string
}

func NewHexLiteral(lit string) *HexLiteral {
	return &HexLiteral{lit: lit}
}

type BytesLiteral struct {
	ASTNode
	lit string
}

func NewBytesLiteral(lit string) *BytesLiteral {
	return &BytesLiteral{lit: lit}
}

type StringLiteral struct {
	ASTNode
	lit string
}

func NewStringLiteral(lit string) *StringLiteral {
	return &StringLiteral{lit: lit}
}

type BoolLiteral struct {
	ASTNode
	value bool
}

func NewBoolLiteral(value bool) *BoolLiteral {
	return &BoolLiteral{value: value}
}

type Variable struct {
	ASTNode
	id *Identifier
}

func NewVariable(id *Identifier) *Variable {
	return &Variable{id: id}
}
func (v *Variable) IsStorage() bool {
	l, ok := v.id.sym.(LValue)
	if !ok {
		v.Err("non-lvalue used in assignment")
	}
	return l.IsStorage()
}
func (v *Variable) Lvalue() {}

type NewExpr struct {
	ASTNode
	qual *QualIdent
	args []*Argument
}

func NewNewExpr(qual *QualIdent, args []*Argument) *NewExpr {
	return &NewExpr{qual: qual, args: args}
}

type AsExpr struct {
	ASTNode
	expr    IExpr
	typeAST ITypeAST
}

func NewAsExpr(expr IExpr, typeAST ITypeAST) *AsExpr {
	return &AsExpr{expr: expr, typeAST: typeAST}
}

type Binary struct {
	ASTNode
	op    token.Token
	left  IExpr
	right IExpr
}

func NewBinary(op token.Token, left IExpr, right IExpr) *Binary {
	return &Binary{op: op, left: left, right: right}
}

type Unary struct {
	ASTNode
	op   token.Token
	expr IExpr
}

func NewUnary(op token.Token, expr IExpr) *Unary {
	return &Unary{op: op, expr: expr}
}

type CallExpr struct {
	ASTNode
	fn   IDesignator
	args []*Argument
}

func NewCallExpr(fn IDesignator, args []*Argument) *CallExpr {
	return &CallExpr{fn: fn, args: args}
}

// exprList = expr { COMMA expr }
// designatorList = designator { COMMA designator }

// simpleStatement = assignOrMoveOrCall | vardecl | constdecl | store | assert | ret | call | emit
// complexStatement = ifstmt | loop | codeblock

// emit = EMIT IDENT LPAREN [ argList ] RPAREN

type EmitStmt struct {
	ASTNode
	id   *Identifier
	args []*Argument
}

func NewEmitStmt(id *Identifier, args []*Argument) *EmitStmt {
	return &EmitStmt{id: id, args: args}
}

// ret = RETURN [ exprList ]
type ReturnStmt struct {
	ASTNode
	args []*Argument
}

func NewReturnStmt(args []*Argument) *ReturnStmt {
	return &ReturnStmt{args: args}
}

// assignOrMoveOrCall = designatorList ( EQ exprList [ FROM expr ] | LARR expr [ FROM expr ] | LPAREN [ argList ] RPAREN | ( PLUSEQ | MINUSEQ ) expr )

type Assignment struct {
	ASTNode
	lhs   []IDesignator
	rhs   []IExpr
	from  *CallExpr
	store bool
}

func NewAssignment(lhs []IDesignator, rhs []IExpr, from *CallExpr) *Assignment {
	return &Assignment{lhs: lhs, rhs: rhs, from: from}
}
func (a *Assignment) SetStore(value bool) {
	a.store = value
}
func (a *Assignment) IsStore() bool {
	return a.store
}

type Move struct {
	ASTNode
	lhs    IDesignator
	amount IExpr
	rhs    IExpr
	store  bool
}

func NewMove(lhs IDesignator, amount IExpr, rhs IExpr) *Move {
	return &Move{lhs: lhs, amount: amount, rhs: rhs}
}
func (m *Move) SetStore(value bool) {
	m.store = value
}
func (m *Move) IsStore() bool {
	return m.store
}

type FunctionCall struct {
	ASTNode
	target    IDesignator
	arguments []*Argument
}

func NewFunctionCall(target IDesignator, arguments []*Argument) *FunctionCall {
	return &FunctionCall{target: target, arguments: arguments}
}

// argList = arg { COMMA arg }
// arg = IDENT COLON expr

type Argument struct {
	ASTNode
	id    *Identifier
	value IExpr
}

func NewArgument(id *Identifier, value IExpr) *Argument {
	return &Argument{id: id, value: value}
}

type AugmentedAssignment struct {
	ASTNode
	lhs   IDesignator
	op    token.Token
	rhs   IExpr
	store bool
}

func NewAugmentedAssignment(lhs IDesignator, op token.Token, rhs IExpr) *AugmentedAssignment {
	return &AugmentedAssignment{lhs: lhs, op: op, rhs: rhs}
}
func (a *AugmentedAssignment) SetStore(value bool) {
	a.store = value
}
func (a *AugmentedAssignment) IsStore() bool {
	return a.store
}

// codeblock = LCURL { simpleStatement SEMI | complexStatement } RCURL

type CodeBlock struct {
	ASTNode
	statements []IStatement
}

func NewCodeBlock(statements []IStatement) *CodeBlock {
	return &CodeBlock{statements: statements}
}

// ifstmt = IF expr codeblock [ ELSE codeblock ]

type IfStmt struct {
	ASTNode
	condition IExpr
	thenpart  *CodeBlock
	elsepart  *CodeBlock
}

func NewIfStmt(condition IExpr, thenpart *CodeBlock, elsepart *CodeBlock) *IfStmt {
	return &IfStmt{condition: condition, thenpart: thenpart, elsepart: elsepart}
}

// forstmt = FOR [ simpleStatement ] SEMI [ expr ] SEMI [ simpleStatement ] codeblock
type ForStmt struct {
	ASTNode
	init IStatement
	cond IExpr
	loop IStatement
	body *CodeBlock
}

func NewForStmt(init IStatement, cond IExpr, loop IStatement, body *CodeBlock) *ForStmt {
	return &ForStmt{init: init, cond: cond, loop: loop, body: body}
}

// while = WHILE expr codeblock
type WhileStmt struct {
	ASTNode
	cond IExpr
	body *CodeBlock
}

func NewWhileStmt(cond IExpr, body *CodeBlock) *WhileStmt {
	return &WhileStmt{cond: cond, body: body}
}

// loop = forstmt | while

// assert = ASSERT expr [ COMMA STRINGLITERAL ]
type AssertStmt struct {
	ASTNode
	cond   IExpr
	reason string
}

func NewAssertStmt(cond IExpr, reason string) *AssertStmt {
	return &AssertStmt{cond: cond, reason: reason}
}

// funcsig = { decorator } FUNC IDENT LPAREN [ paramList ] RPAREN [ type | LPAREN IDENT type { COMMA IDENT type } RPAREN ]
// decorator = PAYABLE | solidity
// solidity = SOLIDITY COLON IDENT LPAREN [ type { COMMA type } ] RPAREN [ type | LPAREN type { COMMA type } RPAREN ]
// funcdef = funcsig codeblock

type PayableDecorator struct {
	ASTNode
}

func NewPayableDecorator() *PayableDecorator {
	return &PayableDecorator{}
}

type SolidityDecorator struct {
	ASTNode
	id          *Identifier
	paramTypes  []ITypeAST
	returnTypes []ITypeAST
}

func NewSolidityDecorator(id *Identifier, paramTypes []ITypeAST, returnTypes []ITypeAST) *SolidityDecorator {
	return &SolidityDecorator{id: id, paramTypes: paramTypes, returnTypes: returnTypes}
}

type InterfaceSig struct {
	ASTNode
	sig *FuncSig
	sym *types.FunctionSym
}

func NewInterfaceSig(sig *FuncSig) *InterfaceSig {
	return &InterfaceSig{sig: sig}
}

type FuncSig struct {
	ASTNode
	id         *Identifier
	params     []*Param
	ret        []*ReturnType
	decorators []IDecorator
	mutating   bool
	sym        *types.FunctionSym
}

func NewFuncSig(id *Identifier, params []*Param, rettype []*ReturnType, decorators []IDecorator, mutating bool) *FuncSig {
	return &FuncSig{id: id, params: params, ret: rettype, decorators: decorators}
}

type ReturnType struct {
	ASTNode
	id      *Identifier
	typeAST ITypeAST
}

func NewReturnType(id *Identifier, typeAST ITypeAST) *ReturnType {
	return &ReturnType{id: id, typeAST: typeAST}
}

// call = CALL expr { COMMA ( GAS | DATA | VALUE ) EQ expr }
type ExternalCall struct {
	ASTNode
	dest  IExpr
	data  IExpr
	gas   IExpr
	value IExpr
}

func NewExternalCall(dest IExpr, data IExpr, gas IExpr, value IExpr) *ExternalCall {
	return &ExternalCall{dest: dest, data: data, gas: gas, value: value}
}

type ReturnValues struct {
	ASTNode
	ids   []*Identifier
	types []ITypeAST
}

func NewReturnValues(ids []*Identifier, types []ITypeAST) *ReturnValues {
	return &ReturnValues{ids: ids, types: types}
}

// interface = INTERFACE IDENT LCURL { funcsig SEMI } RCURL
type Interface struct {
	ASTNode
	id    *Identifier
	decls []*InterfaceSig
	sym   *types.InterfaceSym
}

func NewInterface(id *Identifier, decls []*InterfaceSig) *Interface {
	return &Interface{id: id, decls: decls}
}

type ArrayLiteral struct {
	ASTNode
	typeAST ITypeAST
	values  []IExpr
}

func NewArrayLiteral(typeAST ITypeAST, expressions []IExpr) *ArrayLiteral {
	return &ArrayLiteral{typeAST: typeAST, values: expressions}
}

type Mint struct {
	ASTNode
	expr    IExpr
	typeAST ITypeAST
}

func NewMint(expr IExpr, typeAST ITypeAST) *Mint {
	return &Mint{expr: expr, typeAST: typeAST}
}
