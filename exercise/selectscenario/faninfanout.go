package selectscenario

import (
	"fmt"
	"sort"
	"strings"

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

func longestWords(quit <-chan int, words <-chan string) <-chan string {
	longWords := make(chan string)
	go func() {
		defer close(longWords)
		uniqueWordsMaps := make(map[string]bool)
		uniqueWords := make([]string, 0)
		isOpen, word := true, ""
		for isOpen {
			select {
			case word, isOpen = (<-words):
				if isOpen && !uniqueWordsMaps[word] {
					uniqueWordsMaps[word] = true
					uniqueWords = append(uniqueWords, word)
				}
			case <-quit:
				return
			}
		}
		sort.Slice(uniqueWords, func(a, b int) bool {
			return len(uniqueWords[a]) > len(uniqueWords[b])
		})
		longWords <- strings.Join(uniqueWords[:10], ", ")
	}()
	return longWords
}

func FanInFanOutMainV2() {
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
	results := longestWords(quit, words)

	// consume result channel
	for msg := range results {
		// consume
		fmt.Println(msg)
	}
}
