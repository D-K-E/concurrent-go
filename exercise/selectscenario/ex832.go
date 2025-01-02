package selectscenario

// exercise 8.3.2
import (
	"fmt"
	"math/rand"
	"time"
)

func generateNumbers() chan int {
	output := make(chan int)
	go func() {
		for {
			output <- rand.Intn(10)
			time.Sleep(200 * time.Millisecond)
		}
	}()
	return output
}

func Ex832Main() {
	timeOut := 2 * time.Second
	messages := generateNumbers()
	timeout := time.After(timeOut)
	check := false
	for !check {
		select {
		case msg := (<-messages):
			fmt.Println("Message received:", msg)
		case tNow := (<-timeout):
			fmt.Println("Timed out. Waited until:", tNow.Format("15:04:05"))
			close(messages)
			check = true
		}
	}
}
