package semadowork

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func doworkWaitGroup(id int, waitGroup *sync.WaitGroup) {
	i := rand.Intn(4) // sleeps for random seconds to simulate a random task
	time.Sleep(time.Duration(i) * time.Second)
	fmt.Println(id, "done working after", i, "seconds")
	waitGroup.Done() // signals that goroutine has finished its task
}

func WaitGroupDoWorkMain() {
	//
	wg := sync.WaitGroup{} // create wait group
	nbWaitGroup := 4       // group size
	wg.Add(nbWaitGroup)    // create group with group size
	for i := 0; i < nbWaitGroup; i++ {
		go doworkWaitGroup(i, &wg)
	}
	wg.Wait() // wait until all members of the wait group calls for Done
	fmt.Println("All complete")
}
