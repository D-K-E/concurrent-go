package numfactor

import (
	"fmt"
	"math/rand"
	"sync"
)

func findFactor(number int) []int {
	result := make([]int, 0)
	for i := 1; i < number; i++ {
		if (number % i) == 0 {
			result = append(result, i)
		}
	}
	return result
}

func NumFactorMain() {
	resultCh := make(chan []int)
	go func() {
		resultCh <- findFactor(3129402)
	}()

	nextResult := findFactor(3920147)
	fmt.Println("next result")
	fmt.Println(nextResult)
	firstResult := (<-resultCh)
	fmt.Println(firstResult)
}

func printResult(resultCh <-chan []int) {
	for r := range resultCh {
		fmt.Println(r)
	}
}

func randFactor(resultCh chan<- []int, waitG *sync.WaitGroup) {
	num := rand.Intn(900032)
	resultCh <- findFactor(num)
	waitG.Done()
}

func NumFactorMain2() {
	resultCh := make(chan []int)
	wg := sync.WaitGroup{}
	go printResult(resultCh)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go randFactor(resultCh, &wg)
	}
	wg.Wait()
	close(resultCh)
}
