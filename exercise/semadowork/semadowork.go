package semadowork

import (
	"fmt"

	selfsync "github.com/D-K-E/concurrent-go/selfsync"
)

func doWork(semaphore *selfsync.Semaphore) {
	fmt.Println("Work started")
	fmt.Println("Work finished")
	semaphore.Release()
}

func doWorkWeighted(semaphore *selfsync.WeightedSemaphore, nbChildPermits int) {
	fmt.Println("Work started")
	fmt.Println("Work finished")
	semaphore.ReleasePermit(nbChildPermits)
}

func SemaphoreDoWorkMain() {
	semaphore := selfsync.NewGenericSemaphore[selfsync.Semaphore](0)
	for i := 0; i < 50; i++ {
		go doWork(semaphore) // starts goroutine passing a reference to semaphore
		fmt.Println("Waiting for child goroutine")
		semaphore.Acquire() // waits for available permit on the semaphore
		// indicating the task is complete
		fmt.Println("Child goroutine finished")
	}
}

func WeightedSemaphoreDoWorkMain() {
	semaphore := selfsync.NewGenericSemaphore[selfsync.WeightedSemaphore](0)
	nbChildPermits := 2
	for i := 0; i < 50; i++ {
		go doWorkWeighted(semaphore, nbChildPermits) // starts goroutine passing a reference to semaphore
		fmt.Println("Waiting for child goroutine")
		semaphore.AcquirePermit(nbChildPermits) // waits for available permit on the semaphore
		// indicating the task is complete
		fmt.Println("Child goroutine finished")
	}
}
