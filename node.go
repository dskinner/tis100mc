package tis100mc

import (
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
	Read(Port) Int10
	Write(Port, Int10)
	Step()
}

type CN struct {
	left, right, up, down chan Int10

	// stage writes
	any chan Int10
	// tracks ANY
	last chan Int10
}

func (cn CN) Port(port Port) chan Int10 {
	switch port {
	case LEFT:
		return cn.left
	case RIGHT:
		return cn.right
	case UP:
		return cn.up
	case DOWN:
		return cn.down
	case LAST:
		return cn.last
	default:
		panic(fmt.Sprintf("Port %s is not addressable.", port))
	}
}

// Read returns value from port and will block if there is nothing to read.
//
// TODO There's a fundamental problem with how Nodes currently read in that they
// do not request their neighbors but self-refer, making the defined behavior
// of write on ANY problematic.
func (cn CN) Read(port Port) Int10 {
	switch port {
	case LEFT:
		return <-cn.left
	case RIGHT:
		return <-cn.right
	case UP:
		return <-cn.up
	case DOWN:
		return <-cn.down
	case ANY:
		select {
		case x := <-cn.left:
			(&cn).last = cn.left
			return x
		case x := <-cn.right:
			(&cn).last = cn.right
			return x
		case x := <-cn.up:
			(&cn).last = cn.up
			return x
		case x := <-cn.down:
			(&cn).last = cn.down
			return x
		}
	case LAST:
		if cn.last == nil {
			// TODO return error instead
			panic("Attempt to read LAST before ANY.")
		}
		return <-cn.last
	case NIL:
		return 0
	}
	panic(fmt.Sprintf("Read from %s not supported.", port))
}

// Write writes value to port and will block if port is pending read.
func (cn CN) Write(port Port, x Int10) {
	switch port {
	case LEFT:
		cn.left <- x
	case RIGHT:
		cn.right <- x
	case UP:
		cn.up <- x
	case DOWN:
		cn.down <- x
	case ANY:
		// TODO the docs state when ANY is dest, result will be sent to first
		// node that attempts to read on any port. Behavior is unspecified when
		// multiple nodes are pending reads on src node.
		//
		// Perhaps a RequestRead is necessary to track this behavior for tis-100.
		panic("TODO implement")
	case LAST:
		if cn.last == nil {
			// TODO return error instead
			panic("Attempt to write LAST before ANY.")
		}
		cn.last <- x
	case NIL:
		// TODO check spec
	default:
		panic(fmt.Sprintf("Write to %s not supported.", port))
	}
}

// TODO Join assigns the node's port to the argument node's opposite
// port. If neither node has a port created to share, one is created.
func Join(cn0 CN, port Port, cn1 CN) (CN, CN) {
	ch := make(chan Int10, 1)
	switch port {
	case LEFT:
		cn0.left = ch
		cn1.right = cn0.left
	case RIGHT:
		cn0.right = ch
		cn1.left = cn0.right
	case UP:
		cn0.up = ch
		cn1.down = cn0.up
	case DOWN:
		cn0.down = ch
		cn1.up = cn0.down
	default:
		panic(fmt.Sprintf("Join on %s not supported.", port))
	}
	return cn0, cn1
}

// type PX struct {
// up0, down0, left0, right0 chan Int10
// up1, down1, left1, right1 chan Int10
// }

type BasicNode struct {
	ports CN
	cy    *Cycler

	// TODO make this optional as part of tis-100 emulation
	// step chan struct{}

	// stores last port read or written to when using ANY
	last    Port
	blocked bool
	debug   bool
}

func (n *BasicNode) logf(format string, v ...interface{}) {
	if n.debug {
		log.Printf(format, v...)
	}
}

func (n *BasicNode) Blocked() bool { return n.blocked }

func (n *BasicNode) Step() {}

func (n *BasicNode) Port(port Port) chan Int10 {
	switch port {
	case NIL:
		// TODO implement to spec
		return make(chan Int10, 1)
	case Port(0): // TODO meh
		var ch chan Int10
		return ch
	default:
		return n.ports.Port(port)
	}
}

func (n *BasicNode) Read(port Port) Int10 {
	// n.cy.Cycle()
	return n.ports.Read(port)
}

// Write writes value to port and will block if port is pending read.
func (n *BasicNode) Write(port Port, x Int10) {
	// n.cy.Cycle()
	n.ports.Write(port, x)
}
