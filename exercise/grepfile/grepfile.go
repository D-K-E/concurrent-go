package grepfile

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func read(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	value := string(data[:])
	return value, nil
}

func searchFile(searchTxt string, fileContent string) bool {
	res := strings.Contains(fileContent, searchTxt)
	return res
}

func toConsole(searchTxt string, filename string) {
	str, err := read(filename)
	if err == nil {
		res := searchFile(searchTxt, str)
		fmt.Printf("file %s contains %s: %t\n", filename, searchTxt, res)
	} else {
		fmt.Printf("couldn't found %s", filename)
	}
}

func GrepFileMain() {
	searchTxtVar := flag.String("search", "my", "search text")
	filename1 := flag.String("filename1", "", "a file name")
	filename2 := flag.String("filename2", "", "a file name")
	flag.Parse()
	var files [3]string
	files[0] = *searchTxtVar
	files[1] = *filename1
	files[2] = *filename2
	for i := 1; i < 3; i++ {
		searchTxt := files[0]
		filename := files[i]
		go toConsole(searchTxt, filename)
		time.Sleep(1 * time.Second)
	}
}
