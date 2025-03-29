package patterns

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// fold like loop parallelism with dependent inner loop using wait group
func FoldLikeLoopParallelismWithDependentInnerLoopWaitGroupMain() {
	dir := os.Args[1]
	files, _ := os.ReadDir(dir)
	sha := sha256.New()
	var prev, next *sync.WaitGroup // from solutions
	for _, file := range files {
		if !file.IsDir() {
			next = &sync.WaitGroup{}
			next.Add(1)
			go func(filename string, prev, next *sync.WaitGroup) {
				fp := filepath.Join(dir, filename)
				hash := FHash(fp)
				if prev != nil {
					prev.Wait()
				}
				sha.Write(hash)
				next.Done()
			}(file.Name(), prev, next)
			prev = next // notice we are piping channels successively
		}
	}
	// drain next
	next.Wait()

	fmt.Printf("%x\n", sha.Sum(nil))
}
