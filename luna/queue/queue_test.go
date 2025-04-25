package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	t.Run("should hold a singular value", func(t *testing.T) {
		q := New[int]()
		assert.Equal(t, 0, q.Len(), "Should be empty")

		q.Queue(1)
		assert.NotNil(t, q.First(), "Should have first element")
		assert.Equal(t, 1, q.Len(), "Should have non-zero element count")

		value := q.Dequeue().Value
		assert.Equal(t, 1, value, "Should be equal")
	})
	t.Run("should queue normally", func(t *testing.T) {
		q := New[int]()
		assert.Equal(t, 0, q.Len(), "Should be empty")

		q.Queue(1)
		q.Queue(2)
		q.Queue(3)
		assert.Equal(t, 3, q.Len(), "Should have non-zero element count")

		var value int
		value = q.Dequeue().Value
		assert.Equal(t, 1, value, "Should be equal")
		value = q.Dequeue().Value
		assert.Equal(t, 2, value, "Should be equal")
		value = q.Dequeue().Value
		assert.Equal(t, 3, value, "Should be equal")

		assert.Equal(t, q.Len(), 0, "Should be empty")
	})
	t.Run("should peek normally", func(t *testing.T) {
		q := New[int]()
		assert.Equal(t, 0, q.Len(), "Should be empty")

		q.Queue(1)
		q.Queue(2)
		q.Queue(3)
		assert.Equal(t, 3, q.Len(), "Should have non-zero element count")

		var value int
		value = q.First().Value
		assert.Equal(t, 1, value, "Should be equal")
		value = q.Last().Value
		assert.Equal(t, 3, value, "Should be equal")
		value = q.Peek(1).Value
		assert.Equal(t, 2, value, "Should be equal")

		assert.Equal(t, 3, q.Len(), "Should have non-zero element count")
	})
}
