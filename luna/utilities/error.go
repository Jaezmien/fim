package utilities

import (
	"strings"
)

// An ErrorPair contains details about the origin of an error relative
// to the source code.
type ErrorPair struct {
	// 1-based line number of the error
	Line   int
	// 1-based column number of the error
	Column int
}

// Create an error pair based on a character index.
func GetErrorIndexPair(source string, index int) *ErrorPair {
	content := source[0:min(index+1, len(source))]
	lines := strings.Split(content, "\n")

	return &ErrorPair{
		Line:   len(lines),
		Column: len(lines[len(lines)-1]),
	}
}
