package concurrent

/**** YOU CANNOT MODIFY ANY OF THE FOLLOWING INTERFACES/TYPES ********/
type Task interface{}

type DEQueue interface {
	PushBottom(task Task)
	IsEmpty() bool //returns whether the queue is empty
	PopTop() Task
	PopBottom() Task
}

/******** DO NOT MODIFY ANY OF THE ABOVE INTERFACES/TYPES *********************/

// NewUnBoundedDEQueue returns an empty UnBoundedDEQueue
func NewUnBoundedDEQueue() DEQueue {
	/** TODO: Remove the return nil and implement this function **/
	return nil
}
