package selfsync

import (
	"container/list"
	"sync"
)

type Channel[ChannelType any] struct {
	capacitySema *Semaphore /* capacity semaphore to block sender when buffer is full */
	sizeSema     *Semaphore /* buffer size semaphore to block sender when
	   buffer is empty */
	mutex  sync.Mutex // Mutex protecting our shared data list structure
	buffer *list.List // linked list as a queue data structure
}

func NewChannel[ChannelType any](capacity int) *Channel[ChannelType] {
	channel := Channel[ChannelType]{
		capacitySema: NewSemaphore(capacity),
		sizeSema:     NewSemaphore(0),
		buffer:       list.New(),
	}
	return &channel
}

func (c *Channel[ChannelType]) Send(message ChannelType) {
	c.capacitySema.Acquire() // acquire one permit from the capacity semaphore
	c.mutex.Lock()           /* adds a message to the buffer queue while
	   protecting against race conditions by using a mutex */
	c.buffer.PushBack(message)
	c.mutex.Unlock()
	c.sizeSema.Release() // Release one permit from the buffer size semaphore
}

func (c *Channel[ChannelType]) Receive() ChannelType {
	c.capacitySema.Release() // release one permit from the capacity semaphore
	c.sizeSema.Acquire()     // Acquires one permit from the buffer size semaphore
	c.mutex.Lock()
	v := c.buffer.Remove(c.buffer.Front()).(ChannelType) /*
	   Removes one message from the buffer while protecting against race
	   conditions using the mutex
	*/
	c.mutex.Unlock()
	return v
}

