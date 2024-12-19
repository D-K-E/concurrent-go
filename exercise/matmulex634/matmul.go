package matmul

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	// "time"

	selfsync "github.com/D-K-E/concurrent-go/selfsync"
)

// concurrent matrix
func rowMultiplyWg(matA map[int][]int, matB map[int][]int, row int,
	result map[int][]int,
	waitGroup *sync.WaitGroup,
	mutex *selfsync.ReadWriteMutex,
) error {
	defer waitGroup.Done()
	matrixSize := len(matA)
	if matrixSize != len(matB) {
		return errors.New("Not a square matrix")
	}
	mutex.WriteLock()
	nslice := make([]int, matrixSize)
	result[row] = nslice
	for col := 0; col < matrixSize; col++ {
		sum := 0
		for i := 0; i < matrixSize; i++ {
			sum += matA[row][i] * matB[i][col]
		}
		result[row][col] = sum
	}
	mutex.WriteUnlock()
	return nil
}

// sequential matrix
func rowMultiply(matA map[int][]int, matB map[int][]int, row int,
	result map[int][]int,
) error {
	matrixSize := len(matA)
	if matrixSize != len(matB) {
		return errors.New("Not a square matrix")
	}
	nslice := make([]int, matrixSize)
	result[row] = nslice
	for col := 0; col < matrixSize; col++ {
		sum := 0
		for i := 0; i < matrixSize; i++ {
			sum += matA[row][i] * matB[i][col]
		}
		result[row][col] = sum
	}
	return nil
}

func generateRandMat(size int) map[int][]int {
	m := make(map[int][]int)
	for i := 0; i < size; i++ {
		nslice := make([]int, size)
		for j := 0; j < size; j++ {
			nslice[j] = rand.Intn(10)
		}
		m[i] = nslice
	}
	return m
}

func MatMulMain() {
	size := 5
	matA := generateRandMat(size)
	matB := generateRandMat(size)
	result := make(map[int][]int, size)
	nresult := make(map[int][]int, size)
	wg := sync.WaitGroup{}
	nmutex := selfsync.NewMutex()
	for row := 0; row < size; row++ {
		wg.Add(1)
		go rowMultiplyWg(matA, matB, row, result, &wg, nmutex)
	}
	wg.Wait()

	// print everything to console
	fmt.Println("--------------")
	for row := 0; row < size; row++ {
		fmt.Println(matA[row], matB[row], result[row])
	}
	for row := 0; row < size; row++ {
		rowMultiply(matA, matB, row, nresult)
	}
	fmt.Println("--------------")
	for row := 0; row < size; row++ {
		fmt.Println(matA[row], matB[row], nresult[row])
	}
}
