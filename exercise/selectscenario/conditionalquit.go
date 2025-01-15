package selectscenario

import (
	"fmt"

	selfsync "github.com/D-K-E/concurrent-go/selfsync"
)

/*
A demonstration of conditional quit pattern. We create a counter to stop all
work when a condition on counter (counter == 0) is met
*/

func ConditionalQuitMain() {
	quit := make(chan int)
	quitWords := make(chan int) // create a separate quit channel for condition
	defer close(quit)
	/* notice that we are closing the final quit
	   pipeline explicitly, since the other one is going to be closed by the
	   condition
	*/
	urls := generateUrls(quitWords)
	const nbDownloaders = 20
	pages := make([]<-chan string, nbDownloaders)
	for i := 0; i < nbDownloaders; i++ {
		// download all the pages
		pages[i] = downloadPages(quitWords, urls)
	}

	// join all the pages to a single channel
	merged := selfsync.FanIn(quitWords, pages...)
	words := extractWords(quitWords, merged)
	// take first 10000 words for the pipeline
	myWords := selfsync.Take(quitWords, 10000, words)

	multiWords := selfsync.Broadcast(quit, myWords, 2) // create 2 workers
	longest := longestWords(quit, multiWords[0])
	mostFrequent := frequentWords(quit, multiWords[1])

	// consume result channel
	fmt.Println("longest words are")
	for msg := range longest {
		fmt.Println(msg)
	}
	fmt.Println("most frequent words are")
	for msg := range mostFrequent {
		// consume
		fmt.Println(msg)
	}
}
