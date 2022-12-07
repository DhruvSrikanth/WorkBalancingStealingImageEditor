package vector

import "sync"

type MapReduceContext struct {
	NThreads int // Runs the parallel version of the program with the specified number of threads
	WG       *sync.WaitGroup
}

func mapper(mapperIn <-chan [][]float64, mapperOut chan<- float64, context *MapReduceContext) {
	for vectorPair := range mapperIn {
		mapperOut <- DotP(&vectorPair[0], &vectorPair[1], len(vectorPair[0]))
	}
	close(mapperOut)
	context.WG.Done()
}

func reducer(reducerIn <-chan float64, reducerOut chan<- float64, context *MapReduceContext) {
	globalDotP := 0.0
	for localDotP := range reducerIn {
		globalDotP += localDotP
	}
	reducerOut <- globalDotP
	close(reducerOut)
	context.WG.Done()
}

// MapReduce computes the dot product of two vectors in parallel
func MapReduce(v1, v2 *[]float64, N int, context *MapReduceContext) float64 {
	vector1 := *v1
	vector2 := *v2

	// Create the input channel
	mapperIn := make(chan [][]float64, context.NThreads)

	// Create the output channel
	mapperOut := make(chan float64, context.NThreads)

	// Create the reducer output channel
	reducerOut := make(chan float64, 1)

	// Start the mappers
	for i := 0; i < context.NThreads-1; i++ {
		context.WG.Add(1)
		go mapper(mapperIn, mapperOut, context)
	}

	// Start the reducer
	context.WG.Add(1)
	go reducer(mapperOut, reducerOut, context)

	// Compute the size of the chunks
	chunkSize := N / (context.NThreads - 1)

	// Compute the dot product in parallel
	for i := 0; i < context.NThreads-1; i++ {
		// Compute the start and end indices
		start := i * chunkSize
		end := (i + 1) * chunkSize

		// If this is the last chunk, make sure it includes the remainder
		if i == context.NThreads-2 {
			end = N
		}

		// Create the vector pair
		vectorPair := make([][]float64, 2)
		vectorPair[0] = vector1[start:end]
		vectorPair[1] = vector2[start:end]

		// Send the vector pair to the mapper
		mapperIn <- vectorPair
	}

	// Close the input channel
	close(mapperIn)

	// Return the dot product
	return <-reducerOut
}

// ParallelDotP computes the dot product of two vectors safely in parallel
func ParallelDotP(v1, v2 *[]float64, N int, context *MapReduceContext) float64 {
	// Run the map reduce
	globalDotP := MapReduce(v1, v2, N, context)
	context.WG.Wait()
	return globalDotP

}

// Choose between the parallel and serial vector dot products
func DotPWrapper(v1, v2 *[]float64, N int, context *MapReduceContext) float64 {
	if context.NThreads == 0 {
		return DotP(v1, v2, N)
	} else {
		return ParallelDotP(v1, v2, N, context)
	}
}
