package types

// Types that exist at runtime

type IType interface{}

type Indexer interface {
	GetElementType() IType
	GetIndexType() IType
}

type MapType struct {
	key IType
	val IType
}

func NewMapType(key IType, val IType) *MapType {
	return &MapType{key: key, val: val}
}
func (p *MapType) GetIndexType() IType {
	return p.key
}
func (p *MapType) GetElementType() IType {
	return p.val
}

type ArrayType struct {
	// TODO: what about the bounds?
	base IType
}

func NewArrayType(base IType) *ArrayType {
	return &ArrayType{base: base}
}
func (p *ArrayType) GetIndexType() IType {
	return Uint256
}
func (p *ArrayType) GetElementType() IType {
	return p.base
}

type FunctionType struct {
	Params  []*ParamSym
	Returns []*ReturnSym
}

func NewFunctionType() *FunctionType {
	return &FunctionType{}
}

type EventType struct {
	Params []*ParamSym
}

func NewEventType() *EventType {
	return &EventType{}
}

type ParamType struct {
	name string
	typ  IType
}

func NewParamType(name string, typ IType) *ParamType {
	return &ParamType{name: name, typ: typ}
}

type PoolType struct {
	id string
}

func NewPoolType(id string) *PoolType {
	return &PoolType{id: id}
}

type ContractType struct {
	members map[string]ISymbol
}

func NewContractType() *ContractType {
	return &ContractType{}
}

type ClassType struct {
	members map[string]ISymbol
}

func NewClassType() *ClassType {
	return &ClassType{}
}

type ImportType struct {
	members map[string]ISymbol
}

func NewImportType() *ImportType {
	return &ImportType{}
}

type InterfaceType struct {
	members map[string]ISymbol
}

func NewInterfaceType() *InterfaceType {
	return &InterfaceType{}
}

type ReturnType struct {
	name string
	typ  IType
}

func NewReturnType(name string, typ IType) *ReturnType {
	return &ReturnType{name: name, typ: typ}
}
