package twilight

import (
	"slices"
	"testing"

	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
	"github.com/stretchr/testify/assert"
)

func TestRuneSplittable(t *testing.T) {
	t.Run("should be splittable", func(t *testing.T) {
		assert.True(t, slices.Contains(splittable_runes[:], ' '), "Rune should be splittable")
	})
	t.Run("should not be splittable", func(t *testing.T) {
		assert.False(t, slices.Contains(splittable_runes[:], 'a'), "Rune should not be splittable")
	})
}

func TestPartialTokens(t *testing.T) {
	t.Run("should create one partial token", func(t *testing.T) {
		source := "hello"

		l := createPartialTokens(source)

		assert.Equal(t, 1, l.Len(), "Expected to have generated one partial token")

		token := l.First().Value

		assert.Equal(t, "hello", token.Value, "Expected token value to be same as source")
	})

	t.Run("should create three partial token", func(t *testing.T) {
		source := "hello world"

		l := createPartialTokens(source)

		assert.Equal(t, 3, l.Len(), "Expected to have generated three partial token")

		var token *token.Token
		token = l.Dequeue().Value
		assert.Equal(t, "hello", token.Value, "Expected 'hello'")

		token = l.Dequeue().Value
		assert.Equal(t, " ", token.Value, "Expected whitespace")

		token = l.Dequeue().Value
		assert.Equal(t, "world", token.Value, "Expected 'world'")

		assert.Equal(t, 0, l.Len(), "Expected to have no tokens left")
	})

	t.Run("should clean tokens", func(t *testing.T) {
		source :=
			`hello
			world
		bye`

		l := createPartialTokens(source)
		l = mergePartialTokens(l)

		var token *token.Token
		token = l.Dequeue().Value
		assert.Equal(t, "hello", token.Value, "Expected 'hello'")

		token = l.Dequeue().Value
		assert.Equal(t, "\n", token.Value, "Expected newline")

		token = l.Dequeue().Value
		assert.Equal(t, "world", token.Value, "Expected 'world'")

		token = l.Dequeue().Value
		assert.Equal(t, "\n", token.Value, "Expected newline")

		token = l.Dequeue().Value
		assert.Equal(t, "bye", token.Value, "Expected 'bye'")
	})
}

