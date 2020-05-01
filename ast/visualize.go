package ast

import (
	"fmt"
	"log"
	"reflect"
	"smart/token"
	"strings"
)

type Edge struct {
	from int
	to   int
}

type Node struct {
	label string
}

type Dot struct {
	nodes []Node
	edges []Edge
	stack []int
}

func (a *ASTNode) DumpSymbols() string {
	if a.ctx != nil && a.ctx.symtab != nil {
		keys := make([]string, 0, len(a.ctx.symtab.T))
		for k := range a.ctx.symtab.T {
			keys = append(keys, k)
		}
		return " | Symbols: " + strings.Join(keys, ", ")
	} else {
		return ""
	}
}

func (d *Dot) EnterNode(visitable IVisitor) {
	next := len(d.nodes)
	if next > 0 {
		d.edges = append(d.edges, Edge{d.stack[len(d.stack)-1], next})
	}
	d.stack = append(d.stack, next)

	var label string
	switch v := visitable.(type) {
	case *Program:
		label = "Program"
	case *Contract:
		label = fmt.Sprintf("Contract | %s", v.id.name)
	case *FuncSig:
		label = fmt.Sprintf("FuncSig | %s", v.id.name)
	case *FallbackDecl:
		label = "Fallback"
	case *Constructor:
		label = "Constructor"
	case *PoolDecl:
		label = fmt.Sprintf("Pool | %s", v.id.name)
	case *Param:
		label = fmt.Sprintf("Param | %s | indexed=%t", v.id.name, v.indexed)
	case *StorageDecl:
		label = fmt.Sprintf("StorageDecl | %s", v.id.name)
	case *ConstDecl:
		label = fmt.Sprintf("ConstDecl | %s", v.id.name)
	case *FuncDef:
		label = "FuncDef"
	case *CodeBlock:
		label = "CodeBlock"
	case *Class:
		label = fmt.Sprintf("Class | %s", v.id.name)
	case *ImportDecl:
		label = fmt.Sprintf("ImportDecl | path=\"%s\"", v.path)
		if len(v.hash) > 0 {
			label += fmt.Sprintf(" | hash=0x%x", v.hash)
		}
		if len(v.alias.name) > 0 {
			label += fmt.Sprintf(" | as \"%s\"", v.alias.name)
		}
	case *Interface:
		label = fmt.Sprintf("Interface | %s", v.id.name)
	case *AssertStmt:
		label = fmt.Sprintf("Assert | reason=\"%s\"", strings.ReplaceAll(v.reason, "\"", "\\\\\""))
	case *QualIdent:
		label = "QualIdent | "
		if v.qualifier != nil {
			label += v.qualifier.name + "."
		}
		label += v.id.name
	case *StringLiteral:
		label = fmt.Sprintf("StringLiteral | %s", v.lit)
	case *ArrayType:
		label = "ArrayType"
	case *EventDecl:
		label = fmt.Sprintf("EventDecl | %s", v.id.name)
	case *VarDecl:
		label = fmt.Sprintf("VarDecl | %s", v.id.name)
	case *ExternalCall:
		label = "ExternalCall"
	case *ReturnStmt:
		label = "ReturnStmt"
	case *Argument:
		label = "Argument"
		if v.id != nil {
			label += " | " + v.id.name
		}
	case *Binary:
		label = fmt.Sprintf("Binary \\%s", token.Tokens[v.op])
	case *Variable:
		label = fmt.Sprintf("Variable | %s", v.id.name)
	case *CallExpr:
		label = "CallExpr"
	case *DecLiteral:
		label = fmt.Sprintf("DecLiteral | %s", v.lit)
	case *HexLiteral:
		label = fmt.Sprintf("HexLiteral | %s", v.lit)
	case *BytesLiteral:
		label = fmt.Sprintf("BytesLiteral | %s", v.lit)
	case *FieldAccess:
		label = fmt.Sprintf("FieldAccess | %s", v.field.name)
	case *NewExpr:
		label = "NewExpr"
	case *ArrayAccess:
		label = "ArrayAccess"
	case *Assignment:
		label = "Assignment"
	case *Unary:
		label = fmt.Sprintf("Unary %s", token.Tokens[v.op])
	case *BoolLiteral:
		label = fmt.Sprintf("BoolLiteral | %t", v.value)
	case *EmitStmt:
		label = fmt.Sprintf("EmitStmt | %s", v.id.name)
	case *AugmentedAssignment:
		label = fmt.Sprintf("AugmentedAssignment %s", token.Tokens[v.op])
	case *Move:
		label = "Move"
	case *AsExpr:
		label = "AsExpr"
	case *FunctionCall:
		label = "FunctionCall"
	case *IfStmt:
		label = "IfStmt"
	case *MapType:
		label = "MapType"
	case *ForStmt:
		label = "ForStmt"
	case *WhileStmt:
		label = "WhileStmt"
	case *Mint:
		label = "Mint"
	default:
		log.Fatal("Unimplemented: " + reflect.TypeOf(v).String())
	}
	node := visitable.(IDumper)
	label += node.DumpSymbols()

	d.nodes = append(d.nodes, Node{label})
}

func (d *Dot) ExitNode(v IVisitor) {
	d.stack = d.stack[:len(d.stack)-1]
}

func NewDot() *Dot {
	return &Dot{make([]Node, 0), make([]Edge, 0), make([]int, 0)}
}

func (d *Dot) Print() {
	fmt.Println("digraph ast {")
	fmt.Println("graph [ rankdir = \"LR\" ]")
	for i, node := range d.nodes {
		fmt.Printf("    node%d [ label = \"%s\", shape = \"record\" ];\n", i, strings.ReplaceAll(node.label, "\"", "\\\""))
	}
	for _, edge := range d.edges {
		fmt.Printf("    node%d -> node%d;\n", edge.from, edge.to)
	}
	fmt.Println("}")
}

func Visualize(v IVisitor) {
	d := NewDot()
	v.visit(d.EnterNode, d.ExitNode)
	d.Print()
}
