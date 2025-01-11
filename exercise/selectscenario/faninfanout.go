package selectscenario

import (
	"fmt"

	selfsync "github.com/D-K-E/concurrent-go/selfsync"
)

/*
A demonstration of Fan In - Fan Out pattern. Fan In means we are merging
multiple channel outputs to a single channel. Fan out means, we are
distributing the load to a fixed sized threads/routines.
*/

func FanInFanOutMain() {
	quit := make(chan int)
	defer close(quit)
	urls := generateUrls(quit)
	const nbDownloaders = 20
	pages := make([]<-chan string, nbDownloaders)
	for i := 0; i < nbDownloaders; i++ {
		// download all the pages
		pages[i] = downloadPages(quit, urls)
	}

	// join all the pages to a single channel
	merged := selfsync.FanIn(quit, pages...)
	words := extractWords(quit, merged)

	// consume result channel
	for msg := range words {
		// consume
		fmt.Println(msg)
	}
}
