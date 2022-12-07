package vector

// Compute 1D index from 2D indices
func Compute1DIndex(i, j, n int) int {
	return i*n + j
}

// Perform the operation res = x * y
func DotP(v1, v2 *[]float64, N int) float64 {
	res := 0.0
	for i := 0; i < N; i++ {
		res += (*v1)[i] * (*v2)[i]
	}
	return res
}

// Perform the operation - output = alpha * v1 + beta * v2
func Axpy(output *[]float64, alpha float64, v1 *[]float64, beta float64, v2 *[]float64, N int) {
	for i := 0; i < N; i++ {
		(*output)[i] = alpha*(*v1)[i] + beta*(*v2)[i]
	}
}
