package types

type addressType struct{}
type boolType struct{}
type bytes32Type struct{}
type uint256Type struct{}
type uint8Type struct{}
type int256Type struct{}
type byteType struct{}
type stringType struct{}

var Address *addressType = &addressType{}
var Bool *boolType = &boolType{}
var Bytes32 *bytes32Type = &bytes32Type{}
var Uint8 *uint8Type = &uint8Type{}
var Uint256 *uint256Type = &uint256Type{}
var Int256 *int256Type = &int256Type{}
var Byte *byteType = &byteType{}
var String *stringType = &stringType{}

var addressTypeSym *BuiltinTypeSym = &BuiltinTypeSym{name: "address", typ: Address}
var boolTypeSym *BuiltinTypeSym = &BuiltinTypeSym{name: "bool", typ: Bool}
var bytes32TypeSym *BuiltinTypeSym = &BuiltinTypeSym{name: "bytes32", typ: Bytes32}
var uint256TypeSym *BuiltinTypeSym = &BuiltinTypeSym{name: "uint256", typ: Uint256}
var uint8TypeSym *BuiltinTypeSym = &BuiltinTypeSym{name: "uint8", typ: Uint8}
var int256TypeSym *BuiltinTypeSym = &BuiltinTypeSym{name: "int256", typ: Int256}
var byteTypeSym *BuiltinTypeSym = &BuiltinTypeSym{name: "byte", typ: Byte}
var stringTypeSym *BuiltinTypeSym = &BuiltinTypeSym{name: "string", typ: String}

var BuiltinSymbols map[string]ISymbol = map[string]ISymbol{
	"address": addressTypeSym,
	"bool":    boolTypeSym,
	"bytes32": bytes32TypeSym,
	"string":  stringTypeSym,
	"uint256": uint256TypeSym,
	"uint8":   uint8TypeSym,
	"int256":  int256TypeSym,
	"byte":    byteTypeSym,

	"wei": &BuiltinTypeSym{name: "wei"},

	"msg": &BuiltinTypeSym{
		name: "msg",
		symtab: RootSymTab("msg", map[string]ISymbol{
			"sender": addressTypeSym,
		}),
	},

	"keccak256": &BuiltinTypeSym{name: "keccak256"},
}
