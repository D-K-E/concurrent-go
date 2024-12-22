package hellochannel

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	selfsync "github.com/D-K-E/concurrent-go/selfsync"
)

func receiverStringSelf(messages *selfsync.Channel[string]) {
	msg := ""
	for msg != "STOP" {
		msg = messages.Receive()
		fmt.Println("received", msg, "!")
	}
}

func sendStringSelf(messages *selfsync.Channel[string],
	wg *sync.WaitGroup, nbMessages int,
) {
	for i := 0; i < nbMessages; i++ {
		fmt.Println(time.Now().Format("15:04:05"), "Sending", i)
		istr := strconv.Itoa(i)
		messages.Send(istr)
		time.Sleep(1 * time.Second)
	}
	wg.Done()
}

func HelloSelfChannelMain() {
	//
	msgChannel := selfsync.NewChannel[string](3) // buffer size of 3
	wg := sync.WaitGroup{}
	go receiverStringSelf(msgChannel)

	nbMessages := [2]int{3, 6}
	for _, value := range nbMessages {
		wg.Add(1)
		go sendStringSelf(msgChannel, &wg, value)
	}
	wg.Wait()
	msgChannel.Send("STOP")
	time.Sleep(1 * time.Second)
}
