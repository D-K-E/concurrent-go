package hellochannel

// very simple channel direction program

import (
	"fmt"
	"time"
)

func msgReceiver(messages <-chan int) {
	msg := 0
	for msg != -1 {
		msg = (<-messages)
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
	}
}

func msgSender(messages chan<- int) {
	for i := 1; ; i++ {
		fmt.Println(time.Now().Format("15:04:05"), "Sending", i)
		messages <- i
		time.Sleep(1 * time.Second)
	}
}

func HelloDirectionChannelMain() {
	msgChannel := make(chan int, 3)
	go msgReceiver(msgChannel)
	go msgSender(msgChannel)
	time.Sleep(5 * time.Second)
}
