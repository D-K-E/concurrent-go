package selectscenario

import (
	"fmt"
	"strconv"
	"time"
)

// this snippet shows delaying unblocking the main thread for a specified
// number of time. So far, we had seen blocking/unblocking mechanism that are
// defined based on the completion of various functions like in wait groups
// which waits until certain operations are complete, or barriers which waits
// before certain operations start. This snippet shows how to wait for an
// operation given time using channels.

func sendMsgAfter(seconds time.Duration) <-chan string {
	messages := make(chan string)

	go func() {
		time.Sleep(seconds)
		messages <- "hello"
	}()
	return messages
}

func TimeoutMain() {
	t, _ := strconv.Atoi("4")
	messages := sendMsgAfter(3 * time.Second)
	timeOut := time.Duration(t) * time.Second
	fmt.Printf("Waiting for message for %d seconds...\n", t)
	select {
	case msg := (<-messages):
		fmt.Println("Message received:", msg)
	case tNow := (<-time.After(timeOut)):
		fmt.Println("Timed out. Waited until:", tNow.Format("15:04:05"))
	}
}
