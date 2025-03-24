package parsers

import (
	"slices"

	"git.jaezmien.com/Jaezmien/fim/luna/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

func IsVariableDeclaration(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"Did", " ", "you", " ", "know", " ", "that"}

	if !utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return 0
	}

	return len(ExpectedTokens)
}

func IsVariableModifier(tokens *queue.Queue[*token.Token]) int {
	ExpectedSingleTokens := []string{"becomes", "become", "became"}
	if slices.Contains(ExpectedSingleTokens, tokens.First().Value.Value) {
		return 1
	}

	ExpectedMultiTokens := []string{"is", " ", "now"}
	if utilities.CheckTokenSequence(tokens, ExpectedMultiTokens) {
		return len(ExpectedMultiTokens)
	}
	return 0
}
