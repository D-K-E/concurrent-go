package ex93

import "fmt"

/*
Write a generator goroutine similar to listing 9.2 that, instead of generating
URL strings, generates an infinite stream of square numbers (1, 4, 9, 16, 25 .
. .) on an output channel. Here is the signature:
*/
func GenerateSquares(quit <-chan int) <-chan int {
	//
	squares := make(chan int)
	go func() {
		defer close(squares)
		for i := 0; true; i++ {
			i2 := i * i
			squares <- i2
		}
	}()
	return squares
}

func GenerateSquaresMain() {
	qu := make(chan int)
	defer close(qu)
	resultChannel := GenerateSquares(qu)

	// consume result channel
	for msg := range resultChannel {
		// consume
		fmt.Println(msg)
	}
}