func TestPartialMerging(t *testing.T) {
	t.Run("should merge decimals", func(t *testing.T) {
		source := "1.0"

		l := createPartialTokens(source)
		assert.Equal(t, "1", l.Peek(0).Value.Value, "Should be '1'")
		assert.Equal(t, ".", l.Peek(1).Value.Value, "Should be '.'")
		assert.Equal(t, "0", l.Peek(2).Value.Value, "Should be '0'")

		l = mergePartialTokens(l)

		assert.Equal(t, 1, l.Len(), "Should be one token")

		token := l.Dequeue().Value
		assert.Equal(t, "1.0", token.Value, "Expected '1.0'")
	})
	t.Run("should merge negative decimals", func(t *testing.T) {
		source := "-1.0"

		l := createPartialTokens(source)
		assert.Equal(t, "-1", l.Peek(0).Value.Value, "Should be '-1'")
		assert.Equal(t, ".", l.Peek(1).Value.Value, "Should be '.'")
		assert.Equal(t, "0", l.Peek(2).Value.Value, "Should be '0'")

		l = mergePartialTokens(l)

		assert.Equal(t, 1, l.Len(), "Should be one token")

		token := l.Dequeue().Value
		assert.Equal(t, "-1.0", token.Value, "Expected '-1.0'")
	})
	t.Run("should merge string", func(t *testing.T) {
		source := "\"hello world\""

		l := createPartialTokens(source)
		assert.Equal(t, "\"", l.Peek(0).Value.Value, "Should be '\"'")
		assert.Equal(t, "hello", l.Peek(1).Value.Value, "Should be 'hello'")
		assert.Equal(t, " ", l.Peek(2).Value.Value, "Should be ' '")
		assert.Equal(t, "world", l.Peek(3).Value.Value, "Should be 'world'")
		assert.Equal(t, "\"", l.Peek(4).Value.Value, "Should be '\"'")

		l = mergePartialTokens(l)

		assert.Equal(t, 1, l.Len(), "Should be one token")

		token := l.Dequeue().Value
		assert.Equal(t, "\"hello world\"", token.Value, "Expected '\"hello world\"'")
	})
	t.Run("should merge escaped string", func(t *testing.T) {
		source := `"hello \"world\""`

		l := createPartialTokens(source)
		assert.Equal(t, "\"", l.Peek(0).Value.Value, "Should be '\"'")
		assert.Equal(t, "hello", l.Peek(1).Value.Value, "Should be 'hello'")
		assert.Equal(t, " ", l.Peek(2).Value.Value, "Should be ' '")
		assert.Equal(t, "\\", l.Peek(3).Value.Value, "Should be '\"'")
		assert.Equal(t, "\"", l.Peek(4).Value.Value, "Should be '\"'")
		assert.Equal(t, "world", l.Peek(5).Value.Value, "Should be 'world'")
		assert.Equal(t, "\\", l.Peek(6).Value.Value, "Should be '\"'")
		assert.Equal(t, "\"", l.Peek(7).Value.Value, "Should be '\"'")
		assert.Equal(t, "\"", l.Peek(8).Value.Value, "Should be '\"'")

		l = mergePartialTokens(l)

		assert.Equal(t, 1, l.Len(), "Should be one token")

		token := l.Dequeue().Value
		assert.Equal(t, `"hello \"world\""`, token.Value, `Expected '"hello \"world\""'`)
	})
	t.Run("should merge character", func(t *testing.T) {
		source := "'a'"

		l := createPartialTokens(source)
		assert.Equal(t, "'", l.Peek(0).Value.Value, "Should be '''")
		assert.Equal(t, "a", l.Peek(1).Value.Value, "Should be 'a'")
		assert.Equal(t, "'", l.Peek(4).Value.Value, "Should be '''")

		l = mergePartialTokens(l)

		assert.Equal(t, 1, l.Len(), "Should be one token")

		token := l.Dequeue().Value
		assert.Equal(t, "'a'", token.Value, "Expected ''a''")
	})
	t.Run("should merge delimeters", func(t *testing.T) {
		source := "(hello world)"

		l := createPartialTokens(source)
		assert.Equal(t, "(", l.Peek(0).Value.Value, "Should be '('")
		assert.Equal(t, "hello", l.Peek(1).Value.Value, "Should be 'hello'")
		assert.Equal(t, " ", l.Peek(2).Value.Value, "Should be ' '")
		assert.Equal(t, "world", l.Peek(3).Value.Value, "Should be 'world'")
		assert.Equal(t, ")", l.Peek(4).Value.Value, "Should be ')'")

		l = mergePartialTokens(l)

		assert.Equal(t, 1, l.Len(), "Should be one token")

		token := l.Dequeue().Value
		assert.Equal(t, "(hello world)", token.Value, "Expected '(hello world)'")
	})
}

func TestTokenMerging(t *testing.T) {
	t.Run("should work with one token", func(t *testing.T) {
		source := "hello"
		l := createPartialTokens(source)
		token := utilities.MergeTokens(l, 1)

		assert.Equal(t, token.Value, source, "Should be same as source")
	})
	t.Run("should work with multiple tokens", func(t *testing.T) {
		source := "hello world"
		l := createPartialTokens(source)
		token := utilities.MergeTokens(l, 3)

		assert.Equal(t, token.Value, source, "Should be same as source")
	})
	t.Run("should consume only the given amount", func(t *testing.T) {
		source := "hello world"
		l := createPartialTokens(source)
		token := utilities.MergeTokens(l, 2)

		assert.Equal(t, token.Value, "hello ", "Should be 'hello '")
		assert.Equal(t, l.First().Value.Value, "world", "Should be 'world'")
	})
}
