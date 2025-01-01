package selectscenario

import (
	"fmt"
	"math/rand"
	"time"
)

// this snippet shows how to disable a select case when using select with
// multiple channels. Assuming to channels are closing at different times and
// we want to wait until both them are complete before we move out of select
// statement, we need a mechanism to disable a case without terminating the
// entire select block. Remember that closing a channel does not disable a
// case, it will keep using the default value of the type of the channel, so
// we need something to effectively block the channel

func generateAmounts(n int) <-chan int {
	amounts := make(chan int)
	go func() {
		defer close(amounts)
		for i := 0; i < n; i++ {
			amounts <- rand.Intn(100) + i
			time.Sleep(100 * time.Millisecond) // write some random value at
			// this interval
		}
	}()
	return amounts
}

func SelectNilMain() {
	sales := generateAmounts(50)
	expenses := generateAmounts(40)

	endOfDayAmount := 0
	for sales != nil || expenses != nil {
		select {
		case sale, isOpen := (<-sales):
			if isOpen {
				fmt.Println("sale of", sale)
				endOfDayAmount += sale
			} else {
				// channel is closed so we assigned nil and block this case
				// forever from executing
				sales = nil
			}
		case expense, isOpen := (<-expenses):
			if isOpen {
				fmt.Println("expense of", expense)
				endOfDayAmount -= expense
			} else {
				// channel is closed so we assigned nil and block this case
				// forever from executing
				expenses = nil
			}
		}
	}
}
