package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"smart/ast"
	"smart/parser"
	"smart/scanner"
)

var path string
var typecheck bool

func init() {
	flag.StringVar(&path, "path", "", "Path to .smrt file to parse.")
	flag.BoolVar(&typecheck, "t", false, "Call ast.Check before walking the tree.")
}

func main() {
	flag.Parse()
	if len(path) == 0 {
		log.Fatal("Path argument is required.")
	}
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	s := scanner.New(r)
	p := parser.New(s, path)

	prog := p.Parse_prog()
	if typecheck {
		ast.Check(prog)
	}
	ast.Visualize(prog)
}
