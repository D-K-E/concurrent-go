package selfsync

import "sync"

type ReadWriteMutex struct {
	readerCounter  int        // stores the number of readers currently holding read lock
	writersWaiting int        // stores the number of writers currently waiting
	writerActive   bool       // indicates if a writer is holding the writer lock
	cond           *sync.Cond //
}

func NewMutex() *ReadWriteMutex {
	m := sync.Mutex{}
	rwmutex := ReadWriteMutex{
		readerCounter:  0,
		writersWaiting: 0,
		writerActive:   false,
		cond:           sync.NewCond(&m),
	}
	return &rwmutex
}

func (rw *ReadWriteMutex) ReadLock() {
	rw.cond.L.Lock() // acquires mutex
	for rw.writersWaiting > 0 || rw.writerActive {
		rw.cond.Wait() // wait on condition variable while writers are waiting or active
	}
	rw.readerCounter++ // increment reader counter
	rw.cond.L.Unlock() // release mutex
}

func (rw *ReadWriteMutex) WriteLock() {
	rw.cond.L.Lock()    // acquires mutex
	rw.writersWaiting++ // increments the writer's waiting counter
	for rw.readerCounter > 0 || rw.writerActive {
		rw.cond.Wait() // waits on condition variable as long as there are
		// readers or an active writer
	}
	rw.writersWaiting-- // once wait is over, decrements the writer's waiting
	// counter
	rw.writerActive = true // once the wait is over, marks writer active flag

	rw.cond.L.Unlock() // releases mutex
}

func (rw *ReadWriteMutex) ReadUnlock() {
	rw.cond.L.Lock()   // acquires mutex
	rw.readerCounter-- // decrements reader's counter by 1
	if rw.readerCounter == 0 {
		rw.cond.Broadcast() // sends broadcast if the goroutine is the last
		// remaining reader
	}
	rw.cond.L.Unlock()
}

func (rw *ReadWriteMutex) WriteUnlock() {
	rw.cond.L.Lock()
	rw.writerActive = false // unmarks writer active flag
	rw.cond.Broadcast()     // sends signals to all goroutines
	rw.cond.L.Unlock()
}
