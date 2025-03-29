package deadlock

import (
	"fmt"
	"sync"
)

type Player struct {
	name  string
	score int
	mutex sync.Mutex
}

func incrementScores(players []*Player, increment int, wg *sync.WaitGroup) {
	for _, player := range players {
		player.mutex.Lock()
	}
	for _, player := range players {
		player.score += increment
	}
	for _, player := range players {
		player.mutex.Unlock()
	}
	wg.Done()
}

func incrementScoresV2(players []*Player, increment int, wg *sync.WaitGroup) {
	for _, player := range players {
		player.mutex.Lock()
		player.score += increment
		player.mutex.Unlock()
	}
	wg.Done()
}

func IncrementScoreMain() {
	p1 := Player{name: "map", score: 10, mutex: sync.Mutex{}}
	p2 := Player{name: "nap", score: 20, mutex: sync.Mutex{}}
	p3 := Player{name: "nar", score: 15, mutex: sync.Mutex{}}
	players := []*Player{&p1, &p2, &p3}
	increments := []int{2, 39, 1, 22, 12, 32, 3, 2, 1, 4, 5, 6, 76, 4, 34, 2, 9, 23, 20}
	wg := sync.WaitGroup{}
	for _, increment := range increments {
		wg.Add(1)
		go incrementScores(players, increment, &wg)
	}
	wg.Wait()
	fmt.Println("done")
}
