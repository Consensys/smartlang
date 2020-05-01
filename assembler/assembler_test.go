package asm

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	a := NewAssembler()
	a.Label("hello")
	a.PushLabel("hello")
	a.Push1(5)
	a.Emit(JUMPI)
	a.Comment("this is just a comment")
	a.Label("loop")
	a.PushLabelc("loop", "just jumping to a loop")
	a.Emitc(JUMP, "Go!")
	a.PushBytes([]byte{1, 2, 3, 4, 5})

	fmt.Printf("\n\nBytecode: 0x%X\n", a.Compile())
}
