package concurrent

import (
	"math/rand"
	"sync"
	"time"
)

// Shared Context for Work Balancing
type sharedContextWB struct {
	capacity         int
	thresholdQueue   int
	thresholdBalance int
	queues           []DEQueue
	wg               *sync.WaitGroup
}

// Work Balancing Balancer
type balancer struct {
	workers         []*workerWB
	done            bool
	prevDistributee int
	context         *sharedContextWB
}

// Work Balancing Worker
type workerWB struct {
	id            int
	context       *sharedContextWB
	randGen       *rand.Rand
	workRemaining bool
}

// Returns a new Work Balancing Balancer
func NewWorkerWB(id int, context *sharedContextWB) *workerWB {
	return &workerWB{
		id:            id,
		context:       context,
		randGen:       rand.New(rand.NewSource(time.Now().UnixNano())),
		workRemaining: true,
	}
}

// Check if all queues are empty
func (worker *workerWB) isWorkPoolEmpty() bool {
	// Check if all queues are empty
	for _, queue := range worker.context.queues {
		if !queue.IsEmpty() {
			return false
		}
	}
	return true
}

// Get a random victim
func (worker *workerWB) getVictim() int {
	// Get random victim
	victim := worker.randGen.Intn(worker.context.capacity)

	// Make sure the victim is not the worker
	for victim == worker.id {
		victim = worker.randGen.Intn(worker.context.capacity)
	}
	return victim
}

// Balance policy - returns true if the queues need to be balanced based on the balancing threshold provided
func balancePolicy(smallQueue, largeQueue DEQueue, lambda int) bool {
	// Check if the queues need to be balanced
	return largeQueue.Size()-smallQueue.Size() >= lambda
}

// Balance two queues
func (worker *workerWB) balance() {
	// Get the victim to balance with
	victim := worker.getVictim()

	// Ordering of victim and worker
	var min int
	var max int
	if victim < worker.id {
		min = victim
		max = worker.id
	} else {
		min = worker.id
		max = victim
	}

	// Determine which queue is smaller and which is larger
	var smallQueue DEQueue
	var largeQueue DEQueue

	// Get the queues to balance
	if worker.context.queues[min].Size() < worker.context.queues[max].Size() {
		smallQueue = worker.context.queues[min]
		largeQueue = worker.context.queues[max]
	} else {
		smallQueue = worker.context.queues[max]
		largeQueue = worker.context.queues[min]
	}

	// Check if the queues need to be balanced
	for balancePolicy(smallQueue, largeQueue, worker.context.thresholdBalance) {
		// Get the task to move
		job := largeQueue.PopBottom()

		// Add the task to the small queue
		if job != nil {
			smallQueue.PushBottom(job)
		}
	}
}

// Worker routine
func (worker *workerWB) work() {
	// Worker loops if work is remaining in the overall work pool and if worker's local queue is not empty
	for worker.workRemaining || !worker.isWorkPoolEmpty() {
		// Get the next task
		workerTask := worker.context.queues[worker.id].PopTop()
		if workerTask != nil {
			runnable, ok := workerTask.(Runnable)
			if ok {
				// Run the task
				runnable.Run()
			}
		}

		// Rebalancing is only done if there is more than one worker
		// and at random (i.e. 1/n chance where n is the size of the queue)
		queueSize := worker.context.queues[worker.id].Size()
		if worker.context.capacity > 1 && queueSize == worker.randGen.Intn(queueSize+1) {
			// Balance the queues
			worker.balance()
		}
	}

	// Worker is done
	worker.context.wg.Done()
}

// NewWorkBalancingExecutor returns an ExecutorService that is implemented using the work-balancing algorithm.
// @param capacity - The number of goroutines in the pool
// @param threshold - The number of items that a goroutine in the pool can
// grab from the executor in one time period. For example, if threshold = 10
// this means that a goroutine can grab 10 items from the executor all at
// once to place into their local queue before grabbing more items. It's
// not required that you use this parameter in your implementation.
// @param thresholdBalance - The threshold used to know when to perform
// balancing. Remember, if two local queues are to be balanced the
// difference in the sizes of the queues must be greater than or equal to
// thresholdBalance. You must use this parameter in your implementation.
func NewWorkBalancingExecutor(capacity, thresholdQueue, thresholdBalance int) ExecutorService {
	// Create capacity queues
	queues := []DEQueue{}
	for i := 0; i < capacity; i++ {
		queues = append(queues, NewUnBoundedDEQueue())
	}

	// Create shared context
	context := &sharedContextWB{
		capacity:         capacity,
		thresholdQueue:   thresholdQueue,
		thresholdBalance: thresholdBalance,
		queues:           queues,
		wg:               &sync.WaitGroup{},
	}

	// Create capacity workers
	workers := []*workerWB{}
	for i := 0; i < capacity; i++ {
		workers = append(workers, NewWorkerWB(i, context))
	}

	// Spawn worker routines
	for _, worker := range workers {
		// Increment wait group (Decrement inside the Run() method)
		context.wg.Add(1)
		go worker.work()
	}

	// Create service
	service := &balancer{
		workers:         workers,
		done:            false,
		prevDistributee: capacity - 1,
		context:         context,
	}

	return service
}

// Get the next worker idx to distribute work to
func (service *balancer) nextDistributee() int {
	// Get next distributee and update prevDistributee
	service.prevDistributee = (service.prevDistributee + 1) % service.context.capacity
	return service.prevDistributee
}

// Submit a task to the executor
func (service *balancer) Submit(task interface{}) Future {
	// Check if service is done
	if service.done {
		return nil
	}

	// Get next distributee
	distributee := service.nextDistributee()

	// Add task to distributee's queue
	service.context.queues[distributee].PushBottom(task)

	// Type assertion
	return task.(Future)
}

// Shutdown the executor
func (service *balancer) Shutdown() {
	// Indicate the service is done for all workers
	for _, worker := range service.workers {
		worker.workRemaining = false
	}

	// Indicate the service is done
	service.done = true

	// Wait for all workers to finish
	service.context.wg.Wait()
}
