package asm

import "fmt"

type OpCode uint8

const (
	STOP OpCode = iota
	ADD
	MUL
	SUB
	DIV
	SDIV
	MOD
	SMOD
	ADDMOD
	MULMOD
	EXP
	SIGNEXTEND
)

const (
	LT OpCode = iota + 0x10
	GT
	SLT
	SGT
	EQ
	ISZERO
	AND
	OR
	XOR
	NOT
	BYTE
	SHL
	SHR
	SAR

	SHA3 = 0x20
)

const (
	ADDRESS OpCode = 0x30 + iota
	BALANCE
	ORIGIN
	CALLER
	CALLVALUE
	CALLDATALOAD
	CALLDATASIZE
	CALLDATACOPY
	CODESIZE
	CODECOPY
	GASPRICE
	EXTCODESIZE
	EXTCODECOPY
	RETURNDATASIZE
	RETURNDATACOPY
	EXTCODEHASH
)

const (
	BLOCKHASH OpCode = 0x40 + iota
	COINBASE
	TIMESTAMP
	NUMBER
	DIFFICULTY
	GASLIMIT
	CHAINID     = 0x46
	SELFBALANCE = 0x47
)

const (
	POP OpCode = 0x50 + iota
	MLOAD
	MSTORE
	MSTORE8
	SLOAD
	SSTORE
	JUMP
	JUMPI
	PC
	MSIZE
	GAS
	JUMPDEST
)

const (
	PUSH1 OpCode = 0x60 + iota
	PUSH2
	PUSH3
	PUSH4
	PUSH5
	PUSH6
	PUSH7
	PUSH8
	PUSH9
	PUSH10
	PUSH11
	PUSH12
	PUSH13
	PUSH14
	PUSH15
	PUSH16
	PUSH17
	PUSH18
	PUSH19
	PUSH20
	PUSH21
	PUSH22
	PUSH23
	PUSH24
	PUSH25
	PUSH26
	PUSH27
	PUSH28
	PUSH29
	PUSH30
	PUSH31
	PUSH32
	DUP1
	DUP2
	DUP3
	DUP4
	DUP5
	DUP6
	DUP7
	DUP8
	DUP9
	DUP10
	DUP11
	DUP12
	DUP13
	DUP14
	DUP15
	DUP16
	SWAP1
	SWAP2
	SWAP3
	SWAP4
	SWAP5
	SWAP6
	SWAP7
	SWAP8
	SWAP9
	SWAP10
	SWAP11
	SWAP12
	SWAP13
	SWAP14
	SWAP15
	SWAP16
)

const (
	LOG0 OpCode = 0xa0 + iota
	LOG1
	LOG2
	LOG3
	LOG4
)

const (
	CREATE OpCode = 0xf0 + iota
	CALL
	CALLCODE
	RETURN
	DELEGATECALL
	CREATE2
	STATICCALL = 0xfa

	REVERT       = 0xfd
	SELFDESTRUCT = 0xff
)

