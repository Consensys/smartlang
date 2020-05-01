package parser

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"smart/ast"
	"smart/scanner"
	"smart/token"
	"strings"
)

// prog = { contract | class | importdecl | interface } EOF
func (p *Parser) Parse_prog() *ast.Program {
	var decls []ast.IProgramBodyDecl

	for p.tok == token.IMPORT || p.tok == token.INTERFACE || p.tok == token.CLASS || p.tok == token.CONTRACT {
		switch p.tok {
		case token.CONTRACT:
			decls = append(decls, p.parse_contract())
		case token.CLASS:
			decls = append(decls, p.parse_class())
		case token.IMPORT:
			decls = append(decls, p.parse_importdecl())
		case token.INTERFACE:
			decls = append(decls, p.parse_interface())
		default:
			p.err(p.pos, "syntax error") // PARSE ERROR
		}
	}
	p.expect(token.EOF)
	ret := ast.NewProgram(decls)
	return ret
}

// importdecl = IMPORT [ STRINGLITERAL ] [ LBRACK HEXLITERAL RBRACK ] [ AS IDENT ] SEMI
func (p *Parser) parse_importdecl() *ast.ImportDecl {
	var path string
	var alias *ast.Identifier

	// TODO check for duplicate defintions
	p.expect(token.IMPORT)
	pos := p.GetPos()
	if p.tok == token.STRINGLITERAL {
		path = p.lit[1 : len(p.lit)-1]
		p.expect(token.STRINGLITERAL)
	}

	path = p.ResolvePath(path)
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open(absPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	hash := h.Sum(nil)

	f.Seek(0, 0)
	r := bufio.NewReader(f)
	s := scanner.New(r)
	imported := New(s, path)
	prog := imported.Parse_prog()

	if p.tok == token.LBRACK {
		p.expect(token.LBRACK)
		expectedHash, err := hex.DecodeString(p.lit[2:])
		if err != nil {
			log.Fatal(err)
		}
		if !reflect.DeepEqual(expectedHash, hash) {
			p.err(p.pos, fmt.Sprintf("Imported file hash mismatch: expected 0x%x but computed 0x%x\n", expectedHash, hash))
		}
		p.expect(token.HEXLITERAL)
		p.expect(token.RBRACK)
	}
	// TODO: Make alias optional
	p.expect(token.AS)
	alias = ast.NewIdentifier(p.lit, p.GetPos())
	p.expect(token.IDENT)

	p.expect(token.SEMI)

	ret := ast.NewImportDecl(path, hash, alias, prog)
	ret.SetCoords(pos)
	return ret
}

// contract = CONTRACT IDENT contractBody
func (p *Parser) parse_contract() *ast.Contract {
	p.expect(token.CONTRACT)
	pos := p.GetPos()
	id := ast.NewIdentifier(p.lit, pos)
	p.expect(token.IDENT)
	decls := p.parse_contractBody()
	ret := ast.NewContract(id, decls)
	ret.SetCoords(pos)
	return ret
}

// contractBody = LCURL { classStuff | fallback } RCURL
func (p *Parser) parse_contractBody() []ast.IDecl {
	decls := make([]ast.IDecl, 0)
	p.expect(token.LCURL)
	for p.tok == token.STORAGE || p.tok == token.FALLBACK || p.tok == token.EVENT || p.tok == token.FUNC || p.tok == token.CONST || p.tok == token.POOL || p.tok == token.AT || p.tok == token.CONSTRUCTOR {
		var decl ast.IDecl
		switch p.tok {
		case token.STORAGE, token.EVENT, token.FUNC, token.CONST, token.POOL, token.AT, token.CONSTRUCTOR:
			decl = p.parse_classStuff()
		case token.FALLBACK:
			decl = p.parse_fallback()
		default:
			p.err(p.pos, "syntax error") // PARSE ERROR
		}
		decls = append(decls, decl)
	}
	p.expect(token.RCURL)
	p.expect(token.SEMI)
	return decls
}

// classStuff = constructor | pooldecl | storagedecl | constdecl SEMI | funcdef | eventDecl
func (p *Parser) parse_classStuff() ast.IDecl {
	var d ast.IDecl
	switch p.tok {
	case token.CONSTRUCTOR:
		d = p.parse_constructor()
	case token.POOL:
		d = p.parse_pooldecl()
	case token.STORAGE:
		d = p.parse_storagedecl()
	case token.CONST:
		d = p.parse_constdecl()
		p.expect(token.SEMI)
	case token.FUNC, token.AT:
		d = p.parse_funcdef()
	case token.EVENT:
		d = p.parse_eventDecl()
	default:
		p.err(p.pos, "syntax error") // PARSE ERROR
	}
	return d
}

// fallback = FALLBACK codeblock
func (p *Parser) parse_fallback() *ast.FallbackDecl {
	pos := p.GetPos()
	p.expect(token.FALLBACK)
	c := p.parse_codeblock()
	p.expect(token.SEMI)
	f := ast.NewFallbackDecl(c)
	f.SetCoords(pos)
	return f
}

// eventDecl = EVENT IDENT LPAREN [ IDENT type [ INDEXED ] { COMMA IDENT type [ INDEXED ] } ] RPAREN SEMI

func (p *Parser) parse_eventParam() *ast.Param {
	pId := ast.NewIdentifier(p.lit, p.GetPos())
	p.expect(token.IDENT)
	typ := p.parse_type()
	indexed := false
	if p.tok == token.INDEXED {
		indexed = true
		p.expect(token.INDEXED)
	}
	param := ast.NewParam(pId, typ, indexed)
	return param
}

func (p *Parser) parse_eventDecl() *ast.EventDecl {
	p.expect(token.EVENT)
	pos := p.GetPos()
	eId := ast.NewIdentifier(p.lit, pos)
	p.expect(token.IDENT)
	p.expect(token.LPAREN)
	params := make([]*ast.Param, 0)
	if p.tok == token.IDENT {
		param := p.parse_eventParam()
		params = append(params, param)
		for p.tok == token.COMMA {
			p.expect(token.COMMA)
			param = p.parse_eventParam()
			params = append(params, param)
		}
	}
	p.expect(token.RPAREN)
	p.expect(token.SEMI)
	ret := ast.NewEventDecl(eId, params)
	ret.SetCoords(pos)
	return ret
}

// constructor = CONSTRUCTOR LPAREN [ paramList ] RPAREN codeblock
func (p *Parser) parse_constructor() *ast.Constructor {
	pos := p.GetPos()
	p.expect(token.CONSTRUCTOR)
	p.expect(token.LPAREN)
	var params []*ast.Param
	if p.tok == token.IDENT {
		params = p.parse_paramList()
	}
	p.expect(token.RPAREN)
	code := p.parse_codeblock()
	p.expect(token.SEMI)
	ret := ast.NewConstructor(params, code)
	ret.SetCoords(pos)
	return ret
}

// paramList = IDENT type { COMMA IDENT type }
func (p *Parser) parse_param() *ast.Param {
	pos := p.GetPos()
	id := ast.NewIdentifier(p.lit, pos)
	p.expect(token.IDENT)
	typ := p.parse_type()
	ret := ast.NewParam(id, typ, false)
	ret.SetCoords(pos)
	return ret
}

func (p *Parser) parse_paramList() []*ast.Param {
	params := make([]*ast.Param, 0)
	param := p.parse_param()
	params = append(params, param)
	for p.tok == token.COMMA {
		p.expect(token.COMMA)
		param = p.parse_param()
		params = append(params, param)
	}
	return params
}

// class = CLASS IDENT classBody
func (p *Parser) parse_class() *ast.Class {
	p.expect(token.CLASS)
	pos := p.GetPos()
	id := ast.NewIdentifier(p.lit, pos)
	p.expect(token.IDENT)
	decls := p.parse_classBody()
	ret := ast.NewClass(id, decls)
	ret.SetCoords(pos)
	return ret
}

// classBody = LCURL { classStuff } RCURL
func (p *Parser) parse_classBody() []ast.IClassBodyDecl {
	stuff := make([]ast.IClassBodyDecl, 0)
	p.expect(token.LCURL)
	for p.tok == token.STORAGE || p.tok == token.EVENT || p.tok == token.FUNC || p.tok == token.CONST || p.tok == token.POOL || p.tok == token.AT || p.tok == token.CONSTRUCTOR {
		s := p.parse_classStuff()
		stuff = append(stuff, s)
	}
	p.expect(token.RCURL)
	p.expect(token.SEMI)
	return stuff
}

// pooldecl = POOL IDENT SEMI
func (p *Parser) parse_pooldecl() *ast.PoolDecl {
	p.expect(token.POOL)
	pos := p.GetPos()
	id := ast.NewIdentifier(p.lit, pos)
	p.expect(token.IDENT)
	p.expect(token.SEMI)
	ret := ast.NewPoolDecl(id)
	ret.SetCoords(pos)
	return ret
}

// vardecl = VAR IDENT type [EQ expr]
func (p *Parser) parse_vardecl() *ast.VarDecl {
	p.expect(token.VAR)
	pos := p.GetPos()
	id := ast.NewIdentifier(p.lit, pos)
	p.expect(token.IDENT)
	typ := p.parse_type()
	var value ast.IExpr
	if p.tok == token.EQ {
		p.expect(token.EQ)
		value = p.parse_expr()
	}
	ret := ast.NewVarDecl(id, typ, value)
	ret.SetCoords(pos)
	return ret
}

// storagedecl = STORAGE IDENT type SEMI
func (p *Parser) parse_storagedecl() *ast.StorageDecl {
	p.expect(token.STORAGE)
	pos := p.GetPos()
	id := ast.NewIdentifier(p.lit, pos)
	p.expect(token.IDENT)
	typ := p.parse_type()
	p.expect(token.SEMI)
	ret := ast.NewStorageDecl(id, typ)
	ret.SetCoords(pos)
	return ret
}

// constdecl = CONST IDENT type EQ expr
func (p *Parser) parse_constdecl() *ast.ConstDecl {
	p.expect(token.CONST)
	pos := p.GetPos()
	id := ast.NewIdentifier(p.lit, pos)
	p.expect(token.IDENT)
	typ := p.parse_type()
	p.expect(token.EQ)
	value := p.parse_expr()
	ret := ast.NewConstDecl(id, typ, value)
	ret.SetCoords(pos)
	return ret
}

// type = qualident | arraytype | maptype
func (p *Parser) parse_type() ast.ITypeAST {
	var a ast.ITypeAST
	switch p.tok {
	case token.IDENT:
		a = p.parse_qualident()
	case token.LBRACK:
		a = p.parse_arraytype()
	case token.MAP:
		a = p.parse_maptype()
	default:
		p.err(p.pos, "syntax error") // PARSE ERROR
	}
	return a
}

// qualident = IDENT [ DOT IDENT ]
func (p *Parser) parse_qualident() *ast.QualIdent {
	first := ast.NewIdentifier(p.lit, p.GetPos())
	p.expect(token.IDENT)
	if p.tok == token.DOT {
		p.expect(token.DOT)
		second := ast.NewIdentifier(p.lit, p.GetPos())
		p.expect(token.IDENT)
		return ast.NewQualIdent(second, first) // really!
	} else {
		return ast.NewQualIdent(nil, first) // really!
	}
}

// arraytype = LBRACK [ expr ] RBRACK type
func (p *Parser) parse_arraytype() *ast.ArrayType {
	p.expect(token.LBRACK)
	var e ast.IExpr
	if p.tok == token.LPAREN || p.tok == token.MINUS || p.tok == token.STRINGLITERAL || p.tok == token.BANG || p.tok == token.FALSE || p.tok == token.MINT || p.tok == token.TRUE || p.tok == token.HEXLITERAL || p.tok == token.NEW || p.tok == token.LARR || p.tok == token.IDENT || p.tok == token.BYTESLITERAL || p.tok == token.LBRACK || p.tok == token.BURN || p.tok == token.DECLITERAL {
		e = p.parse_expr()
	}
	p.expect(token.RBRACK)
	typ := p.parse_type()
	ret := ast.NewArrayType(e, typ)
	return ret
}

// maptype = MAP LBRACK type RBRACK type
func (p *Parser) parse_maptype() *ast.MapType {
	p.expect(token.MAP)
	p.expect(token.LBRACK)
	keytype := p.parse_type()
	p.expect(token.RBRACK)
	valuetype := p.parse_type()
	ret := ast.NewMapType(keytype, valuetype)
	return ret
}

// expr  = expr00 { LOGOR expr00 }
func (p *Parser) parse_expr() ast.IExpr {
	e := p.parse_expr00()
	for p.tok == token.LOGOR {
		op := p.tok
		p.expect(token.LOGOR)
		right := p.parse_expr00()
		e = ast.NewBinary(op, e, right)
	}
	return e
}

// expr00 = expr0 { LOGAND expr0 }
func (p *Parser) parse_expr00() ast.IExpr {
	e := p.parse_expr0()
	for p.tok == token.LOGAND {
		op := p.tok
		p.expect(token.LOGAND)
		right := p.parse_expr0()
		e = ast.NewBinary(op, e, right)
	}
	return e
}

// expr0 = expr1 [ relation expr1 ]
func (p *Parser) parse_expr0() ast.IExpr {
	e := p.parse_expr1()
	if p.tok == token.GEQ || p.tok == token.LEQ || p.tok == token.LT || p.tok == token.GT || p.tok == token.EQEQ || p.tok == token.NEQ {
		op := p.parse_relation()
		right := p.parse_expr1()
		e = ast.NewBinary(op, e, right)
	}
	return e
}

// expr1 = expr2 { addop expr2 }
func (p *Parser) parse_expr1() ast.IExpr {
	e := p.parse_expr2()
	for p.tok == token.MINUS || p.tok == token.BITAND || p.tok == token.LSHIFT || p.tok == token.PLUS || p.tok == token.BITOR || p.tok == token.RSHIFT {
		op := p.parse_addop()
		right := p.parse_expr2()
		e = ast.NewBinary(op, e, right)
	}
	return e
}

// expr2 = expr3 { mulop expr3 }
func (p *Parser) parse_expr2() ast.IExpr {
	e := p.parse_expr3()
	for p.tok == token.DIV || p.tok == token.STAR {
		op := p.parse_mulop()
		right := p.parse_expr3()
		e = ast.NewBinary(op, e, right)
	}
	return e
}

// expr3 = expr4 { EXP expr4 }
func (p *Parser) parse_expr3() ast.IExpr {
	e := p.parse_expr4()
	for p.tok == token.EXP { // is this associative correct?
		op := p.tok
		p.expect(token.EXP)
		right := p.parse_expr4()
		e = ast.NewBinary(op, e, right)
	}
	return e
}

// expr4 = expr5 [ AS type ]
func (p *Parser) parse_expr4() ast.IExpr {
	e := p.parse_expr5()
	if p.tok == token.AS {
		p.expect(token.AS)
		pos := p.GetPos()
		typ := p.parse_type()
		ret := ast.NewAsExpr(e, typ)
		ret.SetCoords(pos)
		return ret
	} else {
		return e
	}
}

// expr5 = designator [ LPAREN [ argList ] RPAREN ] | literal | LPAREN expr RPAREN | MINUS expr5 | LARR expr5 | new | MINT BANG expr5 | BURN BANG expr5 | BANG expr5
func (p *Parser) parse_expr5() ast.IExpr {
	var e ast.IExpr
	var tmp ast.IExpr
	var op token.Token
	switch p.tok {
	case token.IDENT:
		e = p.parse_designator()
		if p.tok == token.LPAREN {
			p.expect(token.LPAREN)
			var args []*ast.Argument
			if p.tok == token.IDENT {
				args = p.parse_argList()
			}
			p.expect(token.RPAREN)
			e = ast.NewCallExpr(e, args)
		}
	case token.STRINGLITERAL, token.FALSE, token.TRUE, token.HEXLITERAL, token.BYTESLITERAL, token.LBRACK, token.DECLITERAL:
		e = p.parse_literal()
	case token.LPAREN:
		p.expect(token.LPAREN)
		e = p.parse_expr()
		p.expect(token.RPAREN)
	case token.MINUS:
		op = p.tok
		p.expect(token.MINUS)
		tmp = p.parse_expr5()
		e = ast.NewUnary(op, tmp)
	case token.LARR:
		op = p.tok
		p.expect(token.LARR)
		tmp = p.parse_expr5()
		e = ast.NewUnary(op, tmp)
	case token.NEW:
		e = p.parse_new()
	case token.MINT:
		op = p.tok
		p.expect(token.MINT)
		p.expect(token.BANG)
		tmp = p.parse_expr5()
		typ := p.parse_type()
		e = ast.NewMint(tmp, typ)
	case token.BURN:
		op = p.tok
		p.expect(token.BURN)
		p.expect(token.BANG)
		tmp = p.parse_expr5()
		e = ast.NewUnary(op, tmp)
	case token.BANG:
		op = p.tok
		p.expect(token.BANG)
		tmp = p.parse_expr5()
		e = ast.NewUnary(op, tmp)
	default:
		p.err(p.pos, "syntax error") // PARSE ERROR
	}
	return e
}

// designator = IDENT [ BANG ] { DOT IDENT [ BANG ] | LBRACK expr RBRACK }
func (p *Parser) parse_designator() ast.IDesignator {
	name := p.lit
	pos := p.GetPos()
	p.expect(token.IDENT)
	if p.tok == token.BANG {
		name += p.lit
		p.expect(token.BANG)
	}
	id := ast.NewIdentifier(name, p.GetPos())
	var e ast.IDesignator
	e = ast.NewVariable(id)
	e.SetCoords(pos)
	for p.tok == token.DOT || p.tok == token.LBRACK {
		switch p.tok {
		case token.DOT:
			p.expect(token.DOT)
			name = p.lit
			pos := p.GetPos()
			p.expect(token.IDENT)
			if p.tok == token.BANG {
				name += p.lit
				p.expect(token.BANG)
			}
			field := ast.NewIdentifier(name, p.GetPos())
			e = ast.NewFieldAccess(e, field)
			e.SetCoords(pos)
		case token.LBRACK:
			p.expect(token.LBRACK)
			pos := p.GetPos()
			index := p.parse_expr()
			p.expect(token.RBRACK)
			e = ast.NewArrayAccess(e, index)
			e.SetCoords(pos)
		default:
			p.err(p.pos, "syntax error") // PARSE ERROR
		}
	}
	return e
}

// arrayLiteral = LBRACK RBRACK type LCURL [ exprList ] RCURL
func (p *Parser) parse_arrayLiteral() *ast.ArrayLiteral { // should be array literal type
	var exprList []ast.IExpr
	p.expect(token.LBRACK)
	p.expect(token.RBRACK)
	typ := p.parse_type()
	p.expect(token.LCURL)
	if p.tok == token.LPAREN || p.tok == token.MINUS || p.tok == token.STRINGLITERAL || p.tok == token.BANG || p.tok == token.FALSE || p.tok == token.MINT || p.tok == token.TRUE || p.tok == token.HEXLITERAL || p.tok == token.NEW || p.tok == token.LARR || p.tok == token.IDENT || p.tok == token.BYTESLITERAL || p.tok == token.LBRACK || p.tok == token.BURN || p.tok == token.DECLITERAL {
		exprList = p.parse_exprList()
	}
	p.expect(token.RCURL)
	ret := ast.NewArrayLiteral(typ, exprList)
	return ret
}

// literal = DECLITERAL | HEXLITERAL | BYTESLITERAL | TRUE | FALSE | arrayLiteral | STRINGLITERAL
func (p *Parser) parse_literal() ast.IExpr {
	var e ast.IExpr
	pos := p.GetPos()
	switch p.tok {
	case token.DECLITERAL:
		v := p.lit
		e = ast.NewDecLiteral(v)
		p.expect(token.DECLITERAL)
	case token.HEXLITERAL:
		v := p.lit
		e = ast.NewHexLiteral(v)
		p.expect(token.HEXLITERAL)
	case token.BYTESLITERAL:
		v := p.lit
		e = ast.NewBytesLiteral(v)
		p.expect(token.BYTESLITERAL)
	case token.TRUE:
		e = ast.NewBoolLiteral(true)
		p.expect(token.TRUE)
	case token.FALSE:
		e = ast.NewBoolLiteral(false)
		p.expect(token.FALSE)
	case token.LBRACK:
		e = p.parse_arrayLiteral()
	case token.STRINGLITERAL:
		v := p.lit
		e = ast.NewStringLiteral(v)
		p.expect(token.STRINGLITERAL)
	default:
		p.err(p.pos, "syntax error") // PARSE ERROR
	}
	e.SetCoords(pos)
	return e
}

// relation = EQEQ | NEQ | LT | LEQ | GT | GEQ
func (p *Parser) parse_relation() token.Token {
	op := p.tok
	switch p.tok {
	case token.EQEQ:
		p.expect(token.EQEQ)
	case token.NEQ:
		p.expect(token.NEQ)
	case token.LT:
		p.expect(token.LT)
	case token.LEQ:
		p.expect(token.LEQ)
	case token.GT:
		p.expect(token.GT)
	case token.GEQ:
		p.expect(token.GEQ)
	default:
		p.err(p.pos, "syntax error") // PARSE ERROR
	}
	return op
}

// addop = PLUS | MINUS | BITAND | BITOR | LSHIFT | RSHIFT
func (p *Parser) parse_addop() token.Token {
	op := p.tok
	switch p.tok {
	case token.PLUS:
		p.expect(token.PLUS)
	case token.MINUS:
		p.expect(token.MINUS)
	case token.BITAND:
		p.expect(token.BITAND)
	case token.BITOR:
		p.expect(token.BITOR)
	case token.LSHIFT:
		p.expect(token.LSHIFT)
	case token.RSHIFT:
		p.expect(token.RSHIFT)
	default:
		p.err(p.pos, "syntax error") // PARSE ERROR
	}
	return op
}

// mulop = STAR | DIV
func (p *Parser) parse_mulop() token.Token {
	op := p.tok
	switch p.tok {
	case token.STAR:
		p.expect(token.STAR)
	case token.DIV:
		p.expect(token.DIV)
	default:
		p.err(p.pos, "syntax error") // PARSE ERROR
	}
	return op
}

// new = NEW qualident LPAREN [ argList ] RPAREN
func (p *Parser) parse_new() *ast.NewExpr {
	p.expect(token.NEW)
	qual := p.parse_qualident()
	var args []*ast.Argument
	p.expect(token.LPAREN)
	if p.tok == token.IDENT {
		args = p.parse_argList()
	}
	p.expect(token.RPAREN)
	ret := ast.NewNewExpr(qual, args)
	return ret
}

// exprList = expr { COMMA expr }
func (p *Parser) parse_exprList() []ast.IExpr {
	list := make([]ast.IExpr, 1)
	list[0] = p.parse_expr()
	for p.tok == token.COMMA {
		p.expect(token.COMMA)
		e := p.parse_expr()
		list = append(list, e)
	}
	return list
}

// designatorList = designator { COMMA designator }
func (p *Parser) parse_designatorList() (list []ast.IDesignator) {
	list = make([]ast.IDesignator, 1)
	list[0] = p.parse_designator()
	for p.tok == token.COMMA {
		p.expect(token.COMMA)
		d := p.parse_designator()
		list = append(list, d)
	}
	return list
}

// simpleStatement = assignOrMoveOrCall | vardecl | constdecl | store | assert | call | emit
func (p *Parser) parse_simpleStatement() ast.IStatement {
	var s ast.IStatement
	switch p.tok {
	case token.IDENT:
		s = p.parse_assignOrMoveOrCall()
	case token.VAR:
		s = p.parse_vardecl()
	case token.CONST:
		s = p.parse_constdecl()
	case token.STORE:
		s = p.parse_store()
	case token.ASSERT:
		s = p.parse_assert()
	case token.CALL:
		s = p.parse_call()
	case token.EMIT:
		s = p.parse_emit()
	default:
		p.err(p.pos, "syntax error") // PARSE ERROR
	}
	return s
}

// complexStatement = ifstmt | loop | codeblock | ret
func (p *Parser) parse_complexStatement() (s ast.IStatement) {
	switch p.tok {
	case token.IF:
		s = p.parse_ifstmt()
	case token.FOR, token.WHILE:
		s = p.parse_loop()
	case token.LCURL:
		s = p.parse_codeblock()
		p.expect(token.SEMI)
	case token.RETURN:
		s = p.parse_ret()
	default:
		p.err(p.pos, "syntax error") // PARSE ERROR
	}
	return s
}

// emit = EMIT IDENT LPAREN [ argList ] RPAREN
func (p *Parser) parse_emit() *ast.EmitStmt {
	p.expect(token.EMIT)
	pos := p.GetPos()
	id := ast.NewIdentifier(p.lit, pos)
	p.expect(token.IDENT)
	p.expect(token.LPAREN)
	var args []*ast.Argument
	if p.tok == token.IDENT {
		args = p.parse_argList()
	}
	p.expect(token.RPAREN)
	ret := ast.NewEmitStmt(id, args)
	ret.SetCoords(pos)
	return ret
}

// ret = RETURN [ expr | LCURL argList RCURL ] SEMI
func (p *Parser) parse_ret() *ast.ReturnStmt {
	var argList []*ast.Argument
	pos := p.GetPos()
	p.expect(token.RETURN)
	if p.tok == token.LPAREN || p.tok == token.MINUS || p.tok == token.STRINGLITERAL || p.tok == token.LCURL || p.tok == token.BANG || p.tok == token.FALSE || p.tok == token.MINT || p.tok == token.TRUE || p.tok == token.HEXLITERAL || p.tok == token.NEW || p.tok == token.LARR || p.tok == token.IDENT || p.tok == token.BYTESLITERAL || p.tok == token.LBRACK || p.tok == token.BURN || p.tok == token.DECLITERAL {
		switch p.tok {
		case token.LPAREN, token.MINUS, token.STRINGLITERAL, token.BANG, token.FALSE, token.MINT, token.TRUE, token.HEXLITERAL, token.NEW, token.LARR, token.IDENT, token.BYTESLITERAL, token.LBRACK, token.BURN, token.DECLITERAL:
			pos := p.GetPos()
			expr := p.parse_expr()
			arg := ast.NewArgument(nil, expr)
			arg.SetCoords(pos)
			argList = append(argList, arg)
		case token.LCURL:
			p.expect(token.LCURL)
			argList = p.parse_argList()
			p.expect(token.RCURL)
		default:
			p.err(p.pos, "syntax error") // PARSE ERROR
		}
	}
	p.expect(token.SEMI)

	ret := ast.NewReturnStmt(argList)
	ret.SetCoords(pos)
	return ret
}

// assignOrMoveOrCall = designatorList ( EQ exprList [ FROM expr ] | LARR expr [ FROM expr ] | LPAREN [ argList ] RPAREN | ( PLUSEQ | MINUSEQ ) expr )
func (p *Parser) parse_assignOrMoveOrCall() ast.IStatement {
	var stmt ast.IStatement
	pos := p.GetPos()
	lhsPos := p.pos
	lhs := p.parse_designatorList()
	switch p.tok {
	case token.EQ:
		pos = p.GetPos()
		p.expect(token.EQ)
		rhs := p.parse_exprList()
		var from *ast.CallExpr
		if p.tok == token.FROM {
			p.expect(token.FROM)
			pos := p.pos
			callExpr, ok := p.parse_expr().(*ast.CallExpr)
			if !ok {
				p.err(pos, "from expression must be a function call")
			}
			from = callExpr
		}
		stmt = ast.NewAssignment(lhs, rhs, from)
	case token.LARR:
		if len(lhs) != 1 {
			p.err(p.pos, "too many values on the left-hand side of a move")
		}
		pos = p.GetPos()
		p.expect(token.LARR)
		amount := p.parse_expr()
		var from ast.IExpr
		if p.tok == token.FROM {
			p.expect(token.FROM)
			from = p.parse_expr()
		} else {
			from = amount
			amount = nil
		}
		stmt = ast.NewMove(lhs[0], amount, from)
	case token.LPAREN:
		if len(lhs) != 1 {
			p.err(p.pos, "cannot call more than one target")
		}
		p.expect(token.LPAREN)
		arguments := make([]*ast.Argument, 0)
		if p.tok == token.IDENT {
			arguments = p.parse_argList()
		}
		p.expect(token.RPAREN)
		stmt = ast.NewFunctionCall(lhs[0], arguments)
	case token.PLUSEQ, token.MINUSEQ:
		pos = p.GetPos()
		op := p.tok
		switch p.tok {
		case token.PLUSEQ:
			p.expect(token.PLUSEQ)
		case token.MINUSEQ:
			p.expect(token.MINUSEQ)
		default:
			p.err(p.pos, "syntax error") // PARSE ERROR
		}
		rhs := p.parse_expr()
		if len(lhs) != 1 {
			p.err(lhsPos, "left-hand side of an augmented assignment must be a single lvalue")
		}
		stmt = ast.NewAugmentedAssignment(lhs[0], op, rhs)
	default:
		p.err(p.pos, "syntax error") // PARSE ERROR
	}
	stmt.SetCoords(pos)
	return stmt
}

// argList = arg { COMMA arg }
func (p *Parser) parse_argList() (list []*ast.Argument) {
	list = make([]*ast.Argument, 1)
	pos := p.GetPos()
	arg := p.parse_arg()
	arg.SetCoords(pos)
	list[0] = arg
	for p.tok == token.COMMA {
		p.expect(token.COMMA)
		pos = p.GetPos()
		arg = p.parse_arg()
		arg.SetCoords(pos)
		list = append(list, arg)
	}
	return list
}

// arg = IDENT COLON expr
func (p *Parser) parse_arg() *ast.Argument {
	pos := p.GetPos()
	id := ast.NewIdentifier(p.lit, pos)
	p.expect(token.IDENT)
	p.expect(token.COLON)
	expr := p.parse_expr()
	ret := ast.NewArgument(id, expr)
	ret.SetCoords(pos)
	return ret
}

// codeblock = LCURL { simpleStatement SEMI | complexStatement } RCURL
func (p *Parser) parse_codeblock() *ast.CodeBlock {
	statements := make([]ast.IStatement, 0)
	p.expect(token.LCURL)
	for p.tok == token.IF || p.tok == token.FOR || p.tok == token.LCURL || p.tok == token.VAR || p.tok == token.STORE || p.tok == token.EMIT || p.tok == token.CONST || p.tok == token.ASSERT || p.tok == token.IDENT || p.tok == token.WHILE || p.tok == token.CALL || p.tok == token.RETURN || p.tok == token.SEMI {
		switch p.tok {
		case token.VAR, token.STORE, token.EMIT, token.CONST, token.ASSERT, token.IDENT, token.CALL:
			stmt := p.parse_simpleStatement()
			statements = append(statements, stmt)
			p.expect(token.SEMI)
		case token.IF, token.FOR, token.LCURL, token.WHILE, token.RETURN:
			stmt := p.parse_complexStatement()
			statements = append(statements, stmt)
		default:
			p.err(p.pos, "syntax error") // PARSE ERROR
		}
	}
	p.expect(token.RCURL)
	block := ast.NewCodeBlock(statements)
	return block
}

// ifstmt = IF expr codeblock [ ELSE codeblock ]
func (p *Parser) parse_ifstmt() *ast.IfStmt {
	p.expect(token.IF)
	condition := p.parse_expr()
	thenBlock := p.parse_codeblock()
	var elseBlock *ast.CodeBlock
	if p.tok == token.ELSE {
		p.expect(token.ELSE)
		elseBlock = p.parse_codeblock()
	} else {
		elseBlock = ast.NewCodeBlock(nil)
	}
	p.expect(token.SEMI)
	ret := ast.NewIfStmt(condition, thenBlock, elseBlock)
	return ret
}

// forstmt = FOR [ simpleStatement ] SEMI [ expr ] SEMI [ simpleStatement ] codeblock
func (p *Parser) parse_forstmt() *ast.ForStmt {
	var init ast.IStatement
	var cond ast.IExpr
	var loop ast.IStatement
	var body *ast.CodeBlock
	p.expect(token.FOR)
	if p.tok == token.VAR || p.tok == token.STORE || p.tok == token.EMIT || p.tok == token.CONST || p.tok == token.ASSERT || p.tok == token.IDENT || p.tok == token.CALL {
		init = p.parse_simpleStatement()
	}
	p.expect(token.SEMI)
	if p.tok == token.LPAREN || p.tok == token.MINUS || p.tok == token.STRINGLITERAL || p.tok == token.BANG || p.tok == token.FALSE || p.tok == token.MINT || p.tok == token.TRUE || p.tok == token.HEXLITERAL || p.tok == token.NEW || p.tok == token.LARR || p.tok == token.IDENT || p.tok == token.BYTESLITERAL || p.tok == token.LBRACK || p.tok == token.BURN || p.tok == token.DECLITERAL {
		cond = p.parse_expr()
	}
	p.expect(token.SEMI)
	if p.tok == token.VAR || p.tok == token.STORE || p.tok == token.EMIT || p.tok == token.CONST || p.tok == token.ASSERT || p.tok == token.IDENT || p.tok == token.CALL {
		loop = p.parse_simpleStatement()
	}
	body = p.parse_codeblock()
	p.expect(token.SEMI)
	ret := ast.NewForStmt(init, cond, loop, body)
	return ret
}

// while = WHILE expr codeblock
func (p *Parser) parse_while() *ast.WhileStmt {
	p.expect(token.WHILE)
	expr := p.parse_expr()
	body := p.parse_codeblock()
	p.expect(token.SEMI)
	ret := ast.NewWhileStmt(expr, body)
	return ret
}

// loop = forstmt | while
func (p *Parser) parse_loop() ast.IStatement {
	var s ast.IStatement
	switch p.tok {
	case token.FOR:
		s = p.parse_forstmt()
	case token.WHILE:
		s = p.parse_while()
	default:
		p.err(p.pos, "syntax error") // PARSE ERROR
	}
	return s
}

// store = STORE BANG assignOrMoveOrCall
func (p *Parser) parse_store() ast.IStatement {
	p.expect(token.STORE)
	p.expect(token.BANG)
	tokenPos := p.pos
	stmt := p.parse_assignOrMoveOrCall()
	s, ok := stmt.(ast.Storer)
	if !ok {
		p.err(tokenPos, "only assignments and moves can be inside a store!")
	}
	s.SetStore(true)
	return stmt
}

// assert = ASSERT expr [ COMMA STRINGLITERAL ]
func (p *Parser) parse_assert() *ast.AssertStmt {
	p.expect(token.ASSERT)
	cond := p.parse_expr()
	var reason string
	if p.tok == token.COMMA {
		p.expect(token.COMMA)
		reason = strings.ReplaceAll(p.lit[1:len(p.lit)-1], "\\\"", "\"")
		p.expect(token.STRINGLITERAL)
	}
	ret := ast.NewAssertStmt(cond, reason)
	return ret
}

// funcsig = { decorator } FUNC IDENT [ BANG ] LPAREN [ paramList ] RPAREN [ type | LPAREN IDENT type { COMMA IDENT type } RPAREN ]
func (p *Parser) parse_funcsig() *ast.FuncSig {
	decorators := make([]ast.IDecorator, 0)
	for p.tok == token.AT {
		d := p.parse_decorator()
		decorators = append(decorators, d)
	}
	p.expect(token.FUNC)
	name := p.lit
	pos := p.GetPos()
	p.expect(token.IDENT)
	mutating := false
	if p.tok == token.BANG {
		mutating = true
		name += p.lit
		p.expect(token.BANG)
	}
	id := ast.NewIdentifier(name, pos)
	p.expect(token.LPAREN)
	params := make([]*ast.Param, 0)
	if p.tok == token.IDENT {
		params = p.parse_paramList()
	}
	p.expect(token.RPAREN)
	rets := make([]*ast.ReturnType, 0)
	if p.tok == token.LPAREN || p.tok == token.MAP || p.tok == token.IDENT || p.tok == token.LBRACK {
		switch p.tok {
		case token.MAP, token.IDENT, token.LBRACK:
			typ := p.parse_type()
			rets = append(rets, ast.NewReturnType(nil, typ))
		case token.LPAREN:
			p.expect(token.LPAREN)
			id := ast.NewIdentifier(p.lit, p.GetPos())
			p.expect(token.IDENT)
			typ := p.parse_type()
			rets = append(rets, ast.NewReturnType(id, typ))
			for p.tok == token.COMMA {
				p.expect(token.COMMA)
				id = ast.NewIdentifier(p.lit, p.GetPos())
				p.expect(token.IDENT)
				typ = p.parse_type()
				rets = append(rets, ast.NewReturnType(id, typ))
			}
			p.expect(token.RPAREN)
		default:
			p.err(p.pos, "syntax error") // PARSE ERROR
		}
	}
	ret := ast.NewFuncSig(id, params, rets, decorators, mutating)
	ret.SetCoords(pos)
	return ret
}

// decorator = AT ( PAYABLE | solidity )
func (p *Parser) parse_decorator() (d ast.IDecorator) {
	p.expect(token.AT)
	switch p.tok {
	case token.PAYABLE:
		p.expect(token.PAYABLE)
		d = ast.NewPayableDecorator()
	case token.SOLIDITY:
		d = p.parse_solidity()
	default:
		p.err(p.pos, "syntax error") // PARSE ERROR
	}
	p.expect(token.SEMI)
	return d
}

// solidity = SOLIDITY COLON IDENT LPAREN [ type { COMMA type } ] RPAREN [ type | LPAREN type { COMMA type } RPAREN ]
func (p *Parser) parse_solidity() *ast.SolidityDecorator {
	pos := p.GetPos()
	p.expect(token.SOLIDITY)
	p.expect(token.COLON)
	id := ast.NewIdentifier(p.lit, p.GetPos())
	p.expect(token.IDENT)
	p.expect(token.LPAREN)
	paramTypes := make([]ast.ITypeAST, 0)
	if p.tok == token.MAP || p.tok == token.IDENT || p.tok == token.LBRACK {
		typ := p.parse_type()
		paramTypes = append(paramTypes, typ)
		for p.tok == token.COMMA {
			p.expect(token.COMMA)
			typ = p.parse_type()
			paramTypes = append(paramTypes, typ)
		}
	}
	p.expect(token.RPAREN)
	returnTypes := make([]ast.ITypeAST, 0)
	if p.tok == token.LPAREN || p.tok == token.MAP || p.tok == token.IDENT || p.tok == token.LBRACK {
		switch p.tok {
		case token.MAP, token.IDENT, token.LBRACK:
			typ := p.parse_type()
			returnTypes = append(returnTypes, typ)
		case token.LPAREN:
			p.expect(token.LPAREN)
			typ := p.parse_type()
			returnTypes = append(returnTypes, typ)
			for p.tok == token.COMMA {
				p.expect(token.COMMA)
				typ = p.parse_type()
				returnTypes = append(returnTypes, typ)
			}
			p.expect(token.RPAREN)
		default:
			p.err(p.pos, "syntax error") // PARSE ERROR
		}
	}
	ret := ast.NewSolidityDecorator(id, paramTypes, returnTypes)
	ret.SetCoords(pos)
	return ret
}

// funcdef = funcsig codeblock
func (p *Parser) parse_funcdef() *ast.FuncDef {
	decl := p.parse_funcsig()
	block := p.parse_codeblock()
	p.expect(token.SEMI)
	ret := ast.NewFuncDef(decl, block)
	return ret
}

// call = CALL BANG expr { COMMA ( GAS | DATA | VALUE ) EQ expr }
func (p *Parser) parse_call() *ast.ExternalCall {
	pos := p.GetPos()
	p.expect(token.CALL)
	p.expect(token.BANG)
	dest := p.parse_expr()
	var gas ast.IExpr
	var data ast.IExpr
	var value ast.IExpr
	for p.tok == token.COMMA {
		p.expect(token.COMMA)
		switch p.tok {
		case token.DATA:
			if data != nil {
				p.err(p.pos, "duplicate data")
			}
			p.expect(token.DATA)
			p.expect(token.EQ)
			data = p.parse_expr()
		case token.GAS:
			if gas != nil {
				p.err(p.pos, "duplicate gas")
			}
			p.expect(token.GAS)
			p.expect(token.EQ)
			gas = p.parse_expr()
		case token.VALUE:
			if value != nil {
				p.err(p.pos, "duplicate value")
			}
			p.expect(token.VALUE)
			p.expect(token.EQ)
			value = p.parse_expr()
		default:
			p.err(p.pos, "syntax error") // PARSE ERROR
		}
	}
	ret := ast.NewExternalCall(dest, data, gas, value)
	ret.SetCoords(pos)
	return ret
}

// interface = INTERFACE IDENT LCURL { funcsig SEMI } RCURL
func (p *Parser) parse_interface() *ast.Interface {
	p.expect(token.INTERFACE)
	pos := p.GetPos()
	id := ast.NewIdentifier(p.lit, pos)
	p.expect(token.IDENT)
	p.expect(token.LCURL)
	sigs := make([]*ast.InterfaceSig, 0)
	for p.tok == token.FUNC || p.tok == token.AT {
		sig := p.parse_funcsig()
		inter := ast.NewInterfaceSig(sig)
		sigs = append(sigs, inter)
		p.expect(token.SEMI)
	}
	p.expect(token.RCURL)
	p.expect(token.SEMI)
	ret := ast.NewInterface(id, sigs)
	ret.SetCoords(pos)
	return ret
}
