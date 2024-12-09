package letterfreq

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

func countLettersMutex(url string, frequency []int, mutex *sync.RWMutex) error {
	//
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Server returning error status code: " + resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	mutex.Lock()
	for _, b := range body {
		c := strings.ToLower(string(b))
		cIndex := strings.Index(allLetters, c)
		if cIndex >= 0 {
			frequency[cIndex] += 1
		}
	}
	mutex.Unlock()
	fmt.Println("Completed:", url)
	return nil
}

func LetterFreqMutexMain() {
	freq := make([]int, len(allLetters))

	mutex := sync.RWMutex{}

	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go countLettersMutex(url, freq, &mutex)
	}
	time.Sleep(5 * time.Second)

	for i, c := range allLetters {
		fmt.Printf("%c-%d ", c, freq[i])
	}
	fmt.Println("Done")
}
