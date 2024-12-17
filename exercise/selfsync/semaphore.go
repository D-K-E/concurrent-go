package selfsync

import "sync"

// Semaphore allows a fixed number of permits that enable concurrent
// executions to access shared resources. A mutex ensures that only a single
// goroutine has exclusive access, whereas a semaphore ensures that at most N
// goroutines have access
type Semaphore struct {
	permits int // permits remaining on the semaphore
	cond    *sync.Cond
}

func NewSemaphore(n int) *Semaphore {
	m := sync.Mutex{}
	sema := Semaphore{
		permits: n,
		cond:    sync.NewCond(&m),
	}
	return &sema
}

// organizes access to resources based on permits
func (rw *Semaphore) Acquire() {
	rw.cond.L.Lock()      // acquires mutex to protect permits variable
	for rw.permits <= 0 { // waits until there is a permit available
		rw.cond.Wait()
	}
	rw.permits--
	rw.cond.L.Unlock()
}

func (rw *Semaphore) Release() {
	rw.cond.L.Lock()
	rw.permits++
	rw.cond.Signal()
	// we signal that there is a permit available for possible acquisition
	rw.cond.L.Unlock()
}

// Weighted semaphore is a variation on a semaphore that allows you to acquire
// and release more than one permit at the same time.
type WeightedSemaphore struct {
	permits int // permits remaining on the semaphore
	cond    *sync.Cond
}

func NewGenericSemaphore[SemaphoreType WeightedSemaphore | Semaphore](n int) *SemaphoreType {
	m := sync.Mutex{}
	sema := SemaphoreType{
		permits: n,
		cond:    sync.NewCond(&m),
	}
	return &sema
}

func (rw *WeightedSemaphore) AcquirePermit(n int) {
	rw.cond.L.Lock()
	if rw.permits < n { // wait until we have enough permits
		rw.cond.Wait()
	}
	// we had waited enough, we should have enough permits at this point
	rw.permits -= n
	rw.cond.L.Unlock()
}

func (rw *WeightedSemaphore) ReleasePermit(n int) {
	rw.cond.L.Lock()
	rw.permits += n
	rw.cond.Signal() // we signal that there is n permits available for acquisition
	rw.cond.L.Unlock()
}
