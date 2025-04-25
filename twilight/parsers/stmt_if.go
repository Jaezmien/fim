package parsers

import (
	"slices"

	"git.jaezmien.com/Jaezmien/fim/luna/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

func CheckIfKeyword(tokens *queue.Queue[*token.Token]) int {
	SingleTokens := []string{"If", "When"}

	if !slices.Contains(SingleTokens, tokens.First().Value.Value) {
		return 0
	}

	return 1
}

func CheckThenKeyword(tokens *queue.Queue[*token.Token]) int {
	SingleTokens := []string{"then"}

	if !slices.Contains(SingleTokens, tokens.First().Value.Value) {
		return 0
	}

	return 1
}

func CheckElseKeyword(tokens *queue.Queue[*token.Token]) int {
	// HACK: This is the stupidest way of bypassing
	// Otherwise If (because we're checking tokens instead of values)
	ExpectedMultiTokens := [][]string{
		{"Otherwise", " ", "if"},
		{"Or", " ", "else", " ", "if"},
		{"Or", " ", "else"},
	}
	for _, sequence := range ExpectedMultiTokens {
		if utilities.CheckTokenSequence(tokens, sequence) {
			return len(sequence)
		}
	}

	SingleTokens := []string{"Otherwise"}

	if slices.Contains(SingleTokens, tokens.First().Value.Value) {
		return 1
	}

	return 0
}

func CheckIfEndKeyword(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"That", "'", "s", " ", "what", " ", "I", " ", "would", " ", "do"}

	if !utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return 0
	}

	return len(ExpectedTokens)
}
