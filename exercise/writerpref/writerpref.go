package writerpref

//
import (
	"fmt"
	"time"

	selfsync "github.com/D-K-E/concurrent-go/selfsync"
)

func WriterPrefMain() {
	rwMut := selfsync.NewMutex()
	for i := 0; i < 2; i++ {
		go func() { // starts 2 goroutines
			for { // repeats forever
				rwMut.ReadLock() //
				time.Sleep(1 * time.Second)
				fmt.Println("read done")
				rwMut.ReadUnlock()
			}
		}()
	}
    time.Sleep(1 * time.Second)
    rwMut.WriteLock() // tries to acquire the writer lock from main routine
    fmt.Println("write finished")
}
