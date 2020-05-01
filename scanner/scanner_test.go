package scanner

import (
	"bytes"
	"fmt"
	"smart/token"
	"testing"
)

func TestScanner(t *testing.T) {
	b := []byte("123 this for")
	r := bytes.NewReader(b)
	s := New(r)
	fmt.Println(s)
	var p token.Pos
	var tok token.Token = token.ILLEGAL
	var str string

	for tok != token.EOF {
		p, tok, str = s.Scan()
		fmt.Println(p, tok, token.Tokens[tok], str)
		fmt.Println(s)
	}
}
