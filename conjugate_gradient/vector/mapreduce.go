package vector

import (
	"sync"
)

type vectorTask struct {
	v1 []float64
	v2 []float64
}

type MapReduceContext struct {
	NThreads     int // Runs the parallel version of the program with the specified number of threads
	MapWG        *sync.WaitGroup
	ReduceWG     *sync.WaitGroup
	id           int
	Benchmarking bool
}

func mapper(mapperIn <-chan vectorTask, mapperOut chan float64, reducerOut chan<- float64, context *MapReduceContext) {
	for {
		// Get the vector pair
		vecTask, ok := <-mapperIn
		if !ok {
			break
		}

		// Compute the dot product
		localDotP := DotP(&vecTask.v1, &vecTask.v2, len(vecTask.v1))

		// Send the dot product to the reducer
		mapperOut <- localDotP
	}

	context.MapWG.Done()
}

func reducer(reducerIn <-chan float64, reducerOut chan<- float64, context *MapReduceContext) {
	globalDotP := 0.0
	for localDotP := range reducerIn {
		globalDotP += localDotP
	}
	reducerOut <- globalDotP
	close(reducerOut)
	context.ReduceWG.Done()
}

// MapReduce computes the dot product of two vectors in parallel
func MapReduce(v1, v2 *[]float64, N int, context *MapReduceContext) float64 {
	vector1 := *v1
	vector2 := *v2

	// Create the input channel
	mapperIn := make(chan vectorTask, context.NThreads-1)

	// Create the output channel
	mapperOut := make(chan float64, context.NThreads-1)

	// Create the reducer output channel
	reducerOut := make(chan float64, 1)

	// Start the mappers
	for i := 0; i < context.NThreads-1; i++ {
		context.MapWG.Add(1)
		context.id = i
		go mapper(mapperIn, mapperOut, reducerOut, context)
	}

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
		vecTask := vectorTask{v1: vector1[start:end], v2: vector2[start:end]}

		// Send the vector pair to the mapper
		mapperIn <- vecTask
	}

	// Close the input channel
	close(mapperIn)

	context.MapWG.Wait()

	// Only have the one thread close the output channel
	// Once all the tasks have been added to the output channel, the channel can be closed
	if context.id == context.NThreads-2 {
		close(mapperOut)
		// Once the threads have added values to the output channel, the reducer can start
		context.ReduceWG.Add(1)
		context.id = context.NThreads - 1
		go reducer(mapperOut, reducerOut, context)
	}

	context.ReduceWG.Wait()

	// Return the dot product
	return <-reducerOut
}

// ParallelDotP computes the dot product of two vectors safely in parallel
func ParallelDotP(v1, v2 *[]float64, N int, context *MapReduceContext) float64 {
	// Run the map reduce
	globalDotP := MapReduce(v1, v2, N, context)
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
