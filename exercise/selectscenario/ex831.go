package selectscenario

// exercise 8.3.1

import (
	"fmt"
	"math/rand"
	"time"
)

func generateTemp() chan int {
	out := make(chan int)

	go func() {
		temp := 50
		for {
			out <- temp
			temp += rand.Intn(3) - 1
			time.Sleep(200 * time.Millisecond)
		}
	}()
	return out
}

func outputTemp(input chan int) {
	go func() {
		for {
			msg := (<-input)
			fmt.Println("current temp", msg)
			time.Sleep(2 * time.Second)
		}
	}()
}

func Ex831Main() {
	inputChannel := make(chan int)
	outputChannel := generateTemp()
	msg := (<-outputChannel)
	outputTemp(inputChannel)
	for i := 0; i < 10; i++ {
		select {
		case msg = (<-outputChannel):
			fmt.Println("sent temp", msg)
		case inputChannel <- msg:

		}
		time.Sleep(1 * time.Second)
	}
}
