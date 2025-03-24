package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	t.Run("should hold a singular value", func(t *testing.T) {
		q := New[int]()
		assert.Equal(t, 0, q.Len(), "Should be empty")

		q.Push(1)
		assert.NotNil(t, q.Peek(), "Should have top element")
		assert.Equal(t, 1, q.Len(), "Should have non-zero element count")

		value := *q.Pop()
		assert.Equal(t, 1, value, "Should be equal")
	})
	t.Run("should queue normally", func(t *testing.T) {
		q := New[int]()
		assert.Equal(t, 0, q.Len(), "Should be empty")

		q.Push(1)
		q.Push(2)
		q.Push(3)
		assert.Equal(t, 3, q.Len(), "Should have non-zero element count")

		var value int
		value = *q.Pop()
		assert.Equal(t, 3, value, "Should be equal")
		value = *q.Pop()
		assert.Equal(t, 2, value, "Should be equal")
		value = *q.Pop()
		assert.Equal(t, 1, value, "Should be equal")

		assert.Equal(t, q.Len(), 0, "Should be empty")
	})
	t.Run("should peek normally", func(t *testing.T) {
		q := New[int]()
		assert.Equal(t, 0, q.Len(), "Should be empty")

		q.Push(1)
		q.Push(2)
		q.Push(3)
		assert.Equal(t, 3, q.Len(), "Should have non-zero element count")

		var value int
		value = q.Peek()
		assert.Equal(t, 3, value, "Should be equal")
		value = q.PeekAt(0)
		assert.Equal(t, 1, value, "Should be equal")

		assert.Equal(t, 3, q.Len(), "Should have non-zero element count")
	})
}
