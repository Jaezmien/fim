package errors

// NOTE: "When should I be using an error, and when should I be using panic()?"
// a. If the issue is caused by the reportee, then it should be an error.
// b. If the issue is caused by a fault in the interpreter itself, then it should be a panic.
// For example:
//		A user mistyped an identifier => Throw an error
//		The interpreter somehow popped a variable stack when it's already empty => Panic

import (
	"fmt"
	"strings"
)

type FiMError struct {
	Message string
}

func (e FiMError) Error() string {
	return e.Message
}

func NewFiMError(msg string) FiMError {
	return FiMError{Message: msg}
}

// An ErrorOrigin contains details about the origin of an error relative
// to the source code.
type ErrorOrigin struct {
	// 1-based line number of the error
	Line int
	// 1-based column number of the error
	Column int

	lineContent string
}

// Create an ErrorOrigin based on a character index.
func GetErrorOrigin(source string, index int) ErrorOrigin {
	content := source[0:min(index+1, len(source))]
	lines := strings.Split(content, "\n")

	return ErrorOrigin{
		Line:   len(lines),
		Column: len(lines[len(lines)-1]),

		lineContent: strings.ReplaceAll(strings.Split(source, "\n")[len(lines) - 1], "\t", " "),
	}
}

type ParseError struct {
	FiMError
	ErrorOrigin
}

func (e ParseError) Error() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("[line %d:%d] %s\n", e.Line, e.Column, e.FiMError.Error()))
	sb.WriteString(e.GetErrorLine())

	return sb.String()
}
func (e ParseError) GetErrorLine() string {
	sb := strings.Builder{}

	trimmedContent := strings.TrimLeft(e.lineContent, "\t ")
	sb.WriteString(fmt.Sprintf("%s\n", trimmedContent))
	sb.WriteString(fmt.Sprintf("%s^", strings.Repeat(" ", e.Column - 1 - (len(e.lineContent) - len(trimmedContent)))))

	return sb.String()
}

func NewParseError(msg string, source string, index int) ParseError {
	return ParseError{
		NewFiMError(msg),
		GetErrorOrigin(source, index),
	}
}
