package utilities

import "slices"
import "git.jaezmien.com/Jaezmien/fim/twilight/token"
import "git.jaezmien.com/Jaezmien/fim/twilight/queue"

func IsStringNumber(value string) bool {
	if len(value) == 0 {
		return false
	}

	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
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

func ContainsRune(sampleRune rune, runes []rune) bool {
	return slices.Contains(runes, sampleRune)
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
