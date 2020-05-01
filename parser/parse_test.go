package parser

import (
	"bufio"
	"bytes"
	"os"
	"smart/scanner"
	"testing"
)

func TestParser(t *testing.T) {
	file, err := os.Open("test.smrt")
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(file)
	s := scanner.New(r)
	p := New(s, "test.smrt")
	p.verbose = true

	p.Parse_prog()
}

func skipTestErrorCoordinates(t *testing.T) {
	b := []byte(
		`// line 1

// line 3
  xxx // should report line 4, col 3
`)
	r := bytes.NewReader(b)
	s := scanner.New(r)
	p := New(s, "<unknown file>")

	p.Parse_prog()
}
