package parsers

import (
	"git.jaezmien.com/Jaezmien/fim/luna/queue"
	"git.jaezmien.com/Jaezmien/fim/twilight/token"
	"git.jaezmien.com/Jaezmien/fim/twilight/utilities"
)

func CheckUnaryPrefixIncrement(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"There", " ", "was", " ", "one", " ", "more"}
	if !utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return 0
	}

	return len(ExpectedTokens)
}
func CheckUnaryPostfixIncrement(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"got", " ", "one", " ", "more"}
	if !utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return 0
	}

	return len(ExpectedTokens)
}
func CheckUnaryPrefixDecrement(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"There", " ", "was", " ", "one", " ", "less"}
	if !utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return 0
	}

	return len(ExpectedTokens)
}
func CheckUnaryPostfixDecrement(tokens *queue.Queue[*token.Token]) int {
	ExpectedTokens := []string{"got", " ", "one", " ", "less"}
	if !utilities.CheckTokenSequence(tokens, ExpectedTokens) {
		return 0
	}

	return len(ExpectedTokens)
}
