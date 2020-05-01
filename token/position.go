package token

type Pos int

type Position struct {
	Filename string
	Offset   int
	Line     int
	Column   int
}
