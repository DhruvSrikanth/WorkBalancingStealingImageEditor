package concurrent

import (
	"math/rand"
	"sync"
	"time"
)

type sharedContextWB struct {
	capacity         int
	thresholdQueue   int
	thresholdBalance int
	queues           []DEQueue
	wg               *sync.WaitGroup
}
type balancer struct {
	workers         []*workerWB
	done            bool
	prevDistributee int
	context         *sharedContextWB
}

type workerWB struct {
	id            int
	context       *sharedContextWB
	randGen       *rand.Rand
	workRemaining bool
}

func NewWorkerWB(id int, context *sharedContextWB) *workerWB {
	return &workerWB{
		id:            id,
		context:       context,
		randGen:       rand.New(rand.NewSource(time.Now().UnixNano())),
		workRemaining: true,
	}
}

func (worker *workerWB) isWorkPoolEmpty() bool {
	// Check if all queues are empty
	for _, queue := range worker.context.queues {
		if !queue.IsEmpty() {
			return false
		}
	}
	return true
}

func (worker *workerWB) getVictim() int {
	// Get random victim
	victim := worker.randGen.Intn(worker.context.capacity)

	// Make sure the victim is not the worker
	for victim == worker.id {
		victim = worker.randGen.Intn(worker.context.capacity)
	}
	return victim
}

func balancePolicy(smallQueue, largeQueue DEQueue, lambda int) bool {
	// Check if the queues need to be balanced
	return largeQueue.Size()-smallQueue.Size() >= lambda
}

func (worker *workerWB) balance() {
	// Get the victim to balance with
	victim := worker.getVictim()

	// Canonical ordering of victim and worker
	var min int
	var max int
	if victim < worker.id {
		min = victim
		max = worker.id
	} else {
		min = worker.id
		max = victim
	}

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

func (worker *workerWB) work() {
	// Worker loops if work is remaining in its own queue or the overall work pool
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
		// and at random
		queueSize := worker.context.queues[worker.id].Size()
		if worker.context.capacity > 1 && queueSize == worker.randGen.Intn(queueSize+1) {
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

func (service *balancer) nextDistributee() int {
	// Get next distributee and update prevDistributee
	service.prevDistributee = (service.prevDistributee + 1) % service.context.capacity
	return service.prevDistributee
}

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

func (service *balancer) Shutdown() {
	// Indicate the service is done
	for _, worker := range service.workers {
		worker.workRemaining = false
	}
	service.done = true
	// Wait for all workers to finish
	service.context.wg.Wait()
}
