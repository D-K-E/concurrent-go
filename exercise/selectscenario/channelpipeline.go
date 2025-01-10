package selectscenario

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

/*
this scenario concerns how to use pipeline pattern with channels.
The idea is that upon receiving an quit channel as input we return an output
channel, and the function that consumes the output channel stops consumption
using quit channel. Sounds complicated, but quite easy in code
*/

// generates urls we can think of this as the pipeline channel. Currently it
// has a quit channel and an output channel, but we'll see in subsequent
// versions of this function other channels that are spawned from it
func generateUrls(quit <-chan int) <-chan string {
	urls := make(chan string)
	go func() {
		defer close(urls)
		for i := 100; i <= 130; i++ {
			url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
			select {
			case urls <- url:
			case (<-quit):
				return
			}
		}
	}()
	return urls
}

func ChannelPipelineMainV1() {
	qu := make(chan int)
	defer close(qu)
	resultChannel := generateUrls(qu)

	// consume result channel
	for msg := range resultChannel {
		// consume
		fmt.Println(msg)
	}
}

func downloadPages(quit <-chan int, urls <-chan string) <-chan string {
	pages := make(chan string)
	go func() {
		defer close(pages)
		isOpen := true
		url := ""
		for isOpen {
			select {
			case url, isOpen = (<-urls):
				if isOpen {
					resp, _ := http.Get(url)
					if resp.StatusCode != 200 {
						panic("Server's error: " + resp.Status)
					}
					body, _ := io.ReadAll(resp.Body)
					pages <- string(body)
					resp.Body.Close()
				}
			case <-quit:
				return
			}
		}
	}()
	return pages
}

func ChannelPipelineMainV2() {
	quit := make(chan int)
	defer close(quit)
	urls := generateUrls(quit)
	pages := downloadPages(quit, urls)

	// consume result channel
	for msg := range pages {
		// consume
		fmt.Println(msg)
	}
}

func extractWords(quit <-chan int, pages <-chan string) <-chan string {
	words := make(chan string)
	go func() {
		//
		defer close(words)
		wordRegex := regexp.MustCompile(`[a-zA-Z]+`)
		isOpen := true
		pg := ""
		for isOpen {
			select {
			case pg, isOpen = (<-pages):
				if isOpen {
					for _, word := range wordRegex.FindAllString(pg, -1) {
						words <- strings.ToLower(word)
					}
				}
			case <-quit:
				return
			}
		}
	}()
	return words
}

func ChannelPipelineMainV3() {
	quit := make(chan int)
	defer close(quit)
	urls := generateUrls(quit)
	pages := downloadPages(quit, urls)
	words := extractWords(quit, pages)

	// consume result channel
	for msg := range words {
		// consume
		fmt.Println(msg)
	}
}
