package semadowork

import (
	"fmt"
	"time"

	selfsync "github.com/D-K-E/concurrent-go/selfsync"
)

func doWorkAndWait(name string, timeToWork int, barrier *selfsync.Barrier) {
	start := time.Now()
	for {
		fmt.Println(time.Since(start), name, "is running")
		time.Sleep(time.Duration(timeToWork) * time.Second)
		fmt.Println(time.Since(start), name, "is waiting on barrier")
		barrier.Wait() // waits for others to catch up before proceeding to next line
	}
}

func BarrierMain() {
	nbMembers := 2
	barrier := selfsync.NewBarrier(nbMembers)
	go doWorkAndWait("Mahmut", 3, barrier)
	go doWorkAndWait("Ahmet", 8, barrier)
	time.Sleep(100 * time.Second)
}
