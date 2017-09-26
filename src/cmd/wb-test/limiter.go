package main

import (
	"sync"
)

// limiter
// Limits the number of goroutines
// and allows you to wait
// for the worker group to complete.
type limiter struct {
	waitGroup *sync.WaitGroup
	limit     chan struct{}
}

// Returns new limiter which restrains
// the size of worker group with number
// less or equal k.
func newLimiter(k uint) *limiter {
	return &limiter{
		waitGroup: &sync.WaitGroup{},
		limit:     make(chan struct{}, k),
	}
}

// Adds a new worker to groups
// and returns "done" function,
// which must be run to decrement
// worker counter.
func (lim *limiter) Start() func() {
	lim.waitGroup.Add(1)
	lim.limit <- struct{}{}
	return func() {
		// Important!
		// Do not change order of pull
		// from channel and call of .Done()!
		<-lim.limit
		lim.waitGroup.Done()
	}
}

// Blocks until all workers complete their tasks.
func (lim *limiter) Wait() {
	// You haven't to close channel
	// after limiter and "done" closures dropped
	lim.waitGroup.Wait()
}

// Returns a number of active workers.
func (lim *limiter) Active() int {
	return len(lim.limit)
}
