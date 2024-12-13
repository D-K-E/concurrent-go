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

