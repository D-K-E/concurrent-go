package selectscenario

import (
	"fmt"
	"math"
	"math/rand"
)

/*
The goal of this snippet is to show how select statement can be used to switch
between reading and writing between various channels. This is very similar to
assembly line where you have a stream of input and you want to move certain
items in one line to another line.
*/

func primesOnly(inputs <-chan int) <-chan int {
	results := make(chan int) // will contain prime numbers
	go func() {
		for c := range inputs {
			isPrime := c != 1
			for i := 2; i <= int(math.Sqrt(float64(c))); i++ {
				reminder := c % i
				if reminder == 0 {
					isPrime = false
					break
				}
			}
			if isPrime {
				results <- c // if prime push to results
			}
		}
	}()
	return results
}

func PrimesOnlyMain() {
	numberChannel := make(chan int)
	primeResult := primesOnly(numberChannel)
	for i := 0; i < 100; {
		select {
		case numberChannel <- rand.Intn(100000000) + 1: // adds random number
		case primeNumber := (<-primeResult):
			fmt.Println("found a prime number:", primeNumber)
			i++
		}
	}
}
