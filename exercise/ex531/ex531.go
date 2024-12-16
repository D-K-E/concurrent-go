// exercise 1 at section 5.3
package ex531

// n listing 5.4, Stingy’s goroutine is signaling on the condition variable
// every time we add money to the bank account. Can you change the function so
// that it signals only when there is $50 or more in the account?

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func stingy(money *int, cond *sync.Cond) {
	for i := 0; i < 1000000; i++ {
		cond.L.Lock()
		*money += 10
		if *money >= 50 {
			cond.Signal()
		}
		cond.L.Unlock()
	}
	fmt.Println("Stingy done")
}

func spendy(money *int, cond *sync.Cond) {
	for i := 0; i < 2000000; i++ {
		cond.L.Lock()
		for *money < 50 {
			cond.Wait()
		}
		*money -= 50
		if *money < 0 {
			fmt.Println("Money is negative")
			os.Exit(1)
		}
		cond.L.Unlock()
	}
	fmt.Println("Spendy done")
}

func Ex531Main() {
	money := 100
	mutex := sync.Mutex{}
	cond := sync.NewCond(&mutex)
	go stingy(&money, cond)
	go spendy(&money, cond)
	time.Sleep(2 * time.Second)
	mutex.Lock()
	fmt.Println("money in bank account: ", money)
	mutex.Unlock()
}
