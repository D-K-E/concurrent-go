package patterns

import (
	"bufio"
	"math"
	"os"
	"strings"
	"sync"
)

/*
snippet demonstrating fork-join pattern.

In the loop_level_parallelism we have stated that there are mainly two ways to
decompose a program:
- By data
- By task

The fork-join pattern fits mostly to task decomposition. We have some task or
a suite of tasks and we want to execute simultaneously. In the fork part of
the pattern we spawn a routine for each task, and gather their results in a
common channel in the join part.

Our example program would find the code with the most nested blocks by
searching through all the files.
*/

type FileContent struct {
	file  string
	lines []string
}

// task 1: reading a file content
func readLines(filename string) FileContent {
	f, _ := os.Open(filename)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lines := make([]string, 1)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	fc := FileContent{filename, lines}
	return fc
}

type CodeDepth struct {
	file  string
	level int
}

// task 2: compute code depth for given code text
func deepestNestedBlock(file FileContent) CodeDepth {
	lines := file.lines
	filename := file.file
	maxDepth := 0
	level := 0
	for _, line := range lines {
		for _, c := range line {
			if c == '{' {
				level += 1
				mComp := math.Max(float64(maxDepth), float64(level))
				maxDepth = int(mComp)
			} else if c == '}' {
				level -= 1
			}
		}
	}
	return CodeDepth{filename, maxDepth}
}

type FuncNumber struct {
	filename string
	numFunc  int
}

// task 3: compute number of funcs for given code text
func numberOfFunc(file FileContent) FuncNumber {
	nbFunc := 0
	lines := file.lines
	filename := file.file

	for _, line := range lines {
		if strings.Contains(line, "func ") {
			nbFunc += 1
		}
	}
	return FuncNumber{filename, nbFunc}
}

// now let's see how forking is done
func forkReadLinesIfNeed(path string, info os.FileInfo,
	wg *sync.WaitGroup, fileContents chan FileContent,
) {
	if (!info.IsDir()) && (strings.HasSuffix(path, ".go")) {
		wg.Add(1)
		go func() {
			fc := readLines(path)
			fileContents <- fc
			wg.Done()
		}()
	}
}

func forkFileHandler[FuncOut any](
	wg *sync.WaitGroup, fileContents <-chan FileContent,
	codeDepths chan FuncOut, fileHandler func(FileContent) FuncOut,
) {
	wg.Add(1)
	go func() {
		for fc := range fileContents {
			codeD := fileHandler(fc)
			codeDepths <- codeD
		}
		wg.Done()
	}()
}

func forkDeepestNested(wg *sync.WaitGroup, fileContents <-chan FileContent,
	codeDepths chan CodeDepth,
) {
	forkFileHandler(wg, fileContents, codeDepths, deepestNestedBlock)
}

func forkNumberOfFunc(wg *sync.WaitGroup, fileContents <-chan FileContent,
	numF chan FuncNumber,
) {
	forkFileHandler(wg, fileContents, numF, numberOfFunc)
}

// Now we do the join part
