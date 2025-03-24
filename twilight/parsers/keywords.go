package parsers

import (
	"git.jaezmien.com/Jaezmien/fim/luna/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

func IsConstantKeyword(tokens *queue.Queue[*token.Token]) int {
	if tokens.First().Value.Value != "always" {
		return 0
	}

	return 1
}

func IsAndKeyword(tokens *queue.Queue[*token.Token]) int {
	if tokens.First().Value.Value != "and" {
		return 0
	}

	return 1
}

func IsOrKeyword(tokens *queue.Queue[*token.Token]) int {
	if tokens.First().Value.Value != "or" {
		return 0
	}

	return 1
}

func IsOfKeyword(tokens *queue.Queue[*token.Token]) int {
	if tokens.First().Value.Value != "of" {
		return 0
	}

	return 1
}

func IsStatementEndKeyword(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"That", "'", "s", " ", "what", " ", "I", "did"}

	if !utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return 0
	}

	return len(ExpectedTokens)
}
