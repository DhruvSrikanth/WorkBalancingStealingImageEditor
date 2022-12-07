package rhs

import (
	"math"
	"proj3/conjugate_gradient/vector"
)

// Find the value of b at a given point
func findB(i, j, n int) float64 {
	delta := 1.0 / float64(n)

	x := -0.5 + delta + delta*float64(j)
	y := -0.5 + delta + delta*float64(i)

	// Check if within a circle
	radius := 0.1
	if x*x+y*y < radius*radius {
		return delta * delta / 1.075271758e-02
	} else {
		return 0.0
	}
}

// Fill the b vector with the values of b
func FillB(b *[]float64, N int) {
	n := int(math.Sqrt(float64(N)))
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			(*b)[vector.Compute1DIndex(i, j, n)] = findB(i, j, n)
		}
	}
}
