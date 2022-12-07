package poisson

import (
	"math"
	"proj3/conjugate_gradient/vector"
)

// Compute the poisson equation on the fly (possible because of it being a sparse matrix)
func PoissonOnTheFly(v, w *[]float64, N int) {
	n := int(math.Sqrt(float64(N)))
	var left, right, up, down, here float64
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			// Compute the 1D index of the current point
			idx1D := vector.Compute1DIndex(i, j, n)

			// Initialize the values of the neighbors
			left = 0.0
			right = 0.0
			up = 0.0
			down = 0.0

			here = (*w)[idx1D]

			// Compute the values of the neighbors
			if i > 0 {
				up = (*w)[vector.Compute1DIndex(i-1, j, n)]
			}
			if i < n-1 {
				down = (*w)[vector.Compute1DIndex(i+1, j, n)]
			}
			if j > 0 {
				left = (*w)[vector.Compute1DIndex(i, j-1, n)]
			}
			if j < n-1 {
				right = (*w)[vector.Compute1DIndex(i, j+1, n)]
			}

			// Compute the value of v[i]
			(*v)[idx1D] = 4.0*here - left - right - up - down
		}
	}
}
