package selfsync

import "sync"

type Barrier struct {
	size      int // total number of participants to barrier
	waitCount int /* counter variable representing the number of currently
	   suspended executions */
	cond *sync.Cond
}

// creates a new barrier
func NewBarrier(size int) *Barrier {
	m := sync.Mutex{}
	condVar := sync.NewCond(&m)
	barrier := Barrier{size: size, waitCount: 0, cond: condVar}
	return &barrier
}

func (b *Barrier) Wait() {
	b.cond.L.Lock() // protect access to wait count
	b.waitCount += 1

	// wait count had been reached, so all routines fulfilled their tasks
	if b.waitCount == b.size {
		b.waitCount = 0
		b.cond.Broadcast()
	} else {
		b.cond.Wait() // not reached so we wait
	}
	b.cond.L.Unlock()
}
