package stack

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Len() int {
	return len(s.items)
}

func (s *Stack[T]) IsEmpty() bool {
	return s.Len() == 0
}
func (s *Stack[T]) PeekAt(index int) T {
	return s.items[index]
}
func (s *Stack[T]) Peek() T {
	return s.PeekAt(s.Len() - 1)
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() *T {
	if s.IsEmpty() {
		return nil
	}

	item := s.Peek()
	s.items = s.items[:s.Len() - 1]
	return &item
}

func (q *Stack[T]) Flatten() []T {
	s := make([]T, 0, q.Len())

	for q.Len() > 0 {
		el := *q.Pop()
		s = append(s, el)
	}

	return s
}

func New[T any]() *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0),
	}
}

