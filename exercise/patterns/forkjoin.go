package patterns

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
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

From selectscenarios that were covered, the fork-join pattern resembles to
broadcast pattern.

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

type FileStat interface {
	FileName() string
	Value() int
}

type CodeDepth struct {
	file string
	num  int
}

func (c CodeDepth) FileName() string {
	return c.file
}

func (c CodeDepth) Value() int { return c.num }

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
	cdepth := CodeDepth{filename, maxDepth}
	return cdepth
}

type FuncNumber struct {
	filename string
	numFunc  int
}

func (f FuncNumber) FileName() string { return f.filename }
func (f FuncNumber) Value() int       { return f.numFunc }

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
	fnumber := FuncNumber{filename, nbFunc}
	return fnumber
}

// now let's see how forking is done
func forkReadLinesIfNeed(path string, info os.FileInfo,
	wg *sync.WaitGroup,
	fileContents chan<- FileContent,
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

func forkFileHandler[FuncOut FileStat](wg *sync.WaitGroup,
	fileContents <-chan FileContent,
	codeDepths chan FuncOut,
	fileHandler func(FileContent) FuncOut,
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

func forkDeepestNested(wg *sync.WaitGroup,
	fileContents <-chan FileContent, codeDepths chan CodeDepth,
) {
	forkFileHandler(wg, fileContents, codeDepths, deepestNestedBlock)
}

func forkNumberOfFunc(wg *sync.WaitGroup, fileContents <-chan FileContent, numF chan FuncNumber) {
	forkFileHandler(wg, fileContents, numF, numberOfFunc)
}

// Now we do the join part
func joinFileHandler[FuncOut FileStat](
	wg *sync.WaitGroup,
	partialResult <-chan FuncOut,
	maxDefault FuncOut,
) chan FuncOut {
	finalResult := make(chan FuncOut)
	wg.Add(1)
	go func() {
		localDefault := maxDefault
		for result := range partialResult {
			if result.Value() > localDefault.Value() {
				localDefault = result
			}
		}
		finalResult <- localDefault
		wg.Done()
	}()
	return finalResult
}

func fanInResults(finalCodeDepth chan CodeDepth, finalFuncNumber chan FuncNumber, wg *sync.WaitGroup,
) {
	wg.Add(2)
	go func() {
		isFinalCodeOpen := true
		isFinalFuncOpen := true
		var finalCDepth CodeDepth
		var finalFNumber FuncNumber
		for isFinalCodeOpen || isFinalFuncOpen {
			select {
			case finalCDepth, isFinalCodeOpen = (<-finalCodeDepth):
				if isFinalCodeOpen {
					fmt.Printf("%s has deepest nested code block of %d\n",
						finalCDepth.FileName(), finalCDepth.Value())

					// now that we have the result we can close the channel
					wg.Done()
				}
			case finalFNumber, isFinalFuncOpen = (<-finalFuncNumber):
				if isFinalFuncOpen {
					fmt.Printf("%s has the highest number of func %d\n",
						finalFNumber.FileName(), finalFNumber.Value())

					// now that we have the result we can close the channel
					wg.Done()
				}
			}
		}
	}()
}

// now let's synchronize everything
func ForkJoinMain() {
	dir := os.Args[1]

	// make partial result channels
	fileContents := make(chan FileContent)
	codeDepths := make(chan CodeDepth)
	funcNumbers := make(chan FuncNumber)

	// create the wait group for waiting in the synchronization functions
	wg1 := sync.WaitGroup{}
	wg2 := sync.WaitGroup{}
	wg3 := sync.WaitGroup{}

	// file path walk
	filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			// launch tasks
			forkReadLinesIfNeed(path, info, &wg1, fileContents)
			forkDeepestNested(&wg2, fileContents, codeDepths)
			forkNumberOfFunc(&wg2, fileContents, funcNumbers)
			return nil
		})

	//
	wg1.Wait()
	close(fileContents)

	// join results
	cdepth := CodeDepth{"", 0}
	finalCodeDepth := joinFileHandler(&wg3,
		codeDepths, cdepth)
	funcN := FuncNumber{"", 0}
	finalFuncNumber := joinFileHandler(&wg3,
		funcNumbers, funcN)

	// wait for reading task to complete
	wg2.Wait()

	close(codeDepths)
	close(funcNumbers)

	//
	wg4 := sync.WaitGroup{}
	fanInResults(finalCodeDepth, finalFuncNumber, &wg4)

	wg3.Wait()

	wg4.Wait()
	// close its channel
	close(finalCodeDepth)
	close(finalFuncNumber)
}
