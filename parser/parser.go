package parser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"smart/scanner"
	"smart/token"
)

type Parser struct {
	tok  token.Token
	pos  token.Pos
	lit  string
	scan *scanner.Scanner

	indent  int
	verbose bool

	path string
}

func New(s *scanner.Scanner, path string) *Parser {
	var p Parser
	abspath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	p.path = abspath
	p.scan = s
	p.pos, p.tok, p.lit = s.Scan()
	return &p
}

func (p *Parser) ResolvePath(relpath string) string {
	return filepath.Join(filepath.Dir(p.path), relpath)
}

func (p *Parser) expect(tok token.Token) {
	if tok != p.tok {
		msg := fmt.Sprintf("expecting %s; found %s", token.Tokens[tok], p.lit)
		p.err(p.pos, msg)
	}
	// p.lexeme()
	p.pos, p.tok, p.lit = p.scan.Scan()
}

func (p *Parser) GetPos() token.Position {
	line, col := p.scan.Coordinates(p.pos)
	return token.Position{p.path, int(p.pos), line, col}
}

func (p *Parser) err(pos token.Pos, str string) {
	line, col := p.scan.Coordinates(pos)
	fmt.Fprintf(os.Stderr, "%s line %d, column %d: %s\n", p.path, line, col, str)
	os.Exit(1)
}
