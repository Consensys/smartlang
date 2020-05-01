package types

import "fmt"

type ISymbol interface {
	GetType() IType
}

type Fielder interface {
	GetField(name string) (sym ISymbol, found bool)
}

// Kinds of symbols

type BuiltinTypeSym struct {
	// for uint256, etc.
	name   string
	typ    IType
	symtab *SymTab
}

func NewBuiltinTypeSym(name string, typ IType, symtab *SymTab) *BuiltinTypeSym {
	return &BuiltinTypeSym{name: name, typ: typ, symtab: symtab}
}
func (p *BuiltinTypeSym) GetType() IType {
	return p.typ
}
func (p *BuiltinTypeSym) GetField(name string) (sym ISymbol, found bool) {
	sym, found = p.symtab.Fetch(name)
	return sym, found
}

type FunctionSym struct {
	id  string
	Typ *FunctionType
}

func NewFunctionSym(id string, typ *FunctionType) *FunctionSym {
	return &FunctionSym{id: id, Typ: typ}
}
func (p *FunctionSym) GetType() IType {
	return p.Typ
}

type ContractSym struct {
	id     string
	symtab *SymTab
	typ    *ContractType
}

func NewContractSym(id string, typ *ContractType, symtab *SymTab) *ContractSym {
	return &ContractSym{id: id, typ: typ, symtab: symtab}
}
func (p *ContractSym) GetType() IType {
	return p.typ
}
func (p *ContractSym) GetField(name string) (sym ISymbol, found bool) {
	sym, found = p.symtab.Fetch(name)
	return sym, found
}

type ClassSym struct {
	id     string
	symtab *SymTab
	typ    *ClassType
}

func NewClassSym(id string, typ *ClassType, symtab *SymTab) *ClassSym {
	return &ClassSym{id: id, typ: typ, symtab: symtab}
}
func (p *ClassSym) GetType() IType {
	return p.typ
}
func (p *ClassSym) GetField(name string) (sym ISymbol, found bool) {
	sym, found = p.symtab.Fetch(name)
	return sym, found
}

type InterfaceSym struct {
	id     string
	symtab *SymTab
	typ    *InterfaceType
}

func NewInterfaceSym(id string, typ *InterfaceType, symtab *SymTab) *InterfaceSym {
	return &InterfaceSym{id: id, typ: typ, symtab: symtab}
}
func (p *InterfaceSym) GetType() IType {
	return p.typ
}
func (p *InterfaceSym) GetField(name string) (sym ISymbol, found bool) {
	sym, found = p.symtab.Fetch(name)
	return sym, found
}

type ImportSym struct {
	id     string
	path   string
	hash   []byte
	symtab *SymTab
	typ    *ImportType
}

func NewImportSym(id string, path string, hash []byte, typ *ImportType, symtab *SymTab) *ImportSym {
	return &ImportSym{id: id, path: path, hash: hash, typ: typ, symtab: symtab}
}
func (p *ImportSym) GetType() IType {
	return p.typ
}
func (p *ImportSym) GetField(name string) (sym ISymbol, found bool) {
	sym, found = p.symtab.Fetch(name)
	return sym, found
}

type EventSym struct {
	id  string
	Typ *EventType
}

func NewEventSym(id string, typ *EventType) *EventSym {
	return &EventSym{id: id, Typ: typ}
}
func (p *EventSym) GetType() IType {
	return p.Typ
}

type EventParamSym struct {
	id  string
	typ *ParamType
}

func NewEventParamSym(id string) *EventParamSym {
	return &EventParamSym{id: id}
}
func (p *EventParamSym) GetType() IType {
	return p.typ
}

type PoolSym struct {
	id     string
	typ    *PoolType
	symtab *SymTab
}

func NewPoolSym(id string, typ *PoolType) *PoolSym {
	ret := &PoolSym{id: id, typ: typ}
	ret.symtab = RootSymTab(fmt.Sprintf("pool %s", id), map[string]ISymbol{"sum": uint256TypeSym})
	return ret
}
func (p *PoolSym) GetType() IType {
	return p.typ
}
func (p *PoolSym) GetField(name string) (sym ISymbol, found bool) {
	typ, found := p.symtab.Fetch(name)
	return typ, found
}

type LocalSym struct {
	id  string
	typ IType
}

func NewLocalSym(id string, typ IType) *LocalSym {
	return &LocalSym{id: id, typ: typ}
}
func (p *LocalSym) GetType() IType {
	return p.typ
}
func (p *LocalSym) IsStorage() bool {
	return false
}
func (LocalSym) Lvalue() {}

type StorageSym struct {
	id  string
	Typ IType
}

func NewStorageSym(id string) *StorageSym {
	return &StorageSym{id: id}
}
func (p *StorageSym) GetType() IType {
	return p.Typ
}
func (p *StorageSym) IsStorage() bool {
	return true
}
func (StorageSym) Lvalue() {}

type ReturnSym struct {
	id  string
	typ IType
}

func NewReturnSym(id string, typ IType) *ReturnSym {
	return &ReturnSym{id: id, typ: typ}
}
func (p *ReturnSym) GetType() IType {
	return p.typ
}
func (p *ReturnSym) GetField(name string) (sym ISymbol, found bool) {
	return nil, false
}
func (p *ReturnSym) GetId() string {
	return p.id
}

type ParamSym struct {
	id     string
	typ    IType
	offset int
}

func NewParamSym(id string, typ IType) *ParamSym {
	return &ParamSym{id: id, typ: typ}
}
func (p *ParamSym) GetType() IType {
	return p.typ
}
func (p *ParamSym) GetId() string {
	return p.id
}
func (p *ParamSym) SetOffset(offset int) {
	p.offset = offset
}
func (p *ParamSym) GetOffset() int {
	return p.offset
}
func (p *ParamSym) IsStorage() bool {
	return false
}
func (p *ParamSym) Lvalue() {}

type ConstSym struct {
	id  string
	Typ IType
}

func NewConstSym(id string) *ConstSym {
	return &ConstSym{id: id}
}
func (p *ConstSym) GetType() IType {
	return p.Typ
}
