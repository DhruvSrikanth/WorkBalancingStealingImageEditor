package solver

import (
	"fmt"
	"math"
	"proj3/conjugate_gradient/poisson"
	"proj3/conjugate_gradient/utils"
	"proj3/conjugate_gradient/vector"
)

func ConjugateGradient(b, x *[]float64, n int, context *vector.MapReduceContext) {
	// Initialize variables
	N := n * n

	// Tolerance for the solution to stop after convergence
	tol := 1.0e-10

	r := make([]float64, N)
	p := make([]float64, N)
	z := make([]float64, N)

	// Temporary variables
	Ax := make([]float64, N)

	// r = b - Ax
	poisson.PoissonOnTheFly(&Ax, x, N)
	vector.Axpy(&r, 1.0, b, -1.0, &Ax, N)

	// p = r
	for i := 0; i < N; i++ {
		p[i] = r[i]
	}

	// rsold = rT * r
	rsold := vector.DotPWrapper(&r, &r, N, context)

	for i := 0; i < N; i++ {
		// z = A*p
		poisson.PoissonOnTheFly(&z, &p, N)

		// alpha = rsold / (p*z)
		alpha := rsold / vector.DotPWrapper(&p, &z, N, context)

		// x = x + alpha*p
		vector.Axpy(x, 1.0, x, alpha, &p, N)

		// r = r - alpha*z
		vector.Axpy(&r, 1.0, &r, -alpha, &z, N)

		// rsnew = rT*r
		rsnew := vector.DotPWrapper(&r, &r, N, context)

		// If the residual is small enough, stop
		if math.Sqrt(rsnew) <= tol {
			fmt.Printf("Converged after %d iterations\n", i)
			break
		}

		// p = r + rsnew / rsold * p
		vector.Axpy(&p, 1.0, &r, rsnew/rsold, &p, N)

		// rsold = rsnew
		rsold = rsnew

		// Save x every 10 iterations
		if i%10 == 0 {
			outPath := fmt.Sprintf("./output/x_%d.txt", i)
			utils.WriteToFile(x, outPath, i, N)
		}
	}
}
