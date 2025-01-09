package selectscenario

import (
	"fmt"
)

/*
this scenario concerns how to use pipeline pattern with channels.
The idea is that upon receiving an quit channel as input we return an output
channel, and the function that consumes the output channel stops consumption
using quit channel. Sounds complicated, but quite easy in code
*/

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
