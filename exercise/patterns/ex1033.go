package patterns

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func WebDownloadSequentialMain() {
	const pagesToDownload = 12
	totalLines := 0
	for i := 1000; i < 1000+pagesToDownload; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		fmt.Println("Downloading", url)
		resp, _ := http.Get(url)
		if resp.StatusCode != 200 {
			panic("Server error: " + resp.Status)
		}
		bodyBytes, _ := io.ReadAll(resp.Body)
		totalLines += strings.Count(string(bodyBytes), "\n")
		resp.Body.Close()
	}
	fmt.Println("Total lines", totalLines)
}

// task 1: prepare url
func prepUrl(i int) string {
	url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
	return url
}

// task 2: download url
func downloadPage(url string) *http.Response {
	fmt.Println("downloading", url)
	resp, _ := http.Get(url)
	if resp.StatusCode != 200 {
		return nil
	}
	return resp
}

// task 3: read response
func readResponse(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	return string(bodyBytes)
}

// task 4: count lines
func countLines(body string) int {
	lines := strings.Count(body, "\n")
	return lines
}

func WebDownloadConcurrentMain() {
	const pagesToDownload = 12
	input := make(chan int)
	quit := make(chan int)
	urlOut := AddOnPipe(quit, prepUrl, input)
	downloadOut := AddOnPipe(quit, downloadPage, urlOut)
	respOut := AddOnPipe(quit, readResponse, downloadOut)
	countOut := AddOnPipe(quit, countLines, respOut)
	go func() {
		for i := 1000; i < 1000+pagesToDownload; i++ {
			input <- i
		}
	}()
	totalLines := 0
	for j := 0; j < pagesToDownload; j++ {
		totalLines += (<-countOut)
	}
	quit <- 1
	fmt.Println("Total lines", totalLines)
}
