package utilities

import (
	"strings"

	"git.jaezmien.com/Jaezmien/fim/luna/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
)

func IsStringNumber(value string) bool {
	// XXX: "Why not use strconv?"
	// Well, not exactly handing hexadecimal/octal notation right now.
	// So, we're making exactly sure that what we're getting is an actual decimal.

	if len(value) == 0 {
		return false
	}

	isNegative := strings.HasPrefix(value, "-")
	hasNumber := false

	for idx, r := range value {
		if isNegative && idx == 0 {
			continue
		}

		if r < '0' || r > '9' {
			return false
		}

		hasNumber = true
	}

	return hasNumber
}
func IsIndentCharacter(r rune) bool {
	return r == ' ' || r == '\t'
}
func IsIndentString(s string) bool {
	return len(s) == 1 && IsIndentCharacter(rune(s[0]))
}

func MergeTokens(q *queue.Queue[*token.Token], amount int) *token.Token {
	token := q.Dequeue()

	if token.Value == nil {
		return nil
	}

	amount--

	for amount > 0 {
		t := q.Dequeue()

		if t == nil {
			break
		}

		token.Value.Append(t.Value)

		amount--
	}

	return token.Value
}

func CheckTokenSequence(tokens *queue.Queue[*token.Token], sequence []string) bool {
	if tokens.Len() < len(sequence) {
		return false
	}

	for idx, value := range sequence {
		if idx >= tokens.Len() {
			return false
		}
		if tokens.Peek(idx).Value.Value != value {
			return false
		}
	}

	return true
}
