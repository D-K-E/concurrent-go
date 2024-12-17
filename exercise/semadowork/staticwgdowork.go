package semadowork

import (
	"fmt"
	"math/rand"
	"time"

	selfsync "github.com/D-K-E/concurrent-go/selfsync"
)

func doWorkStaticWaitGroup(id int, waitGroup *selfsync.StaticWaitGroup) {
	i := rand.Intn(4) // sleeps for random seconds to simulate a random task
	time.Sleep(time.Duration(i) * time.Second)
	fmt.Println(id, "done working after", i, "seconds")
	waitGroup.Done()
}

func StaticWaitGroupDoWorkMain() {
	//
	nbWaitGroup := 4                           // group size
	wg := selfsync.NewStaticGroup(nbWaitGroup) // create wait group
	for i := 0; i < nbWaitGroup; i++ {
		go doWorkStaticWaitGroup(i, wg)
	}
	wg.Wait() // wait until all members of the wait group calls for Done
	fmt.Println("All complete")
}
