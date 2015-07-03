//go:generate stringer -type=Instruction,Port -output=types_string.go

package tis100mc

import (
	"bytes"
	"fmt"
	"strconv"
)

type Instruction Int10

type Port Int10

// TODO representing multiple ops where each requires either 0, 1, or 2
// args in a single struct is finicky, where currently nil represents
// no arg.
type Operation struct {
	Instruction

	A *Int10
	B *Port
}

func NewOperation(ns Instruction, a Int10, b Port) *Operation {
	return &Operation{ns, &a, &b}
}

func (op Operation) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s(", op.Instruction))
	if op.A != nil {
		if (*op.A).IsPort() {
			buf.WriteString(fmt.Sprintf("%v", Port(*op.A)))
		} else {
			buf.WriteString(fmt.Sprintf("%v", *op.A))
		}
	}
	if op.B != nil {
		buf.WriteString(fmt.Sprintf(", %v", *op.B))
	}
	buf.WriteString(")")
	return buf.String()
}

const (
	NOP Instruction = iota + MaxInt10 + 1
	MOV
	SWP
	SAV
	ADD
	SUB
	NEG
	JMP
	JEZ
	JNZ
	JGZ
	JLZ
	JRO

	ACC Port = iota + MaxInt10 + 1
	BAK
	NIL
	LEFT
	RIGHT
	UP
	DOWN
	ANY
	LAST
)

var InstructionMap = map[string]Instruction{
	"NOP": NOP,
	"MOV": MOV,
	"SWP": SWP,
	"SAV": SAV,
	"ADD": ADD,
	"SUB": SUB,
	"NEG": NEG,
	"JMP": JMP,
	"JEZ": JEZ,
	"JNZ": JNZ,
	"JGZ": JGZ,
	"JLZ": JLZ,
	"JRO": JRO,
}

var PortMap = map[string]Port{
	"ACC":   ACC,
	"LEFT":  LEFT,
	"RIGHT": RIGHT,
	"UP":    UP,
	"DOWN":  DOWN,
	"ANY":   ANY,
	"LAST":  LAST,
	"NIL":   NIL,
}

func ParseInt10(x string) Int10 {
	if i, ok := PortMap[x]; ok {
		return Int10(i)
	}
	i, err := strconv.ParseInt(x, 10, 16)
	if err != nil {
		panic(err)
	}
	return Int10(i)
}
