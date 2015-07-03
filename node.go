package tis100mc

import (
	"errors"
	"fmt"
	"log"
)

func init() {
	_ = Node(&BasicNode{})
}

// Int10 allows for valid numerical ranges between -999, 999 and reserves values outside
// of those bounds for constants. It is also an exercise in sillyness.
type Int10 int16

const MaxInt10 = 1<<10 - 25

// IsPort determines if value is outside of valid numerical range -999, 999.
func (n Int10) IsPort() bool {
	return -MaxInt10 > n || n > MaxInt10
}

// clamp returns value clamped to Int10 numerical range.
func clamp(x Int10) Int10 {
	if x > MaxInt10 {
		x = MaxInt10
	} else if x < -MaxInt10 {
		x = -MaxInt10
	}
	return x
}

type Node interface {
	Connect(Port, Node) error

	Read(Port) Int10
	Write(Port, Int10)

	Port(Port) chan Int10
	// TODO objectional as can break system
	SetPort(Port, chan Int10)

	Step()
}

type BasicNode struct {
	up, down, left, right chan Int10

	blocked bool

	// TODO make this optional as part of tis-100 emulation
	step chan struct{}

	// stores last port read or written to when using ANY
	last Port

	debug bool
}

func (n *BasicNode) logf(format string, v ...interface{}) {
	if n.debug {
		log.Printf(format, v...)
	}
}

func (n *BasicNode) Blocked() bool { return n.blocked }

func (n *BasicNode) Step() {}

// Connect assigns the node's port to the argument node's opposite
// port. If neither node has a port created to share, one is created.
func (n *BasicNode) Connect(port Port, cn Node) error { // TODO accept Node interface arg
	switch port {
	case UP:
		return cn.Connect(DOWN, n)
	case DOWN:
		if n.down != nil || cn.Port(UP) != nil {
			return errors.New("Attempting to connect up/down but chan already exists.")
		}
		n.down = make(chan Int10, 1)
		cn.SetPort(UP, n.down)
	case LEFT:
		return cn.Connect(RIGHT, n)
	case RIGHT:
		if n.right != nil || cn.Port(LEFT) != nil {
			return errors.New("Attempting to connect left/right but chan already exists.")
		}
		n.right = make(chan Int10, 1)
		cn.SetPort(LEFT, n.right)
	default:
		return errors.New(fmt.Sprintf("Unknown port: %s", port))
	}
	return nil
}

func (n *BasicNode) Port(port Port) chan Int10 {
	switch port {
	case UP:
		return n.up
	case DOWN:
		return n.down
	case LEFT:
		return n.left
	case RIGHT:
		return n.right
	case ACC:
		panic("ACC not supported")
	case NIL:
		// TODO implement to spec
		return make(chan Int10, 1)
	default:
		var ch chan Int10
		return ch
	}
}

func (n *BasicNode) SetPort(port Port, ch chan Int10) {
	switch port {
	case UP:
		n.up = ch
	case DOWN:
		n.down = ch
	case LEFT:
		n.left = ch
	case RIGHT:
		n.right = ch
	default:
		panic(fmt.Sprintf("SetPort does not support %s", port))
	}
}

// Read returns value from port and will block if there is nothing to read.
func (n *BasicNode) Read(port Port) Int10 {
	switch port {
	case UP:
		return <-n.up
	case DOWN:
		return <-n.down
	case LEFT:
		return <-n.left
	case RIGHT:
		return <-n.right
	case ANY:
		select {
		case x := <-n.up:
			n.last = UP
			return x
		case x := <-n.down:
			n.last = DOWN
			return x
		case x := <-n.left:
			n.last = LEFT
			return x
		case x := <-n.right:
			n.last = RIGHT
			return x
		}
	case NIL:
		// panic("TODO implement")
		// TODO check spec
		return 0
		// TODO handle in ExecNode
		// case ACC:
		// return n.acc
	}
	panic(fmt.Sprintf("Read from unknown port %s", port))
}

// Write writes value to port and will block if port is pending read.
func (n *BasicNode) Write(port Port, val Int10) {
	switch port {
	case UP:
		n.up <- val
	case DOWN:
		n.down <- val
	case LEFT:
		n.left <- val
	case RIGHT:
		n.right <- val
	case ANY:
		// TODO the docs state when ANY is dest, result will be sent to first
		// node that attempts to read on any port. Behavior is unspecified when
		// multiple nodes are pending reads on src node.
		//
		// Perhaps a RequestRead is necessary to track this behavior for tis-100.
		panic("TODO implement")
	case ACC:
		// n.acc = val
		// TODO handle in ExecNode
		panic("ACC not supported")
	}
}
