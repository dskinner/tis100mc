package tis100mc

import "sync"

type ExecNode struct {
	*BasicNode

	acc, bak Int10

	ops []*Operation
	pos int

	dr DeferRead
	dw DeferWrite
}

func (n *ExecNode) ACC() Int10 { return n.acc }
func (n *ExecNode) BAK() Int10 { return n.bak }

func (n *ExecNode) LoadProgram(prg []byte) {
	n.ops = Parse(prg)
	n.pos = 0
}

func (n *ExecNode) Operation() (int, *Operation) {
	if n.dr.P == 0 && n.dw.P == 0 {
		return n.pos, n.ops[n.pos]
	}
	pos := n.pos - 1
	if pos < 0 {
		pos = 0
	}
	return pos, n.ops[pos]
}

// Step executes the next operation as specified by current position.
func (n *ExecNode) Step() {
	n.logf("n.Step()")

	if n.dr.P == 0 && n.dw.P == 0 {
		n.exec()
	}

	// check deferred read off registers/pseudo-ports
	switch n.dr.P {
	case ACC:
		n.logf("deferred read from ACC")
		n.dr.F(n.acc)
		n.dr = DeferRead{}
	case NIL:
		n.logf("deferred read from NIL")
		n.dr.F(0)
		n.dr = DeferRead{}
	}

	// check deferred write off registers/pseudo-ports
	switch n.dw.P {
	case ACC:
		n.logf("deferred write from ACC")
		n.acc = n.dw.X
		n.dw = DeferWrite{}
	case NIL:
		n.logf("deferred write from NIL")
		n.dw = DeferWrite{}
	}

	// check deferred read off chans
	select {
	case x := <-n.Port(n.dr.P):
		n.dr.F(x)
		n.logf("performed deferred read")
		n.dr = DeferRead{}
	default:
	}

	// check deferred write off chans
	select {
	case n.Port(n.dw.P) <- n.dw.X:
		n.logf("performed deferred write")
		n.dw = DeferWrite{}
	default:
	}
}

func (n *ExecNode) Start(stop chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case <-stop:
			return
		case <-n.step:
			n.Step()
		}
		wg.Done()
	}
}

func (n *ExecNode) exec() {
	if len(n.ops) == 0 {
		return
	}
	op := n.ops[n.pos]
	n.logf("executing: %s\n", op)
	switch op.Instruction {
	case NOP:
	case MOV:
		if (*op.A).IsPort() {
			// n.Write(*op.B, n.Read(Port(*op.A)))
			n.dr = DeferMov(n, Port(*op.A), *op.B)
		} else {
			// n.Write(*op.B, *op.A)
			n.dw = DeferWrite{P: *op.B, X: *op.A}
		}
	case SWP:
		n.acc, n.bak = n.bak, n.acc
	case SAV:
		n.bak = n.acc
	case ADD:
		if (*op.A).IsPort() {
			// n.acc = clamp(n.acc + n.Read(Port(*op.A)))
			n.dr = DeferAdd(n, Port(*op.A))
		} else {
			n.acc = clamp(n.acc + *op.A)
		}
	case SUB:
		if (*op.A).IsPort() {
			// n.acc = clamp(n.acc - n.Read(Port(*op.A)))
			n.dr = DeferSub(n, Port(*op.A))
		} else {
			n.acc = clamp(n.acc - *op.A)
		}
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
}

type DeferRead struct {
	P Port
	F func(Int10)
}

type DeferWrite struct {
	P Port
	X Int10
}

// TODO document defer methods in relation to TIS-100 emu. Defers should
// be avoided for MC.

func DeferAdd(n *ExecNode, p Port) DeferRead {
	return DeferRead{
		P: p,
		F: func(x Int10) { n.acc = clamp(n.acc + x) },
	}
}

func DeferSub(n *ExecNode, p Port) DeferRead {
	return DeferRead{
		P: p,
		F: func(x Int10) { n.acc = clamp(n.acc - x) },
	}
}

func DeferMov(n *ExecNode, src, dest Port) DeferRead {
	return DeferRead{
		P: src,
		F: func(x Int10) {
			n.dw = DeferWrite{
				P: dest,
				X: x,
			}
		},
	}
}
