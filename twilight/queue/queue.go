package queue

type Queue[T any] struct {
	front *QueueItem[T]
}
type QueueItem[T any] struct {
	previous *QueueItem[T]
	next     *QueueItem[T]

	Value T
}

func (q *Queue[T]) First() *QueueItem[T] {
	if q.front == nil {
		return nil
	}

	return q.front
}
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

func (q *Queue[T]) Dequeue() *QueueItem[T] {
	if q.front == nil {
		return nil
	}

	item := q.front
	item.previous = nil
	q.front = item.next

	return item
}
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

func New[T any]() *Queue[T] {
	return &Queue[T]{}
}
