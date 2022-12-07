package concurrent

import (
	"math/rand"
	"sync"
	"time"
)

// Shared Context for Work Stealing
type sharedContextST struct {
	capacity int
	queues   []DEQueue
	wg       *sync.WaitGroup
}

// Work Stealing Stealer
type stealer struct {
	workers         []*workerST
	done            bool
	prevDistributee int
	context         *sharedContextST
}

// Work Stealing Worker
type workerST struct {
	id            int
	context       *sharedContextST
	randGen       *rand.Rand
	workRemaining bool
	victimOptions []int
}

// Returns a new Work Stealing Stealer
func NewWorkerST(id int, context *sharedContextST) *workerST {
	// Add the appropriate stealing options
	victims := []int{}
	for i := 0; i < context.capacity; i++ {
		// Worker cant steal from itself
		if i != id {
			victims = append(victims, i)
		}
	}

	return &workerST{
		id:            id,
		context:       context,
		randGen:       rand.New(rand.NewSource(time.Now().UnixNano())),
		workRemaining: true,
		victimOptions: victims,
	}
}

// Check if all queues are empty
func (worker *workerST) isWorkPoolEmpty() bool {
	// Check if all queues are empty
	for _, queue := range worker.context.queues {
		if !queue.IsEmpty() {
			return false
		}
	}
	return true
}

// Get a random victim
func (worker *workerST) getVictim() int {
	// Get random victim
	victimIdx := worker.randGen.Intn(len(worker.victimOptions))
	// Victim will not be the worker itself since we have removed it from the options
	// Remove the victim from the options
	victim := worker.victimOptions[victimIdx]
	worker.victimOptions = append(worker.victimOptions[:victimIdx], worker.victimOptions[victimIdx+1:]...)
	return victim
}

func stealingPolicy(smallQueue, largeQueue DEQueue) bool {
	// Check if the size of the smaller queue is 0 and the larger queue has at least 1 element
	return smallQueue.IsEmpty() && !largeQueue.IsEmpty()
}

// Steal a task from the victim
func (worker *workerST) steal() {
	// Get the victim to balance with
	victimIdx := worker.getVictim()
	victim := worker.context.queues[victimIdx]
	workerQueue := worker.context.queues[worker.id]

	// Check if the queues need to be balanced
	if stealingPolicy(workerQueue, victim) {
		// Get the task to move
		job := victim.PopBottom()

		// Add the task to the small queue
		if job != nil {
			workerQueue.PushBottom(job)
			// Add all options back to the stealing options
			for i := 0; i < worker.context.capacity; i++ {
				if i != worker.id {
					worker.victimOptions = append(worker.victimOptions, i)
				}
			}
		}
	}
}

// Worker routine
func (worker *workerST) work() {
	// Worker loops if work is remaining in its own queue or the overall work pool
	for worker.workRemaining || !worker.isWorkPoolEmpty() {
		// Finish all of your own tasks before stealing
		workerTask := worker.context.queues[worker.id].PopTop()
		for workerTask != nil {
			runnable, ok := workerTask.(Runnable)
			if ok {
				// Run the task
				runnable.Run()
			}
			// Get the next task
			workerTask = worker.context.queues[worker.id].PopTop()
		}

		if worker.context.capacity > 1 && len(worker.victimOptions) > 0 {
			worker.steal()
		}
	}

	// Worker is done
	worker.context.wg.Done()
}

// NewWorkStealingExecutor returns an ExecutorService that is implemented using the work-stealing algorithm.
// @param capacity - The number of goroutines in the pool
// @param threshold - The number of items that a goroutine in the pool can
// grab from the executor in one time period. For example, if threshold = 10
// this means that a goroutine can grab 10 items from the executor all at
// once to place into their local queue before grabbing more items. It's
// not required that you use this parameter in your implementation.
func NewWorkStealingExecutor(capacity, threshold int) ExecutorService {
	// Create capacity queues
	queues := []DEQueue{}
	for i := 0; i < capacity; i++ {
		queues = append(queues, NewUnBoundedDEQueue())
	}

	// Create shared context
	context := &sharedContextST{
		capacity: capacity,
		queues:   queues,
		wg:       &sync.WaitGroup{},
	}

	// Create capacity workers
	workers := []*workerST{}
	for i := 0; i < capacity; i++ {
		workers = append(workers, NewWorkerST(i, context))
	}

	// Spawn worker routines
	for _, worker := range workers {
		// Increment wait group (Decrement inside the Run() method)
		context.wg.Add(1)
		go worker.work()
	}

	// Create service
	service := &stealer{
		workers:         workers,
		done:            false,
		prevDistributee: capacity - 1,
		context:         context,
	}

	return service
}

// Get index of the next worker to distribute work to
func (service *stealer) nextDistributee() int {
	// Get next distributee and update prevDistributee
	service.prevDistributee = (service.prevDistributee + 1) % service.context.capacity
	return service.prevDistributee
}

// Submit a task to the executor
func (service *stealer) Submit(task interface{}) Future {
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
func (service *stealer) Shutdown() {
	// Indicate no more work is remaining
	for _, worker := range service.workers {
		worker.workRemaining = false
	}

	// Indicate the service is done
	service.done = true

	// Wait for all workers to finish
	service.context.wg.Wait()
}
