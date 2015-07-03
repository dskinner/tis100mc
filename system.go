package tis100mc

import "sync"

// TODO setup simple way to connect multiple systems to connect
// boundary nodes.
type System struct {
	nodes [12]*ExecNode
	stop  chan struct{}
	wg    *sync.WaitGroup

	// toggles tis-100 emulation
	emu bool
}

type connection struct {
	A int
	P Port
	B int
}

var defaultConnections = []connection{
	{0, RIGHT, 1},
	{0, DOWN, 4},
	{1, RIGHT, 2},
	{1, DOWN, 5},
	{2, RIGHT, 3},
	{2, DOWN, 6},
	{3, DOWN, 7},

	{4, RIGHT, 5},
	{4, DOWN, 8},
	{5, RIGHT, 6},
	{5, DOWN, 9},
	{6, RIGHT, 7},
	{6, DOWN, 10},
	{7, DOWN, 11},

	{8, RIGHT, 9},
	{9, RIGHT, 10},
	{10, RIGHT, 11},
}

func NewSystem() *System {
	sys := &System{
		stop: make(chan struct{}),
		wg:   new(sync.WaitGroup),
	}

	for i := range sys.nodes {
		sys.nodes[i] = &ExecNode{BasicNode: &BasicNode{step: make(chan struct{})}}
	}

	for _, cn := range defaultConnections {
		sys.nodes[cn.A].Connect(cn.P, sys.nodes[cn.B].BasicNode) // TODO yuck .Node
	}

	for _, n := range sys.nodes {
		go n.Start(sys.stop, sys.wg)
	}

	return sys
}

// TODO return a Node interface
func (sys *System) Node(i int) *ExecNode {
	return sys.nodes[i]
}

// TODO check the integrity of system nodes. Given the public api
// may allow breaking the system, run check before stepping through
// programs. Better yet, don't allow breaking the system ..
func (sys *System) Check() error {
	return nil
}

func (sys *System) Start() {
	if err := sys.Check(); err != nil {
		panic(err)
	}
	go func() {
		for {
			select {
			case <-sys.stop:
				return
			default:
				sys.Step()
			}
		}
	}()
}

func (sys *System) Stop() {
	close(sys.stop)
	sys.stop = make(chan struct{})
}

func (sys *System) Step() {
	// TODO check integrity every step? depends on use

	sys.wg.Add(len(sys.nodes))
	for _, n := range sys.nodes {
		n.step <- struct{}{}
	}
	sys.wg.Wait()
}
