package ex93

import (
	"fmt"

	selfsync "github.com/D-K-E/concurrent-go/selfsync"
)

/*
Connect the components developed in exercises 1 to 4 together in a main()
function using the following pseudocode:
*/

func Ex935Main() {
	quit := make(chan int)
	defer close(quit)
	squares := GenerateSquares(quit)
	taken := selfsync.TakeUntil(func(s int) bool { return s <= 10000 },
		quit, squares)
	//
	printed := selfsync.Print(quit, taken)
	<-selfsync.Drain(quit, printed)
}
