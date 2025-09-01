package twilight

import (
	"slices"
	"strings"

	"git.jaezmien.com/Jaezmien/fim/luna/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

var splittable_runes = [...]rune{'.', '!', '?', ':', ',', '(', ')', '"', '\'', ' ', '\t', '\\', '\n'}

// Split the source string into a queue separated by a splittable rune
// This extra step is needed to handle the multi-word tokens that FiM++ uses.
//
// See:
//
//	var splittable_runes
func createPartialTokens(source string) *queue.Queue[*token.Token] {
	l := queue.New[*token.Token]()

	start := 0
	for idx, r := range source {
		if slices.Contains(splittable_runes[:], r) || len(source) == idx {
			length := idx - start

			if length > 0 {
				token := &token.Token{
					Start:  start,
					Length: length,
					Type:   token.TokenType_Unknown,
					Value:  source[start : start+length],
				}
				l.Queue(token)
			}

			token := &token.Token{
				Start:  idx,
				Length: 1,
				Type:   token.TokenType_Unknown,
				Value:  source[idx : idx+1],
			}
			l.Queue(token)

			start = idx + 1
			continue
		}
	}

	// Fallback for empty sources
	if start < len(source) {
		length := len(source) - start

		token := &token.Token{
			Start:  start,
			Length: length,
			Type:   token.TokenType_Unknown,
			Value:  source[start : start+length],
		}

		l.Queue(token)
	}

	return l
}

// Merge basic tokens that are split across multiple partial tokens.
//
// e.g.: '0' '.' '0' becomes '0.0'
func mergePartialTokens(tokens *queue.Queue[*token.Token]) *queue.Queue[*token.Token] {
	l := queue.New[*token.Token]()

	partialTokensProcessor := []struct {
		process func(tokens *queue.Queue[*token.Token]) int
	}{
		{process: checkDecimalTokens},
		{process: checkStringTokens},
		{process: checkCharacterTokens},
		{process: checkDelimiters},
	}

	newline := false
	for tokens.Len() > 0 {
		if newline && utilities.IsIndentString(tokens.First().Value.Value) {
			tokens.Dequeue()
			continue
		}
		newline = false

		for _, processor := range partialTokensProcessor {
			mergeAmount := processor.process(tokens)
			if mergeAmount <= 0 {
				continue
			}

			token := utilities.MergeTokens(tokens, mergeAmount)
			tokens.QueueFront(token)
		}

		if tokens.First().Value.Length == 1 && tokens.First().Value.Value == "\n" {
			newline = true
		}

		l.Queue(tokens.Dequeue().Value)
	}

	return l
}

// Check for tokens that has the decimal pattern (\d+.\d+).
// If it is a decimal token, return the amount of tokens it consumes.
// Otherwise, return 0
func checkDecimalTokens(tokens *queue.Queue[*token.Token]) int {
	left := tokens.Peek(0)
	if !utilities.IsStringNumber(left.Value.Value) {
		return 0
	}

	decim := tokens.Peek(1)
	if decim.Value.Value != "." {
		return 0
	}

	right := tokens.Peek(2)
	if !utilities.IsStringPositiveNumber(right.Value.Value) {
		return 0
	}

	return 3
}

// Check for tokens that has a string pattern.
// If it is a string token, return the amount of tokens it consumes.
// Otherwise, return 0.
//
// This function also handles escaped quotes.
func checkStringTokens(tokens *queue.Queue[*token.Token]) int {
	const Delimeter = "\""
	const EscapeToken = "\\"

	current := tokens.Peek(0)
	if current.Value.Value != Delimeter {
		return 0
	}

	ignoreNextToken := false
	endIndex := -1
	for idx := 1; idx < tokens.Len(); idx++ {
		token := tokens.Peek(idx)
		if token.Value.Value == "\n" || idx == tokens.Len() {
			return 0
		}

		if ignoreNextToken {
			ignoreNextToken = false
			continue
		}

		if token.Value.Value == EscapeToken {
			ignoreNextToken = true
			continue
		}

		if token.Value.Value == Delimeter {
			endIndex = idx
			break
		}
	}

	if endIndex == -1 {
		return 0
	}

	return endIndex + 1
}

// Check for tokens that has a character pattern.
// If it is a character token, return the amount of tokens it consumes.
// Otherwise, return 0.
//
// This function also supports special characters.
func checkCharacterTokens(tokens *queue.Queue[*token.Token]) int {
	const Delimeter = "'"
	const EscapeToken = "\\"

	current := tokens.Peek(0)
	if current.Value.Value != Delimeter {
		return 0
	}

	mergeAmount := -1

	// FIXME: This doesn't handle Unicode characters well (e.g. �)

	if tokens.Len() >= 4 {
		if tokens.Peek(1).Value.Value == EscapeToken &&
			len(tokens.Peek(2).Value.Value) == 1 &&
			tokens.Peek(3).Value.Value == Delimeter {
			mergeAmount = 4
		}
	}

	if tokens.Len() >= 3 {
		if len(tokens.Peek(1).Value.Value) == 1 &&
			tokens.Peek(2).Value.Value == Delimeter {
			mergeAmount = 3
		}
	}

	if mergeAmount == -1 {
		return 0
	}

	return mergeAmount
}

// Check for tokens that matches against generic delimiter patterns.
// If it matches, return the amount of tokens it consumes.
// Otherwise, return 0.
func checkDelimiters(tokens *queue.Queue[*token.Token]) int {
	StartDelimeters := []rune{'('}
	EndDelimeters := []rune{')'}

	current := tokens.Peek(0)

	startDelimeter := strings.IndexAny(current.Value.Value, string(StartDelimeters))
	if startDelimeter == -1 {
		return 0
	}

	endDelimeterIndex := -1
	for idx := 1; idx < tokens.Len(); idx++ {
		token := tokens.Peek(idx)
		if token.Value.Value == "\n" || idx == tokens.Len() {
			return 0
		}

		endDelimeter := strings.IndexAny(token.Value.Value, string(EndDelimeters))
		if startDelimeter != endDelimeter {
			continue
		}

		endDelimeterIndex = idx
		break
	}

	if endDelimeterIndex == -1 {
		return 0
	}

	return endDelimeterIndex + 1
}
