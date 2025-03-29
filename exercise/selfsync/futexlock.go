package selfsync

import (
	"sync/atomic"
	"syscall"
)

/*
Futex stands for fast user space mutex. Although it is not really a mutex, it
is a wait queue primitive that can be accessed from user space. By user space,
we basically mean by our program, not necessarily by the operating system's
scheduler. It gives us the ability to suspend or awaken executions on certain
addresses.

There are two functions involved with this capacity:
- futex_wait
- futex_wake

Usually these are named somewhat differently depending on the operating
system, but their functionality is mostly same.
The wait suspends the execution of the caller and sends it to the back of
the wait queue. The wake wakes up a given number of suspended executions that
are waiting on a memory address
*/

/*
futex_wait:
*/
func futex_wait(address *int32, value int32) {
	//
	var val_ptr uintptr = uintptr(value)
	address_ptr := (uintptr)(address)
	syscall.Syscall(uintptr(syscall.SYS_FUTEX), address, val_ptr)
}
