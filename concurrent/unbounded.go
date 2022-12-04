package concurrent

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

/**** YOU CANNOT MODIFY ANY OF THE FOLLOWING INTERFACES/TYPES ********/
type Task interface{}

type DEQueue interface {
	PushBottom(task Task)
	IsEmpty() bool //returns whether the queue is empty
	PopTop() Task
	PopBottom() Task
	Size() int
}

/******** DO NOT MODIFY ANY OF THE ABOVE INTERFACES/TYPES *********************/

// Node in the double ended unbounded queue
type Node struct {
	task Task
	next unsafe.Pointer
	prev unsafe.Pointer
}

func newNode(task Task) *Node {
	node := unsafe.Pointer(&Node{
		task: task,
		next: nil,
		prev: nil,
	})
	return (*Node)(node)

}

// UnBoundedDEQueue is a double ended unbounded queue
type UnBoundedDEQueue struct {
	head *Node // bottom part of the queue
	tail *Node // top part of the queue
	size int64
	lock *sync.Mutex // Lock for the queue that is used by the executor
}

// Visualized representation of the queue
// (top) prev -> tail -> head -> next (bottom)

// NewUnBoundedDEQueue returns an empty UnBoundedDEQueue
func NewUnBoundedDEQueue() DEQueue {
	return &UnBoundedDEQueue{
		head: nil,
		tail: nil,
		size: 0,
		lock: &sync.Mutex{},
	}
}

// PushBottom adds a task to the bottom of the queue
func (q *UnBoundedDEQueue) PushBottom(task Task) {
	// Lock the queue and defer the unlock
	q.lock.Lock()
	defer q.lock.Unlock()

	// Create a new node
	node := newNode(task)

	// Increase the size of the queue
	q.size += 1

	// Check if the queue is empty
	if q.head == nil {
		// Set the head and tail to the new node
		q.head = node
		q.tail = node
	} else {
		// Set the next node of the head to the new node
		q.head.next = unsafe.Pointer(node)
		// Set the prev node of the new node to the head
		node.prev = unsafe.Pointer(q.head)
		// Set the head to the new node
		q.head = node
	}
}

// PopBottom removes a task from the bottom of the queue
func (q *UnBoundedDEQueue) PopBottom() Task {
	// Lock the queue and defer the unlock
	q.lock.Lock()
	defer q.lock.Unlock()

	// Check if the queue is empty
	if q.size == 0 {
		return nil
	}

	// Decrease the size of the queue
	q.size -= 1

	// Check if the queue has only one element
	if q.head == q.tail {
		// Get the task
		task := q.head.task
		// Reset the deque
		q.head = nil
		q.tail = nil
		return task
	}

	// Get the task from the head
	task := q.head.task
	// Set the head to the prev node
	q.head = (*Node)(q.head.prev)
	// Set the next node to nil
	q.head.next = unsafe.Pointer(nil)

	return task

}

// PopTop removes a task from the top of the queue
func (q *UnBoundedDEQueue) PopTop() Task {
	// Lock the queue and defer the unlock
	q.lock.Lock()
	defer q.lock.Unlock()

	// Check if the queue is empty
	if q.size == 0 {
		return nil
	}

	// Decrease the size of the queue
	q.size -= 1

	// Check if the queue has only one element
	if q.head == q.tail {
		// Get the task
		task := q.tail.task
		// Reset the deque
		q.head = nil
		q.tail = nil
		return task
	}

	// Get the task from the tail
	task := q.tail.task
	// Set the tail to the next node
	q.tail = (*Node)(q.tail.next)
	// Set the prev node to nil
	q.tail.prev = unsafe.Pointer(nil)

	return task
}

// IsEmpty returns whether the queue is empty
func (q *UnBoundedDEQueue) IsEmpty() bool {
	// Lock the queue and defer the unlock
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.size == 0
}

// Size returns the size of the queue
func (q *UnBoundedDEQueue) Size() int {
	return int(atomic.LoadInt64(&q.size))
}
