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

func SwapIfEqualBool(ptr *atomic.Bool, oldValue, newValue bool) bool {
	result := ptr.CompareAndSwap(oldValue, newValue)
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

type SpinLockBool atomic.Bool

func (s *SpinLockBool) Lock() {
	lockAsBool := (*atomic.Bool)(s) // cast to bool pointer

	// if swap occurs we call the scheduler to give execution time
	// to other schedules
	for !SwapIfEqualBool(lockAsBool, false, true) {
		runtime.Gosched()
	}
}

func (s *SpinLockBool) Unlock() {
	lockAsBool := (*atomic.Bool)(s) // cast to bool pointer
	lockAsBool.Store(false)
}

func (s *SpinLockBool) TryLock() bool {
	lockAsBool := (*atomic.Bool)(s) // cast to bool pointer
	isLocked := SwapIfEqualBool(lockAsBool, false, true)
	return isLocked
}
