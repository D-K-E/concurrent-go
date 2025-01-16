package selectscenario

import (
	"fmt"
)

/*
A demonstration of first class channels pattern. Basically we can use channels
recursively to grow the pipeline dynamically.
*/

func primeMultiFilter(numbers <-chan int, quit chan<- int) {
	var right chan int
	p := (<-numbers) // get the first number from the channel
	fmt.Println("number received:", p)
	for n := range numbers { // read next numbers from the channel

		/* make sure that the next number is not multiple of
		   previous */
		if (n % p) != 0 {
			// we found a number but we don't have an output channel
			if right == nil {
				right = make(chan int)

				// notice that we found a number, now we create a new channel
				// to push the number and continue with the search
				go primeMultiFilter(right, quit)
			}
			// now we push the number
			right <- n
		}
	}
	if right == nil {
		// finished going through all the numbers, yet we were unable to
		// instantiate the channel, meaning that all the subsequent numbers
		// are multiple of the previous number
		close(quit)
	} else {
		//
		// finished going through all the numbers, and we were able to
		// instantiate it, so now it is being consumed by another go routine
		// we can safely close it, since there will be no more numbers going
		// through it
		close(right)
	}
}

func FirstClassChannelMain() {
	numbers := make(chan int)
	quit := make(chan int)
	go primeMultiFilter(numbers, quit)
	for i := 2; i < 10000; i++ {
		numbers <- i
	}
	close(numbers)
	// now wait for quit signal
	<-quit
	for range quit {
		fmt.Println("quit signal received!")
	}
}
