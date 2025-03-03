package utilities

import "strings"

type ErrorPair struct {
	Line int
	Column int
}
func GetErrorIndexPair(source string, index int) *ErrorPair {
	content := source[0 : index + 1]
	lines := strings.Split(content, "\n")

	return &ErrorPair{
		Line: len(lines),
		Column: len(lines[len(lines) - 1]),
	}
}
