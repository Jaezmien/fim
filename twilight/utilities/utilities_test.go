package utilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringNumber(t *testing.T) {
	assert.True(t, IsStringNumber("1"), "Should be a number")
	assert.True(t, IsStringNumber("10"), "Should be a number")
	assert.False(t, IsStringNumber("a"), "Should not be a number")
	assert.False(t, IsStringNumber("Zz"), "Should not be a number")
}
