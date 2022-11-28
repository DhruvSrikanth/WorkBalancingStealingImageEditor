package concurrent

/**** YOU CANNOT MODIFY ANY OF THE FOLLOWING INTERFACES ********/

// Runnable represents a task that does not return a value.
type Runnable interface {
	Run() // Starts the execution of a Runnable
}

// Callable represents a task that will return a value.
type Callable interface {
	Call() interface{} // Starts the execution of a Callable
}

// Future represents the value that is returned after executing a Runnable or Callable task.
type Future interface {
	// Get waits (if necessary) for the task to complete. If the task associated with the Future is a Callable Task then it will return the value returned by the Call method. If the task associated with the Future is a Runnable then it must return nil once the task is complete.
	Get() interface{}
}

// ExecutorService represents a service that can run om Runnable and/or Callable tasks concurrently.
type ExecutorService interface {

	// Submits a task for execution and returns a Future representing that task.
	Submit(task interface{}) Future

	// Shutdown initiates a shutdown of the service. It is unsafe to call Shutdown at the same time as the Submit method. All tasks must be submitted before calling Shutdown. All Submit calls during and after the call to the Shutdown method will be ignored. A goroutine that calls Shutdown is blocked until the service is completely shutdown (i.e., no more pending tasks and all goroutines spawned by the service are terminated).
	Shutdown()
}

/******** DO NOT MODIFY ANY OF THE ABOVE INTERFACES *********************/
