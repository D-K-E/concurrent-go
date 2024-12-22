package selfsync

// condition variable based channel implementation
import (
	"container/list"
	"sync"
)

type CondChannel[ChannelType any] struct {
	capacityCounter int        // counts remaining capacity until full if full block sender
	capacityCond    *sync.Cond //
	sizeCounter     int        /* buffer size counter to block sender when
	   buffer is empty */
	sizeCond *sync.Cond
	buffer   *list.List // linked list as a queue data structure
	mutex    *sync.Mutex
}

func NewCondChannel[ChannelType any](capacity int) *CondChannel[ChannelType] {
	mut := sync.Mutex{}
	cc := sync.NewCond(&mut)
	sc := sync.NewCond(&mut)
	channel := CondChannel[ChannelType]{
		capacityCounter: capacity,
		capacityCond:    cc,
		sizeCounter:     0,
		sizeCond:        sc,
		buffer:          list.New(),
		mutex:           &mut,
	}
	return &channel
}

func (c *CondChannel[ChannelType]) Send(message ChannelType) {
	c.capacityCond.L.Lock() /* adds a message to the buffer queue while
	   protecting against race conditions by using a mutex */
	for c.capacityCounter <= 0 {
		c.capacityCond.Wait()
	}
	c.capacityCounter-- // acquire one permit from the capacity semaphore
	c.capacityCond.L.Unlock()
	c.mutex.Lock()
	c.buffer.PushBack(message)
	c.mutex.Unlock()

	c.sizeCond.L.Lock()
	c.sizeCounter++ // Release one permit from the buffer size semaphore
	c.sizeCond.Signal()
	c.sizeCond.L.Unlock()
}

func (c *CondChannel[ChannelType]) Receive() ChannelType {
	c.capacityCond.L.Lock()
	c.capacityCounter++ // release one permit from the capacity semaphore
	c.capacityCond.Signal()
	c.capacityCond.L.Unlock()

	c.sizeCond.L.Lock()
	for c.sizeCounter <= 0 {
		c.sizeCond.Wait()
	}
	c.sizeCounter-- // Acquires one permit from the buffer size semaphore
	c.sizeCond.L.Unlock()
	c.mutex.Lock()
	v := c.buffer.Remove(c.buffer.Front()).(ChannelType) /*
	   Removes one message from the buffer while protecting against race
	   conditions using the mutex
	*/
	c.mutex.Unlock()
	return v
}
