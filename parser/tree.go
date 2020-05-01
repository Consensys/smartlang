package parser

import (
	"fmt"
	"smart/token"
)

func (p *Parser) printIndent() {
	for i := 0; i < p.indent; i++ {
		fmt.Print(". ")
	}
}

func (p *Parser) enter(element string) {
	if p.verbose {
		p.printIndent()
		fmt.Println("[", element)
		p.indent += 1
	}
}

func (p *Parser) leave(element string) {
	if p.verbose {
		p.indent -= 1
		p.printIndent()
		fmt.Println("]", element)
	}
}

func (p *Parser) lexeme() {
	if p.verbose {
		p.printIndent()
		tok := token.Tokens[p.tok]
		fmt.Print("lexeme: ", tok)
		if tok != p.lit {
			fmt.Print("[", p.lit, "]")
		}
		fmt.Println()
	}
}
