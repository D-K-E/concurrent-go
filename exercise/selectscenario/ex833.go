package selectscenario

// exercise 8.3.3

import (
	"fmt"
	"math/rand"
	"time"
)

func player() chan string {
	output := make(chan string)
	count := rand.Intn(100)
	move := []string{"UP", "DOWN", "LEFT", "RIGHT"}
	go func() {
		defer close(output)
		for i := 0; i < count; i++ {
			output <- move[rand.Intn(4)]
			d := time.Duration(rand.Intn(200))
			time.Sleep(d * time.Millisecond)
		}
	}()
	return output
}

func Ex833Main() {
	outCh1 := player()
	outCh2 := player()
	outCh3 := player()
	outCh4 := player()
	remainingPlayers := 4
	for {
		select {
		case msg, isOpen := (<-outCh1):
			if !isOpen {
				outCh1 = nil
				remainingPlayers--
				if remainingPlayers == 1 {
					return
				}
			} else {
				fmt.Println("player 1", msg)
			}
		case msg, isOpen := (<-outCh2):
			if !isOpen {
				outCh2 = nil
				remainingPlayers--
				if remainingPlayers == 1 {
					return
				}
			} else {
				fmt.Println("player 2", msg)
			}
		case msg, isOpen := (<-outCh3):
			if !isOpen {
				outCh3 = nil
				remainingPlayers--
				if remainingPlayers == 1 {
					return
				}
			} else {
				fmt.Println("player 3", msg)
			}
		case msg, isOpen := (<-outCh4):
			if !isOpen {
				outCh4 = nil
				remainingPlayers--
				if remainingPlayers == 1 {
					return
				}
			} else {
				fmt.Println("player 4", msg)
			}
		}
	}
}
