package tis100mc

import "fmt"

func init() {
	_ = Node(&StackNode{})
}

func consume(n *StackNode, ch chan Int10) {
	if n.length != len(n.mem) {
		select {
		case x := <-ch:
			n.mem[n.length] = x
			n.length++
		default:
		}
	}
}

type StackNode struct {
	// TODO replace with PX
	source CN
	sink   CN // TODO assure not buffered

	read chan struct{}

	mem    [9]Int10
	length int
}

func NewStackNode() *StackNode {
	return &StackNode{
		read: make(chan struct{}, 1),
	}
}

// Perhaps a RequestRead is necessary to track this behavior for tis-100.

// TODO Step across all nodes is inherently racey which performing read/write
// at same time. This is not applicable for MC implementation where there is no
// "steps", but for TIS-100 emu, will likely need some kind of bridge where each
// node either signals write or no-write, and each neighbor blocks until receiving
// signal before completing step.
//
// TODO if there's no room on stack, don't read source. This simulates a block for
// TIS-100. For MC, this should actually block.
func (n *StackNode) Step() {

	// can consume multiple ports during a single cycle.
	consume(n, n.source.left)
	consume(n, n.source.right)
	consume(n, n.source.up)
	consume(n, n.source.down)

	// may only produce for one port during a single cycle.
	if n.length > 0 {
		x := n.mem[n.length-1] // peek
		select {
		case <-n.read:
			select {
			case n.sink.left <- x:
				n.length--
			case n.sink.right <- x:
				n.length--
			case n.sink.up <- x:
				n.length--
			case n.sink.down <- x:
				n.length--
			default:
				// TODO racey and only use if System.emu == true
				panic("Received read signal but sink not waiting on write.")
			}
		default:
		}
	}
}

func (n *StackNode) Read(port Port) Int10 {
	return 0
}

func (n *StackNode) Write(port Port, x Int10) {

}

func (n *StackNode) Connect(port Port, cn Node) error {
	return nil
}

func (n *StackNode) SetPort(port Port, ch chan Int10) {
	switch port {
	case UP:
		n.source.up = ch
	case DOWN:
		n.source.down = ch
	case LEFT:
		n.source.left = ch
	case RIGHT:
		n.source.right = ch
	default:
		panic(fmt.Sprintf("SetPort does not support %s", port))
	}
}

/*
n := StackNode

n0 := ExecNode
n0.Connect(RIGHT, n) //internally calls n.SetPort

MOV 1, RIGHT // writes 1 to n0.right which == n.source.left

// due to this instruction, all stacknode connections will have
// to originate from the stacknode to replace the basicnode's chan
// appropriately, unless something else can be worked out with type CN
MOV RIGHT, ACC

*/

func (n *StackNode) Port(port Port) chan Int10 {
	var ch chan Int10
	return ch
	// if port == ANY {

	// TODO this needs to be re-evaluated
	// To keep ANY dest in same select statement, n.Port(ANY) returns
	// nil channel when n.length == 0 so it always blocks.

	// if n.length == 0 {
	// var ch chan Int10
	// return ch
	// }

	// TODO don't spin off go routine just to do this
	// ch := make(chan Int10, 1)
	// go func() {
	// select {
	// case x := <-n.left:
	// log.Println("any left")
	// ch <- x
	// case x := <-n.right:
	// log.Println("any right")
	// ch <- x
	// case x := <-n.up:
	// log.Println("any up")
	// ch <- x
	// case x := <-n.down:
	// log.Println("any down")
	// ch <- x
	// case x := <-ch:
	// log.Printf("X: %v\n", x)
	// }
	// }()
	// return ch
	// }
	// return n.BasicNode.Port(port)
}
