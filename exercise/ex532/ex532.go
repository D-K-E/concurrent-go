package ex532

// Change the game-sync listings 5.8 and 5.9 so that, still using condition
// variables, the players wait for a fixed number of seconds. If the players
// haven’t all joined within this time, the goroutines should stop waiting and
// let the game start without all the players. Hint: try using another
// goroutine with an expiry timer.

import (
	"fmt"
	"sync"
	"time"
)

func playerHandler(cond *sync.Cond, playersRemaining *int, playerId int,
	gameHasStarted *bool,
) {
	// lock the mutex of the condition variable to avoid race condition
	cond.L.Lock() // lock since gameHasStarted and playersRemaining are shared
	if !(*gameHasStarted) {
		fmt.Println(playerId, ": Connected")
		remaining := *playersRemaining
		remaining--
		*playersRemaining = remaining
	}
	// no other player are remaining time to unblock all the waiting threads
	if *playersRemaining == 0 {
		cond.Broadcast()
	}

	// wait until all players have connected or time has elapsed
	for *playersRemaining > 0 && !(*gameHasStarted) {
		fmt.Println(playerId, ": Waiting for more players")
		cond.Wait()
	}
	// unlock all goroutines and resume execution
	cond.L.Unlock()
	if *gameHasStarted {
		fmt.Println("Game has started: ", playerId)
	} else {
		fmt.Println("All players are connected. Ready player", playerId)
	}
}

func isExpired(start *time.Time, threshold float32) bool {
	elapsed := time.Since(*start).Seconds()
	return float32(elapsed) >= threshold
}

func expiryTimer(cond *sync.Cond, start *time.Time,
	threshold float32, gameHasStarted *bool,
) {
	for !isExpired(start, threshold) {
		// while not expired let other goroutine's work
	}

	// time is expired time to move on to game
	cond.L.Lock() // lock since gameHasStarted is shared
	fmt.Println("time's up! Starting game!")
	*gameHasStarted = true
	cond.Broadcast() // signal all waiting threads to proceed with their business
	cond.L.Unlock()
}

func Ex532Main() {
	varMutex := sync.Mutex{}

	// declaring conditional variable
	cond := sync.NewCond(&varMutex)

	start := time.Now()
	var threshold float32 = 3.2
	gameHasStarted := false
	go expiryTimer(cond, &start, threshold, &gameHasStarted)
	time.Sleep(2 * time.Second)

	playersInGame := 4
	for playerId := 0; playerId < 4; playerId++ {
		go playerHandler(cond, &playersInGame, playerId, &gameHasStarted)
		time.Sleep(1 * time.Second)
	}
}
