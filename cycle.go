package tis100mc

import "sync"

// TODO completely remove node.step and replace with some type of
// Cycler interface. The implementation will track cycle counts and
// also cause blocks within the go routine for emulation purposes.

type Cycler struct {
	sync.WaitGroup

	count int

	step   chan struct{}
	waitio chan struct{}
}

func NewCycler(emu bool) *Cycler {
	cy := &Cycler{
		step:   make(chan struct{}, 1),
		waitio: make(chan struct{}),
	}
	close(cy.waitio)

	// if !emu {
	// close(cy.step)
	// }
	return cy
}

func (cy *Cycler) Cycle() {
	<-cy.step
	cy.count++
}

func (cy *Cycler) Add(delta int) {
	select {
	case <-cy.waitio:
		cy.WaitGroup.Add(delta)
		cy.step <- struct{}{}
	default:
	}
}

func (cy *Cycler) WaitIO() {
	cy.Wait()
	<-cy.waitio
}
