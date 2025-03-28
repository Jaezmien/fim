package stack

// The Stack represents a Stack data structure
type Stack[T any] struct {
	items []T
}

// Return the amount of items in the queue
func (s *Stack[T]) Len() int {
	return len(s.items)
}

// Checks if the stack is empty
func (s *Stack[T]) IsEmpty() bool {
	return s.Len() == 0
}

// Returns the item at the n-th index
func (s *Stack[T]) PeekAt(index int) T {
	return s.items[index]
}

// Returns the item at the top of the stack
// without removing it
func (s *Stack[T]) Peek() T {
	return s.PeekAt(s.Len() - 1)
}

// Inserts an item at the bottom of the stack
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Returns the item at the top of the stack
func (s *Stack[T]) Pop() *T {
	if s.IsEmpty() {
		return nil
	}

	item := s.Peek()
	s.items = s.items[:s.Len()-1]
	return &item
}

// Return the queue as an array
func (q *Stack[T]) Flatten() []T {
	s := make([]T, 0, q.Len())

	for q.Len() > 0 {
		el := *q.Pop()
		s = append(s, el)
	}

	return s
}

// Create a new Stack
func New[T any]() *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0),
	}
}
