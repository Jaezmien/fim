package queue

// The Queue represents a Queue data structure
type Queue[T any] struct {
	front *QueueItem[T]
}

// The QueueItem is a wrapper to allow double-linked list 
type QueueItem[T any] struct {
	previous *QueueItem[T]
	next     *QueueItem[T]

	Value T
}

// Get the first item inserted into the queue
func (q *Queue[T]) First() *QueueItem[T] {
	if q.front == nil {
		return nil
	}

	return q.front
}

// Get the last item inserted into the queue
func (q *Queue[T]) Last() *QueueItem[T] {
	if q.front == nil {
		return nil
	}

	item := q.front
	for item.next != nil {
		item = item.next
	}
	return item
}

// Peek the n-th item inserted into the queue
func (q *Queue[T]) Peek(index int) *QueueItem[T] {
	if q.front == nil {
		return nil
	}

	item := q.front
	for item.next != nil && index > 0 {
		item = item.next
		index -= 1
	}
	if item == nil {
		return nil
	}
	return item
}

// Return the amount of items in the queue
func (q *Queue[T]) Len() int {
	if q.front == nil {
		return 0
	}

	length := 0

	item := q.front
	length += 1
	for item.next != nil {
		item = item.next
		length += 1
	}

	return length
}

// Remove the first item from the start of the queue
func (q *Queue[T]) Dequeue() *QueueItem[T] {
	if q.front == nil {
		return nil
	}

	item := q.front
	item.previous = nil
	q.front = item.next

	return item
}

// Insert an item into the back of the queue
func (q *Queue[T]) Queue(value T) {
	item := &QueueItem[T]{
		Value: value,
	}

	if q.front == nil {
		q.front = item
		return
	}

	last := q.front
	for last.next != nil {
		last = last.next
	}

	last.next = item
	item.previous = last
}

// Insert an item into the front of the queue.
// Preposterous!
func (q *Queue[T]) QueueFront(value T) {
	item := &QueueItem[T]{
		Value: value,
	}

	if q.front == nil {
		q.front = item
		return
	}

	item.next = q.front
	q.front.previous = item

	q.front = item
}

// Turn the queue into an array
func (q *Queue[T]) Flatten() []T {
	s := make([]T, 0, q.Len())

	idx := 0
	for q.Len() > 0 {
		el := q.Dequeue()
		s = append(s, el.Value)
		idx += 1
	}

	return s
}

// Create a new Queue
func New[T any]() *Queue[T] {
	return &Queue[T]{}
}
