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

/*
Given these task let's first see the synchronized version of this code
*/
func ForkJoinSynchronizedMain() {
	// get root directory
	dir := os.Args[1]

	//
	cdepth_ := CodeDepth{"", 0}
	funcN := FuncNumber{"", 0}
	// file path walk
	filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			// launch tasks
			fcontent := readLines(path)
			cdepth := deepestNestedBlock(fcontent)
			numFunc := numberOfFunc(fcontent)
			if cdepth.Value() > cdepth_.Value() {
				cdepth_ = cdepth
			}
			if numFunc.Value() > funcN.Value() {
				funcN = numFunc
			}
			return nil
		})

	// print results
	fmt.Printf("%s has deepest nested code block of %d\n",
		cdepth_.FileName(), cdepth_.Value())

	fmt.Printf("%s has highest number of func of %d\n",
		funcN.FileName(), funcN.Value())
}

// now let's see how forking is done
func forkReadLinesIfNeed(path string, info os.FileInfo,
	wg *sync.WaitGroup,
	codeDepths chan<- CodeDepth,
	funcNums chan<- FuncNumber,
) {
	if (!info.IsDir()) && (strings.HasSuffix(path, ".go")) {
		wg.Add(2)
		fc := readLines(path)

		// launch task 1
		go func() {
			cdepth := deepestNestedBlock(fc)
			codeDepths <- cdepth
			wg.Done()
		}()
		// launch task 2
		go func() {
			numf := numberOfFunc(fc)
			funcNums <- numf
			wg.Done()
		}()
	}
}

type Pair[FirstType, SecondType any] struct {
	first  FirstType
	second SecondType
}

// Now we do the join part
func joinFileHandler(
	codeDepths <-chan CodeDepth,
	funcNums <-chan FuncNumber,
) Pair[chan CodeDepth, chan FuncNumber] {
	finalCodeDepth := make(chan CodeDepth)
	mxDepth := CodeDepth{"", 0}
	go func() {
		for result := range codeDepths {
			if result.Value() > mxDepth.Value() {
				mxDepth = result
			}
		}
		finalCodeDepth <- mxDepth
	}()
	finalFuncNum := make(chan FuncNumber)
	mxFunc := FuncNumber{"", 0}
	go func() {
		for result := range funcNums {
			if result.Value() > mxFunc.Value() {
				mxFunc = result
			}
		}
		finalFuncNum <- mxFunc
	}()
	p := Pair[chan CodeDepth, chan FuncNumber]{finalCodeDepth, finalFuncNum}
	return p
}

// now let's synchronize everything

// fork join main
func ForkJoinMain() {
	dir := os.Args[1]

	// make partial result channels
	codeDepths := make(chan CodeDepth)
	funcNumbers := make(chan FuncNumber)

	// create the wait group for waiting in the synchronization functions
	wg1 := sync.WaitGroup{}

	// file path walk
	filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			// launch tasks with forks
			forkReadLinesIfNeed(path, info, &wg1, codeDepths, funcNumbers)
			return nil
		})

	//
	// join results
	pair := joinFileHandler(
		codeDepths, funcNumbers)

	wg1.Wait()
	close(codeDepths)
	close(funcNumbers)

	// now let's print the values
	isCodeDepthOpen, isFuncOpen := true, true
	var finalCodeDepth CodeDepth
	var finalFuncNum FuncNumber
	for isCodeDepthOpen || isFuncOpen {
		select {
		case finalCodeDepth, isCodeDepthOpen = (<-pair.first):
			if isCodeDepthOpen {
				fmt.Printf("%s has deepest nested code block of %d\n",
					finalCodeDepth.FileName(), finalCodeDepth.Value())
				close(pair.first)
			}
		case finalFuncNum, isFuncOpen = (<-pair.second):
			if isFuncOpen {
				fmt.Printf("%s has the highest number of func %d\n",
					finalFuncNum.FileName(), finalFuncNum.Value())
				close(pair.second)
			}
		}
	}
}
