package scanner

import (
	"bytes"
	"io"
	"smart/token"
	"unicode/utf8"
)

type Scanner struct {
	src []byte

	offset int

	prev *token.Token
}

func (s *Scanner) init(r io.Reader) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		panic(err)
	}
	s.src = buf.Bytes()
}

func New(r io.Reader) *Scanner {
	var s Scanner

	s.init(r)
	return &s
}

func (s *Scanner) Scan() (pos token.Pos, tok token.Token, lit string) {
	if s.skipWhitespace() {
		s.prev = nil
		return pos - 1, token.SEMI, ""
	}
	tok, lit = s.Scan0(s.src[s.offset:])
	if tok == token.IDENT {
		tok = token.Lookup(lit)
	}
	pos = token.Pos(s.offset)
	s.offset += len(lit)
	s.prev = &tok
	return pos, tok, lit
}

func (s *Scanner) Coordinates(pos token.Pos) (line int, col int) {
	line = 1
	col = 1
	for i := token.Pos(0); i < pos; i++ {
		col += 1
		if rune(s.src[i]) == '\n' {
			line += 1
			col = 1
		}
	}
	return line, col
}

func (s *Scanner) skipWhitespace() bool {
	for s.offset < len(s.src) {
		s.skipComment()
		if s.offset >= len(s.src) {
			break
		}
		ch := rune(s.src[s.offset])
		if ch > utf8.RuneSelf {
			// ERROR
		}
		switch ch {
		default:
			return false
		case ' ', '\t':
			s.offset += 1
		case '\n':
			s.offset += 1
			if s.prev != nil && (*s.prev == token.IDENT || *s.prev == token.BYTESLITERAL || *s.prev == token.DECLITERAL || *s.prev == token.HEXLITERAL || *s.prev == token.STRINGLITERAL || *s.prev == token.RPAREN || *s.prev == token.RBRACK || *s.prev == token.RCURL || *s.prev == token.PAYABLE || *s.prev == token.TRUE || *s.prev == token.FALSE) {
				return true
			}
		}
	}
	return false
}

func (s *Scanner) skipComment() {
	if bytes.HasPrefix(s.src[s.offset:], []byte("//")) {
		s.offset += 2
		for s.offset < len(s.src) && s.src[s.offset] != '\n' {
			s.offset += 1
		}
	}
}
