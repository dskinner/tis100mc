package tis100mc

func init() {
	_ = Node(&StackNode{})
}

type StackNode struct {
	*BasicNode

	mem    [9]Int10
	length int
}

// TODO the docs state when ANY is dest, result will be sent to first
// node that attempts to read on any port. Behavior is unspecified when
// multiple nodes are pending reads on src node.
//
// Perhaps a RequestRead is necessary to track this behavior for tis-100.

// TODO Step across all nodes is inherently racey which performing read/write
// at same time. This is not applicable for MC implementation where there is no
// "steps", but for TIS-100 emu, will likely need some kind of bridge where each
// node either signals write or no-write, and each neighbor blocks until receiving
// signal before completing step.
func (n *StackNode) Step() {

	// To keep ANY dest in same select statement, n.Port(ANY) returns
	// nil channel when n.length == 0 so it always blocks.
	select {
	case n.Port(ANY) <- n.mem[n.length-1]:
		n.length--
	case x := <-n.Port(ANY):
		n.mem[n.length] = x
		n.length++
	default:
	}
}

func (n *StackNode) Port(port Port) chan Int10 {
	if port == ANY && n.length == 0 {
		var ch chan Int10
		return ch
	}
	return n.BasicNode.Port(port)
}
