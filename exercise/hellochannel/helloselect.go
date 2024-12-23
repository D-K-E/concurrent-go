package hellochannel

import (
	"fmt"
	"time"
)

func sendMsgAfter(seconds time.Duration) <-chan string {
	msg := make(chan string)
	go func() {
		time.Sleep(seconds)
		msg <- "hello"
	}()
	return msg
}

func HelloSelectMain() {
	msg := sendMsgAfter(3 * time.Second)
	for {
		select { // picks whichever case is possible
		// if none are available applies default branch
		case m := (<-msg):
			fmt.Println("received", m)
			return
		default:
			fmt.Println("no message received")
			time.Sleep(1 * time.Second)
		}
	}
}
