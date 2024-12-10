package gamesync

import (
	"fmt"
	"sync"
	"time"
)

func playerHandler(cond *sync.Cond, playersRemaining *int, playerId int) {
	// lock the mutex of the condition variable to avoid race condition
	cond.L.Lock()
	fmt.Println(playerId, ": Connected")
	remaining := *playersRemaining
	remaining--
	*playersRemaining = remaining
	// no other player are remaining time to unblock all the waiting threads
	if *playersRemaining == 0 {
		cond.Broadcast()
	}

	// wait until all players have connected
	for *playersRemaining > 0 {
		fmt.Println(playerId, ": Waiting for more players")
		cond.Wait()
	}
	// unlock all goroutines and resume execution
	cond.L.Unlock()
	fmt.Println("All players are connected. Ready player", playerId)
}

func GameSyncMain() {
	varMutex := sync.Mutex{}

	// declaring conditional variable
	cond := sync.NewCond(&varMutex)

	playersInGame := 4
	for playerId := 0; playerId < 4; playerId++ {
		go playerHandler(cond, &playersInGame, playerId)
		time.Sleep(1 * time.Second)
	}
}
