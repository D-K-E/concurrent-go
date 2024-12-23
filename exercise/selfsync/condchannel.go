package selfsync

// condition variable based channel implementation
import (
	"container/list"
	"sync"
)

type CondChannel[ChannelType any] struct {
	cond    *sync.Cond //
	maxSize int
	buffer  *list.List // linked list as a queue data structure
}

func NewCondChannel[ChannelType any](capacity int) *CondChannel[ChannelType] {
	mut := sync.Mutex{}
	cc := sync.NewCond(&mut)
	channel := CondChannel[ChannelType]{
		maxSize: capacity,
		cond:    cc,
		buffer:  list.New(),
	}
	return &channel
}

func (c *CondChannel[ChannelType]) Send(message ChannelType) {
	c.cond.L.Lock() /* adds a message to the buffer queue while
	   protecting against race conditions by using a mutex */
	for c.buffer.Len() == c.maxSize {
		c.cond.Wait()
	}
	c.buffer.PushBack(message)
	c.cond.Broadcast()
	c.cond.L.Unlock()
}

func (c *CondChannel[ChannelType]) Receive() ChannelType {
	c.cond.L.Lock()
	c.maxSize++
	c.cond.Broadcast()
	for c.buffer.Len() == 0 {
		c.cond.Wait()
	}
	c.maxSize--
	v := c.buffer.Remove(c.buffer.Front()).(ChannelType)
	c.cond.L.Unlock()
	return v
}
