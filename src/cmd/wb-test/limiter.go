package main

import (
	"sync"
)

type limiter struct {
	waitGroup *sync.WaitGroup
	limit     chan struct{}
}

func newLimiter(k uint) *limiter {
	return &limiter{
		waitGroup: &sync.WaitGroup{},
		limit:     make(chan struct{}, k),
	}
}

func (lim *limiter) Start() func() {
	lim.waitGroup.Add(1)
	lim.limit <- struct{}{}
	return func() {
		<-lim.limit
		lim.waitGroup.Done()
	}
}

func (lim *limiter) Wait() {
	lim.waitGroup.Wait()
	close(lim.limit)
}

func (lim *limiter) Active() int {
	return len(lim.limit)
}