var codeToString = map[OpCode]string{
	STOP:       "STOP",
	ADD:        "ADD",
	MUL:        "MUL",
	SUB:        "SUB",
	DIV:        "DIV",
	SDIV:       "SDIV",
	MOD:        "MOD",
	SMOD:       "SMOD",
	ADDMOD:     "ADDMOD",
	MULMOD:     "MULMOD",
	EXP:        "EXP",
	SIGNEXTEND: "SIGNEXTEND",

	LT:     "LT",
	GT:     "GT",
	SLT:    "SLT",
	SGT:    "SGT",
	EQ:     "EQ",
	ISZERO: "ISZERO",
	AND:    "AND",
	OR:     "OR",
	XOR:    "XOR",
	NOT:    "NOT",
	BYTE:   "BYTE",
	SHL:    "SHL",
	SHR:    "SHR",
	SAR:    "SAR",

	SHA3: "SHA3",

	ADDRESS:        "ADDRESS",
	BALANCE:        "BALANCE",
	ORIGIN:         "ORIGIN",
	CALLER:         "CALLER",
	CALLVALUE:      "CALLVALUE",
	CALLDATALOAD:   "CALLDATALOAD",
	CALLDATASIZE:   "CALLDATASIZE",
	CALLDATACOPY:   "CALLDATACOPY",
	CODESIZE:       "CODESIZE",
	CODECOPY:       "CODECOPY",
	GASPRICE:       "GASPRICE",
	EXTCODESIZE:    "EXTCODESIZE",
	EXTCODECOPY:    "EXTCODECOPY",
	RETURNDATASIZE: "RETURNDATASIZE",
	RETURNDATACOPY: "RETURNDATACOPY",
	EXTCODEHASH:    "EXTCODEHASH",

	BLOCKHASH:   "BLOCKHASH",
	COINBASE:    "COINBASE",
	TIMESTAMP:   "TIMESTAMP",
	NUMBER:      "NUMBER",
	DIFFICULTY:  "DIFFICULTY",
	GASLIMIT:    "GASLIMIT",
	CHAINID:     "CHAINID",
	SELFBALANCE: "SELFBALANCE",

	POP:      "POP",
	MLOAD:    "MLOAD",
	MSTORE:   "MSTORE",
	MSTORE8:  "MSTORE8",
	SLOAD:    "SLOAD",
	SSTORE:   "SSTORE",
	JUMP:     "JUMP",
	JUMPI:    "JUMPI",
	PC:       "PC",
	MSIZE:    "MSIZE",
	GAS:      "GAS",
	JUMPDEST: "JUMPDEST",

	PUSH1:  "PUSH1",
	PUSH2:  "PUSH2",
	PUSH3:  "PUSH3",
	PUSH4:  "PUSH4",
	PUSH5:  "PUSH5",
	PUSH6:  "PUSH6",
	PUSH7:  "PUSH7",
	PUSH8:  "PUSH8",
	PUSH9:  "PUSH9",
	PUSH10: "PUSH10",
	PUSH11: "PUSH11",
	PUSH12: "PUSH12",
	PUSH13: "PUSH13",
	PUSH14: "PUSH14",
	PUSH15: "PUSH15",
	PUSH16: "PUSH16",
	PUSH17: "PUSH17",
	PUSH18: "PUSH18",
	PUSH19: "PUSH19",
	PUSH20: "PUSH20",
	PUSH21: "PUSH21",
	PUSH22: "PUSH22",
	PUSH23: "PUSH23",
	PUSH24: "PUSH24",
	PUSH25: "PUSH25",
	PUSH26: "PUSH26",
	PUSH27: "PUSH27",
	PUSH28: "PUSH28",
	PUSH29: "PUSH29",
	PUSH30: "PUSH30",
	PUSH31: "PUSH31",
	PUSH32: "PUSH32",
	DUP1:   "DUP1",
	DUP2:   "DUP2",
	DUP3:   "DUP3",
	DUP4:   "DUP4",
	DUP5:   "DUP5",
	DUP6:   "DUP6",
	DUP7:   "DUP7",
	DUP8:   "DUP8",
	DUP9:   "DUP9",
	DUP10:  "DUP10",
	DUP11:  "DUP11",
	DUP12:  "DUP12",
	DUP13:  "DUP13",
	DUP14:  "DUP14",
	DUP15:  "DUP15",
	DUP16:  "DUP16",
	SWAP1:  "SWAP1",
	SWAP2:  "SWAP2",
	SWAP3:  "SWAP3",
	SWAP4:  "SWAP4",
	SWAP5:  "SWAP5",
	SWAP6:  "SWAP6",
	SWAP7:  "SWAP7",
	SWAP8:  "SWAP8",
	SWAP9:  "SWAP9",
	SWAP10: "SWAP10",
	SWAP11: "SWAP11",
	SWAP12: "SWAP12",
	SWAP13: "SWAP13",
	SWAP14: "SWAP14",
	SWAP15: "SWAP15",
	SWAP16: "SWAP16",

	LOG0: "LOG0",
	LOG1: "LOG1",
	LOG2: "LOG2",
	LOG3: "LOG3",
	LOG4: "LOG4",

	CREATE:       "CREATE",
	CALL:         "CALL",
	CALLCODE:     "CALLCODE",
	RETURN:       "RETURN",
	DELEGATECALL: "DELEGATECALL",
	CREATE2:      "CREATE2",
	STATICCALL:   "STATICCALL",

	REVERT:       "REVERT",
	SELFDESTRUCT: "SELFDESTRUCT",
}

type Comment struct {
	comment string
}

func (c Comment) Length() int {
	return 0
}

func (c Comment) GetComment() string {
	return c.comment
}

type Instruction interface {
	Length() int
	GetComment() string
}

type PushByte struct {
	Comment
	value uint8
}

func (PushByte) Length() int {
	return 2
}

type PushBytes struct {
	Comment
	bytes []byte
}

func (p PushBytes) Length() int {
	return 1 + len(p.bytes)
}

type PushLabel struct {
	Comment
	label string
}

func (PushLabel) Length() int {
	return 3
}

type SimpleOp struct {
	Comment
	code OpCode
}

func (SimpleOp) Length() int {
	return 1
}

type LabelOp struct {
	Comment
	label string
}

func (LabelOp) Length() int {
	return 1
}

type Assembler struct {
	instructions []Instruction
}

func NewAssembler() *Assembler {
	return &Assembler{make([]Instruction, 0)}
}

func (a *Assembler) Emit(code OpCode) {
	a.instructions = append(a.instructions, SimpleOp{code: code})
}
func (a *Assembler) Emitc(code OpCode, comment string) {
	a.instructions = append(a.instructions, SimpleOp{code: code, Comment: Comment{comment}})
}

func (a *Assembler) Push1(value uint8) {
	a.instructions = append(a.instructions, PushByte{value: value})
}
func (a *Assembler) Pushc(value uint8, comment string) {
	a.instructions = append(a.instructions, PushByte{value: value, Comment: Comment{comment}})
}

