package selectscenario

import (
	"fmt"
)

/* this scenario concerns how to use quit channel pattern.
Quit channel pattern implies that we have a function that consumes from
multiple channels until the quit channel signals that all consumption must be
terminated
*/

// we have a function that prints generated numbers 10 times. After that we
// call the quits
func printNumbers(numbers <-chan int, quit chan int) {
	go func() {
		for i := 0; i < 10; i++ {
			number := (<-numbers)
			fmt.Println(number)
		}
		close(quit)
	}()
}

func QuitChannelMain() {
	numbers := make(chan int)
	quit := make(chan int)
	printNumbers(numbers, quit)

	next := 0
	for i := 1; ; i++ {
		next += i
		select {
		case numbers <- next:
		case <-quit:
			fmt.Println("Quitting number generation")
			return
		}
	}
}
