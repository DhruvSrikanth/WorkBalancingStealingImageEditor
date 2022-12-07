package main

import (
	"fmt"
	"os"
	"proj3/conjugate_gradient/rhs"
	"proj3/conjugate_gradient/solver"
	"proj3/conjugate_gradient/utils"
	"strconv"
	"time"
)

const usage = "Usage: go run conjugate_reduce.go <N>\n<N> = number of grid points in each dimension\n"

func main() {

	if len(os.Args) != 2 {
		fmt.Println(usage)
		return
	}
	// Initialize variables
	n, _ := strconv.Atoi(os.Args[1])
	N := n * n

	fmt.Println("Simulation Parameters:")
	fmt.Println("n = ", n)

	// Start timer
	s := time.Now()

	b := make([]float64, N)
	rhs.FillB(&b, N)

	// Result vector
	x := make([]float64, N)

	// Run simulation
	solver.ConjugateGradient(&b, &x, n)

	// Stop timer
	e := time.Now()

	// Write RHS to file
	utils.WriteToFile(&b, "./output/b.txt", 0, N)

	// Write solution to file
	utils.WriteToFile(&x, "./output/x.txt", N, N)

	// Print time taken
	fmt.Println("Time taken = ", e.Sub(s).Seconds(), " seconds")
}
