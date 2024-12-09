package letterfreq

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

const allLetters = "qwertyuiop@asdfghjkl;:]zxcvbnm,./\\1234567890-^ふあうわん"

func countLetters(url string, frequency []int) error {
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

	for _, b := range body {
		c := strings.ToLower(string(b))
		cIndex := strings.Index(allLetters, c)
		if cIndex >= 0 {
			frequency[cIndex] += 1
		}
	}
	fmt.Println("Completed:", url)
	return nil
}

func LetterFreqMain() {
	freq := make([]int, len(allLetters))

	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		countLetters(url, freq)
	}

	for i, c := range allLetters {
		fmt.Printf("%c-%d ", c, freq[i])
	}
	fmt.Println("Done")
}
