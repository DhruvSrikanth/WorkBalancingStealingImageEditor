package utils

import (
	"fmt"
	"math"
	"os"
	"proj3/conjugate_gradient/vector"
)

// Write the result to a file
func WriteToFile(result *[]float64, filename string, n_iter, N int) {
	// Create the file
	file, _ := os.Create(filename)
	defer file.Close()

	// Compute dimensions of the 1D array to represent as 2D array
	n := int(math.Sqrt(float64(N)))

	// Write the timestep to the file
	file.WriteString(fmt.Sprintf("[%d], ", n_iter))

	// Write the data to the file
	file.WriteString("[")
	for i := 0; i < n; i++ {
		if i == 0 {
			file.WriteString("[")
		} else {
			file.WriteString(", [")
		}
		for j := 0; j < n; j++ {
			res := (*result)[vector.Compute1DIndex(i, j, n)]
			if j == N-1 {
				file.WriteString(fmt.Sprintf("%f", res))
			} else {
				file.WriteString(fmt.Sprintf("%f, ", res))
			}
		}
		file.WriteString("]")
	}
	file.WriteString("]")

	// Release the memory for the file
	file.Close()
}
