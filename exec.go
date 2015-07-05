package tis100mc

import "sync"

type ExecNode struct {
	*BasicNode

	acc, bak Int10

	ops []*Operation
	pos int
}

func (n *ExecNode) ACC() Int10 { return n.acc }
func (n *ExecNode) BAK() Int10 { return n.bak }

func (n *ExecNode) LoadProgram(prg []byte) {
	n.ops = Parse(prg)
	n.pos = 0
}

func (n *ExecNode) Operation() (int, *Operation) {
	return n.pos, n.ops[n.pos]
}

func (n *ExecNode) Read(port Port) Int10 {
	if port == ACC {
		return n.acc
	}
	return n.BasicNode.Read(port)
}

func (n *ExecNode) Write(port Port, x Int10) {
	if port == ACC {
		n.acc = x
	} else {
		n.BasicNode.Write(port, x)
	}
}

func (n *ExecNode) Start(stop chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case <-stop:
			return
			// case <-n.step:
		default:
			n.Step()
		}
	}
}

// Step executes the next operation as specified by current position.
func (n *ExecNode) Step() {
	n.cy.Cycle()
	if len(n.ops) == 0 {
		// TODO stop endless for
		n.cy.Done()
		return
	}
	<-n.cy.waitio // better name than done,step

	op := n.ops[n.pos]
	n.logf("executing: %s\n", op)
	switch op.Instruction {
	case NOP:
	case MOV:
		n.cy.waitio = make(chan struct{})
		go func() {
			if (*op.A).IsPort() {
				n.Write(*op.B, n.Read(Port(*op.A)))
			} else {
				n.Write(*op.B, *op.A)
			}
			close(n.cy.waitio)
		}()
	case SWP:
		n.acc, n.bak = n.bak, n.acc
	case SAV:
		n.bak = n.acc
	case ADD:
		n.cy.waitio = make(chan struct{})
		go func() {
			if (*op.A).IsPort() {
				n.acc = clamp(n.acc + n.Read(Port(*op.A)))
			} else {
				n.acc = clamp(n.acc + *op.A)
			}
			close(n.cy.waitio)
		}()
	case SUB:
		n.cy.waitio = make(chan struct{})
		go func() {
			if (*op.A).IsPort() {
				n.acc = clamp(n.acc - n.Read(Port(*op.A)))
			} else {
				n.acc = clamp(n.acc - *op.A)
			}
			close(n.cy.waitio)
		}()
	case NEG:
		n.acc = -n.acc
	case JMP:
		n.pos = int(*op.A) - 1
	case JEZ:
		if n.acc == 0 {
			n.pos = int(*op.A) - 1
		}
	case JNZ:
		if n.acc != 0 {
			n.pos = int(*op.A) - 1
		}
	case JGZ:
		if n.acc > 0 {
			n.pos = int(*op.A) - 1
		}
	case JLZ:
		if n.acc < 0 {
			n.pos = int(*op.A) - 1
		}
	case JRO:
		n.pos += int(*op.A) - 1
	default:
		panic("Unknown instruction")
	}

	n.pos++
	if n.pos >= len(n.ops) {
		n.pos = 0
	}

	n.cy.Done()
}
