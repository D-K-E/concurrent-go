package main

import (
	"flag"
	"fmt"
	"time"
)

func doWork(id int) {
	fmt.Printf("work %d started at %s\n", id, time.Now().Format("15:04:05"))
	time.Sleep(1 * time.Second)
	fmt.Printf("work %d finished at %s\n", id, time.Now().Format("15:04:05"))
}

func seqMain() {
	for i := 0; i < 5; i++ {
		doWork(i)
	}
}

func parallelMain() {
	for i := 0; i < 5; i++ {
		go doWork(i)
	}
	time.Sleep(2 * time.Second)
}

func main() {
	var isParallel *bool = flag.Bool("is_parallel", false, "a bool var")
	flag.Parse()
	if *isParallel {
		parallelMain()
	} else {
		seqMain()
	}
}
