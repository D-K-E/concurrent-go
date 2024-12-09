package catfile

import (
	"flag"
	"fmt"
	"os"
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

func toConsole(filename string) {
	str, err := read(filename)
	if err == nil {
		fmt.Printf("file content:\n %s", str)
	} else {
		fmt.Printf("couldn't found %s", filename)
	}
}

func CatFileMain() {
	var filename1 *string = flag.String("filename1", "", "a file name")
	var filename2 *string = flag.String("filename2", "", "a file name")
	flag.Parse()
	var files [2]string
	files[0] = *filename1
	files[1] = *filename2
	for i := 0; i < 2; i++ {
		filename := files[i]
		go toConsole(filename)
		time.Sleep(1 * time.Second)
	}
}
