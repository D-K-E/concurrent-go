package selectscenario

import (
	"fmt"
	"math"
	"time"
)

/*
this scenario concerns performing concurrent computations on default case
and stopping all concurrent actions using select statement
*/

const (
	passToGuess = "nopass"
	chars       = " abcdefghijklmnopqrstuvwxyz"
	charlen     = len(chars)
)

// helper function to convert integer to string. By working with integers we
// can divide the search into multiple ranges
func toBaseChar(n int) string {
	result := ""
	for n > 0 {
		result = string(chars[n%charlen]) + result
		n /= charlen
	}
	return result
}

func guessPass(from int, upto int, stop chan int, result chan string,
) {
	// iterate from variable to upto variable
	for guessNb := from; guessNb < upto; guessNb++ {
		select {
		case _, isOpen := (<-stop):
			fmt.Printf("Stopped at %d [%d,%d)\n", guessNb, from, upto)
			if isOpen {
				close(stop)
			}
			return
		default:

			// transform number to string
			guessedPass := toBaseChar(guessNb)

			// check if it matches to the password we are trying to guess
			// normally we would not know passToGuess but would attempt to
			// access to the resource with the guessedPass
			if guessedPass == passToGuess {
				result <- guessedPass
				// send signal to stop channel
				close(stop)
				return
			}
		}
	}
	fmt.Printf("Not found between [%d,%d)\n", from, upto)
}

func SelectScenarioDefaultCaseMain() {
	stopChannel := make(chan int)

	passChannel := make(chan string)

	stepSize := 10000000

	maxSize := int(math.Pow(float64(len(chars)), float64(len(passToGuess))))

	for i := 1; i < maxSize; i += stepSize {
		go guessPass(i, i+stepSize, stopChannel, passChannel)
		//)
	}

	// notice that this blocks main thread until passChannel is populated with
	// a message. If the pass channel doesn't get anything, all threads go to
	// sleep and we have a deadlock
	for msg := range passChannel {
		fmt.Println("password found", msg)
		close(passChannel)
	}
	time.Sleep(5 * time.Second)
}
