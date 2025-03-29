package selfsync

// spin lock implementation
import (
	"runtime"
	"sync"
	"sync/atomic"
)

/*
This is an atomic function that checks whether swapping the value of `ptr` is
possible. If `ptr` value is same as old value, then it is possible to swap.
The function would swap the value of `ptr` with `newvalue` and outputs true.
If `ptr` value is not same as old value, then the swap is not possible and
function outputs false. We use a CompareAndSwapInt32 from atomic
*/
func SwapIfEqualInt32(ptr *int32, oldValue, newValue int32) bool {
	result := atomic.CompareAndSwapInt32(ptr, oldValue, newValue)
	return result
}

type SpinLockInt32 int32

func (s *SpinLockInt32) Lock() {
	lockAsInt := (*int32)(s) // cast to int32 pointer

	// if swap occurs we call the scheduler to give execution time
	// to other schedules
	for !SwapIfEqualInt32(lockAsInt, 0, 1) {
		runtime.Gosched()
	}
}

func (s *SpinLockInt32) Unlock() {
	lockAsInt := (*int32)(s) // cast to int32 pointer
	atomic.StoreInt32(lockAsInt, 0)
}

func NewSpinLock() sync.Locker {
	var lock SpinLockInt32
	return &lock
}
