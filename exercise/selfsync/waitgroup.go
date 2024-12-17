package selfsync

import "sync"

type StaticWaitGroup struct {
	semaphore *Semaphore
}

func NewStaticGroup(size int) *StaticWaitGroup {
	sema := NewSemaphore(1 - size)
	wg := StaticWaitGroup{semaphore: sema}
	return &wg
}

func (wg *StaticWaitGroup) Wait() {
	wg.semaphore.Acquire()
}

func (wg *StaticWaitGroup) Done() {
	wg.semaphore.Release()
}

type DynamicWaitGroup struct {
	groupSize int
	cond      *sync.Cond
}

func NewDynamicGroup() *DynamicWaitGroup {
	mut := sync.Mutex{}
	cond := sync.NewCond(&mut)
	wg := DynamicWaitGroup{cond: cond, groupSize: 0}
	return &wg
}

func (wg *DynamicWaitGroup) AddMember(delta int) {
	wg.cond.L.Lock()
	wg.groupSize += delta
	wg.cond.L.Unlock()
}

func (wg *DynamicWaitGroup) WaitMembers() {
	wg.cond.L.Lock()
	for wg.groupSize > 0 {
		wg.cond.Wait()
	}
	wg.cond.L.Unlock()
}

func (wg *DynamicWaitGroup) TaskDone() {
	wg.cond.L.Lock()
	wg.groupSize--
	if wg.groupSize == 0 {
		wg.cond.Broadcast()
	}
	wg.cond.L.Unlock()
}
