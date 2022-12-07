package main

import (
	"fmt"
	"os"
	"proj3/conjugate_gradient/rhs"
	"proj3/conjugate_gradient/solver"
	"proj3/conjugate_gradient/utils"
	"proj3/conjugate_gradient/vector"
	"strconv"
	"sync"
	"time"
)

const usage = "Usage: go run conjugate_reduce.go <imageHeight> <nThreads\n  <image_height> = number of grid points in each dimension\n<imageHeight> = number of threads to use for parallelization\n"

func main() {
	// Flag set manually by user
	benchmarking := false

	if len(os.Args) < 2 {
		fmt.Println(usage)
		return
	} else if len(os.Args) > 3 {
		fmt.Println(usage)
		return
	}

	// Initialize variables
	n, _ := strconv.Atoi(os.Args[1])
	nThreads := 0
	if len(os.Args) == 3 {
		nThreads, _ = strconv.Atoi(os.Args[2])
		if nThreads < 2 {
			fmt.Println("Number of threads should be greater than or equal to 2 since there must be atleast 1 mapper and 1 reducer.\n")
			return
		}
	}

	context := vector.MapReduceContext{
		NThreads:     nThreads,
		MapWG:        &sync.WaitGroup{},
		ReduceWG:     &sync.WaitGroup{},
		Benchmarking: benchmarking,
	}

	N := n * n

	if !benchmarking {
		fmt.Println("Simulation Parameters:")
		fmt.Println("n =", n)
		if nThreads > 0 {
			fmt.Println("Number of threads =", nThreads)
		}
	}
	// Start timer
	s := time.Now()

	b := make([]float64, N)
	rhs.FillB(&b, N)

	// Result vector
	x := make([]float64, N)

	// Run simulation
	solver.ConjugateGradient(&b, &x, n, &context)

	// Stop timer
	e := time.Now()

	if !benchmarking {
		// Write RHS to file
		utils.WriteToFile(&b, "./output/b.txt", 0, N)

		// Write solution to file
		utils.WriteToFile(&x, "./output/x.txt", N, N)
	}

	// Print time taken
	if !benchmarking {
		fmt.Println("Time taken =", e.Sub(s).Seconds(), "seconds")
	} else {
		fmt.Println(e.Sub(s).Seconds())
	}
}
