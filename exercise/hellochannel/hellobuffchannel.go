package hellochannel

// very simple channel example
import (
	"fmt"
	"sync"
	"time"
)

func receiverInt(messages chan int, wg *sync.WaitGroup) {
	msg := 0
	for msg != -1 {
		time.Sleep(1 * time.Second)
		msg = (<-messages)
		fmt.Println("Received:", msg)
	}
	wg.Done()
}

func HelloBufferedChannelMain() {
	//
	msgChannel := make(chan int, 3) // create a new channel with buffer of 3 messages
	wg := sync.WaitGroup{}
	wg.Add(1)
	go receiverInt(msgChannel, &wg)
	//
	for i := 1; i <= 6; i++ {
		size := len(msgChannel)
		fmt.Printf("%s Sending: %d. Buffer size: %d\n",
			time.Now().Format("15:04:05"), i, size)
		msgChannel <- i
	}
	msgChannel <- (-1)
	wg.Wait() // wait until the other wait groups are done
}
