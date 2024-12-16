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

func SemaphoreDoWorkMain() {
	semaphore := selfsync.NewSemaphore(0)
	for i := 0; i < 50; i++ {
		go doWork(semaphore) // starts goroutine passing a reference to semaphore
		fmt.Println("Waiting for child goroutine")
		semaphore.Acquire() // waits for available permit on the semaphore
		// indicating the task is complete
		fmt.Println("Child goroutine finished")
	}
}
