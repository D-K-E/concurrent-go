package selectscenario

import (
	"fmt"
	"sort"
	"strings"

	selfsync "github.com/D-K-E/concurrent-go/selfsync"
)

/*
A demonstration of Broadcast/Worker pattern. We create multiple workers to
work on a single output. Each worker does a different task
*/

func frequentWords(quit <-chan int, words <-chan string) <-chan string {
	mostFrequentWords := make(chan string)
	go func() {
		defer close(mostFrequentWords)
		freqMap := make(map[string]int)
		freqList := make([]string, 0)
		isOpen, word := true, ""
		for isOpen {
			select {
			case word, isOpen = (<-words):
				if isOpen {
					if freqMap[word] == 0 {
						freqList = append(freqList, word)
					}
					freqMap[word] += 1
				}
			case <-quit:
				return
			}
		}
		sort.Slice(freqList, func(a, b int) bool {
			return freqMap[freqList[a]] > freqMap[freqList[b]]
		})
		mostFrequentWords <- strings.Join(freqList[:10], ", ")
	}()
	return mostFrequentWords
}

func BroadcastMain() {
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
	multiWords := selfsync.Broadcast(quit, words, 2) // create 2 workers
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