func (a *Assembler) PushLabel(label string) {
	a.instructions = append(a.instructions, PushLabel{label: label})
}
func (a *Assembler) PushLabelc(label string, comment string) {
	a.instructions = append(a.instructions, PushLabel{label: label, Comment: Comment{comment}})
}

func (a *Assembler) PushBytes(bytes []byte) {
	if len(bytes) > 32 {
		panic("Maximum amount of data to PUSH is 32 bytes.")
	}
	a.instructions = append(a.instructions, PushBytes{bytes: bytes})
}
func (a *Assembler) PushBytesc(bytes []byte, comment string) {
	if len(bytes) > 32 {
		panic("Maximum amount of data to PUSH is 32 bytes.")
	}
	a.instructions = append(a.instructions, PushBytes{bytes: bytes, Comment: Comment{comment}})
}

func (a *Assembler) Label(label string) {
	a.instructions = append(a.instructions, LabelOp{label: label})
}
func (a *Assembler) Labelc(label string, comment string) {
	a.instructions = append(a.instructions, LabelOp{label: label, Comment: Comment{comment}})
}

func (a *Assembler) Comment(comment string) {
	a.instructions = append(a.instructions, Comment{comment: comment})
}

func (a *Assembler) Print() {
	pc := 0
	labels := make(map[string]int)
	for _, i := range a.instructions {
		switch v := i.(type) {
		case LabelOp:
			labels[v.label] = pc
		}
		pc += i.Length()
	}

	pc = 0
	for _, i := range a.instructions {
		commentSuffix := ""
		if len(i.GetComment()) > 0 {
			commentSuffix = " // " + i.GetComment()
		}
		switch v := i.(type) {
		case SimpleOp:
			fmt.Printf("0x%-6.2X 0x%X    (%s)%s\n", pc, v.code, codeToString[v.code], commentSuffix)
		case PushByte:
			fmt.Printf("0x%-6.2X 0x%X %2d (PUSH1)%s\n", pc, PUSH1, v.value, commentSuffix)
		case PushLabel:
			pos := labels[v.label]
			fmt.Printf("0x%-6.2X 0x%X %2d (PUSH2 @%s)%s\n", pc, PUSH2, pos, v.label, commentSuffix)
		case LabelOp:
			fmt.Printf("%s:\n", v.label)
			fmt.Printf("0x%-6.2X 0x%X    (JUMPDEST)%s\n", pc, JUMPDEST, commentSuffix)
		case PushBytes:
			op := OpCode(uint8(PUSH1) + uint8(len(v.bytes)) - 1)
			fmt.Printf("0x%-6.2X 0x%X 0x%X%s\n", pc, op, v.bytes, commentSuffix)
		case Comment:
			fmt.Printf("// %s\n", v.GetComment())
		default:
			panic("Unknown instruction type.")
		}
		pc += i.Length()
	}
}

func (a *Assembler) Compile() []byte {
	out := make([]byte, 0)

	pc := 0
	labels := make(map[string]int)
	for _, i := range a.instructions {
		switch v := i.(type) {
		case LabelOp:
			labels[v.label] = pc
		}
		pc += i.Length()
	}

	pc = 0
	for _, i := range a.instructions {
		// commentSuffix := ""
		// if len(i.GetComment()) > 0 {
		// 	commentSuffix = " // " + i.GetComment()
		// }
		switch v := i.(type) {
		case SimpleOp:
			// fmt.Printf("0x%-6.2X 0x%X    (%s)%s\n", pc, v.code, codeToString[v.code], commentSuffix)
			out = append(out, byte(v.code))
		case PushByte:
			// fmt.Printf("0x%-6.2X 0x%X %2d (PUSH1)%s\n", pc, PUSH1, v.value, commentSuffix)
			out = append(out, byte(PUSH1))
			out = append(out, byte(v.value))
		case PushLabel:
			pos, ok := labels[v.label]
			if !ok {
				panic("Label not found: " + v.label)
			}
			// fmt.Printf("0x%-6.2X 0x%X %2d (PUSH2 @%s)%s\n", pc, PUSH2, pos, v.label, commentSuffix)
			out = append(out, byte(PUSH2))
			out = append(out, byte(pos/256))
			out = append(out, byte(pos%256))
		case LabelOp:
			// fmt.Printf("%s:\n", v.label)
			// fmt.Printf("0x%-6.2X 0x%X    (JUMPDEST)%s\n", pc, JUMPDEST, commentSuffix)
			out = append(out, byte(JUMPDEST))
		case PushBytes:
			op := OpCode(uint8(PUSH1) + uint8(len(v.bytes)) - 1)
			// fmt.Printf("0x%-6.2X 0x%X 0x%X%s\n", pc, op, v.bytes, commentSuffix)
			out = append(out, byte(op))
			for _, b := range v.bytes {
				out = append(out, b)
			}
		case Comment:
			// fmt.Printf("// %s\n", v.GetComment())
		default:
			panic("Unknown instruction type.")
		}
		pc += i.Length()
	}

	return out
}
