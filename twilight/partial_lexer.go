package twilight

import (
	"strings"

	"git.jaezmien.com/Jaezmien/fim/twilight/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

var splittable_runes = [...]rune{'.', '!', '?', ':', ',', '(', ')', '"', '\'', ' ', '\t', '\\', '\n'}
func isRuneSplittable(r rune) bool {
	return utilities.ContainsRune(r, splittable_runes[:])
}

func createPartialTokens(source string) *queue.Queue[*token.Token] {
	l := queue.New[*token.Token]()

	start := 0
	for idx, r := range source {
		if isRuneSplittable(r) || len(source) == idx {
			length := idx - start

			if length > 0 {
				token := &token.Token{
					Start:  start,
					Length: length,
					Type:   token.TokenType_Unknown,
					Value: source[start: start + length],
				}
				l.Queue(token)
			}

			token := &token.Token{
				Start:  idx,
				Length: 1,
				Type:   token.TokenType_Unknown,
				Value: source[idx: idx + 1],
			}
			l.Queue(token)

			start = idx + 1
			continue
		}
	}

	if start < len(source) {
		length := len(source) - start

		token := &token.Token{
			Start:  start,
			Length: length,
			Type:   token.TokenType_Unknown,
			Value: source[start: start + length],
		}

		l.Queue(token)
	}

	return l
}

type mergePartialTokensResult = func(tokens *queue.Queue[*token.Token]) int
func processPartialTokens(tokens *queue.Queue[*token.Token], process mergePartialTokensResult) {
	mergeAmount := process(tokens)
	if mergeAmount <= 0 { return }

	token := utilities.MergeTokens(tokens, mergeAmount)
	tokens.QueueFront(token)
}
func mergePartialTokens(tokens *queue.Queue[*token.Token]) *queue.Queue[*token.Token] {
	l := queue.New[*token.Token]()

	newline := false
	for tokens.Len() > 0 {
		if newline && utilities.IsIndentString(tokens.First().Value.Value) {
			tokens.Dequeue()
			continue
		}
		newline = false

		processPartialTokens(tokens, mergeDecimalTokens)
		processPartialTokens(tokens, mergeStringTokens)
		processPartialTokens(tokens, mergeCharacterTokens)
		processPartialTokens(tokens, mergeDelimiters)

		if tokens.First().Value.Length == 1 && tokens.First().Value.Value == "\n" {
			newline = true
		}

		l.Queue(tokens.Dequeue().Value)
	}

	return l
}

func mergeDecimalTokens(tokens *queue.Queue[*token.Token]) int {
	left := tokens.Peek(0)
	if !utilities.IsStringNumber(left.Value.Value) {
		return 0
	}

	decim := tokens.Peek(1)
	if decim.Value.Value != "." {
		return 0
	}

	right := tokens.Peek(2)
	if !utilities.IsStringNumber(right.Value.Value) {
		return 0
	}
	
	return 3
}

func mergeStringTokens(tokens *queue.Queue[*token.Token]) int {
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
func mergeCharacterTokens(tokens *queue.Queue[*token.Token]) int {
	const Delimeter = "'"
	const EscapeToken = "\\"

	current := tokens.Peek(0)
	if current.Value.Value != Delimeter {
		return 0
	}

	mergeAmount := -1

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

	if mergeAmount == -1 { return 0  }

	return mergeAmount
}
func mergeDelimiters(tokens *queue.Queue[*token.Token]) int {
	StartDelimeters := []rune{'('}
	EndDelimeters := []rune{')'}

	current := tokens.Peek(0)

	startDelimeter := strings.IndexAny(current.Value.Value, string(StartDelimeters))
	if startDelimeter == -1 { return 0 }

	endDelimeterIndex := -1
	for idx := 1; idx < tokens.Len(); idx++ {
		token := tokens.Peek(idx)
		if token.Value.Value == "\n" || idx == tokens.Len() {
			return 0 
		}

		endDelimeter := strings.IndexAny(token.Value.Value, string(EndDelimeters))
		if startDelimeter != endDelimeter { continue }
		
		endDelimeterIndex = idx
		break
	}
	
	if endDelimeterIndex == -1 { return 0 }

	return endDelimeterIndex + 1
}
