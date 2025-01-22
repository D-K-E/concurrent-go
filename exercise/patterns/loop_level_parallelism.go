package patterns

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

/*
snippet demonstrating how to do loop level parallelism pattern.

There are two common ways to decompose a program:
- By data
- By task

Loop level parallelism concerns mostly data decomposition. You have some kind
of container and you want to apply some sort of function to it and create a
new container with the transformed data. `xs = Map(f, ys)` kind of scenario
where you map function `f` to iterable `ys`. Optionally you can have some kind
of fold scenario where the next iteration depends on the previous one:
`xs = Fold(Acc, f, ys)` where `Acc` is an accumulation variable (it can be an
integer or another list entirely), `f` is a function and `ys` is the
container.
*/

// first we'll look at a map scenario, where `f` will compute the hash of a
// file, given its path:

func FHash(filepath string) []byte {
	file, _ := os.Open(filepath)
	defer file.Close()
	sha := sha256.New()
	io.Copy(sha, file)
	shaSum := sha.Sum(nil)
	return shaSum
}

func MapLikeLoopParallelismWithIndependentInnerLoopMain() {
	dir := os.Args[1]
	files, _ := os.ReadDir(dir)
	wg := sync.WaitGroup{}
	for _, file := range files {
		if !file.IsDir() {
			wg.Add(1)
			go func(filename string) {
				fp := filepath.Join(dir, filename)
				hash := FHash(fp)
				fmt.Printf("%s - %x\n", filename, hash)
				wg.Done()
			}(file.Name())
		}
	}
	wg.Wait()
}

/*
The second scenario is fold like scenario. This time we'll compute the hash of
a folder using the hash of files inside it. The tricky bit is the fact that
order of iteration would affect the folder hash. The solution to this is to
run independent parts concurrently and compute dependent parts using
synchronisation mechanisms.
*/

func FoldLikeLoopParallelismWithDependentInnerLoopMain() {
	dir := os.Args[1]
	files, _ := os.ReadDir(dir)
	sha := sha256.New()
	var prev, next chan int
	for _, file := range files {
		if !file.IsDir() {
			next = make(chan int)
			go func(filename string, prev, next chan int) {
				fp := filepath.Join(dir, filename)
				hash := FHash(fp)
				if prev != nil {
					<-prev // if not first go routine wait for previous
				}
				sha.Write(hash)
				next <- 0
			}(file.Name(), prev, next)
			prev = next // notice we are piping channels successively
		}
	}
	// drain next
	<-next
	/*for n := range next {
		_ = n
	}*/
	fmt.Printf("%x\n", sha.Sum(nil))
}
