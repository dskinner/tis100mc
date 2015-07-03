package tis100mc

import (
	"errors"
	"math"
	"testing"
	"time"
)

func newNode() *ExecNode {
	return &ExecNode{
		BasicNode: &BasicNode{
			up:    make(chan Int10, 1),
			down:  make(chan Int10, 1),
			left:  make(chan Int10, 1),
			right: make(chan Int10, 1),
		},
	}
}

func TestIsPort(t *testing.T) {
	for i := -MaxInt10; i <= MaxInt10; i++ {
		if Int10(i).IsPort() {
			t.Fatal(i, "is not a port.")
		}
	}
	for i := -MaxInt10 - 1; i >= -math.MaxInt16; i-- {
		if !Int10(i).IsPort() {
			t.Fatal(i, "is a port.")
		}
	}
	for i := MaxInt10 + 1; i <= math.MaxInt16; i++ {
		if !Int10(i).IsPort() {
			t.Fatal(i, "is a port.")
		}
	}
}

func TestConnect(t *testing.T) {
	n0 := &BasicNode{}
	n1 := &BasicNode{}

	ps := []struct {
		A, B Port
	}{
		{UP, DOWN},
		{DOWN, UP},
		{LEFT, RIGHT},
		{RIGHT, LEFT},
	}

	for _, p := range ps {
		if err := n0.Connect(p.A, n1); err != nil {
			t.Fatal(err)
		}
		if n0.Port(p.A) == nil || n1.Port(p.B) == nil || n0.Port(p.A) != n1.Port(p.B) {
			t.Fatalf("Node connect failure.\nn0: %+v\nn1: %+v", n0, n1)
		}
		if err := n1.Connect(p.B, n0); err == nil {
			t.Fatal("Node connect should return error, already connected.")
		}
	}

	if err := n0.Connect(Port(0), n1); err == nil {
		t.Fatal("Node connect should return error, unknown port.")
	}
}

func TestReadWrite(t *testing.T) {
	n := newNode()

	ms := []struct {
		C chan Int10
		P Port
		X Int10
	}{
		{n.up, UP, 0},
		{n.down, DOWN, 1},
		{n.left, LEFT, 2},
		{n.right, RIGHT, 3},
	}

	for _, m := range ms {
		// Read
		m.C <- m.X
		if n.Read(m.P) != m.X {
			t.Fatalf("Read failed, value doesn't match: %+v\n", m)
		}
		// Write
		n.Write(m.P, m.X)
		if <-m.C != m.X {
			t.Fatalf("Write failed, value doesn't match: %+v\n", m)
		}
	}
}

func read(ch chan Int10) (Int10, error) {
	select {
	case <-time.After(100 * time.Millisecond):
		return 0, errors.New("Timeout on read")
	case x := <-ch:
		return x, nil
	}
}

func write(ch chan Int10, x Int10) error {
	select {
	case <-time.After(100 * time.Millisecond):
		return errors.New("Timeout on write")
	case ch <- x:
	}
	return nil
}

func TestMOV(t *testing.T) {
	n := newNode()
	n.debug = true

	n.ops = []*Operation{
		NewOperation(MOV, 1, ACC),
		// NewOperation(MOV, Int10(ACC), ACC),
		NewOperation(MOV, Int10(ACC), RIGHT),
		NewOperation(MOV, Int10(RIGHT), LEFT),
		NewOperation(MOV, Int10(LEFT), LEFT),
		NewOperation(MOV, Int10(LEFT), NIL),
	}

	// expects := []struct {
	// C chan Int10
	// X Int10
	// }{
	// {n.right, 1},
	// {n.left, 1},
	// {n.left, 1},
	// }

	// for _, e := range expects {
	// n.Step()
	// if x, err := read(e.C); x != 1 || err != nil {
	// t.Fatalf("Failed: error %s: %s - %+v", err, n.ops[n.pos-1], n)
	// } else {
	// e.C <- x
	// }
	// }

	n.logf("exec op 0")
	n.Step()
	if n.acc != 1 {
		t.Fatalf("Failed: %s - %+v", n.ops[0], n)
	}
	// n.Step()
	// if n.acc != 1 {
	// t.Fatal(fmt.Sprintf("Failed: %s - %+v", n.ops[0], n))
	// }

	n.logf("exec op 1")
	n.Step()
	if x, err := read(n.right); x != 1 || err != nil {
		t.Fatalf("Failed: %s - %+v", n.ops[1], n)
	} else {
		n.right <- x
	}

	n.logf("exec op 2")
	n.Step()
	if x, err := read(n.left); x != 1 || err != nil {
		t.Fatalf("Failed: %s - %+v", n.ops[2], n)
	} else {
		n.left <- x
	}
	n.Step()
	if x, err := read(n.left); x != 1 || err != nil {
		t.Fatalf("Failed: %s - %+v", n.ops[3], n)
	} else {
		n.left <- x
	}

	n.Step()
	if err := write(n.left, 1); err != nil || n.dw.P != 0 {
		t.Fatalf("Failed (pending write): %s - %+v", n.ops[4], n)
	}
}

func TestSWP(t *testing.T) {
	n := newNode()
	n.ops = []*Operation{
		{SWP, nil, nil},
	}
	n.acc = 1
	n.Step()
	if n.acc != 0 || n.bak != 1 {
		t.Fatal("Failed: SWP")
	}
	n.Step()
	if n.acc != 1 || n.bak != 0 {
		t.Fatal("Failed: SWP")
	}
}

func TestSAV(t *testing.T) {
	n := newNode()
	n.ops = []*Operation{
		{SAV, nil, nil},
	}
	n.acc = 1
	n.Step()
	if n.acc != 1 || n.bak != 1 {
		t.Fatal("Failed: SAV")
	}
}

func TestADD(t *testing.T) {
	n := newNode()
	n.ops = []*Operation{
		NewOperation(ADD, 5, 0),
		NewOperation(ADD, -7, 0),
	}
	n.Step()
	if n.acc != 5 {
		t.Fatal()
	}
	n.Step()
	if n.acc != -2 {
		t.Fatal()
	}
}

func TestSUB(t *testing.T) {
	n := newNode()
	n.ops = []*Operation{
		NewOperation(SUB, 5, 0),
		NewOperation(SUB, -7, 0),
	}
	n.Step()
	if n.acc != -5 {
		t.Fatal()
	}
	n.Step()
	if n.acc != 2 {
		t.Fatal()
	}
}

func TestNEG(t *testing.T) {
	n := newNode()
	n.ops = []*Operation{
		{NEG, nil, nil},
	}
	n.acc = 5
	n.Step()
	if n.acc != -5 {
		t.Fatal()
	}
	n.Step()
	if n.acc != 5 {
		t.Fatal()
	}
}

func testStep(t *testing.T) {
	// JMP
	// JEZ
	// JNZ
	// JGZ
	// JLZ
	// JRO
}
