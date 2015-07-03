package tis100mc

import "fmt"

type parser struct {
	lex    *lexer
	ops    []*Operation
	labels map[string]int
}

type labelsReceiver struct {
	labels map[string]int
	index  int
}

func (r *labelsReceiver) Receive(tkn Token) {
	switch tkn.typ {
	case TokenLabel:
		r.labels[tkn.val] = r.index
	case TokenInstruction:
		r.index++
	}
}

func Parse(bytes []byte) []*Operation {
	p := &parser{}

	lbls := &labelsReceiver{make(map[string]int), 0}
	lex := NewLexer(lbls)
	lex.bytes = bytes
	lex.Run()

	p.labels = lbls.labels
	lex = NewLexer(p)
	lex.bytes = bytes
	lex.Run()

	return p.ops
}

func (p *parser) Receive(tkn Token) {
	switch tkn.typ {
	// case TokenLabel:
	// p.labels[tkn.val] = len(p.ops)
	case TokenInstruction:
		if _, ok := InstructionMap[tkn.val]; !ok {
			panic("Invalid instruction - " + tkn.val)
		}
		p.ops = append(p.ops, &Operation{Instruction: InstructionMap[tkn.val]})
	case TokenArgument:
		op := p.ops[len(p.ops)-1]
		var arg Int10
		if i, ok := p.labels[tkn.val]; ok {
			arg = Int10(i)
		} else if op.Instruction == JEZ {

		} else {
			arg = ParseInt10(tkn.val)
		}
		if op.A == nil {
			op.A = &arg
		} else if op.B == nil {
			if !arg.IsPort() {
				panic(fmt.Sprintf("Invalid port value - %v\n", arg))
			}
			p := Port(arg)
			op.B = &p
		} else {
			panic("Received too many arguments.")
		}
	default:
	}
}
