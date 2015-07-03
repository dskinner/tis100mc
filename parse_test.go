package tis100mc

import (
	"testing"
)

func TestInstruction(t *testing.T) {
	t.Log(NOP, MOV)
	var x Int10 = 22
	_ = Instruction(x)
}

func TestParser(t *testing.T) {
	Parse([]byte(example))
}

var funcMap = map[Instruction]func(int) int{
	ADD: func(x int) int { return x * x },
}

func makeFunc(x int) func() int {
	return func() int { return x * x }
}

var result int
var op Instruction = ADD

type Adder struct {
	Instruction

	Arg int
}

func (a *Adder) Exec() int { return a.Arg * a.Arg }

func BenchmarkClosure(b *testing.B) {
	for n := 0; n < b.N; n++ {
		switch op {
		case ADD:
			fn := makeFunc(20)
			result = fn()
		}
	}
}

func BenchmarkMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fn := funcMap[ADD]
		result = fn(20)
	}
}

func BenchmarkStruct(b *testing.B) {
	for n := 0; n < b.N; n++ {
		a := &Adder{ADD, 20}
		result = a.Exec()
	}
}
