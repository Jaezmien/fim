package utilities

import (
	"errors"
	"fmt"
	"strings"
)

// An ErrorPair contains details about the origin of an error relative
// to the source code.
type ErrorPair struct {
	// 1-based line number of the error
	Line int
	// 1-based column number of the error
	Column int
}

// Create an ErrorPair based on a character index.
func GetErrorIndexPair(source string, index int) *ErrorPair {
	content := source[0:min(index+1, len(source))]
	lines := strings.Split(content, "\n")

	return &ErrorPair{
		Line:   len(lines),
		Column: len(lines[len(lines)-1]),
	}
}

// Create an error based on a character index from the source code.
func CreateErrorFromIndex(source string, index int, errorMessage string) error {
	pair := GetErrorIndexPair(source, index)
	return errors.New(fmt.Sprintf("[line %d:%d] %s", pair.Line, pair.Column, errorMessage))
}
